#!/bin/bash

export PGHOST=192.168.99.100
export PGUSER=one
export PGPASSWORD=two
export PGDATABASE=bench

go build
./go_db_bench