# NanoPi M4

## ffmpeg hardware accel

Looks like Rockchip has its own hardware accel type, [MPP](https://github.com/rockchip-linux/mpp).
See references to it in [ffmpeg](https://github.com/FFmpeg/FFmpeg/blob/master/configure#L336).  I can't find a Debian Stretch
package, so probably have to build from source.


Neither mmal or rkmpp seen in ffmpeg installed from Armbian package:
```
jeremys@nanopim4> dpkg -l | grep ffmpeg
ii  ffmpeg                          7:3.2.12-1~deb9u1                 arm64        Tools for transcoding, streaming and playing of multimedia files

jeremys@nanopim4> ffmpeg 2>&1 | perl -pe 's/\s+/\n/g' | egrep "mmal|omx|mpp"
--enable-omx
```
