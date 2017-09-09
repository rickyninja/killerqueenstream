#!/usr/bin/perl

use strict;
use warnings;

# video mosaic
# https://trac.ffmpeg.org/wiki/Create%20a%20mosaic%20out%20of%20several%20input%20videos
# https://superuser.com/questions/1223900/overlay-two-live-streams-with-ffmpeg

# command line play stream via ffmpeg tools
# ffplay -i "rtmp://localhost:1935/live/aaa live=1"

# send file to nginx rtmp
# ffmpeg -re -i /media/velociraptor1/tivo/minecraft/Blaze\ Farm.mp4 -c:v libx264 -c:a libmp3lame -ar 44100 -ac 1 -f flv rtmp://localhost/myapp/blaze


my $file = shift;

stream_file($file);
#capture_screen($file);
#capture_cam();
#capture_cam_video_only();

sub stream_file {
    my $file = shift || die 'missing file';
    my @args = (
        "ffmpeg",
        "-i", $file,
        "-f", "flv",
        q{-metadata}, q{streamName="aaa"},
        "tcp://localhost:6666"
    );
    if (system(@args) != 0) {
        die "system @args failed: $?";
    }
}

# screen capture.  my native resolution is 2560x1600
sub capture_screen {
    my $file = shift || die 'missing file';
    my @args = (
        qw(
            ffmpeg
            -thread_queue_size 10485760
            -f alsa
            -ac 2
            -i hw:0
            -video_size 2560x1600
            -framerate 60
            -f x11grab
            -i :0.0
            -preset superfast
            -vcodec libx264
            -acodec aac
            -strict -2
            -ar 44100
            -ab 96000
        ),
        "-f", "flv",
        q{-metadata}, q{streamName="aaa"},
        "tcp://localhost:6666"
    );
    if (system(@args) != 0) {
        die "system @args failed: $?";
    }
}

# working cam capture audio & video (used to work)
sub capture_cam {
    my @args = (
        "ffmpeg",
        "-thread_queue_size", "10485760",
        "-f", "alsa",
        "-ac", "1",
        "-i", "hw:3",
        "-f", "v4l2",
        "-i", "/dev/video0",
        "-preset", "veryfast",
        "-threads", "0",
        "-framerate", "60",
        "-s", "1680x1050",
        "-vcodec", "libx264",
        "-acodec", "aac",
        "-strict", "-2",
        "-ar", "44100",
        "-ab", "96000",
        "-f", "flv",
        q{-metadata}, q{streamName=aaa},
        "tcp://localhost:6666"
    );
    if (system(@args) != 0) {
        die "system @args failed: $?";
    }
}

# works no audio
sub capture_cam_video_only {
    my @args = qw(
      ffmpeg
      -i /dev/video0
      -f v4l2
      -framerate 60
      -s 1680x1050
      -f flv
      -metadata streamName=aaa
      tcp://localhost:6666
    );
    if (system(@args) != 0) {
        die "system @args failed: $?";
    }
}


=pod example I began with, doesn't work

ffmpeg \
-f alsa \
-i pulse \
-f x11grab \
-s 1680x1050 \
-r 30 \
-i :0.0+0,0 \
-vf "movie=/dev/video0:f=video4linux2, scale=240:-1, fps, setpts=PTS-STARTPTS [movie]; [in][movie] overlay=main_w-overlay_w-2:main_h-overlay_h-2 [out]" \
-vcodec libx264 \
-crf 20 \
-preset veryfast \
-minrate 150k \
-maxrate 500k \
-s 960x540 \
-acodec libfaac \
-ar 44100 \
-ab 96000 \
-threads 0 \
-f flv rtmp://localhost:1935/vod/inboundLiveFlv
#-f flv - | tee name.flv | ffmpeg -i - -codec copy -f flv -metadata streamName=livestream tcp://localhost:1935

=cut
