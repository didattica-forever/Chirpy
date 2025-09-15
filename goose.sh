#! /bin/bash

# $1 must be up or down
cd sql/schema || { echo "Wrong directory, aborted!"; exit 1; }
goose postgres  postgres://postgres:postgres@localhost:5432/chirpy $1
