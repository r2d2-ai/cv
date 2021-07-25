package ipcam

import (
	"context"
	"net/url"
	"time"

	"github.com/r2d2-ai/core/data/metadata"
	"github.com/r2d2-ai/core/support/log"
	"github.com/r2d2-ai/core/trigger"
	"gocv.io/x/gocv"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	return &Trigger{settings: s}, nil
}

type Trigger struct {
	// metadata       *trigger.Metadata
	settings       *Settings
	cameraHandlers []*CameraHandler
	logger         log.Logger
}

type CameraHandler struct {
	handler  trigger.Handler
	settings *HandlerSettings
	cap      *gocv.VideoCapture
	done     chan bool
	logger   log.Logger
	id       string
}

func (t *Trigger) Initialize(ctx trigger.InitContext) error {

	t.logger = log.ChildLogger(ctx.Logger(), "ipcam-feed-listener")
	for _, handler := range ctx.GetHandlers() {
		s := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), s, true)
		if err != nil {
			return err
		}
		camHnd := &CameraHandler{}
		camHnd.settings = s
		camHnd.handler = handler
		camHnd.logger = t.logger
		// camHnd.done = make(chan bool)
		// camHnd.done <- false
		t.cameraHandlers = append(t.cameraHandlers, camHnd)
	}
	return nil
}

// Start implements util.Managed.Start
func (t *Trigger) Start() error {

	for _, camHdl := range t.cameraHandlers {
		err := camHdl.startStream()
		if err != nil {
			return err
		}
		go camHdl.run()
	}
	return nil
}

func (camHnd *CameraHandler) startStream() error {
	//protocol := camHnd.settings.Protocol

	host := camHnd.settings.Host
	user := camHnd.settings.User
	password := camHnd.settings.Password
	videoUri := camHnd.settings.VideoURI
	camHnd.logger.Infof("Start IP Cam %v stream", host)
	cap, err := gocv.OpenVideoCapture("rtsp://" + url.QueryEscape(user) + ":" + url.QueryEscape(password) + "@" + host + "/" + url.QueryEscape(videoUri))

	if err != nil {
		return err
	}

	camHnd.cap = cap

	camHnd.id = "SomeId"
	return nil
}

type FPSCounter []int64

func (counter *FPSCounter) FPS() float64 {
	var total int64 = 0
	slice := *counter
	if len(slice) > 1000 {
		*counter = slice[len(*counter)-1000:]
	}

	for _, val := range *counter {
		total += val
	}

	return float64(total) / float64(len(*counter))
}

func (camHnd *CameraHandler) run() {
	var err error
	var counter *FPSCounter = new(FPSCounter)
	*counter = make([]int64, 1100)

	img := gocv.NewMat()
	host := camHnd.settings.Host

	camHnd.logger.Infof("Running IP Cam %v stream", host)
	for {
		start := time.Now()
		// if <-camHnd.done {
		// 	break
		// }
		camHnd.cap.Read(&img)
		if img.Empty() {
			camHnd.logger.Errorf("Received blank frame IP Cam %v", host)
			continue
		}

		duration := time.Since(start).Milliseconds()
		*counter = append(*counter, duration)
		image := img //.ToBytes()
		output := &Output{}
		output.Image = image
		output.FPS = 1000. / counter.FPS()
		_, err = camHnd.handler.Handle(context.Background(), output)

		if err != nil {
			camHnd.logger.Errorf("Failed to process frame for IP Cam %v", host)
		}
	}
}

// Stop implements util.Managed.Stop
func (t *Trigger) Stop() error {
	for _, camHdl := range t.cameraHandlers {
		camHdl.done <- true
	}
	return nil
}
