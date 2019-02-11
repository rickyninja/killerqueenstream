package main

import (
	"testing"

	kq "github.com/rickyninja/killerqueenstream"
)

func TestSetDefaultStreams(t *testing.T) {
	setDefaultStreams()
	want := &kq.Stream{
		Application:      "killerqueen",
		Name:             "game",
		Host:             "localhost",
		VideoEncoding:    "libx264",
		VideoInputFormat: "video4linux2",
		VideoResolution:  "",
		AudioDisabled:    false,
		AudioEncoding:    "aac",
		AudioInput:       "alsa",
		AudioRate:        44100,
		AverageBitRate:   96000,
		OutputEncoding:   "flv",
		Live:             false,
		Camera: kq.Camera{
			Serial: "C204170601172",
			Model:  "USB_Capture_HDMI+",
			Vendor: "Magewell",
			IdPath: "platform-xhci-hcd.6.auto-usb-0:1.2:1.0",
			Device: "/dev/video0",
		},
	}
	running.Lock()
	got, ok := running.Stream["game"]
	running.Unlock()
	if !ok {
		t.Fatal("expected game to be configured")
	}
	if got.Application != want.Application {
		t.Errorf("wrong Application, got %s want %s", got.Application, want.Application)
	}
	if got.Name != want.Name {
		t.Errorf("wrong Name, got %s want %s", got.Name, want.Name)
	}
	if got.Host != want.Host {
		t.Errorf("wrong Host, got %s want %s", got.Host, want.Host)
	}
	if got.VideoEncoding != want.VideoEncoding {
		t.Errorf("wrong VideoEncoding, got %s want %s", got.VideoEncoding, want.VideoEncoding)
	}
	if got.VideoInputFormat != want.VideoInputFormat {
		t.Errorf("wrong VideoInputFormat, got %s want %s", got.VideoInputFormat, want.VideoInputFormat)
	}
	if got.VideoResolution != want.VideoResolution {
		t.Errorf("wrong VideoResolution, got %s want %s", got.VideoResolution, want.VideoResolution)
	}
	if got.AudioDisabled != want.AudioDisabled {
		t.Errorf("wrong AudioDisabled, got %t want %t", got.AudioDisabled, want.AudioDisabled)
	}
	if got.AudioEncoding != want.AudioEncoding {
		t.Errorf("wrong AudioEncoding, got %s want %s", got.AudioEncoding, want.AudioEncoding)
	}
	if got.AudioInput != want.AudioInput {
		t.Errorf("wrong AudioInput, got %s want %s", got.AudioInput, want.AudioInput)
	}
	if got.AudioRate != want.AudioRate {
		t.Errorf("wrong AudioRate, got %d want %d", got.AudioRate, want.AudioRate)
	}
	if got.AverageBitRate != want.AverageBitRate {
		t.Errorf("wrong AverageBitRate, got %d want %d", got.AverageBitRate, want.AverageBitRate)
	}
	if got.OutputEncoding != want.OutputEncoding {
		t.Errorf("wrong OutputEncoding, got %s want %s", got.OutputEncoding, want.OutputEncoding)
	}
	if got.Live != want.Live {
		t.Errorf("wrong Live, got %t want %t", got.Live, want.Live)
	}
	if got.Camera.Serial != want.Camera.Serial {
		t.Errorf("wrong Serial, got %s want %s", got.Camera.Serial, want.Camera.Serial)
	}
	if got.Camera.Model != want.Camera.Model {
		t.Errorf("wrong Model, got %s want %s", got.Camera.Model, want.Camera.Model)
	}
	if got.Camera.Vendor != want.Camera.Vendor {
		t.Errorf("wrong Vendor, got %s want %s", got.Camera.Vendor, want.Camera.Vendor)
	}
	if got.Camera.IdPath != want.Camera.IdPath {
		t.Errorf("wrong IdPath, got %s want %s", got.Camera.IdPath, want.Camera.IdPath)
	}
	if got.Camera.Device != want.Camera.Device {
		t.Errorf("wrong Device, got %s want %s", got.Camera.Device, want.Camera.Device)
	}
	if got.Camera.Device != want.Camera.Device {
		t.Errorf("wrong Device, got %s want %s", got.Camera.Device, want.Camera.Device)
	}
}
