# Battlesnake Rules CLI

This tool allows running a Battlesnake game locally. There are several command-line options for running games, including the ability to send Battlesnake requests sequentially or concurrently, set a custom timeout, etc.

### Installation

Download precompiled binaries here: <br>
[https://github.com/BattlesnakeOfficial/rules/releases](https://github.com/BattlesnakeOfficial/rules/releases)

Install as a Go package. Requires Go 1.13 or higher. [[Download](https://golang.org/dl/)]
```
go get github.com/BattlesnakeOfficial/rules/cli/battlesnake
```

Compile from source. Also requires Go 1.13 or higher.
```
git clone git@github.com:BattlesnakeOfficial/rules.git
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
  -g, --gametype string     Type of Game Rules (default "standard")
  -H, --height int32        Height of Board (default 11)
  -h, --help                help for play
  -n, --name stringArray    Name of Snake
  -r, --seed int            Random Seed (default 1607708568137187300)
  -s, --sequential          Use Sequential Processing
  -S, --squad stringArray   Squad of Snake
  -t, --timeout int32       Request Timeout (default 500)
  -u, --url stringArray     URL of Snake
  -v, --viewmap             View the Map Each Turn
  -W, --width int32         Width of Board (default 11)

Global Flags:
      --config string   config file (default is $HOME/.battlesnake.yaml)
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

### Sample Output
```
$ battlesnake play --width 3 --height 3 --url http://redacted:4567/ --url http://redacted:4568/  --name Bob --name Sue
2020/10/31 22:05:56 [1]: State: &{3 3 [{1 0}] [{e74892ba-9f0c-4e96-9bde-1a9efaff0ab4 [{0 1} {0 2} {0 2} {0 2}] 100  } {89e20d26-7da7-4964-b0ae-148c8f60f7ee [{2 1} {2 2} {2 2} {2 2}] 100  }]} OutOfBounds: []
2020/10/31 22:05:56 [2]: State: &{3 3 [{1 0}] [{e74892ba-9f0c-4e96-9bde-1a9efaff0ab4 [{0 0} {0 1} {0 2} {0 2}] 99  } {89e20d26-7da7-4964-b0ae-148c8f60f7ee [{2 0} {2 1} {2 2} {2 2}] 99  }]} OutOfBounds: []
2020/10/31 22:05:56 [3]: State: &{3 3 [{1 2}] [{e74892ba-9f0c-4e96-9bde-1a9efaff0ab4 [{1 0} {0 0} {0 1} {0 2} {0 2}] 100 head-collision 89e20d26-7da7-4964-b0ae-148c8f60f7ee} {89e20d26-7da7-4964-b0ae-148c8f60f7ee [{1 0} {2 0} {2 1} {2 2} {2 2}] 100 head-collision e74892ba-9f0c-4e96-9bde-1a9efaff0ab4}]} OutOfBounds: []
2020/10/31 22:05:56 [DONE]: Game completed after 3 turns. It was a draw.
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
