package options

import (
	"bytes"
	"io"
)

// IOStreams provides the standard names for iostreams.  This is useful for embedding and for unit testing.
// Inconsistent and different names make it hard to read and review code.
type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

// NewTestIOStreamsDiscard returns a valid IOStreams that just discards.
func NewTestIOStreams() (IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}, in, out, errOut
}

func NewTestIOStreamsDiscard() IOStreams {
	in := &bytes.Buffer{}

	return IOStreams{
		In:     in,
		Out:    io.Discard,
		ErrOut: io.Discard,
	}
}
