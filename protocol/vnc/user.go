package vnc

import "github.com/changkun/occamy/lib"

// Client ...
type Client struct {
}

// JoinHanlder ...
func (c *Client) JoinHanlder(u *lib.User, params []interface{}) bool {

}

// LeaveHandler ...
func (c *Client) LeaveHandler(u *lib.User) bool {
	if u.Client.Display {
		// Update shared cursor state
		C.guac_common_cursor_remove_user(u.Client.display.cursor, u)
	}
}
