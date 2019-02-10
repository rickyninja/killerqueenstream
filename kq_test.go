package kq

import (
	"io"
	"strings"
	"testing"
)

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
