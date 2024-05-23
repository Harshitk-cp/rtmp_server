package media

import (
	"fmt"

	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/go-gst/go-gst/gst"
	"github.com/sirupsen/logrus"
)

type Transcoder struct {
	rtmpServer     *rtmp.RTMPServer
	resourceID     string
	gstPipeline    *gst.Pipeline
	audioQueue     *gst.Element
	videoQueue     *gst.Element
	audioEncoder   *gst.Element
	videoEncoder   *gst.Element
	muxer          *gst.Element
	sink           *gst.Element
	pipelinePaused bool
}

func NewTranscoder() *Transcoder {
	return &Transcoder{}
}

func (t *Transcoder) Start(rtmpServer *rtmp.RTMPServer, resourceID string) error {
	t.rtmpServer = rtmpServer
	t.resourceID = resourceID

	pipeline, err := gst.NewPipeline("transcoder-pipeline")
	if err != nil {
		return fmt.Errorf("failed to create GStreamer pipeline: %w", err)
	}
	t.gstPipeline = pipeline

	if !t.pipelinePaused {
		err := t.gstPipeline.SetState(gst.StatePlaying)
		if err != nil {
			return fmt.Errorf("failed to start GStreamer pipeline: %w", err)
		}
		t.pipelinePaused = false
	}

	audioQueue, err := gst.NewElement("queue")
	if err != nil {
		return fmt.Errorf("failed to create audio queue element: %w", err)
	}
	t.audioQueue = audioQueue
	t.gstPipeline.Add(audioQueue)

	videoQueue, err := gst.NewElement("queue")
	if err != nil {
		return fmt.Errorf("failed to create video queue element: %w", err)
	}
	t.videoQueue = videoQueue
	t.gstPipeline.Add(videoQueue)

	rtmpsrc, err := gst.NewElement("rtmpsrc")
	if err != nil {
		return fmt.Errorf("failed to create rtmpsrc element: %w", err)
	}
	t.gstPipeline.Add(rtmpsrc)
	rtph264pay, err := gst.NewElementWithName("rtph264pay", "rtph264pay")
	if err != nil {
		return fmt.Errorf("failed to create rtph264pay element: %w", err)
	}
	t.gstPipeline.Add(rtph264pay)

	aacpay, err := gst.NewElement("faac")
	if err != nil {
		return fmt.Errorf("failed to create aacpay element: %w", err)
	}
	t.gstPipeline.Add(aacpay)
	rtmpsrc.Link(rtph264pay)
	rtph264pay.Link(aacpay)

	appName := "live"

	if err := rtmpsrc.SetProperty("location", fmt.Sprintf("rtmp://localhost:1935/%s/%s", appName, resourceID)); err != nil {
		return fmt.Errorf("failed to set rtmpsrc location: %w", err)
	}

	sink, err := gst.NewElement("fakesink")
	if err != nil {
		return fmt.Errorf("failed to create fakesink element: %w", err)
	}

	t.gstPipeline.Add(sink)
	aacpay.Link(sink)

	err = t.gstPipeline.SetState(gst.StatePlaying)
	if err != nil {
		t.gstPipeline.SetState(gst.StateNull)
		return fmt.Errorf("failed to start GStreamer pipeline: %w", err)
	}

	// go func() {
	// 	bus := t.gstPipeline.GetBus()
	// 	for {
	// 		msg := bus.Pop()
	// 		if msg == nil {
	// 			break
	// 		}

	// 		switch msg.Type() {
	// 		case gst.MessageEOS:
	// 			logrus.Info("Transcoder pipeline reached EOS (End of Stream)")
	// 			t.gstPipeline.SetState(gst.StateNull)
	// 			return
	// 		default:
	// 		}
	// 	}
	// }()

	return nil
}

func (t *Transcoder) Stop() {

	if t.gstPipeline == nil {
		return
	}

	err := t.gstPipeline.SetState(gst.StateNull)
	if err != nil {
		logrus.Errorf("Error stopping GStreamer pipeline: %v", err)
	}

	t.gstPipeline.Unref()
	t.gstPipeline = nil
	t.audioQueue = nil
	t.videoQueue = nil
	t.audioEncoder = nil
	t.videoEncoder = nil
	t.muxer = nil
	t.sink = nil

	t.pipelinePaused = false
}
