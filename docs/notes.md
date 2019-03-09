# Notes

- HDMI 1.0 does up to 3.96 Gbs, rPi usb 2.0 max is 480 Mbs
  https://superuser.com/questions/1118496/can-raspberry-pi-capture-hdmi-input/1118643

- looks interesting, iffy though
  https://github.com/AdamLaurie/hdmi-rip

- raspberry pi only seeing usb 1.0 when device connected
  https://stackoverflow.com/questions/49867887/usb-serial-communication-slow-on-raspberry-pi-3

- an overview of linux usb
  https://www.linuxjournal.com/article/8093

# USB power ports on/off

- http://embeddedapocalypse.blogspot.com/2016/10/how-to-power-off-raspberry-pi-3-usb-or.html
  - https://github.com/codazoda/hub-ctrl.c

# duplicate /dev/video devices

I'm seeing this on Armbian on the nanopi m4.  If you select the wrong device, ffmpeg commands fail with
obscure ioctl error.
