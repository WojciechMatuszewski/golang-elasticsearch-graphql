#!/usr/bin/env bash

if [ "" == "$(docker ps | grep dynamodb-local)" ]; then
    echo "Starting DynamoDB..."
    docker network create kuna
    docker run --net kuna -d --rm --name dynamodb-local -p 8000:8000 amazon/dynamodb-local
    while ! nc -z localhost 8000; do
        echo "Waiting for services to come up..."
        sleep 2;
    done
    echo "DynamoDB started."
    docker ps
else
    echo "DynamoDB already running, deleting tables..."
    go run $(pwd)/tool/clear_local_dynamo.go
fi

