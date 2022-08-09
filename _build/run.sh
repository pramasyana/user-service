INPUT_STRING=user
while [ "$INPUT_STRING" != "bye" ]
do
 PID=$(ps -ef | grep 'rsyslogd' | grep -v 'grep' | awk '{ printf $1 }')
 if [ -n "$PID" -a -e /proc/$PID ]; then
    echo "process exists"
 else
    /usr/sbin/rsyslogd
fi
 pr=$(ps -ef | grep 'user-service-linux' | grep -v 'grep' | awk '{ printf $1 }')
#echo $processId
 if [ -n "$pr" -a -e /proc/$PID ]; then
    echo "process exists"
 else
    /go/src/github.com/Bhinneka/user-service/user-service-linux &
fi
 sleep 60
 done
