#!/bin/bash



start=8090
for i in {1..30}; do
    mkdir -p ./tmp/musician${i};

    ENV_REST_ENDPOINTADDRESS=0.0.0.0:$start \
	ENV_CONDUCTOR_ADVERTISEADDR=http://localhost:${start} \
	ENV_CONDUCTOR_CONDUCTORADDR=http://localhost:8080 \
	./build/musician -d=tmp/musician${i} > tmp/musician${i}/logs.txt &

    start=$((start + 1))
done
