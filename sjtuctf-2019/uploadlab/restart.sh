#!/bin/bash
while true
do
    docker-compose down
    docker-compose up -d
    sleep 900
done