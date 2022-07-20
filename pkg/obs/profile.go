package obs

import (
	"fmt"
)

var (
	ProfileTemplate = `[General]
Name=zixun

[Video]
BaseCX=%d
BaseCY=%d
OutputCX=%d
OutputCY=%d
FPSCommon=25 PAL
ScaleType=bicubic

[Panels]
CookieId=CCE0DA3FA51A8247

[Audio]
MonitoringDeviceName=%s
MonitoringDeviceId=%s

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
RescaleRes=%s
RecRescaleRes=%s
FFRescaleRes=%s

[Stats]
geometry=AdnQywADAAAAAAIZAAABPgAABToAAAJ1AAACGgAAAV0AAAU5AAACdAAAAAAAAAAAB1QAAAIaAAABXQAABTkAAAJ0

[WebsocketAPI]
ServerEnabled=true
ServerPort=%d
LockToIPv4=false
DebugEnabled=true
AlertsEnabled=false
AuthRequired=true
AuthSecret=%s
AuthSalt=%s
`
)

func GetProfileContent(config *ProfileConfig) string {
	return fmt.Sprintf(ProfileTemplate,
		config.Size.X, config.Size.Y,
		config.Size.X, config.Size.Y,
		config.Mic.Name, config.Mic.ID,
		config.Size.GetSize(), config.Size.GetSize(), config.Size.GetSize(),
		config.Port, config.Secret, config.Salt,
	)
}

type ProfileConfig struct {
	Size     *VideoSize
	Salt     string
	Secret   string
	Port     uint
	Password string
	Mic      *VirtualMicInfo
}

func NewProfileConfig(size *VideoSize, salt string, secret string, port uint, mic *VirtualMicInfo) *ProfileConfig {
	return &ProfileConfig{
		Size:   size,
		Salt:   salt,
		Secret: secret,
		Port:   port,
		Mic:    mic,
	}
}
