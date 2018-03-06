#!/usr/bin/env sh
echo waiting for elastic to start
while ! curl --output /dev/null --silent --head --fail http://elasticsearch:9200 -u elastic:changeme; do sleep 1 && echo -n .; done;
echo waiting elastic started
makelogs --host "elasticsearch:9200" --auth "elastic:changeme"
echo finished creating logs
