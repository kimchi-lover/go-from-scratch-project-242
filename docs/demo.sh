#!/bin/sh
# Demo script for the asciinema recording: prints each command as a prompt
# line, runs it, and pauses so the cast reads like a live session.
run() {
  printf '$ %s\n' "$1"
  sleep 0.6
  sh -c "$1"
  sleep 1.4
}

run 'make build'
run './bin/hexlet-path-size -h'
run './bin/hexlet-path-size testdata/test.txt'
run './bin/hexlet-path-size testdata/dir'
run './bin/hexlet-path-size -r testdata/dir'
run './bin/hexlet-path-size testdata/hidden'
run './bin/hexlet-path-size -r -a testdata/hidden'
run './bin/hexlet-path-size -H bin/hexlet-path-size'
run './bin/hexlet-path-size missing.txt'
sleep 1
