package vnc

/*
#cgo LDFLAGS: -L/usr/local/lib -lrbf

#include <rfb/rfbclient.h>
#include <rfb/rfbproto.h>

char* VNC_CLIENT_KEY = "OccamyVNC"
char* vnc_get_password(rfbClient* client) {
    guac_client* gc = rfbClientGetClientData(client, VNC_CLIENT_KEY);
    return ((guac_vnc_client*) gc->data)->settings->password;
}
*/
import "C"

// RfbClient ...
type RfbClient struct {
	client *C.rfbClient
}

const vncClientKey = "OccamyVNC"

// GetPassword ...
func (c *RfbClient) GetPassword() string {
	return C.GoString(C.vnc_get_password(c.client))
}
