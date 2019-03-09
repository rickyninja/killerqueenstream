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
## compile ffmpeg with rkmpp

### mpp

I wasn't able to build a debian package based on the [reference doc](http://opensource.rock-chips.com/wiki_Mpp#Unix.2FLinux).
Had better luck installing with make.

```
jeremys@nanopim4> sudo apt install -y cmake fakeroot debhelper
jeremys@nanopim4> git clone -b release https://github.com/rockchip-linux/mpp.git
jeremys@nanopim4> cd mpp
jeremys@nanopim4> make
jeremys@nanopim4> sudo make install
```

### ffmpeg

```
jeremys@nanopim4> git clone https://git.ffmpeg.org/ffmpeg.git
jeremys@nanopim4> cd ffmpeg
jeremys@nanopim4> sudo apt install -y libdrm-dev
jeremys@nanopim4> ./configure --enable-rkmpp --enable-libdrm --enable-libx264 --enable-version3 --enable-gpl
jeremys@nanopim4> make -j `nproc`
jeremys@nanopim4> sudo make install
```
