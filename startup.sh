#!/usr/bin/env bash

process="ipserving"
echo ${process} $1 $2

start(){
    if [ "$2"x = ""x ];then
        echo "miss conf file"
    else
        pid=`pgrep ${process}`
        if [ "${pid}"x = ""x ];then
            echo "start new process..."
            nohup ./${process} -c $2 >/dev/null 2>&1 &
        else
            for i in ${pid}
            do
                echo "reload the process [ $i ]"
                kill -USR2 $i
            done
	    sleep 2
        fi
	pid=`pgrep ${process}`
	echo "new process id: ${pid}"
    fi
}

stop(){
    pid=`pgrep ${process}`
    echo ${pid}
    for i in ${pid}
    do
        echo "kill the process [ $i ]"
	kill -9 $i
    done
}

status(){
    ps aux | grep -w ${process} | grep -v 'grep'
}


case "$1" in
    start)
	start $1 $2;;
    reload)
	start $1 $2;;
    stop)
	stop ;;
    status)
	status ;;
    *)
	echo "Usage: $0 {start|stop|reload|status}"

esac
