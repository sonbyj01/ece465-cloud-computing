#!/bin/sh
# Simple way to run multiple clients and a single server on the same host
# not intended to get any speedup, mostly to ensure that overall handshake
# and general algorithm work. This uses a simple four-node network on predefined
# ports on localhost.

# usage: ./sample_run.sh [GRAPH_FILE]

GRAPH_FILE=$1
CONFIG_NODES_FILE=/tmp/config.nodes

# cleanup from previous run; this should be run after this script
# if client_latest and/or server_latest did not properly terminate
pkill client_latest
pkill server_latest
rm -f $CONFIG_NODES_FILE

sleep 0.25

# write config nodes
cat <<EOF >$CONFIG_NODES_FILE
127.0.0.1:5000
127.0.0.1:5001
127.0.0.1:5002
127.0.0.1:5003
EOF

# start clients
target/client_latest --port 5000 &
target/client_latest --port 5001 &
target/client_latest --port 5002 &
target/client_latest --port 5003 &

# wait a little...
sleep 0.25

# start server
target/server_latest --config $CONFIG_NODES_FILE --graph $GRAPH_FILE &