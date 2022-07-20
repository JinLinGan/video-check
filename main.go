package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.ziniao.com/environment/client/video-check/pkg/obs"
	"github.com/andreykaipov/goobs/api/events"
)

var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

func main() {
	o := obs.NewObsManager(&obs.VideoSource{
		Size: &obs.VideoSize{
			X: 1876,
			Y: 1017,
		},
		Url: "https://vdo.ganjl.top/?view=uWMvcn8",
	}, 4444, "4444")
	if err := o.StartOBS(); err != nil {
		log.Panic(err)
	}

	client := o.Client
	go func() {
		for event := range client.IncomingEvents {
			switch e := event.(type) {
			case *events.SourceVolumeChanged:
				fmt.Printf("Volume changed for %-25q: %f\n", e.SourceName, e.Volume)
			default:
				log.Printf("Unhandled event: %#v", e.GetUpdateType())
			}
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, terminationSignals...)
	<-ch
	stop(o)

}

func stop(o *obs.OBSManager) {
	err := o.StopOBS()
	if err != nil {
		log.Printf("Close OBS fail: %s \n", err)
	}
}
