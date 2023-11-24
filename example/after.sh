#!/bin/sh
## Example: a typical script with several problems
for f in ./*.m3u
do
  grep -qi "hq.*mp3" "$f" \
    && echo "Playlist $f contains a HQ file in mp3 format"
done