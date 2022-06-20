#!/usr/bin/python
# Same as ttv.sh but in Python
# Use through ttv wrapper

import os
import requests
import threading

URL = "https://www.twitch.tv/"
DEFAULT_RC_PATH = "~/.config/ttv/ttv.rc"


def parse_config():
    TTV_RC_PATH = os.path.expanduser(DEFAULT_RC_PATH)

    if os.path.exists(TTV_RC_PATH):
        with open(TTV_RC_PATH, "r") as TTV_RC:
            channel_names = TTV_RC.readlines()
        if len(channel_names) == 0:
            print("Config file empty")
            exit()
        return channel_names
    else:
        print(f"Config file not found ({TTV_RC_PATH})")
        exit()


def is_live(name):
    page = requests.get(URL + name).content.decode("utf-8")
    if "isLiveBroadcast" in page:
        print(name)


def main():
    threads = []
    channels = parse_config()

    for channel in channels:
        # name[:-1] to strip "\n" in link
        thread = threading.Thread(target=is_live, args=(channel[:-1],))
        thread.start()
        threads.append(thread)

    for thread in threads:
        thread.join()


main()
