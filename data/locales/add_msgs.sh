#!/bin/bash
while IFS= read -r string; do
    for file in *.po*; do
        printf '\nmsgid "%s"\nmsgstr ""\n' "$string" >> "$file"
    done
    printf 'Added "%s"\n' "$string"
done