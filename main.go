package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"git.ziniao.com/environment/client/video-check/pkg/obs"
	"github.com/andreykaipov/goobs/api/events"
)

// https://app.chime.aws/meetings/3095770894
// video-check.exe -url "https://vdo.ninja/?view=ganjl" -port 4444 -mid "{0.0.0.00000000}.{951fe9f9-630c-4984-997f-ca1850d38080}" -mname "Line 1 (Virtual Audio Cable)"
// video-check.exe -url "https://vdo.ninja/?view=ganjl" -port 4444 -mid {0.0.0.00000000}.{e31d0034-6676-4f21-9170-dbaa12b1ee03}" -mname "Line 2 (Virtual Audio Cable)"
var terminationSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

var url = flag.String("url", "https://vdo.ninja/?view=ganjl", "URL to video")
var port = flag.Int("port", 4444, "Port to listen")
var mName = flag.String("mname", "Line 1 (Virtual Audio Cable)", "Mic ID \"Line 1 (Virtual Audio Cable)\" \"Line 2 (Virtual Audio Cable)\"")
var mID = flag.String("mid", "{0.0.0.00000000}.{951fe9f9-630c-4984-997f-ca1850d38080}", "Mic Name \"{0.0.0.00000000}.{951fe9f9-630c-4984-997f-ca1850d38080}\" {0.0.0.00000000}.{e31d0034-6676-4f21-9170-dbaa12b1ee03}")

func main() {

	flag.Parse()
	o := obs.NewObsManager(&obs.VideoSource{
		Size: &obs.VideoSize{
			X: 1876,
			Y: 1017,
		},
		Url: *url,
	}, uint(*port), "4444", *mID, *mName)
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
