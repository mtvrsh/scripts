#!/bin/python
import argparse
import json
import sys


DEDUP_STORE = set()


def walk_bookmarks_tree(bookmark, opts):
    # reversed because deleting is not allowed while iterating
    # https://code.whatever.social/questions/18418/elegant-way-to-remove-items-from-sequence-in-python#181062
    for child in reversed(bookmark["children"]):
        child["title"] = child["title"].strip()
        if child["type"] == "text/x-moz-place-container":
            if not "children" in child or len(child["children"]) == 0:
                bookmark["children"].remove(child)
                continue
            walk_bookmarks_tree(child, opts)
        elif child["type"] == "text/x-moz-place":
            child["uri"] = child["uri"].strip()
            if child["uri"].startswith("about:newtab"):
                bookmark["children"].remove(child)
                continue
            if child["title"] == "":
                child["title"] = child["uri"]
            if opts.d:
                if child["uri"] in DEDUP_STORE:
                    bookmark["children"].remove(child)
                    continue
                else:
                    DEDUP_STORE.add(child["uri"])
            if not opts.json:
                if opts.u:
                    print(child["uri"], file=opts.o)
                else:
                    # remove delimeter from title
                    title = "".join(filter(lambda c: c != "|", child["title"]))
                    print(f"{title}|{child['uri']}", file=opts.o)


def main():
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawDescriptionHelpFormatter,
        description="""\
CLI utility to manage bookmarks exported from firefox.

Examples:
$ ff-bookmarks.py bookmarks-*.json | fuzzel -d | cut -d"|" -f2- | wl-copy
$ ff-bookmarks.py -u bookmarks-*.json | fzf
$ ff-bookmarks.py bookmarks.json -d --json -o clean.json
""",
    )
    parser.add_argument(
        "input",
        nargs="?",
        type=argparse.FileType("r"),
        default=sys.stdin,
        help="path to bookmarks-*.json",
    )
    parser.add_argument(
        "-o",
        nargs="?",
        type=argparse.FileType("w"),
        default=sys.stdout,
        help="file to write modified bookmarks (default: stdout)",
        metavar="output",
    )
    # parser.add_argument(
    # "-t",
    # help="add tags to bookmarks",
    # action="store_true",
    # )
    parser.add_argument(
        "-u",
        help="print only URI",
        action="store_true",
    )
    parser.add_argument(
        "-d",
        help="remove bookmarks with duplicated URIs",
        action="store_true",
    )
    parser.add_argument(
        "--json",
        help="output JSON",
        action="store_true",
    )
    opts = parser.parse_args()

    bookmarks = json.loads(opts.input.read())
    walk_bookmarks_tree(bookmarks, opts)
    if opts.json:
        print(json.dumps(bookmarks), file=opts.o)


if __name__ == "__main__":
    sys.exit(main())
