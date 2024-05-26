// pkg/utils/appsinkwriter.go
package utils

import (
	"github.com/go-gst/go-gst/gst"
)

type AppSinkWriter struct {
	appsink *gst.Element
}

func NewAppSinkWriter(appsink *gst.Element) *AppSinkWriter {
	return &AppSinkWriter{appsink: appsink}
}

func (w *AppSinkWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (w *AppSinkWriter) Close() error {

	return nil
}
