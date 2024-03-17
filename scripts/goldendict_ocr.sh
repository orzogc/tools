#!/usr/bin/env bash

set -e

TEMP_FILE_NAME="/tmp/goldendict_ocr_tmp"
TEMP_SCREENSHOT_IMAGE="${TEMP_FILE_NAME}_screenshot.png"
TEMP_IMAGE="${TEMP_FILE_NAME}.png"
TEMP_TXT="${TEMP_FILE_NAME}.txt"

while getopts ":l:p:" flag; do
    case "${flag}" in
        l)
            # -l 指定OCR的语言，具体看`man tesseract`的`LANGUAGES AND SCRIPTS`部分
            lang="${OPTARG}"
            ;;
        p)
            # -p 指定goldendict popup的group
            group="${OPTARG}"
            ;;
        *)
            exit 1
            ;;
    esac
done

# 默认翻译英文
if [[ -z "${lang}" ]]; then
    lang="eng"
fi

# 选择取词范围
range=$(slurp)
if [[ -z "${range}" ]]; then
    exit 1
fi

# 截图整个屏幕
spectacle --fullscreen --background --nonotify --output $TEMP_SCREENSHOT_IMAGE

if [[ -f $TEMP_SCREENSHOT_IMAGE ]]; then
    # 拿到截图范围
    read x y w h <<< ${range//[^0-9]/ }
    # 根据截图范围裁剪截图
    crop_png -x $x -y $y -width $w -height $h -input $TEMP_SCREENSHOT_IMAGE -output $TEMP_IMAGE

    if [[ -f $TEMP_IMAGE ]]; then
        # OCR裁剪后的截图
        tesseract $TEMP_IMAGE $TEMP_FILE_NAME -l "${lang}" --oem 1

        # 发送OCR后的文字到goldendict
        if [[ -z "${group}" ]]
        then
            goldendict -s "$(/usr/bin/cat ${TEMP_TXT})"
        else
            goldendict -s -p "${group}" "$(/usr/bin/cat ${TEMP_TXT})"
        fi

        # 清理临时文件
        rm $TEMP_SCREENSHOT_IMAGE $TEMP_IMAGE $TEMP_TXT
    fi
fi
