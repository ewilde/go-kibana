#!/usr/bin/env bash
docker run -d --name elasticsearch -p 9200:9200 elastic-local:5.5.3
echo "waiting for 15 seconds for elastic to finish starting up"
sleep 15
docker run --name kibana --link elasticsearch -p 5601:5601 docker.elastic.co/kibana/kibana:5.5.3
