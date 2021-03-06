package kq

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

// use raspberry pi hardware encoder
// sudo ffmpeg -f video4linux2 -i /dev/video0 -preset veryfast -tune zerolatency -c:v h264_omx -f flv rtmp://10.0.0.146/killerqueen/blue

// probe camera's resolutions
// ffmpeg -f v4l2 -list_formats all -i /dev/video0

// Start begins streaming a camera to an rtmp server.
func (s *Stream) Start() startResponse {
	uri := "rtmp://" + s.Host + "/" + s.Application + "/" + s.Name
	resp := startResponse{
		URL: uri,
	}
	args := []string{}
	args = append(args, "-hide_banner")
	if !s.AudioDisabled {
		cardId := s.Camera.CardId()
		args = append(args, "-f", s.AudioInput)
		args = append(args, "-i", fmt.Sprintf("hw:%d", cardId))
		args = append(args, "-ac", strconv.Itoa(s.Camera.NAudioChannels(cardId)))
	}
	args = append(args, "-f", s.VideoInputFormat)
	args = append(args, "-r", "60")
	args = append(args, "-i", s.Camera.Device)
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
	cmd.Stderr = new(bytes.Buffer)
	s.cmd = cmd
	go func() {
		defer func() { s.Live = false }()
		s.Live = true
		err := cmd.Run()
		if err != nil && !strings.Contains(err.Error(), "signal: killed") {
			log.Printf("Failed to start stream: %s\n%s\n", strings.Join(append([]string{"ffmpeg"}, args...), " "), cmd.Stderr)
		}
	}()
	return resp
}

// Stop halts streaming a camera to an rtmp server.
func (s *Stream) Stop() error {
	return s.cmd.Process.Kill()
}

type Camera struct {
	Serial    string `json:"serial"`
	Model     string `json:"model"`
	Vendor    string `json:"vendor"`
	IdPath    string `json:"id_path"`
	Device    string `json:"device"`
	asound    func() io.ReadCloser
	asoundPcm func() io.ReadCloser
}

func (c *Camera) NAudioChannels(cardId int) int {
	if c.asoundPcm == nil {
		c.asoundPcm = func() io.ReadCloser {
			fd, err := os.Open("/proc/asound/pcm")
			if err != nil {
				return ioutil.NopCloser(strings.NewReader(""))
			}
			return fd
		}
	}
	rc := c.asoundPcm()
	defer rc.Close()
	scanner := bufio.NewScanner(rc)
	var n int
	for scanner.Scan() {
		line := scanner.Text()
		// line len check and comment check are for tests.
		// Not expected from /proc/asound/pcm.
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, fmt.Sprintf("%02d", cardId)) {
			continue
		}
		if strings.Contains(line, "Audio") && strings.Contains(line, "capture") &&
			!strings.Contains(line, "playback") {
			n++
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Scanner failed: %s\n", err)
		return 0
	}
	return n
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
	path = strings.TrimPrefix(path, "platform-")
	path = strings.Replace(path, "usb-0:", "", 1)
	path = path[:strings.LastIndex(path, ":")] // omit trailing :1.2 etc.

	if c.asound == nil {
		c.asound = func() io.ReadCloser {
			fd, err := os.Open("/proc/asound/cards")
			if err != nil {
				return ioutil.NopCloser(strings.NewReader(""))
			}
			return fd
		}
	}
	fd := c.asound()
	defer fd.Close()
	cardmap, err := getCardmap(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 0
	}
	id, ok := cardmap[path]
	if !ok {
		fmt.Fprintf(os.Stderr, "Failed to find card for path: %s\n", path)
		fmt.Fprintf(os.Stderr, "Full path is %s\n", c.IdPath)
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
func getCardmap(r io.Reader) (map[string]int, error) {
	m := make(map[string]int)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// line len check and comment check are for tests.
		// Not expected from /proc/asound/cards.
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
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
