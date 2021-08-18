#!/bin/bash

git tag -a "$1" -m "$2"

git push origin main

goreleaser --rm-dist
