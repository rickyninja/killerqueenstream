# raspberry pi performance insights

I began doing streaming on a full desktop PC, with intel core i7, 12GB ram.  I wasn't encountering any performance
bottlenecks at this point because the hardware is so beefy compared to a pi.  I wanted to do streaming on a pi to
cut costs and save space inside the arcade cabinets.

hardware used in testing:

* logitech C615 x 2
* raspberry pi 3 model B.

## 32 bit operating system vs 64 bit operating system

**Use a 32 bit operating system**

I began with a vanilla Debian 64 bit operating system image because the pi 3 is 64 bit.  I was
having a lot of problems with extraordinarily high CPU usage (350% - 400% CPU) when using ffmpeg
to grab a USB camera stream.  This was all of the available CPU on the pi 3.  Stream frames per
second and latency was awful.

After a lot of reading online, I found that raspbian is still only supporting 32 bit, mainly
because the pi's GPU is still only 32 bit, while the pi 3 CPU is 64 bit.  This will require some way to convert
kernel memory (64 bit) to GPU memory (32 bit), which nobody has taken on yet.  In retrospect, I'd bet all the CPU
was context switching to convert bit width between CPU and GPU.

After switching to 32 bit raspbian image, CPU usage dropped to below 50% for a single USB camera in raw mode.

## camera raw output vs compressed output

**Take advantage of compression if your camera is capable**

This project was my first experience doing video streaming with a tool like ffmpeg.  Who knew cameras could do hardware compression?
When I tried to start streaming from a 2nd camera, both camera streams would get bogged down until both stopped updating.  I found that
ffmpeg seems to default to requesting raw output from the camera.  Once I configured ffmpeg to request a compressed output (mjpeg), I was
able to see a live stream from both cameras.

### query camera compression capability

I used 2 logitech C615 cameras, which is what I happen to have at home, to begin trying to stream with a raspberry pi.

#### with ffmpeg

I find the v4l2-ctl output a bit easier to read, but ffmpeg includes resolution capabilities.

```
jeremys@raspberrypi> sudo ffmpeg -f v4l2 -list_formats all -i /dev/video0 -hide_banner
[video4linux2,v4l2 @ 0xb73ab0] Raw       :     yuyv422 :           YUYV 4:2:2 : 640x480 160x120 176x144 320x240 432x240 352x288 640x360 800x448 864x480 1024x576 800x600 960x720 1280x720 1600x896 1920x1080
[video4linux2,v4l2 @ 0xb73ab0] Compressed:       mjpeg :          Motion-JPEG : 640x480 160x120 176x144 320x240 432x240 352x288 640x360 800x448 864x480 1024x576 800x600 960x720 1280x720 1600x896 1920x1080
/dev/video0: Immediate exit requested
```

#### with v4l2-ctl

YUYV is raw output, as noted in the 'ffmpeg -list_formats` output.

```
jeremys@raspberrypi> sudo v4l2-ctl --list-formats
ioctl: VIDIOC_ENUM_FMT
        Index       : 0
        Type        : Video Capture
        Pixel Format: 'YUYV'
        Name        : YUYV 4:2:2

        Index       : 1
        Type        : Video Capture
        Pixel Format: 'MJPG' (compressed)
        Name        : Motion-JPEG
```


## hardware acceleration for h264

Raspberry pi's have hardware capabilities for h264 encoding and decoding.  This is one of the things that makes them great for being
a media center computer with software like Kodi.

**Use hardware to encode h264 output**

This cut pi 3 CPU usage per camera to less than 10%.  With so much CPU remaining, I hope to be able to run the capture card on the pi too.

The ffmpeg syntax to do hardware encoding for output stream is `-c:v h264_omx`.
In context of a full command:
`sudo ffmpeg -f video4linux2 -input_format mjpeg -i /dev/video0 -preset veryfast -tune zerolatency -c:v h264_omx -f flv rtmp://localhost/killerqueen/blue`
