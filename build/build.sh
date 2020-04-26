#!/bin/bash

go build -o teel -ldflags="-X main.version=$(git describe --always --long)" cmd/teel/main.go
