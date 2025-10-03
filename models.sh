#!/bin/sh

rm -r internal/handler/http/api/v1/models/*
/usr/bin/swagger generate model -f ./api/http/v1/swagger.yaml -t . -m internal/handler/http/api/v1/models