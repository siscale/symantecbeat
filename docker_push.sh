#!/bin/bash
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
make
docker build -t mariancraciunescu/symantecbeat .
docker push mariancraciunescu/symantecbeat:latest
