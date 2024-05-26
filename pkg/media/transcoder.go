package media

import (
	"fmt"

	"github.com/Harshitk-cp/rtmp_server/pkg/rtmp"
	"github.com/Harshitk-cp/rtmp_server/pkg/utils"
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

func (t *Transcoder) Start(rtmpServer *rtmp.RTMPServer, resourceID string, mediaBuffer *utils.PrerollBuffer) error {
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

	flvdemux, err := gst.NewElementWithName("flvdemux", "flv-demux")
	if err != nil {
		return fmt.Errorf("failed to create flvdemux element: %w", err)
	}
	t.gstPipeline.Add(flvdemux)

	rtmpsrc.Link(flvdemux)

	rtmpsrc.Link(rtph264pay)
	rtph264pay.Link(aacpay)

	appName := "live"
	if err := rtmpsrc.SetProperty("location", fmt.Sprintf("rtmp://localhost:1935/%s/%s", appName, "ky")); err != nil {
		return fmt.Errorf("failed to set rtmpsrc location: %w", err)
	}

	muxer, err := gst.NewElementWithName("qtmux", "muxer")
	if err != nil {
		return fmt.Errorf("failed to create muxer element: %w", err)
	}
	t.muxer = muxer
	t.gstPipeline.Add(muxer)

	aacpay.Link(muxer)
	rtph264pay.Link(muxer)

	appsink, err := gst.NewElementWithName("appsink", "app-sink")
	if err != nil {
		return fmt.Errorf("failed to create appsink element: %w", err)
	}
	t.gstPipeline.Add(appsink)

	appsink.SetProperty("sync", false)
	appsink.SetProperty("emit-signals", true)
	appsink.SetProperty("drop", false)

	appSinkWriter := utils.NewAppSinkWriter(appsink)

	if err := mediaBuffer.SetWriter(appSinkWriter); err != nil {
		return fmt.Errorf("failed to set PrerollBuffer writer: %w", err)
	}

	muxer.Link(appsink)

	err = t.gstPipeline.SetState(gst.StatePlaying)
	if err != nil {
		t.gstPipeline.SetState(gst.StateNull)
		return fmt.Errorf("failed to start GStreamer pipeline: %w", err)
	}

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
