#!/bin/bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

STANDALONE_DIR=".next/standalone"
STATIC_SRC=".next/static"
STATIC_DEST="$STANDALONE_DIR/.next/static"
PUBLIC_SRC="public"
PUBLIC_DEST="$STANDALONE_DIR/public"

if [ ! -d "$STANDALONE_DIR" ]; then
  echo "skip: standalone output not found"
  exit 0
fi

rm -rf "$STATIC_DEST" "$PUBLIC_DEST"
mkdir -p "$STATIC_DEST" "$PUBLIC_DEST"

cp -r "$STATIC_SRC"/. "$STATIC_DEST"/
cp -r "$PUBLIC_SRC"/. "$PUBLIC_DEST"/

test -d "$PUBLIC_DEST/images"
test -d "$PUBLIC_DEST/videos"
test -f "$PUBLIC_DEST/images/favicon.webp"
test -f "$PUBLIC_DEST/videos/LearnWeaver-Seamless-Loop.mp4"
test -d "$STATIC_DEST"

echo "standalone assets prepared"
