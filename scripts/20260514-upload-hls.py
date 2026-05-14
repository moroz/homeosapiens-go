#!/usr/bin/env -S uv run
import os
from os import path
from pathlib import Path
import psycopg
import boto3

bucket = "homeosapiens-staging-assets"
s3 = boto3.client('s3')
db = psycopg.connect(os.environ['DATABASE_URL'])

videos_dir = "/Users/karol/working/hs"

videos = {
    '019e26f5-94ff-738f-b892-f25c3ceaa231': 'tptaoh_d2p1.mp4',
    '019e26f8-414c-7032-b627-77d1107b558d': 'tptaoh_d2p2.mp4',
    '019e26f8-a4d3-729b-872e-9d66242969ef': 'tptaoh_d2p3.mp4',
}

for video_id in videos.keys():
    dirname = Path(videos[video_id]).stem
    dirname = path.join(videos_dir, dirname)

    print(f"Uploading segments for video {video_id}...")

    for segment in Path(dirname).glob("*.ts"):
        filename = path.basename(segment)
        object_key = f"videos/{video_id}/hls/{filename}"

        s3.upload_file(
            segment,
            bucket,
            object_key,
            ExtraArgs={
                "ContentType": "video/mp2t",
                "CacheControl": "public, max-age=31536000, immutable",
            }
        )

    playlist_key = f"videos/{video_id}/hls/index.m3u8"
    print(f"Uploading playlist for video {video_id}...")
    s3.upload_file(
        path.join(dirname, "index.m3u8"),
        bucket,
        playlist_key,
        ExtraArgs={
            "ContentType": "application/vnd.apple.mpegurl",
            "CacheControl": "no-cache",
        }
    )

    with db.cursor() as cur:
        cur.execute(
            """
            insert into video_sources (content_type, codec, video_id, object_key, priority)
            values ('application/vnd.apple.mpegurl', null, %s, %s, 0)
            """,
            params=(video_id, f"/{playlist_key}")
        )

db.commit()
db.close()
