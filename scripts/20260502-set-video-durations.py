import os
import shutil
import subprocess
import urllib.parse

import psycopg
from psycopg.rows import namedtuple_row

database_url = os.environ["DATABASE_URL"]
cloudfront_url = "https://d3n1g0yg3ja4p3.cloudfront.net"
conn = psycopg.connect(database_url, row_factory=namedtuple_row)
cur = conn.cursor()

ffprobe_bin = shutil.which("ffprobe") or "/home/linuxbrew/.linuxbrew/bin/ffprobe"

cur.execute("""
            select v.id as video_id, vs.object_key, vs.content_type
            from videos v
                     join video_sources vs on vs.video_id = v.id
            where vs.content_type = 'video/mp4'
              and v.duration_seconds is null
            """)

rows = list(cur)

for row in rows:
    url = urllib.parse.urljoin(cloudfront_url, row.object_key)
    duration = subprocess.check_output(
        f"{ffprobe_bin} -v quiet -show_entries format=duration -of csv=p=0 \"{url}\"",
        shell=True,
    )
    duration = int(float(duration.decode('utf-8').strip()))

    cur.execute("update videos set duration_seconds = %s where id = %s", params=(duration, row.video_id))

conn.commit()
conn.close()
