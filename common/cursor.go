package common

/*
#include <cairo/cairo.h>
*/
import "C"
import (
	"time"

	"github.com/changkun/occamy/lib"
)

// Layer ...
type Layer int

// Cursor ...
type Cursor struct {
	client      *lib.Client
	buffer      *Layer
	width       int
	height      int
	imageBuffer [CursorDefaultSize]byte

	surface    *C.cairo_surface_t
	hotspotX   int
	hotspotY   int
	user       *lib.User
	x          int
	y          int
	buttonMask lib.ClientMouse
	timestamp  time.Time
}

// NewCursor ...
func NewCursor(c *lib.Client) *Cursor {
	return &Cursor{
		client:    c,
		timestamp: time.Now().UTC(),
	}
}

func (c *Cursor) Dup() {

}

func (c *Cursor) Update() {

}

func (c *Cursor) SetARGB() {

}

func (c *Cursor) SetSurface() {

}

func (c *Cursor) SetPointer() {

}

func (c *Cursor) SetDot() {

}

func (c *Cursor) SetIBar() {

}

func (c *Cursor) SetBlank() {

}

func (c *Cursor) RemoveUser() {

}
