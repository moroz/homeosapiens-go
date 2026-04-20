# Preparing videos for Safari

Start with a HVC1-encoded video (neither Chrome nor Safari will play HEVC-encoded video):

```
$ ffprobe ../hvc1_1080.mp4
ffprobe version 8.1 Copyright (c) 2007-2026 the FFmpeg developers
  built with Apple clang version 17.0.0 (clang-1700.6.4.2)
  configuration: --prefix=/opt/homebrew/Cellar/ffmpeg/8.1_1 --enable-shared --enable-pthreads --enable-version3 --cc=clang --host-cflags= --host-ldflags= --enable-ffplay --enable-gpl --enable-libsvtav1 --enable-libopus --enable-libx264 --enable-libmp3lame --enable-libdav1d --enable-libvmaf --enable-libvpx --enable-libx265 --enable-openssl --enable-videotoolbox --enable-audiotoolbox --enable-neon
  libavutil      60. 26.100 / 60. 26.100
  libavcodec     62. 28.100 / 62. 28.100
  libavformat    62. 12.100 / 62. 12.100
  libavdevice    62.  3.100 / 62.  3.100
  libavfilter    11. 14.100 / 11. 14.100
  libswscale      9.  5.100 /  9.  5.100
  libswresample   6.  3.100 /  6.  3.100
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from '../hvc1_1080.mp4':
  Metadata:
    major_brand     : isom
    minor_version   : 512
    compatible_brands: isomiso2mp41
    encoder         : Lavf62.3.100
  Duration: 02:19:04.85, start: 0.000000, bitrate: 574 kb/s
  Stream #0:0[0x1](und): Video: hevc (Main) (hvc1 / 0x31637668), yuv420p(tv, bt709/bt709/unknown, progressive), 1920x1080 [SAR 1:1 DAR 16:9], 440 kb/s, 25 fps, 25 tbr, 12800 tbn (default)
    Metadata:
      handler_name    : VideoHandler
      encoder         : Lavc62.11.100 libx265
      timecode        : 01:00:00:00
  Stream #0:1[0x2](und): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 128 kb/s (default)
    Metadata:
      handler_name    : SoundHandler
  Stream #0:2[0x3](eng): Data: none (tmcd / 0x64636D74)
    Metadata:
      handler_name    : VideoHandler
      timecode        : 01:00:00:00
```

Repackage the file as M3U8 (VOD playlist):

```
ffmpeg -i hvc1_1080.mp4 -c:v copy -c:a copy -tag:v hvc1 -hls_time 6 -hls_playlist_type vod -hls_segment_type fmp4 -hls_fmp4_init_filename "init.mp4" -hls_segment_filename "seg%03d.m4s" p1_hls.m3u8
```

TODO: Variable bitrates.

Upload the file to S3 with the correct headers:

```
video_id="019a8668-bb4f-7c9c-b9b8-3f274de96566"

aws s3 cp . s3://homeosapiens-staging-assets/videos/$video_id/hls/ \
  --recursive \
  --exclude "*" \
  --include "*.m4s" \
  --content-type "video/iso.segment" \
  --cache-control "public, max-age=31536000, immutable"

aws s3 cp . s3://homeosapiens-staging-assets/videos/$video_id/hls/ \
  --recursive \
  --exclude "*" \
  --include "init.mp4" \
  --content-type "video/mp4" \
  --cache-control "public, max-age=31536000, immutable"

aws s3 cp . s3://homeosapiens-staging-assets/videos/$video_id/hls/ \
  --recursive \
  --exclude "*" \
  --include "*.m3u8" \
  --content-type "application/vnd.apple.mpegurl" \
  --cache-control "public, max-age=3600"
```

If requests show stale cache you may need to create a CloudFront invalidation.
