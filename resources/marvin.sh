#!/bin/sh
# Startup script for Debian derivatives
#
# the skeleton script in debian 6 does not work properly in package scripts.
# the return/exit codes cause {pre|post}{inst|rm} to fail regardless of the
# script completion status.  this script exits explicitly.
#
# the skeleton script also does not work properly with python applications,
# as the lsb tools cannot distinguish between the python interpreter and
# the python code that was invoked.  this script uses ps and grep to look
# for the application signature instead of using the lsb tools to determine
# whether the app is running.
#
### BEGIN INIT INFO
# Provides:          marvin
# Required-Start:    $local_fs $syslog networking
# Required-Stop:     $local_fs $syslog networking
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: marvin server
# Description:       marvin server for alexa powered smarthome
### END INIT INFO

# Do NOT "set -e"

# PATH should only include /usr/* if it runs after the mountnfs.sh script
PATH=/sbin:/usr/sbin:/bin:/usr/bin
NAME=marvin
HOME=/home/pi/$NAME
USER="pi"
DESC="marvin server for alexa powered smarthome"
SCRIPTNAME=/etc/init.d/$NAME
PIDFILE=$HOME/$NAME.pid
DAEMON=$HOME/$NAME
DAEMON_ARGS="-config $HOME/config.json" 

# Exit if the package is not installed
[ -x "$DAEMON" ] || exit 0

# Load the VERBOSE setting and other rcS variables
. /lib/init/vars.sh

# Define LSB log_* functions.
# Depend on lsb-base (>= 3.0-6) to ensure that this file is present.
. /lib/lsb/init-functions

# start the daemon/service
#   0 if daemon has been started
#   1 if daemon was already running
#   2 if daemon could not be started
do_start() {
    NPROC=$(count_procs)
    if [ $NPROC != 0 ]; then
	return 1
    fi
    start-stop-daemon --start --background --chuid $USER --chdir $HOME --pidfile $PIDFILE --make-pidfile --exec $DAEMON -- $DAEMON_ARGS || return 2
    return 0
}

# stop the daemon/service
#   0 if daemon has been stopped
#   1 if daemon was already stopped
#   2 if daemon could not be stopped
#   other if a failure occurred
do_stop() {
    # bail out if the app is not running
    NPROC=$(count_procs)
    if [ $NPROC = 0 ]; then
        return 1
    fi
    # bail out if there is no pid file
    if [ ! -f $PIDFILE ]; then
        return 1
    fi
    start-stop-daemon --stop --pidfile $PIDFILE
    # we cannot trust the return value from start-stop-daemon
    RETVAL=2
    c=0
    while [ $c -lt 24 -a "$RETVAL" = "2" ]; do
        c=`expr $c + 1`
        # process may not really have completed, so check it
        NPROC=$(count_procs)
        if [ $NPROC = 0 ]; then
	    RETVAL=0
        else
            echo -n "."
            sleep 5
        fi
    done
    if [ "$RETVAL" = "0" -o "$RETVAL" = "1" ]; then
        # delete the pid file just in case
        rm -f $PIDFILE
    fi
    return "$RETVAL"
}

count_procs() {
    NPROC=`ps ax | grep -v grep | grep $DAEMON | wc -l`
    echo $NPROC
}

RETVAL=0
case "$1" in
    start)
	log_daemon_msg "Starting $DESC" "$NAME"
	do_start
	case "$?" in
	    0) log_end_msg 0; RETVAL=0 ;;
	    1) log_action_cont_msg " already running" && log_end_msg 0; RETVAL=0 ;;
	    2) log_end_msg 1; RETVAL=1 ;;
	esac
	;;
    stop)
	log_daemon_msg "Stopping $DESC" "$NAME"
	do_stop
	case "$?" in
	    0) log_end_msg 0; RETVAL=0 ;;
	    1) log_action_cont_msg " not running" && log_end_msg 0; RETVAL=0 ;;
	    2) log_end_msg 1; RETVAL=1 ;;
	esac
	;;
    status)
        NPROC=$(count_procs)
        if [ $NPROC -gt 1 ]; then
            MSG="running multiple times"
	elif [ $NPROC = 1 ]; then
	    MSG="running"
	else
	    MSG="not running"
	fi
	log_daemon_msg "Status of $DESC" "$MSG"
	log_end_msg 0
	RETVAL=0
	;;
    restart)
	log_daemon_msg "Restarting $DESC" "$NAME"
	do_stop
	case "$?" in
	    0|1)
		do_start
		case "$?" in
		    0) log_end_msg 0; RETVAL=0 ;;
		    1) log_end_msg 1; RETVAL=1 ;; # Old process still running
		    *) log_end_msg 1; RETVAL=1 ;; # Failed to start
		esac
		;;
	    *)
	  	# Failed to stop
		log_end_msg 1
		RETVAL=1
		;;
	esac
	;;
    *)
	echo "Usage: $SCRIPTNAME {start|stop|status|restart}"
	exit 3
	;;
esac

exit $RETVAL

