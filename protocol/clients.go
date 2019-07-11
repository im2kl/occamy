package protocol

import "github.com/changkun/occamy/lib"

type Client interface {
	JoinHanlder(user *lib.User, params []interface{}) bool
	LeaveHandler(user *lib.User) bool
}
