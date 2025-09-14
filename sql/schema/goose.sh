#! /bin/bash

# $1 must be up or down
goose postgres  postgres://postgres:postgres@localhost:5432/chirpy $1
