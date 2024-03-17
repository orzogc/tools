#!/usr/bin/env bash

# 第一个参数是goldendict popup的group
if [[ -z "$1" ]]
then
    goldendict -s "$(wl-paste -p)"
else
    goldendict -s -p "$1" "$(wl-paste -p)"
fi
