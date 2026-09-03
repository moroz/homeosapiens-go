#!/usr/bin/env uv run

import csv
from uuid import uuid7
import os
import re
from slugify import slugify
import textwrap

inputfile = open('sources.csv', newline='')
outfile = open('output.csv', mode = "w")
script = open('importer.sh', mode = "w")

playlist_name = "Dr. Asher Shaikh Seminar: "

reader = csv.reader(inputfile)
writer = csv.writer(outfile)

for row in reader:
    print(row)
    id = uuid7()

    source = row[0]
    title = row[1]

    full_title = playlist_name+ title
    writer.writerow([id, source, full_title, slugify(full_title)])

    script.write(textwrap.dedent(f"""\
        mkdir -p {id}
        ffmpeg -i {source} -c copy {id}/avc1.mp4
    """).rstrip() + "\n\n")
