#!/bin/sh

git show --format="format:PROGRAM_COMMIT_HASH=%h%n" -s --output aa2.env >& /dev/null
echo "READER_NAME=${READER_NAME:-$1}" >> aa2.env

