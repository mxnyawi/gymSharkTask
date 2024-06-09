#!/bin/sh

# Create a log file
LOGFILE=/root/startup.log

# Load environment variables from .env file
. /root/config.env >> $LOGFILE 2>&1

# Trim leading and trailing whitespace from variables
USERNAME=$(echo "$USERNAME" | tr -d '[:space:]')
PASSWORD=$(echo "$PASSWORD" | tr -d '[:space:]')

# Trim leading and trailing whitespace from variables
USERNAME=$(echo -n "$USERNAME" | xargs)
PASSWORD=$(echo -n "$PASSWORD" | xargs)
# Wait for the Couchbase Server to start
echo "Waiting for the Couchbase Server to start..." >> $LOGFILE
until $(curl --output /dev/null --silent --head --fail http://db:8091); do
  printf '.'
  sleep 5
done

echo "User: $USERNAME" >> $LOGFILE
echo "Password: $PASSWORD" >> $LOGFILE

# Initialize the cluster
echo "Initializing the cluster..." >> $LOGFILE
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST http://db:8091/clusterInit \
  -d username=$USERNAME \
  -d password=$PASSWORD \
  -d clusterName=myCluster \
  -d services=kv,index,n1ql,fts \
  -d port=SAME \
  -d allowedHosts=*)
echo "Response: $RESPONSE" >> $LOGFILE


echo "Startup script finished." >> $LOGFILE