#!/bin/bash
host=$1
port=$2
remote="circleci@${host}"
echo "deploying..."
ssh -p $port $remote sudo systemctl stop hook2xmpp;
RETVAL=$?
[ $RETVAL -ne 0 ] && exit 1
scp -q -P $port ~/.go_workspace/bin/hook2xmpp $remote:~/bin/hook2xmpp;
RETVAL=$?
ssh -p $port $remote sudo systemctl start hook2xmpp;
[ $RETVAL -eq 0 ] && RETVAL=$?
[ $RETVAL -ne 0 ] && exit 1
echo "deployed"
