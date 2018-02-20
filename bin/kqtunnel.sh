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
sshopts="$sshopts -i /home/jeremys/.ssh/brpi_rsa"
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

# tcp        0      0 10.0.0.118:22           10.0.0.146:54388        ESTABLISHED

tunnelpid=0
while true; do
    netstat -an | awk '$6 == "ESTABLISHED" {print $5}' | sort -u | grep -qs "${remoteip}:22"
    if [ $? -ne 0 ]; then
        if [ ${tunnelpid} -ne 0 ]; then
            [ ${verbose} -eq 1 ] && echo killing prior tunnel, pid ${tunnelpid}
            kill ${tunnelpid}
        fi
        [ ${verbose} -eq 1 ] && echo no tunnel seen
        ssh -N -R ${remoteport}:localhost:22 ${username}@${remotehost} ${sshopts} &
        if [ $? -eq 0 ]; then
	    tunnelpid=$!
	    trap "{ kill ${tunnelpid}; exit; }" EXIT SIGINT SIGTERM
	    [ ${verbose} -eq 1 ] && echo tunnelpid is $tunnelpid
        else
            [ ${verbose} -eq 1 ] && echo failed to establish tunnel
        fi
    else
        [ ${verbose} -eq 1 ] && echo tunnel to ${remoteip} is established
    fi
    [ ${verbose} -eq 1 ] && echo sleeping
    sleep ${duration}
done
