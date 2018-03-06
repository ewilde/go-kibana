#!/usr/bin/env bash
docker run -d --name elasticsearch -e ELASTIC_PASSWORD=changeme -p 9200:9200 elastic-local-platinum:6.2.2
echo "waiting for 15 seconds for elastic to finish starting up"
sleep 15
docker run --name kibana --link elasticsearch -p 5601:5601 docker.elastic.co/kibana/kibana-x-pack:6.2.2
