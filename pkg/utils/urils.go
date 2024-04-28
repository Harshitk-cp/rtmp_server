package utils

import "io"

type PrerollBuffer struct {
	writer io.WriteCloser
	// Implement other necessary methods and fields if needed
}

func NewPrerollBuffer(resetCallback func() error) *PrerollBuffer {
	return &PrerollBuffer{}
}

func (b *PrerollBuffer) SetWriter(w io.WriteCloser) error {
	b.writer = w
	return nil
}

func (b *PrerollBuffer) Close() {
	// Implement close logic if needed
}
