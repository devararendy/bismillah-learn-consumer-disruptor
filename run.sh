#!/bin/bash

set -o allexport
source consumer.env
set +o allexport

go run . $@