#!/usr/bin/env bash

set -e

TEMP_FILE_NAME="/tmp/goldendict_ocr_tmp"
TEMP_IMAGE="${TEMP_FILE_NAME}.png"
TEMP_TXT="${TEMP_FILE_NAME}.txt"

while getopts ":l:p:" flag; do
    case "${flag}" in
        l)
            lang="${OPTARG}"
            ;;
        p)
            group="${OPTARG}"
            ;;
        *)
            exit 1
            ;;
    esac
done

if [[ -z "${lang}" ]]; then
    lang="eng"
fi

spectacle --region --background --nonotify --output $TEMP_IMAGE

if [[ -f $TEMP_IMAGE ]]; then
    tesseract $TEMP_IMAGE $TEMP_FILE_NAME -l "${lang}" --oem 1

    if [[ -z "${group}" ]]
    then
        goldendict -s "$(/usr/bin/cat ${TEMP_TXT})"
    else
        goldendict -s -p "${group}" "$(/usr/bin/cat ${TEMP_TXT})"
    fi

    rm $TEMP_IMAGE $TEMP_TXT
fi