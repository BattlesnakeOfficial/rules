# Battlesnake Rules CLI

This tool allows running a Battlesnake game locally. There are several command-line options for running games, including the ability to send Battlesnake requests sequentially or concurrently, set a custom timeout, etc.

### Installation

Download precompiled binaries here: <br>
[https://github.com/BattlesnakeOfficial/rules/releases](https://github.com/BattlesnakeOfficial/rules/releases)

Install as a Go package. Requires Go 1.18 or higher. [[Download](https://golang.org/dl/)]
```
go install github.com/BattlesnakeOfficial/rules/cli/battlesnake@latest
```

Compile from source. Also requires Go 1.18 or higher.
```
git clone https://github.com/BattlesnakeOfficial/rules.git
cd rules
go build -o battlesnake ./cli/battlesnake/main.go
```

### Usage

Example command to run a game locally:
```
battlesnake play -W 11 -H 11 --name <SNAKE_NAME> --url <SNAKE_URL> -g solo -v
```

Complete usage documentation:
```
Usage:
  battlesnake play [flags]

Flags:
  -W, --width int                 Width of Board (default 11)
  -H, --height int                Height of Board (default 11)
  -n, --name stringArray          Name of Snake
  -u, --url stringArray           URL of Snake
  -t, --timeout int               Request Timeout (default 500)
  -s, --sequential                Use Sequential Processing
  -g, --gametype string           Type of Game Rules (default "standard")
  -m, --map string                Game map to use to populate the board (default "standard")
  -v, --viewmap                   View the Map Each Turn
  -c, --color                     Use color to draw the map
  -r, --seed int                  Random Seed (default 1656460409268690000)
  -d, --delay int                 Turn Delay in Milliseconds
  -D, --duration int              Minimum Turn Duration in Milliseconds
  -o, --output string             File path to output game state to. Existing files will be overwritten
      --browser                   View the game in the browser using the Battlesnake game board
      --board-url string          Base URL for the game board when using --browser (default "https://board.battlesnake.com")
      --foodSpawnChance int       Percentage chance of spawning a new food every round (default 15)
      --minimumFood int           Minimum food to keep on the board every turn (default 1)
      --hazardDamagePerTurn int   Health damage a snake will take when ending its turn in a hazard (default 14)
      --shrinkEveryNTurns int     In Royale mode, the number of turns between generating new hazards (default 25)
  -h, --help                      help for play

Global Flags:
      --config string   config file (default is $HOME/.battlesnake.yaml)
      --verbose         Enable debug logging
```

Battlesnake names and URLs will be paired together in sequence, for example:

```
battlesnake play --name Snake1 --url http://snake1-url-whatever --name Snake2 --url http://snake2-url-whatever
```

This will create a game with the following Battlesnakes:
* Snake1, http://snake1-url-whatever
* Snake2, http://snake2-url-whatever

Names are optional, and if you don't provide them UUIDs will be generated instead. However names are way easier to read and highly recommended!

URLs are technically optional too, but your Battlesnake will lose if the server is only sending move requests to http://example.com.

Example creating a 7x7 Standard game with two Battlesnakes:
```
battlesnake play --width 7 --height 7 --name Snake1 --url http://snake1-url-whatever --name Snake2 --url http://snake2-url-whatever
```

### Maps
The `map` command provides map information for use with the `play` command.

List all available maps using the `list` subcommand:
```
battlesnake map list
```
Display map information using the `info` subcommand:
```
battlesnake map info standard
Name: Standard
Author: Battlesnake
Description: Standard snake placement and food spawning
Version: 2
Min Players: 1
Max Players: 16
Board Sizes (WxH): 7x7 9x9 11x11 13x13 15x15 17x17 19x19 21x21 23x23 25x25
```

### Sample Output
```
$ battlesnake play --width 3 --height 3 --url http://redacted:4567/ --url http://redacted:4568/  --name Bob --name Sue
2022/04/10 04:16:08 [1]: State: &{1 3 3 [{2 2} {2 0}] [{8d8e31df-a1c1-4cac-a981-2453febe76ae [{1 0} {0 0} {0 0}] 99  0 } {cb5f5f33-3b35-477d-9c32-3e3fecd0e3c2 [{1 2} {0 2} {0 2}] 99  0 }] []}
2022/04/10 04:16:08 [2]: State: &{2 3 3 [{1 1}] [{8d8e31df-a1c1-4cac-a981-2453febe76ae [{2 0} {1 0} {0 0} {0 0}] 100  0 } {cb5f5f33-3b35-477d-9c32-3e3fecd0e3c2 [{2 2} {1 2} {0 2} {0 2}] 100  0 }] []}
2022/04/10 04:16:08 [3]: State: &{3 3 3 [{1 1}] [{8d8e31df-a1c1-4cac-a981-2453febe76ae [{2 1} {2 0} {1 0} {0 0}] 99 head-collision 0 cb5f5f33-3b35-477d-9c32-3e3fecd0e3c2} {cb5f5f33-3b35-477d-9c32-3e3fecd0e3c2 [{2 1} {2 2} {1 2} {0 2}] 99 head-collision 0 8d8e31df-a1c1-4cac-a981-2453febe76ae}] []}
2022/04/10 04:16:08 [DONE]: Game completed after 3 turns. It was a draw.
```

### Sample JSON Output

The default output logs the BoardState struct each turn, but is not well suited for parsing with a script.

Use the `--output` flag to write the game state in JSON format to a log file, for example:
```
battlesnake play --output out.log --name Snake1 --url http://snake1-url-whatever --name Snake2 --url http://snake2-url-whatever
```

The above command will write the game state each turn as a single line to `out.log` file. An example line:
```
{"game":{"id":"202b0f42-8d66-4adf-b29c-5ae1afd4c3cf","ruleset":{"name":"standard","version":"cli","settings":{"foodSpawnChance":15,"minimumFood":1,"hazardDamagePerTurn":14,"hazardMap":"","hazardMapAuthor":"","royale":{"shrinkEveryNTurns":0},"squad":{"allowBodyCollisions":false,"sharedElimination":false,"sharedHealth":false,"sharedLength":false}}},"timeout":500,"source":""},"turn":60,"board":{"height":11,"width":11,"snakes":[{"id":"55860e87-7c39-4911-8b67-aea861f27af6","name":"Snake1","latency":"0","health":98,"body":[{"x":10,"y":8},{"x":10,"y":9},{"x":10,"y":10},{"x":9,"y":10},{"x":9,"y":9},{"x":9,"y":8},{"x":8,"y":8},{"x":7,"y":8}],"head":{"x":10,"y":8},"length":8,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}},{"id":"fdb00735-1602-4a4c-bf23-2b704f80bbeb","name":"Snake2","latency":"0","health":92,"body":[{"x":9,"y":7},{"x":8,"y":7},{"x":7,"y":7},{"x":7,"y":6},{"x":7,"y":5},{"x":8,"y":5},{"x":9,"y":5},{"x":9,"y":4},{"x":8,"y":4}],"head":{"x":9,"y":7},"length":9,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}}],"food":[{"x":4,"y":6},{"x":0,"y":9},{"x":4,"y":5}],"hazards":[]},"you":{"id":"55860e87-7c39-4911-8b67-aea861f27af6","name":"Snake1","latency":"0","health":98,"body":[{"x":10,"y":8},{"x":10,"y":9},{"x":10,"y":10},{"x":9,"y":10},{"x":9,"y":9},{"x":9,"y":8},{"x":8,"y":8},{"x":7,"y":8}],"head":{"x":10,"y":8},"length":8,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}}}
```

To get the request data sent to each snake, use the `--debug-requests` flag (note this contains the `you` field which is missing in data generated using the `--output` flag):
```
2022/04/10 04:41:16 POST http://localhost:8080/move: {"game":{"id":"0baa4367-b1ee-40c7-96c8-34227b88af24","ruleset":{"name":"standard","version":"cli","settings":{"foodSpawnChance":15,"minimumFood":1,"hazardDamagePerTurn":14,"hazardMap":"","hazardMapAuthor":"","royale":{"shrinkEveryNTurns":0},"squad":{"allowBodyCollisions":false,"sharedElimination":false,"sharedHealth":false,"sharedLength":false}}},"timeout":500,"source":""},"turn":5,"board":{"height":11,"width":11,"snakes":[{"id":"5bddff9f-d3ff-458c-b0f5-df81a830b5d8","name":"Snake1","latency":"0","health":96,"body":[{"x":5,"y":7},{"x":4,"y":7},{"x":4,"y":8}],"head":{"x":5,"y":7},"length":3,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}},{"id":"f76e8994-6457-49f0-9102-6a1bcfee5695","name":"Snake2","latency":"0","health":96,"body":[{"x":6,"y":6},{"x":7,"y":6},{"x":7,"y":5}],"head":{"x":6,"y":6},"length":3,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}}],"food":[{"x":6,"y":10},{"x":10,"y":4},{"x":5,"y":5},{"x":9,"y":0}],"hazards":[]},"you":{"id":"f76e8994-6457-49f0-9102-6a1bcfee5695","name":"Snake2","latency":"0","health":96,"body":[{"x":6,"y":6},{"x":7,"y":6},{"x":7,"y":5}],"head":{"x":6,"y":6},"length":3,"shout":"","squad":"","customizations":{"color":"#03d3fc","head":"beluga","tail":"bolt"}}}
```

### Sample Output (With ASCII Board)
```
$ battlesnake play --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --url http://redacted:4567/ --name Snake1 --name Snake2 --name Snake3 --name Snake4 --name Snake5 --name Snake6 --name Snake7 --name Snake8 --width 13 --height 13 --timeout 1000 --viewmap
2020/11/01 21:56:50 [1]
Hazards ░: []
Food ⚕: [{12 10} {8 4} {10 10} {9 11} {8 2} {9 6} {1 11} {9 12}]
Snake1 ■: {cca4652d-26b5-4c09-9d05-dbe01d24626c [{0 3} {0 2} {0 2}] 99  }
Snake2 ⌀: {aff9c973-fc49-4b1e-b219-1a4d2023d76b [{8 1} {8 0} {8 0}] 99  }
Snake3 ●: {03c90cd1-62dc-4393-8c1c-185601cfe00a [{8 11} {7 11} {7 11}] 99  }
Snake4 ⍟: {c112965a-0b5a-45f6-b4de-88a68f3373e3 [{3 2} {3 1} {3 1}] 99  }
Snake5 ◘: {f4810018-cd5e-44bd-b871-3f6afd84250f [{5 2} {5 1} {5 1}] 99  }
Snake6 ☺: {50c2933a-c4e4-4727-bc2e-e54778129308 [{7 4} {6 4} {6 4}] 99  }
Snake7 □: {f760d89c-e503-45c4-9453-0284ed172120 [{1 12} {2 12} {2 12}] 99  }
Snake8 ☻: {8e42531e-bd55-4d76-8d3a-e0eda0578812 [{4 7} {4 6} {4 6}] 99  }
◦□□◦◦◦◦◦◦⚕◦◦◦
◦⚕◦◦◦◦◦●●⚕◦◦◦
◦◦◦◦◦◦◦◦◦◦⚕◦⚕
◦◦◦◦◦◦◦◦◦◦◦◦◦
◦◦◦◦◦◦◦◦◦◦◦◦◦
◦◦◦◦☻◦◦◦◦◦◦◦◦
◦◦◦◦☻◦◦◦◦⚕◦◦◦
◦◦◦◦◦◦◦◦◦◦◦◦◦
◦◦◦◦◦◦☺☺⚕◦◦◦◦
■◦◦◦◦◦◦◦◦◦◦◦◦
■◦◦⍟◦◘◦◦⚕◦◦◦◦
◦◦◦⍟◦◘◦◦⌀◦◦◦◦
◦◦◦◦◦◦◦◦⌀◦◦◦◦
```
