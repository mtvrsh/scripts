#!/bin/bash
# helper for ttv.go and ttv.py
set -e

if [ "$1" = "go" ]; then
	echo go
	channel="$(go run ttv.go | fzf)"
else
	echo python
	channel="$(./ttv.py | fzf)"
fi

mpv --profile=low-latency "$@" https://twitch.tv/"$channel"
