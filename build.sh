#!/bin/sh

git show --format="format:PROGRAM_COMMIT_HASH=%h%n" -s --output .env >& /dev/null

