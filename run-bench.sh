#!/bin/bash

export PGHOST=192.168.99.100
export PGUSER=one
export PGPASSWORD=two
export PGDATABASE=bench

go test -test.bench=. -test.benchmem