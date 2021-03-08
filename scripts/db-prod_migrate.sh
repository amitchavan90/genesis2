#!/bin/bash
source ./init/genesis-prod.env
./bin/migrate -database "postgres://$GENESIS_DATABASE_USER:$GENESIS_DATABASE_PASS@$GENESIS_DATABASE_HOST:$GENESIS_DATABASE_PORT/$GENESIS_DATABASE_NAME?sslmode=disable" -path ./migrations up