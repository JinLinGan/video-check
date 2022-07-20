package obs

import "testing"

// import (
// 	. "github.com/smartystreets/goconvey/convey"
// )

// // TestGetProfileConfig
// //  @param t
// func TestGetProfileConfig(t *testing.T) {

// 	tb := struct {
// 		x        uint16
// 		y        uint16
// 		template string
// 		out      string
// 	}{
// 		x:        0,
// 		y:        0,
// 		template: "",
// 		out:      "",
// 	}
// 	Convey("get Profile config", t, func() {
// 		p, err := GetProfileConfig(NewScenseSize(100, 200))
// 		So(err, ShouldBeNil)
// 	})

// }

func TestGetProfileContent(t *testing.T) {
	type args struct {
		config *ProfileConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				config: &ProfileConfig{
					Size: &VideoSize{
						X: 100,
						Y: 200,
					},
					Salt:   "aaa",
					Secret: "bbb",
					Port:   1111,
					Mic: &VirtualMicInfo{
						Name: "mmmm",
						ID:   "ididid",
					},
				},
			},
			want: `[General]
Name=zixun

[Video]
BaseCX=100
BaseCY=200
OutputCX=100
OutputCY=200
FPSCommon=25 PAL
ScaleType=bicubic

[Panels]
CookieId=CCE0DA3FA51A8247

[Audio]
MonitoringDeviceName=mmmm
MonitoringDeviceId=ididid

[Output]
Mode=Simple

[SimpleOutput]
VBitrate=300
UseAdvanced=true
Preset=ultrafast

[AdvOut]
TrackIndex=1
RecType=Standard
RecTracks=1
FLVTrack=1
FFOutputToFile=true
FFFormat=
FFFormatMimeType=
FFVEncoderId=0
FFVEncoder=
FFAEncoderId=0
FFAEncoder=
FFAudioMixes=1
VodTrackIndex=2
RescaleRes=100x200
RecRescaleRes=100x200
FFRescaleRes=100x200

[Stats]
geometry=AdnQywADAAAAAAIZAAABPgAABToAAAJ1AAACGgAAAV0AAAU5AAACdAAAAAAAAAAAB1QAAAIaAAABXQAABTkAAAJ0

[WebsocketAPI]
ServerEnabled=true
ServerPort=1111
LockToIPv4=false
DebugEnabled=true
AlertsEnabled=false
AuthRequired=true
AuthSecret=bbb
AuthSalt=aaa
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetProfileContent(tt.args.config)
			if got != tt.want {
				t.Errorf("GetProfileContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
