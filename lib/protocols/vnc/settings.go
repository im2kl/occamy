package vnc

// DefaultRecordingName is the filename to use
// for the screen recording, if not specified.
const DefaultRecordingName = "recording"

// Settings .
type Settings struct {
	hostname          string
	port              int
	password          string
	encoding          string
	swapRedBlue       bool
	colorDepth        int
	readOnly          bool
	retries           int
	clipboardEncoding string
}

func (s *Settings) Parse() {
}
