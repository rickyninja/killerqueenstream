#!/bin/bash
# simple wrapper script to test ffmpeg cli flags

#-thread_queue_size 10485760 \

ffmpeg \
-thread_queue_size 512 \
-f alsa \
-ac 2 \
-i hw:0 \
-video_size 2560x1600 \
-framerate 60 \
-f x11grab \
-i :0.0 \
-preset superfast \
-tune zerolatency \
-vcodec libx264 \
-acodec aac \
-strict -2 \
-ar 44100 \
-ab 96000 \
-f flv rtmp://localhost/killerqueen/minecraft
