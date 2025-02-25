#! /usr/bin/env bash
if [ -f "/tmp/stderr.log" ]; then
    sudo chown root /tmp/stderr.log
fi

if [ -f "/tmp/stdout.log" ]; then
    sudo chown root /tmp/stdout.log
fi

(tail -f /tmp/stdout.log) &
pid=$!
(tail -f /tmp/stderr.log) &
pid="$pid $1"

trap "eval kill -9 $pid" EXIT TERM
/start.sh $@ > /tmp/stdout.log 2> /tmp/stderr.log
