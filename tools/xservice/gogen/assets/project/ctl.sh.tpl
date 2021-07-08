#!/bin/bash
set -e

# uncomment & config for service discovery
# export XSERVICE_ETCD=127.0.0.1:2379
# export XSERVICE_ETCD_USER=root
# export XSERVICE_ETCD_PASSWORD=123456

appName="{{.Name}}"

##########################

cd $(dirname $0)

getpid() {
  echo $(ps -ef | grep -E "\s\.?\/${appName}" | awk '{print $2}')
}

status() {
  local pid=$(getpid)
  if [ ! -z "$pid" ]; then
    echo "$appName is runing pid: $pid"

    echo ""
    echo "ps status"
    ps -p "$pid" -o "user,pid,ppid,lstart,etime,rss,%mem,%cpu,command"
  else
    echo "$appName is not runing"
  fi
}

start() {
  local pid=$(getpid)
  if [ -z $pid ]; then
    echo "starting $appName"
    # disable stdlog, bug keep err log for track panic issues.
    XSERVICE_DISABLE_STDOUT=true ./$appName &> .err.log &
    echo "$appName is runing pid: $!"
  else
    echo "$appName is already runing pid:$pid"
  fi
}

stop() {
  echo "stopping $appName"
  local pid=$(getpid)
  if [ ! -z "$pid" ]; then
    kill "$pid"
    sleep 2s
    pid=$(getpid)
    if [ ! -z "$pid" ]; then
      echo "$appName is still runing, try force stop!"
      kill -9 "$pid"
      sleep 2s
    fi
  fi
  echo "$appName stopped"
}

reload() {
  local pid=$(getpid)
  if [ ! -z "$pid" ]; then
    kill -USR2 "$pid"
    echo "$appName reloaded"
  else
    echo "$appName is not runing"
  fi
}

startOrReload() {
  local pid=$(getpid)
  if [ -z $pid ]; then
    start ${@:2}
  else
    echo "reloading $appName"
    reload
  fi
}

version() {
  ./$appName -version
}

help() {
  ./$appName -h
}

case "$1" in
status)
  status
  ;;
start)
  start ${@:2}
  ;;
stop)
  stop
  ;;
restart)
  stop
  start
  ;;
reload)
  reload
  ;;
startOrReload)
  startOrReload ${@:2}
  ;;
version)
  version
  ;;
help)
  help
  ;;
*)
  cat <<EOF

usage: $0 action [args]

  $0 script for simple controlling $appName

  actions supported:
    status          show $appName current status
    start           start $appName
                      see more arguments please execute: $0 help
    stop            stop $appName
    restart         stop then start
    reload          send reload signal (SIGUSR2) via kill command
    startOrReload   start $appName or just reload it if it is already runing
    version         show version
    help            show $appName help

EOF
  ;;
esac

exit 0
