#!/bin/zsh
cd backend
go build .
./fgo-calc-backend -config config.dev.json

