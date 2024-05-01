package utils

import (
	"bytes"
	"io"
	"sync"

	"github.com/Harshitk-cp/rtmp_server/pkg/errors"
)

const (
	maxBufferSize = 10000000
)

type PrerollBuffer struct {
	lock   sync.Mutex
	buffer *bytes.Buffer
	w      io.WriteCloser

	onBufferReset func() error
}

func NewPrerollBuffer(onBufferReset func() error) *PrerollBuffer {
	return &PrerollBuffer{
		buffer:        &bytes.Buffer{},
		onBufferReset: onBufferReset,
	}
}

func (pb *PrerollBuffer) SetWriter(w io.WriteCloser) error {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	pb.w = w
	if pb.w == nil {
		// Send preroll buffer reset event
		pb.buffer.Reset()
		if pb.onBufferReset != nil {
			if err := pb.onBufferReset(); err != nil {
				return err
			}
		}
	} else {
		_, err := io.Copy(pb.w, pb.buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pb *PrerollBuffer) Write(p []byte) (int, error) {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.w == nil {
		if len(p)+pb.buffer.Len() > maxBufferSize {
			// We would overflow the max allowed buffer size. Reset th buffer state
			pb.buffer.Reset()
			if pb.onBufferReset != nil {
				if err := pb.onBufferReset(); err != nil {
					return 0, err
				}
			}
			return 0, errors.ErrPrerollBufferReset
		}
		return pb.buffer.Write(p)
	}

	n, err := pb.w.Write(p)
	if err == io.ErrClosedPipe {
		// Do not return errors caused by a consuming pipe getting closed
		err = nil
	}
	return n, err
}

func (pb *PrerollBuffer) Close() error {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	if pb.w != nil {
		return pb.w.Close()
	}

	return nil
}
