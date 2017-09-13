package kq

import (
	"log"
	"os/exec"
	"strings"
)

// Stream represents configuration to send a stream to an rtmp server.
type Stream struct {
	Application    string `json:"application"`
	Name           string `json:"name"`
	Host           string `json:"host"`
	VideoEncoding  string `json:"video_encoding"`
	AudioEncoding  string `json:"audio_encoding"`
	OutputEncoding string `json:"output_encoding"`
	Live           bool   `json:"live"`
	Camera         Camera `json:"camera"`
	cmd            *exec.Cmd
}

// NewStream constructs a *Stream with sensible defaults.
func NewStream() *Stream {
	return &Stream{
		Application:    "killerqueen",
		Name:           "",
		Host:           "localhost",
		VideoEncoding:  "libx264",
		AudioEncoding:  "libmp3lame",
		OutputEncoding: "flv",
	}
}

// ffmpeg -i /dev/video1 -c:v libx264 -c:a libmp3lame -f flv rtmp://localhost/killerqueen/capture

// Start begins streaming a camera to an rtmp server.
func (s *Stream) Start() string {
	uri := "rtmp://" + s.Host + "/" + s.Application + "/" + s.Name
	go func() {
		defer func() { s.Live = false }()
		s.Live = true
		cmd := exec.Command("ffmpeg", "-i", s.Camera.Device, "-c:v", s.VideoEncoding, "-c:a", s.AudioEncoding, "-f", s.OutputEncoding, uri)
		s.cmd = cmd
		err := cmd.Run()
		if err != nil && !strings.Contains(err.Error(), "signal: killed") {
			log.Printf("Failed to start stream: %s\n", err)
		}
	}()
	return uri
}

// Stop halts streaming a camera to an rtmp server.
func (s *Stream) Stop() error {
	return s.cmd.Process.Kill()
}

type Camera struct {
	Serial string `json:"serial"`
	Model  string `json:"model"`
	Vendor string `json:"vendor"`
	Device string `json:"device"`
}
