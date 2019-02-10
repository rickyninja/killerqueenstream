package kq

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNChannels(t *testing.T) {
	cases := []struct {
		cardId int
		want   int
		pcm    func() io.ReadCloser
	}{
		{2, 2, func() io.ReadCloser {
			r := strings.NewReader(`# cat /proc/asound/pcm
00-00: ff890000.i2s-rt5651-aif1 rt5651-aif1-0 :  : playback 1 : capture 1
01-00: ff8a0000.i2s-i2s-hifi i2s-hifi-0 :  : playback 1
02-00: USB Audio : USB Audio : playback 1 : capture 1
02-01: USB Audio : USB Audio #1 : capture 1
02-02: USB Audio : USB Audio #2 : capture 1
`)
			return ioutil.NopCloser(r)
		}},
		{2, 1, func() io.ReadCloser {
			r := strings.NewReader(`# cat /proc/asound/pcm
00-03: HDMI 0 : HDMI 0 : playback 1
00-07: HDMI 1 : HDMI 1 : playback 1
00-08: HDMI 2 : HDMI 2 : playback 1
00-09: HDMI 3 : HDMI 3 : playback 1
01-00: Generic Analog : Generic Analog : playback 1 : capture 1
01-01: Generic Digital : Generic Digital : playback 1
01-02: Generic Alt Analog : Generic Alt Analog : capture 1
02-00: USB Audio : USB Audio : capture 1
`)
			return ioutil.NopCloser(r)
		}},
	}
	for _, tc := range cases {
		cam := Camera{
			asoundPcm: tc.pcm,
		}
		got := cam.NAudioChannels(tc.cardId)
		if got != tc.want {
			t.Errorf("wrong number of channels, got %d want %d", got, tc.want)
			t.Logf("%#v", tc)
		}
	}
}

func TestCamera_CardId(t *testing.T) {
	cases := []struct {
		idpath string
		want   int
		asound io.ReadCloser
	}{
		{"platform-xhci-hcd.6.auto-usb-0:1.2:1.0", 2,
			ioutil.NopCloser(strings.NewReader(`#
 0 [realtekrt5651co]: realtek_rt5651- - realtek,rt5651-codec
                      realtek,rt5651-codec
 1 [rockchiphdmi   ]: rockchip_hdmi - rockchip,hdmi
                      rockchip,hdmi
 2 [HDMI           ]: USB-Audio - USB Capture HDMI+
                      Magewell USB Capture HDMI+ at usb-xhci-hcd.6.auto-1.2, super speed
`)),
		},
		{"pci-0000:01:00.0-usb-0:2.1.2:1.2", 3,
			ioutil.NopCloser(strings.NewReader(`#
 0 [NVidia         ]: HDA-Intel - HDA NVidia
                      HDA NVidia at 0xef080000 irq 73
 1 [Generic        ]: HDA-Intel - HD-Audio Generic
                      HD-Audio Generic at 0xef900000 irq 75
 2 [U0x46d0x80a    ]: USB-Audio - USB Device 0x46d:0x80a
                      USB Device 0x46d:0x80a at usb-0000:00:1d.7-3, high speed
 3 [C615           ]: USB-Audio - HD Webcam C615
                      HD Webcam C615 at usb-0000:01:00.0-2.1.2, high speed
 4 [C615           ]: USB-Audio - HD Webcam C615
                      HD Webcam C615 at usb-0000:00:1a.7-6.1.2, high speed
`)),
		},
	}
	for _, tc := range cases {
		cam := Camera{
			IdPath: tc.idpath,
			asound: func() io.ReadCloser { return tc.asound },
		}
		got := cam.CardId()
		if got != tc.want {
			t.Errorf("wrong Camera Id, got: %d want %d", got, tc.want)
		}
	}
}

func TestGetCardmap(t *testing.T) {
	cases := []struct {
		input io.Reader
		want  map[string]int
	}{
		{
			strings.NewReader(`#
 0 [realtekrt5651co]: realtek_rt5651- - realtek,rt5651-codec
                      realtek,rt5651-codec
 1 [rockchiphdmi   ]: rockchip_hdmi - rockchip,hdmi
                      rockchip,hdmi
 2 [HDMI           ]: USB-Audio - USB Capture HDMI+
                      Magewell USB Capture HDMI+ at usb-xhci-hcd.6.auto-1.2, super speed
`),
			// The first 2 device indexes aren't in the map because they don't match the expected pattern.
			map[string]int{
				"xhci-hcd.6.auto-1.2": 2,
			}},
		{
			strings.NewReader(`#
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
`),
			map[string]int{
				"0000:00:1d.7-3":     3,
				"0000:00:1a.7-6.1.2": 4,
			}},
	}
	for _, tc := range cases {
		got, err := getCardmap(tc.input)
		if err != nil {
			t.Fatal(err)
		}
		if !cardMapEqual(got, tc.want) {
			t.Errorf("wrong card map")
			t.Errorf(" got: %#v", got)
			t.Errorf("want: %#v", tc.want)
		}
	}
}

func cardMapEqual(got, want map[string]int) bool {
	for k, w := range want {
		g, ok := got[k]
		if !ok {
			return false
		}
		if g != w {
			return false
		}
	}
	return true
}
