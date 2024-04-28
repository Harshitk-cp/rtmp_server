package errors

import "errors"

var (
	ErrNoConfig                     = errors.New("missing config")
	ErrInvalidAudioOptions          = errors.New("invalid audio options")
	ErrInvalidVideoOptions          = errors.New("invalid video options")
	ErrInvalidAudioPreset           = errors.New("invalid audio encoding preset")
	ErrInvalidVideoPreset           = errors.New("invalid video encoding preset")
	ErrSourceNotReady               = errors.New("source encoder not ready")
	ErrUnsupportedDecodeFormat      = errors.New("unsupported format for the source media")
	ErrUnsupportedEncodeFormat      = errors.New("unsupported mime type for encoder")
	ErrUnsupportedURLFormat         = errors.New("unsupported URL type")
	ErrDuplicateTrack               = errors.New("more than 1 track with given media kind")
	ErrUnableToAddPad               = errors.New("could not add pads to bin")
	ErrMissingResourceId            = errors.New("missing resource ID")
	ErrInvalidRelayToken            = errors.New("invalid token")
	ErrIngressNotFound              = errors.New("ingress not found")
	ErrServerCapacityExceeded       = errors.New("server capacity exceeded")
	ErrServerShuttingDown           = errors.New("server shutting down")
	ErrIngressClosing               = errors.New("ingress closing")
	ErrMissingStreamKey             = errors.New("missing stream key")
	ErrPrerollBufferReset           = errors.New("preroll buffer reset")
	ErrInvalidSimulcast             = errors.New("invalid simulcast configuration")
	ErrSimulcastTranscode           = errors.New("simulcast is not supported when transcoding")
	ErrRoomDisconnected             = errors.New("room disconnected")
	ErrRoomDisconnectedUnexpectedly = errors.New("room disconnected unexpectedly")
)