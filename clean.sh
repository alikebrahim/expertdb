#!/bin/bash

rm -rf ./logs/
rm -rf ./tmp/
rm ./db/sqlite/main.db
goose -dir ./db/migrations/sqlite/ sqlite3 ./db/sqlite/main.db up
