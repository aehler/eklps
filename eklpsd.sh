#!/bin/sh

PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/opt/eklps/bin/
DAEMON=/opt/eklps/bin/eklpsd
NAME=eklpsdaemon
DESC=eklpsdaemon
PID=/run/eklpsd.pid
LOGFILE=/home/admor/src/eklps/eklpsd.log

echo "$NAME..." >> $LOGFILE


do_start()
{
        # Return
        #   0 if daemon has been started
        #   1 if daemon was already running
        #   2 if daemon could not be started
        start-stop-daemon --start --quiet --pidfile $PID --exec $DAEMON --test > /dev/null \
                || return 1
	start-stop-daemon --start --pidfile /run/eklpsd.pid --exec $DAEMON --make-pidfile $PID --verbose >> $LOGFILE & \
        return 0 
}

#
# Function that stops the daemon/service
#
do_stop()
{
        # Return
        #   0 if daemon has been stopped
        #   1 if daemon was already stopped
        #   2 if daemon could not be stopped
        #   other if a failure occurred
        start-stop-daemon --stop --quiet --retry=TERM/30/KILL/5 --pidfile $PID
        RETVAL="$?"

	rm $PID
        sleep 1
        return "$RETVAL"
}

case "$1" in
        start)
                do_start
                case "$?" in
                        0|1) echo " started." && echo " started" >> $LOGFILE ;;
                        2) echo " couldn't start." && echo " couldn't start" >> $LOGFILE ;;
                esac
                ;;
        stop)
                do_stop
                case "$?" in
                        0|1) echo " stopped." && echo " stopped" >> $LOGFILE ;;
                        2) echo " couldn't stop." && echo " couldn't stop" >> $LOGFILE ;;
                esac
                ;;
        *)
                echo "Usage: eklpsd.sh  {start|stop}" >&2
                exit 3
                ;;
esac

exit 0
