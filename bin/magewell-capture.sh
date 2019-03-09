#!/bin/bash
# example script that gets input from magewell capture card

#-thread_queue_size 10485760 \

#res=1920x1080
res='640x360'
#-s $res \
hw=2


ffmpeg -hide_banner \
-f alsa \
-i hw:$hw \
-ac 2 \
-f video4linux2 \
-r 30 \
-i /dev/video0 \
-vf scale=1280:720 \
-preset ultrafast \
-c:a aac \
-c:v libx264 \
-ar 44100 \
-ab 96000 \
-f flv \
rtmp://localhost/killerqueen/game

#ffmpeg -f video4linux2 -input_format mjpeg -i /dev/video0 -preset veryfast -tune zerolatency -c:v h264_omx -f flv rtmp://localhost/killerqueen/blue

#jeremys@nanopim4> cat /proc/asound/cards
# 0 [realtekrt5651co]: realtek_rt5651- - realtek,rt5651-codec
#                      realtek,rt5651-codec
# 1 [rockchiphdmi   ]: rockchip_hdmi - rockchip,hdmi
#                      rockchip,hdmi
# 2 [HDMI           ]: USB-Audio - USB Capture HDMI+
#                      Magewell USB Capture HDMI+ at usb-xhci-hcd.6.auto-1.2, super speed

