#!/usr/bin/env -S uv run

import os
from os import path
import boto3
import psycopg

database_url = os.environ['DATABASE_URL']
db = psycopg.connect(database_url)

s3 = boto3.client('s3')

bucket = "homeosapiens-staging-assets"

videos = {
    '019e26f5-94ff-738f-b892-f25c3ceaa231': 'tptaoh_d2p1.mp4',
    '019e26f8-414c-7032-b627-77d1107b558d': 'tptaoh_d2p2.mp4',
    '019e26f8-a4d3-729b-872e-9d66242969ef': 'tptaoh_d2p3.mp4',
}

videos_dir = path.abspath('../..')

for video_id in videos.keys():
    object_key = f"videos/{video_id}/avc1_1080.mp4"
    video_path = path.join(videos_dir, videos[video_id])
    print(f"Uploading file {video_path} to {object_key}...")
    s3.upload_file(
        video_path,
        bucket,
        object_key,
        ExtraArgs={
            "ContentType": "video/mp4",
            "CacheControl": "public, max-age=31536000, immutable",
        }
    )
    with db.cursor() as cur:
        cur.execute("""
                    insert into video_sources (content_type, codec, video_id, object_key, priority)
                    values ('video/mp4', 'avc1.640028,mp4a.40.2', %s, %s, 1)
                    """, params=(video_id, object_key))

db.commit()
db.close()
