#!/bin/sh
set -ex

go build --ldflags '-linkmode external -extldflags "-static -s -w"' -v ./
