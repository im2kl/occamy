package common

import (
	"sync"

	"github.com/changkun/occamy/lib"
)

// Clipboard defines a generic clipboard structure
type Clipboard struct {
	Mimetype string
	Buffer   []byte
	MaxSize  int
	mu       sync.Mutex
}

// NewClipboard creates a new clipboard having the given initial size.
func NewClipboard(size int) *Clipboard {
	return &Clipboard{MaxSize: size}
}

// Send sends the contents of the clipboard along the given client, splitting
// the contents as necessary.
func (c *Clipboard) Send(client *lib.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client.ForEachUser(sendUserClipboard, c)
}

// Reset clears the clipboard contents and assigns a new mimetype for future data.
func (c *Clipboard) Reset(mimetype string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Buffer = []byte{}
	c.Mimetype = mimetype
}

// Append appends the given data to the current clipboard contents. The data must
// match the mimetype chosen for the clipboard data by
// guac_common_clipboard_reset().
func (c *Clipboard) Append(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Buffer) < c.MaxSize {
		c.Buffer = append(c.Buffer, data[:MaxSize-len(c.Buffer)])
	}
}

func sendUserClipboard(u *lib.User, data interface{}) {
	clipboard := data.(*Clipboard)

	// TODO: send data
}
