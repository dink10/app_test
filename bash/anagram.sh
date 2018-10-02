#!/bin/sh
tr -c [:alpha:] \\n | grep . | while read n; do echo "$n" | sed 's/\(.\)/\1 /g' | tr ' ' \\n | sort | tr -d \\n; echo " $n"; done | sort -k1 | awk '{dups[$1]++} END{for (num in dups) {print num,dups[num]}}';
