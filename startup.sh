#!/usr/bin/env bash

process="syncdemo"
echo ${process} $1

start(){
      pid=`pgrep ${process}`
      if [ "${pid}"x = ""x ];then
          echo "start new process..."
          nohup ./${process} >/dev/null 2>&1 &
      else
          for i in ${pid}
          do
              echo "reload the process [ $i ]"
              kill -SIGUSR2 $i
          done
      fi
      sleep 1
      pid=`pgrep ${process}`
      echo "new process id: ${pid}"
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
	start $1;;
    reload)
	start $1;;
    stop)
	stop ;;
    status)
	status ;;
    *)
	echo "Usage: $0 {start|stop|reload|status}"

esac
