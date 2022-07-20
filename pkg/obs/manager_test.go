package obs

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

var fackNow = func() time.Time {
	return time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
}

func TestGenerateSalts(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "seed 1974-05-19 01:02:03.00000004 +0000 UTC",
			want: "NEhPY1pqYjIzaVlYNjlvdGc0b3VQVW9wR3NkYTR3R3U=",
		},
	}

	now = fackNow

	Convey("GenerateSalt", t, func() {
		for _, tt := range tests {
			Convey(fmt.Sprintf("GenerateSalt %s", tt.name), func() {
				salt := GenerateSalt()
				So(salt, ShouldResemble, tt.want)
			})
		}
	})
}

func TestGenerateSecret(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "obs password 4444",
			args: args{
				password: "4444",
				salt:     "17wUISsvrMkgxckSrCu5G0KZJANTSaf+FS2/qIPh6Zk=",
			},
			want: "Wn78WPhpk7mSmH0XU33gyCm4NKKtD0LSWt4Ek6GhZow=",
		},
		{
			name: "password zixun seed 1974-05-19 01:02:03.00000004 +0000 UTC",
			args: args{
				password: "zixun",
				salt:     "NEhPY1pqYjIzaVlYNjlvdGc0b3VQVW9wR3NkYTR3R3U=",
			},
			want: "MNgmOPC8L08AHnHmLgnPjj1t564ecNtFTeeS0bSdpck=",
		},
	}
	Convey("GenerateSecret", t, func() {
		for _, tt := range tests {
			Convey(fmt.Sprintf("GenerateSecret(%v)", tt.name), func() {
				got := GenerateSecret(tt.args.password, tt.args.salt)
				So(got, ShouldEqual, tt.want)
			})
		}
	})
}
