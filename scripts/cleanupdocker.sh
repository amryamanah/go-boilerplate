#!/bin/bash

docker rm -f $(docker ps -a -q)
docker volume rm boilerplate-volume