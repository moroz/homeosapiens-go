#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEGMENT_DURATION=6

for input in "$SCRIPT_DIR"/*.mp4; do
    name="$(basename "$input" .mp4)"
    outdir="$SCRIPT_DIR/$name"

    echo "[$(date +%H:%M:%S)] Remuxing $name ..."
    mkdir -p "$outdir"

    ffmpeg -y -i "$input" \
        -c copy \
        -hls_time "$SEGMENT_DURATION" \
        -hls_list_size 0 \
        -hls_segment_filename "$outdir/seg_%04d.ts" \
        "$outdir/index.m3u8"

    echo "[$(date +%H:%M:%S)] Done -> $outdir/index.m3u8"
done

echo "All done."
