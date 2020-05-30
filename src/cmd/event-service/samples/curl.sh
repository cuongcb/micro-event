#!bin/bash

curl -X POST -H "Content-type: application/json" --data @samples/events.json localhost:8080/v1/event/
