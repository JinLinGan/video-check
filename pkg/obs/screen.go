package obs

import "fmt"

type VideoSize struct {
	X, Y uint16
}

func NewScenseSize(x, y uint16) *VideoSize {
	return &VideoSize{
		X: x,
		Y: y,
	}
}

func (s *VideoSize) GetSize() string {
	return fmt.Sprintf("%dx%d", s.X, s.Y)
}
