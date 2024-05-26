package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"sync"
)

const (
	maxBufferSize = 100000000
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
			excessBytes := len(p) + pb.buffer.Len() - maxBufferSize
			if excessBytes >= pb.buffer.Len() {
				pb.buffer.Reset()
			} else {
				pb.buffer.Next(excessBytes)
			}
		}
		return pb.buffer.Write(p)
	}

	// if pb.w == nil {
	// 	if len(p)+pb.buffer.Len() > maxBufferSize {
	// 		pb.buffer.Reset()
	// 		if pb.onBufferReset != nil {
	// 			if err := pb.onBufferReset(); err != nil {
	// 				return 0, err
	// 			}
	// 		}
	// 		return 0, errors.ErrPrerollBufferReset
	// 	}
	// 	return pb.buffer.Write(p)
	// }

	n, err := pb.w.Write(p)
	if err == io.ErrClosedPipe {
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

func (pb *PrerollBuffer) Read(p []byte) (int, error) {
	pb.lock.Lock()
	defer pb.lock.Unlock()

	return pb.buffer.Read(p)
}

func NewGuid(prefix string) string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	id := fmt.Sprintf("%s%x-%x-%x-%x-%x", prefix, bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
	return id
}
