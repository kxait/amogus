#!/bin/sh

if [ -z ${SLAVES+x} ]; then
    echo "set SLAVES"
    exit 1
fi

cat << EOF
length_start: 6
length_end: 6
characters: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPRSTUVWXYZ012345678
mode: sha512
slaves: $SLAVES
chunk_size: 50000
test_suite_sample_size: 12