#!/usr/bin/env bash

if [ "" == "$(docker ps | grep elasticsearch)" ]; then
    echo "Starting Elasticsearch..."
    docker network create kuna
    docker run -d  --name elasticsearch --net kuna -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" --rm elasticsearch:7.6.0
    while ! nc -z localhost 9200; do
        echo "Waiting for services to come up..."
        sleep 2;
    done
    echo "Elasticsearch started."
    docker ps
else
    echo "Elasticsearch already running"
fi

