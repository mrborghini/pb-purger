#!/bin/bash

go build -v -ldflags "-s -w"
./pb-purger