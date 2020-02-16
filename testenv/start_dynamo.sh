#!/usr/bin/env bash

if [ "" == "$(docker ps | grep dynamodb-local)" ]; then
    echo "Starting DynamoDB..."
    docker network create kuna
    docker run --net kuna -d --rm --name dynamodb-local -p 8000:8000 -v $(pwd)/local/dynamodb:/data/ amazon/dynamodb-local -jar DynamoDBLocal.jar -sharedDb -inMemory
    while ! nc -z localhost 8000; do
        echo "Waiting for services to come up..."
        sleep 2;
    done
    echo "DynamoDB started."
    docker ps
else
    echo "DynamoDB already running"
fi

