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

	rtmpsrc, err := gst.NewElementWithName("rtmpsrc", "rtmpsrc")
	if err != nil {
		return fmt.Errorf("failed to create rtmpsrc element: %w", err)
	}
	t.gstPipeline.Add(rtmpsrc)
	rtph264pay, err := gst.NewElementWithName("rtph264pay", "rtph264pay")
	if err != nil {
		return fmt.Errorf("failed to create rtph264pay element: %w", err)
	}
	t.gstPipeline.Add(rtph264pay)

	aacpay, err := gst.NewElementWithName("aacpay", "aacpay")
	if err != nil {
		return fmt.Errorf("failed to create aacpay element: %w", err)
	}
	t.gstPipeline.Add(aacpay)
	rtmpsrc.Link(rtph264pay)
	rtph264pay.Link(aacpay)

	if err := rtmpsrc.SetProperty("location", fmt.Sprintf("rtmp://localhost:1935/live/%s", resourceID)); err != nil {
		return fmt.Errorf("failed to set rtmpsrc location: %w", err)
	}

	sink, err := gst.NewElementWithName("fakesink", "fakesink")
	if err != nil {
		return fmt.Errorf("failed to create fakesink element: %w", err)
	}
	t.gstPipeline.Add(sink)
	aacpay.Link(sink)

	err = t.gstPipeline.SetState(gst.StatePlaying)
	if err != nil {
		return fmt.Errorf("failed to start GStreamer pipeline: %w", err)
	}

	go func() {
		bus := t.gstPipeline.GetBus()
		for {
			msg := bus.Pop()
			if msg == nil {
				break
			}

			switch msg.Type() {
			case gst.MessageEOS:
				logrus.Info("Transcoder pipeline reached EOS (End of Stream)")
				t.gstPipeline.SetState(gst.StateNull)
				return
			default:
			}
		}
	}()

	return nil
}

func (t *Transcoder) Stop() {

	if t.pipelinePaused || t.gstPipeline == nil {
		return
	}

	if t.gstPipeline == nil || t.gstPipeline.GetCurrentState() == gst.StateNull {
		return
	}

	err := t.gstPipeline.SetState(gst.StateNull)
	if err != nil {
		logrus.Error("Error stopping GStreamer pipeline: %w", err)
	}

	t.pipelinePaused = true
}
