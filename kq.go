package kq

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// starts playing stream quickly:
// ffplay -probesize 32 -sync ext rtmp://localhost/killerqueen/blue

// Stream represents configuration to send a stream to an rtmp server.
type Stream struct {
	Application      string    `json:"application"`
	Name             string    `json:"name"`
	Host             string    `json:"host"`
	ThreadQueueSize  int       `json:"thread_queue_size"`
	Preset           string    `json:"preset"`
	VideoEncoding    string    `json:"video_encoding"`
	VideoInputFormat string    `json:"video_input_format"`
	VideoResolution  string    `json:"video_resolution"`
	AudioDisabled    bool      `json:"audio_disabled"`
	AudioEncoding    string    `json:"audio_encoding"`
	AudioInput       string    `json:"audio_input"`
	AudioRate        int       `json:"audio_rate"`
	AverageBitRate   int       `json:"average_bit_rate"`
	OutputEncoding   string    `json:"output_encoding"`
	Live             bool      `json:"live"`
	Camera           Camera    `json:"camera"`
	cmd              *exec.Cmd `json:"cmd"`
}

// NewStream constructs a *Stream with sensible defaults.
func NewStream() *Stream {
	return &Stream{
		Application:      "killerqueen",
		Name:             "",
		Host:             "localhost",
		ThreadQueueSize:  1024, // max frames in buffer
		Preset:           "veryfast",
		VideoEncoding:    "libx264",
		VideoInputFormat: "video4linux2",
		VideoResolution:  "",
		AudioEncoding:    "aac",
		AudioInput:       "alsa",
		AudioRate:        44100,
		AverageBitRate:   96000,
		OutputEncoding:   "flv",
	}
}

type startResponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
	Cmd    string `json:"cmd"`
}

// Start begins streaming a camera to an rtmp server.
func (s *Stream) Start() startResponse {
	uri := "rtmp://" + s.Host + "/" + s.Application + "/" + s.Name
	resp := startResponse{
		URL: uri,
	}
	args := []string{}
	//args = append(args, "-thread_queue_size", strconv.Itoa(s.ThreadQueueSize))
	if !s.AudioDisabled {
		args = append(args, "-f", s.AudioInput)
		args = append(args, "-ac", strconv.Itoa(s.Camera.NAudioChannels))
		args = append(args, "-i", fmt.Sprintf("hw:%d", s.Camera.CardId()))
	}
	args = append(args, "-f", s.VideoInputFormat)
	args = append(args, "-i", s.Camera.Device)
	args = append(args, "-preset", s.Preset)
	args = append(args, "-tune", "zerolatency")
	if s.VideoResolution != "" {
		args = append(args, "-vf", "scale="+s.VideoResolution)
	}
	args = append(args, "-c:v", s.VideoEncoding)
	if !s.AudioDisabled {
		args = append(args, "-c:a", s.AudioEncoding)
		args = append(args, "-ar", strconv.Itoa(s.AudioRate))
		args = append(args, "-ab", strconv.Itoa(s.AverageBitRate))
	}
	args = append(args, "-f", s.OutputEncoding, uri)
	resp.Cmd = strings.Join(append([]string{"ffmpeg"}, args...), " ")
	resp.Status = "on"
	cmd := exec.Command("ffmpeg", args...)
	s.cmd = cmd
	go func() {
		defer func() { s.Live = false }()
		s.Live = true
		err := cmd.Run()
		if err != nil && !strings.Contains(err.Error(), "signal: killed") {
			log.Printf("Failed to start stream: %s\n%s\n", strings.Join(append([]string{"ffmpeg"}, args...), " "), err)
		}
	}()
	return resp
}

// Stop halts streaming a camera to an rtmp server.
func (s *Stream) Stop() error {
	return s.cmd.Process.Kill()
}

type Camera struct {
	Serial         string `json:"serial"`
	Model          string `json:"model"`
	Vendor         string `json:"vendor"`
	IdPath         string `json:"id_path"`
	Device         string `json:"device"`
	NAudioChannels int    `json:"n_audio_channels"`
}

// Load hardware specific settings.
func (c *Camera) LoadHardware() {
	if c.Model == "HD_Pro_Webcam_C920" {
		c.NAudioChannels = 2
		// set default resolution?
	}
	// 1 may not be a valid value, but zero is never a valid value.
	if c.NAudioChannels == 0 {
		c.NAudioChannels = 1
	}
}

/*
jeremys@jeremys-desktop> udevadm info --query=all --name=/dev/video0 | grep ID_PATH
E: ID_PATH=pci-0000:00:1d.7-usb-0:3:1.0

jeremys@jeremys-desktop> udevadm info --query=all --name=/dev/video1 | grep ID_PATH
E: ID_PATH=pci-0000:00:1a.7-usb-0:6.1.2:1.2
*/

// CardId attempts to correlate this video device to a card ID for use with ffmpeg.
func (c Camera) CardId() int {
	path := strings.TrimPrefix(c.IdPath, "pci-")
	path = strings.Replace(path, "usb-0:", "", 1)
	path = path[:strings.LastIndex(path, ":")] // omit trailing :1.2 etc.

	cardmap, err := getCardmap()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0
	}
	id, ok := cardmap[path]
	if !ok {
		fmt.Fprintln(os.Stderr, err)
		return 0
	}
	return id
}

/*
jeremys@jeremys-desktop> cat /proc/asound/cards
 0 [XFi            ]: SB-XFi - Creative X-Fi
                      Creative X-Fi 20K1 Unknown
 1 [Intel          ]: HDA-Intel - HDA Intel
                      HDA Intel at 0xf1df8000 irq 40
 2 [NVidia         ]: HDA-Intel - HDA NVidia
                      HDA NVidia at 0xf3dfc000 irq 41
 3 [U0x46d0x80a    ]: USB-Audio - USB Device 0x46d:0x80a
                      USB Device 0x46d:0x80a at usb-0000:00:1d.7-3, high speed
 4 [C615           ]: USB-Audio - HD Webcam C615
                      HD Webcam C615 at usb-0000:00:1a.7-6.1.2, high speed
*/

// getCardmap builds a mapping of the device string (as seen in udevadm info and /proc/asound/cards) to card ID.
func getCardmap() (map[string]int, error) {
	m := make(map[string]int)
	fd, err := os.Open("/proc/asound/cards")
	if err != nil {
		return m, err
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		cardId, err := strconv.Atoi(fields[0])
		if err != nil {
			return m, err
		}
		if !scanner.Scan() {
			break
		}
		usbDevice := scanner.Text()
		if !strings.Contains(usbDevice, " at usb-") {
			continue
		}
		afterStr := " at usb-"
		a := strings.Index(usbDevice, afterStr) + len(afterStr)
		b := strings.LastIndex(usbDevice, ",")

		m[usbDevice[a:b]] = cardId
	}
	if err := scanner.Err(); err != nil {
		return m, err
	}
	return m, nil
}
