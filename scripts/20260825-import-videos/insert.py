#!/usr/bin/env uv run

import csv, psycopg, re
import shutil
import subprocess
import uuid

from psycopg.rows import namedtuple_row

conn = psycopg.connect("postgres://postgres:postgres@localhost/homeosapiens_dev")

inputfile = open('output.csv')

ffprobe = shutil.which('ffprobe')

INSERT_SOURCE_QUERY = """
                      insert into video_sources (video_id, object_key, priority, content_type, codec)
                      values (%s, %s, %s, %s, %s) \
                      """

ASHER_ID = "019beef9-4287-714f-982b-2524fdef7063"

with conn.cursor(row_factory=namedtuple_row) as cur:
    playlist_id = uuid.uuid7()
    cur.execute(
        "insert into video_groups (id, title_en, title_pl, slug) values (%s, %s, %s, %s)",
        (
            playlist_id,
            "Dr. Asher Shaikh Seminar",
            "Seminarium z drem Asherem Shaikh",
            "dr-asher-shaikh-seminar"
        )
    )

    for index, row in enumerate(csv.reader(inputfile)):
        id = row[0]
        title = row[2]
        slug = row[3]

        mp4_file = f"{id}/avc1_1080.mp4"
        duration = subprocess.check_output(
            f"{ffprobe} -v quiet -show_entries format=duration -of csv=p=0 \"{mp4_file}\"",
            shell=True,
        ).decode('utf-8').strip()

        day, part = re.findall(r"(\d+)", title)
        title_pl = f"Seminarium z drem Asherem Shaikh: Dzień {day}, część {part}"
        cur.execute(
            """
            insert into videos (id, title_pl, title_en, provider, slug, duration_seconds, recorded_on, host_id)
            values (%s, %s, %s, 'cloudfront', %s, %s, '2026-06-04'::date + %s * interval '1 day', %s)
            """,
            (id, title_pl, title, slug, int(float(duration)), int(day), ASHER_ID),
        )

        cur.execute(
            """
            insert into video_groups_videos ("position", video_id, video_group_id)
            values (%s, %s, %s)
            """,
            (index, id, playlist_id)
        )

        cur.execute(INSERT_SOURCE_QUERY, (id, f"/videos/{id}/hls/index.m3u8", 0, "application/vnd.apple.mpegurl", None))
        cur.execute(INSERT_SOURCE_QUERY, (id, f"/videos/{id}/avc1_1080.mp4", 1, "video/mp4", "avc1.640028,mp4a.40.2"))

conn.commit()
