#!/usr/bin/env uv run

import os, glob, uuid, subprocess


def is_uuid(value: str) -> bool:
    try:
        uuid.UUID(value)
        return True
    except ValueError:
        return False


files = glob.glob('*/*.mp4')

for source in files:
    outdir = os.path.dirname(source)
    os.makedirs(f'{outdir}/hls', exist_ok=True)
    cmd = f'ffmpeg -i {source} -c copy -f hls -hls_time 6 -hls_list_size 0 -hls_segment_filename "{outdir}/hls/seg_%04d.ts" {outdir}/hls/index.m3u8'
    subprocess.run(cmd, shell=True)
