package dtos

type FrameRequestDTO struct {
	PID     int   `json:"pid"`
	Entries []int `json:"entries"`
}
