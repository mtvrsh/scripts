#!/usr/bin/bash
# ttv.sh - query twitch.tv for online streams and open selected one in mpv
# ttv.sh <mpv options>
# List channels in ttv.rc, default is ~/.config/ttv/ttv.rc, 1 name per line.
set -e pipefail

TTV_RC="${TTV_RC:=$HOME/.config/ttv/ttv.rc}"

URL="https://www.twitch.tv/"

# Another function because idk how to send it to background otherwise
is_live() {
    if curl -4 -s "$URL""$1" | rg -o "isLiveBroadcast" &>/dev/null; then
        echo "$1"
    fi
}

live_channels() {
    while read -r channel; do
        is_live "$channel" &
    done <"$TTV_RC"
}

if [ ! -f "$TTV_RC" ]; then
    echo Config file not found.
elif [ -s "$TTV_RC" ]; then
    # Assign to variable to avoid passing empty strings to mpv (thanks to set -e)
    name="$(live_channels | fzf)"
    mpv --profile=low-latency "$@" "${URL}${name}"
else
    echo Config file empty.
fi
