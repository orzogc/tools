#!/usr/bin/env bash

# todo

set -e

TEMP_FILE_NAME="/tmp/goldendict_ocr_tmp"
TEMP_IMAGE="${TEMP_FILE_NAME}.png"
TEMP_TXT="${TEMP_FILE_NAME}.txt"

spectacle --region --background --nonotify --output $TEMP_IMAGE

if [[ -f $TEMP_IMAGE ]]; then
    tesseract $TEMP_IMAGE $TEMP_FILE_NAME -l eng --oem 1
    goldendict -s "$(/usr/bin/cat ${TEMP_TXT})"
    rm $TEMP_IMAGE $TEMP_TXT
fi
