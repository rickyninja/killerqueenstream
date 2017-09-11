#!/bin/bash

if [ ! -f /usr/bin/dig ]; then
    echo dig needs to be installed
    exit 1
fi

PATH=/usr/bin:$HOME/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
remoteport=19999
username=killerqueen
remotehost=rickyninja.net
remoteip=$(dig +short ${remotehost})
sshopts="-o ExitOnForwardFailure=yes -o ServerAliveInterval=60"
verbose=0
duration=60

usage() {
    echo <<EOF
usage: $0 [options]
-v output verbose debug info
EOF

    exit 1
}

while getopts 'v' flag; do
    case "${flag}" in
        v) verbose=1 ;;
        h) usage ;;
        *) usage ;;
    esac
done

while true; do
    netstat -an | awk '{print $5}' | sort -u | grep -qs "${remoteip}:22"
    if [ $? -ne 0 ]; then
        [ ${verbose} -eq 1 ] && echo no tunnel seen
        ssh -fN -R ${remoteport}:localhost:22 ${username}@${remotehost} ${sshopts}
        if [ $? -ne 0 ]; then
            [ ${verbose} -eq 1 ] && echo failed to establish tunnel
        fi
    fi
    [ ${verbose} -eq 1 ] && echo sleeping
    sleep ${duration}
done
