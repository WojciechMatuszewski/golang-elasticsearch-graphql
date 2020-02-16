#!/usr/bin/env bash

echo "starting up dynamo"
sh ./testenv/start_dynamo.sh
echo "starting elasticsearch"
sh ./testenv/start_elasticsearch.sh
