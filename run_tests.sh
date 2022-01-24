#/bin/sh
. ./set_env_variables.sh
go test -parallel=1 -coverprofile=coverage.out
go tool cover -html=coverage.out
