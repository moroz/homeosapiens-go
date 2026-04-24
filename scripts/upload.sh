#!/usr/bin/env -S bash -euo pipefail

base_dir="/home/karol/working/hs/thumbs"
out_dir="$base_dir/out"

videos="$(find $base_dir -mindepth 1 -maxdepth 1 -type f -name "*.mp4")"

for video in $videos; do
  base="$(basename $video ".mp4")"
  uuid="$(npx -y uuid v7 | tr -d '\n')"
  mkdir -p $out_dir/$uuid

  resolution=1080
  if [[ "$base" =~ ^modi_webinar ]]; then
    resolution=720
  fi

  mv $video $out_dir/$uuid/avc1_$resolution.mp4
  mkdir -p $out_dir/$uuid/hls/
  mv $base_dir/$base/* $out_dir/$uuid/hls/
done

BUCKET="s3://homeosapiens-staging-assets/videos/"
VAULT="aws-vault exec medic --"

# HLS playlists — no caching (dynamic index files)
$VAULT aws s3 cp --recursive \
  --exclude "*" --include "*.m3u8" \
  --content-type "application/vnd.apple.mpegurl" \
  --cache-control "no-cache" \
  $out_dir/ "$BUCKET"

# HLS segments — immutable
$VAULT aws s3 cp --recursive \
  --exclude "*" --include "*.ts" \
  --content-type "video/mp2t" \
  --cache-control "public, max-age=31536000, immutable" \
  $out_dir/ "$BUCKET"

# MP4
$VAULT aws s3 cp --recursive \
  --exclude "*" --include "*.mp4" \
  --content-type "video/mp4" \
  --cache-control "public, max-age=31536000, immutable" \
  $out_dir/ "$BUCKET"