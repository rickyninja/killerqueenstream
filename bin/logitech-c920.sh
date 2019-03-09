#!/bin/bash

ffmpeg -hide_banner \
-f alsa \
-i hw:3 \
-ac 2 \
-f video4linux2 \
-r 60 \
-i /dev/video2 \
-vf scale=1920x1080 \
-c:v libx264 \
-c:a aac \
-ar 44100 \
-ab 96000 \
-f flv rtmp://127.0.0.1/killerqueen/cam1
