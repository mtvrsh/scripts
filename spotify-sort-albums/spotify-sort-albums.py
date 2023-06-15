#!/usr/bin/python3
# spt-sort-albums.py - backup, restore and sort spotify albums
# Saves all albums, deletes them and adds sorted so they appear sorted in
# in applications that don't support sorting.
#
# Fill ./env with appropriate values and source it.
# $ vi env; source ./env

import argparse
import os
import pickle
import sys
from datetime import date

import spotipy

scope = "user-library-read user-library-modify"

username = os.getenv("username", default="")

if len(username) < 1:
    print("username not set")
    sys.exit()

parser = argparse.ArgumentParser()
parser.add_argument("-r", help="restore saved albums", dest="filename")
parser.add_argument("-b", help="backup albums", action="store_true")

args = parser.parse_args()

token = spotipy.util.prompt_for_user_token(username, scope)

if token:
    sp = spotipy.Spotify(auth=token)

    online_albums = {}
    sorted_album_ids = []

    results = sp.current_user_saved_albums()
    albums = results["items"]

    while results["next"]:
        results = sp.next(results)
        albums.extend(results["items"])

    print("Albums in online library:")

    for album in albums:
        name = album["album"]["name"]
        idn = album["album"]["id"]

        print("{} - {}".format(name, idn))
        online_albums[name] = album
    print()

    sorted_album_names = sorted(online_albums)

    for album in sorted_album_names:
        sorted_album_ids.append(online_albums[album]["album"]["id"])

    if args.b:
        try:
            with open(username + "_ids_" + str(date.today()) + ".bak",
                      "xb") as f:
                pickle.dump(sorted_album_ids, f)
                print("Backup created succesfully", "\n")
        except FileExistsError:
            print("Cannot create backup, file already exists.")
            sys.exit()

    if args.filename:
        try:
            with open(args.filename, "rb") as f:
                sorted_album_ids = pickle.load(f)
                print("Backup restored succesfully")
        except FileNotFoundError:
            print("Cannot restore backup, file not found.")
            sys.exit()

    print(
        "Albums prepared to add:\n",
        sorted_album_names,
        "\n",
        sorted_album_ids,
        "\n",
        sep="",
    )

    sorted_album_ids_copy = sorted_album_ids.copy()

    if online_albums and input("Delete online albums? [y/n] ") == "y":
        while len(sorted_album_ids):
            sp.current_user_saved_albums_delete(albums=sorted_album_ids[:50])
            del sorted_album_ids[:50]
        print("Done")

    if sorted_album_ids_copy and input("Add albums? [y/n] ") == "y":
        for a in sorted_album_ids:
            sp.current_user_saved_albums_add(
                albums=[sorted_album_ids_copy.pop()])
        print("Done")

else:
    print("Can't get token for", username)
