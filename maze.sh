#!/bin/bash

set -ex

go build -o battlesnake ./cli/battlesnake/main.go
./battlesnake play -W 25 -H 25 --name Frank --url http://0:3000/famished-frank --gametype solo -v -d 150 --map coreyja_maze --hazardDamagePerTurn 100
