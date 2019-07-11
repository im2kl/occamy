// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lib

/*
#cgo LDFLAGS: -L/usr/local/lib -lguac
#include <stdlib.h>
#include "../guacamole/src/libguac/guacamole/user.h"
#include "../guacamole/src/libguac/guacamole/client.h"

#include "../guacamole/src/libguac/guacamole/parser.h"
 const char *mimetypes[] = {"", NULL};
 void set_user_info(guac_user* user) {
	user->info.optimal_width = 1024;
	user->info.optimal_height = 768;
	user->info.optimal_resolution = 96;
	user->info.audio_mimetypes = (const char**) mimetypes;
	user->info.video_mimetypes = (const char**) mimetypes;
	user->info.image_mimetypes = (const char**) mimetypes;
}
int get_args_length(const char** args) {
	int i = 0;
	int argc = 0;
	for (i=0; args[i] != NULL; i++) {
		argc++;
	}
	return argc;
}
 static char** makeCharArray(int size) {
	return calloc(sizeof(char*), size);
}
 static void setArrayString(char **a, char *s, int n) {
	a[n] = s;
}
 static void freeCharArray(char **a, int size) {
	int i;
	for (i = 0; i < size; i++)
		free(a[i]);
	free(a);
}
*/
import "C"
import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/changkun/occamy/config"
	"github.com/sirupsen/logrus"
)

// UserCallback ...
type UserCallback func(u *User, data interface{}) interface{}

// User is the representation of a physical connection within a larger logical connection
// which may be shared. Logical connections are represented by guac_client.
type User struct {
	guacUser *C.struct_guac_user
	Client   *Client
	userid   string
	owner    bool
	active   int32 // atomic
	params   UserParams
	once     sync.Once

	Next *User // points to next connected user
}

// UserParams contains all parameters for establish a connection
type UserParams struct {
	Host     string
	Port     string
	Username string
	Password string
}

// NewUser creates a user and associate the user with any specific client
func NewUser(s *Socket, c *Client, owner bool, jwt *config.JWT) (*User, error) {
	u := C.guac_user_alloc()
	if u == nil {
		return nil, errors.New(errorStatus())
	}
	u.socket = s.guacSocket
	if owner {
		u.owner = C.int(1)
	} else {
		u.owner = C.int(0)
	}

	host, port, err := net.SplitHostPort(jwt.Host)
	if err != nil {
		return nil, err
	}
	return &User{
		guacUser: u,
		client:   c,
		owner:    owner,
		active:   1,
		params: UserParams{
			Host:     host,
			Port:     port,
			Username: jwt.Username,
			Password: jwt.Password,
		},
	}, nil
}

func (u *User) IsActive() bool {
	if atomic.LoadInt32(&u.active) == 1 {
		return true
	}
	return false
}

func (u *User) Stop() {
	C.guac_user_stop(u.guacUser)
	atomic.StoreInt32(&u.active, 0)
}

// Close frees the user and detach the association to the attached client
func (u *User) Close() {
	u.once.Do(func() {
		C.guac_user_free(u.guacUser)
	})
}

const usecTimeout time.Duration = 15 * time.Millisecond

// HandleConnection handles all I/O for the portion of a user's Guacamole connection
// following the initial "select" instruction, including the rest of the handshake.
// The handshake-related properties of the given guac_user are automatically
// populated, and HandleConnection() is invoked for all instructions received after
// the handshake has completed. This function blocks until the connection/user is aborted
// or the user disconnects.
func (u *User) HandleConnection() error {
	if int(C.guac_user_handle_connection(u.guacUser, C.int(int(usecTimeout)))) != 0 {
		return errors.New(errorStatus())
	}
	return nil
}

// HandleConnectionWithHandshake ...
func (u *User) HandleConnectionWithHandshake() error {
	// general args
	C.set_user_info(u.guacUser)

	// client args
	length := int(C.get_args_length(u.Client.guacClient.args))
	tmpslice := (*[1 << 30]*C.char)(unsafe.Pointer(u.Client.guacClient.args))[:length:length]
	args := make([]string, length)
	for i, s := range tmpslice {
		args[i] = C.GoString(s)
	}
	for i := range args {
		switch args[i] {
		case "hostname":
			args[i] = u.params.Host
		case "port":
			args[i] = u.params.Port
		case "username":
			args[i] = u.params.Username
		case "password":
			args[i] = u.params.Password
		default:
			args[i] = ""
		}
	}
	cargs := C.makeCharArray(C.int(len(args)))
	defer C.freeCharArray(cargs, C.int(len(args)))

	// FIXME: why NULL in cargs?
	if int(C.guac_client_add_user(u.Client.guacClient, u.guacUser, C.int(len(args)), cargs)) != 0 {
		logrus.Errorf("User %s could NOT join connection %s", u.guacUser.user_id, u.Client.guacClient.connected_users)
	}
	u.start() // block here
	return nil
}

func (u *User) start() {

	parser := C.guac_parser_alloc()
	defer C.guac_parser_free(parser)

	logrus.Info("user_start: ready for parser read loop!!!!")

	for u.Client.Running() && u.IsActive() {

		// 1. read instruction, stop on error
		if int(C.guac_parser_read(parser, u.guacUser.socket, C.int(int(usecTimeout)))) != 0 {
			errString := errorStatus()
			if errString != statusString[statusTimeout] {
				logrus.Errorf("User is not responding.")
			} else if errString != statusString[statusClosed] {
				logrus.Errorf("Guacamole connection failure")
				u.Stop()
			}

			return
		}

		// 2. reset guac_error and guac_error_message (user/client handlers are not guaranteed to set these)

		// 3. call handler, stop on error
		if int(C.guac_user_handle_instruction(u.guacUser, parser.opcode, parser.argc, parser.argv)) < 0 {
			logrus.Errorf("User connection aborted.")
			u.Stop()
			return
		}
	}

	C.guac_client_remove_user(u.Client.guacClient, u.guacUser)
	logrus.Infof("User %s disconnected (%i users remain)", u.guacUser.user_id, u.Client.guacClient.connected_users)
}
