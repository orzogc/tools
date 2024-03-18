#!/usr/bin/env bash

set -e

TEMP_FILE_NAME="/tmp/goldendict_ocr_tmp"
TEMP_SCREENSHOT_IMAGE="${TEMP_FILE_NAME}_screenshot.png"
TEMP_IMAGE="${TEMP_FILE_NAME}.png"
TEMP_TXT="${TEMP_FILE_NAME}.txt"

while getopts ":l:g:s" flag; do
    case "${flag}" in
        l)
            # -l 指定OCR的语言，具体看`man tesseract`的`LANGUAGES AND SCRIPTS`部分
            lang="${OPTARG}"
            ;;
        g)
            # -g 指定goldendict popup的group
            group="${OPTARG}"
            ;;
        s)
            # -s 只扫描鼠标指针后的单词，鼠标指针要在单词的左下角附近，目前只支持OCR英文单词
            single_word=true
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

if [[ -z "${single_word}" ]]
then
    # 选择OCR范围，需要模拟tty否则`slurp`可能不会正常运行
    range=$(script -qefc "slurp -f \"%x %y %w %h\"" /dev/null)
    if [[ -z "${range}" ]]; then
        exit 1
    fi
    # 拿到截图范围
    read x y w h <<< ${range}
else
    # 只扫描鼠标指针后的单词
    if [[ $(kdotool getmouselocation) =~ x:([0-9]+)\ y:([0-9]+) ]]; then
        x=${BASH_REMATCH[1]}
        y=$((${BASH_REMATCH[2]}+8)) # 向下偏移一些像素方便OCR单词
        w=250
        h=-50
    fi
fi

# 截图整个屏幕
spectacle --fullscreen --background --nonotify --output $TEMP_SCREENSHOT_IMAGE

if [[ -f $TEMP_SCREENSHOT_IMAGE ]]; then
    # 根据截图范围裁剪截图（`crop_png`必须在`$PATH`范围里）
    crop_png -x $x -y $y -width $w -height $h -input $TEMP_SCREENSHOT_IMAGE -output $TEMP_IMAGE

    if [[ -f $TEMP_IMAGE ]]; then
        # OCR裁剪后的截图
        tesseract $TEMP_IMAGE $TEMP_FILE_NAME -l "${lang}" --oem 1

        if [[ -z "${single_word}" ]]
        then
            ocr_result=$(cat ${TEMP_TXT})
        else
            # 只扫描最后一行的第一个单词
            if [[ $(cat ${TEMP_TXT} | tail -n1) =~ ([a-zA-Z]+)[^a-zA-Z]? ]]; then
                ocr_result=${BASH_REMATCH[1]}
            fi
        fi

        if [[ ! -z "${ocr_result}" ]]; then
            # 发送OCR得到的文字到goldendict
            if [[ -z "${group}" ]]
            then
                goldendict -s "${ocr_result}"
            else
                goldendict -s -p "${group}" "${ocr_result}"
            fi
        fi
    fi

    # 清理临时文件
    rm -f $TEMP_SCREENSHOT_IMAGE $TEMP_IMAGE $TEMP_TXT
fi
