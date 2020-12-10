# Battlesnake CLI

This tool allows running a Battlesnake game locally. There are several command-line options for running games, including the ability to send Battlesnake requests sequentially or concurrently, set a custom timeout, etc.

## Installation

The Battlesnake CLI requires Go 1.13 or later.

```
go get github.com/BattlesnakeOfficial/rules/cli/battlesnake
```

## Usage

```
Usage:
  battlesnake play [flags]

Flags:
  -g, --gametype string     Type of Game Rules (default "standard")
  -H, --height int32        Height of Board (default 11)
  -h, --help                help for play
  -n, --name stringArray    Name of Snake
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

### Sample Map Output
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

### Sample Solo Game
```
$ battlesnake play --url http://redacted:4567/ --name Bob --width 3 --height 3 --timeout 500 --gametype solo
2020/10/31 22:02:58 [1]: State: &{3 3 [{2 2}] [{cc8831e8-d517-4216-a8d8-a64243decada [{1 2} {0 2} {0 2}] 99  }]} OutOfBounds: []
2020/10/31 22:02:58 [2]: State: &{3 3 [{2 1}] [{cc8831e8-d517-4216-a8d8-a64243decada [{2 2} {1 2} {0 2} {0 2}] 100  }]} OutOfBounds: []
2020/10/31 22:02:59 [3]: State: &{3 3 [{0 1}] [{cc8831e8-d517-4216-a8d8-a64243decada [{2 1} {2 2} {1 2} {0 2} {0 2}] 100  }]} OutOfBounds: []
2020/10/31 22:02:59 [4]: State: &{3 3 [{0 1}] [{cc8831e8-d517-4216-a8d8-a64243decada [{1 1} {2 1} {2 2} {1 2} {0 2}] 99  }]} OutOfBounds: []
2020/10/31 22:02:59 [5]: State: &{3 3 [{0 2}] [{cc8831e8-d517-4216-a8d8-a64243decada [{0 1} {1 1} {2 1} {2 2} {1 2} {1 2}] 100  }]} OutOfBounds: []
2020/10/31 22:02:59 [6]: State: &{3 3 [{2 0}] [{cc8831e8-d517-4216-a8d8-a64243decada [{0 2} {0 1} {1 1} {2 1} {2 2} {1 2} {1 2}] 100  }]} OutOfBounds: []
2020/10/31 22:02:59 [7]: State: &{3 3 [{2 0} {0 0}] [{cc8831e8-d517-4216-a8d8-a64243decada [{0 1} {0 2} {0 1} {1 1} {2 1} {2 2} {1 2}] 99 snake-self-collision cc8831e8-d517-4216-a8d8-a64243decada}]} OutOfBounds: []
2020/10/31 22:02:59 [DONE]: Game completed after 7 turns. It was a draw.
```

### Sample Squad Game
```
-$ battlesnake play --url http://redacted:4567/ --name Bob --squad A --url http://redacted:4567/ --name Sue --squad A --url http://redacted:4567/ --name Jim --squad B --url http://redacted:4567/ --name Francine --squad B --width 5 --height 5 --gametype squad
2020/10/31 22:14:27 [1]: State: &{5 5 [{2 4} {4 1} {4 3} {1 4} {0 2}] [{92a1bd60-8f8d-4adb-8468-e8eb1028b7f0 [{3 0} {4 0} {4 0}] 99  } {25c5607c-a2da-421e-84c3-e2a040cffae5 [{1 2} {1 1} {1 1}] 99  } {9dc22d73-3631-43cc-9472-a2ff074bc4a1 [{3 2} {4 2} {4 2}] 99  } {54157a58-2e07-4f84-b035-6d6df73d751a [{3 4} {4 4} {4 4}] 99  }]} OutOfBounds: []
2020/10/31 22:14:28 [2]: State: &{5 5 [{4 1} {4 3} {1 4}] [{92a1bd60-8f8d-4adb-8468-e8eb1028b7f0 [{2 0} {3 0} {4 0} {4 0}] 100  } {25c5607c-a2da-421e-84c3-e2a040cffae5 [{0 2} {1 2} {1 1} {1 1}] 100  } {9dc22d73-3631-43cc-9472-a2ff074bc4a1 [{3 3} {3 2} {4 2} {4 2}] 100  } {54157a58-2e07-4f84-b035-6d6df73d751a [{2 4} {3 4} {4 4} {4 4}] 100  }]} OutOfBounds: []
2020/10/31 22:14:28 [3]: State: &{5 5 [{4 1}] [{92a1bd60-8f8d-4adb-8468-e8eb1028b7f0 [{2 1} {2 0} {3 0} {4 0}] 99  } {25c5607c-a2da-421e-84c3-e2a040cffae5 [{0 3} {0 2} {1 2} {1 1}] 99  } {9dc22d73-3631-43cc-9472-a2ff074bc4a1 [{4 3} {3 3} {3 2} {4 2} {4 2}] 100  } {54157a58-2e07-4f84-b035-6d6df73d751a [{1 4} {2 4} {3 4} {4 4} {4 4}] 100  }]} OutOfBounds: []
2020/10/31 22:14:28 [4]: State: &{5 5 [{4 1}] [{92a1bd60-8f8d-4adb-8468-e8eb1028b7f0 [{3 1} {2 1} {2 0} {3 0}] 98  } {25c5607c-a2da-421e-84c3-e2a040cffae5 [{0 4} {0 3} {0 2} {1 2}] 98  } {9dc22d73-3631-43cc-9472-a2ff074bc4a1 [{4 4} {4 3} {3 3} {3 2} {4 2}] 99  } {54157a58-2e07-4f84-b035-6d6df73d751a [{1 3} {1 4} {2 4} {3 4} {4 4}] 99  }]} OutOfBounds: []
2020/10/31 22:14:28 [5]: State: &{5 5 [{1 0}] [{92a1bd60-8f8d-4adb-8468-e8eb1028b7f0 [{4 1} {3 1} {2 1} {2 0} {2 0}] 100 squad-eliminated } {25c5607c-a2da-421e-84c3-e2a040cffae5 [{0 3} {0 4} {0 3} {0 2}] 97 snake-self-collision 25c5607c-a2da-421e-84c3-e2a040cffae5} {9dc22d73-3631-43cc-9472-a2ff074bc4a1 [{3 4} {4 4} {4 3} {3 3} {3 2}] 98  } {54157a58-2e07-4f84-b035-6d6df73d751a [{2 3} {1 3} {1 4} {2 4} {3 4}] 98  }]} OutOfBounds: []
2020/10/31 22:14:28 [DONE]: Game completed after 5 turns. Francine is the winner.
```

### Sample Royale Game
```
$ battlesnake play --url http://redacted:4567/ --url http://redacted:4567/ --name Bob --name Sue --width 7 --height 7 --timeout 800 --gametype royale
2020/10/31 22:16:44 [1]: State: &{7 7 [{4 0} {0 0} {3 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{4 1} {5 1} {5 1}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{0 1} {1 1} {1 1}] 99  }]} OutOfBounds: []
2020/10/31 22:16:44 [2]: State: &{7 7 [{3 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{4 0} {4 1} {5 1} {5 1}] 100  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{0 0} {0 1} {1 1} {1 1}] 100  }]} OutOfBounds: []
2020/10/31 22:16:45 [3]: State: &{7 7 [{3 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 0} {4 0} {4 1} {5 1}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 0} {0 0} {0 1} {1 1}] 99  }]} OutOfBounds: []
2020/10/31 22:16:45 [4]: State: &{7 7 [{3 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 1} {3 0} {4 0} {4 1}] 98  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 1} {1 0} {0 0} {0 1}] 98  }]} OutOfBounds: []
2020/10/31 22:16:45 [5]: State: &{7 7 [{3 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 2} {3 1} {3 0} {4 0}] 97  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 2} {1 1} {1 0} {0 0}] 97  }]} OutOfBounds: []
2020/10/31 22:16:45 [6]: State: &{7 7 [{0 4}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 3} {3 2} {3 1} {3 0} {3 0}] 100  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 3} {1 2} {1 1} {1 0}] 96  }]} OutOfBounds: []
2020/10/31 22:16:45 [7]: State: &{7 7 [{0 4}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 3} {3 3} {3 2} {3 1} {3 0}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 4} {1 3} {1 2} {1 1}] 95  }]} OutOfBounds: []
2020/10/31 22:16:45 [8]: State: &{7 7 [{1 1}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 4} {2 3} {3 3} {3 2} {3 1}] 98  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{0 4} {1 4} {1 3} {1 2} {1 2}] 100  }]} OutOfBounds: []
2020/10/31 22:16:45 [9]: State: &{7 7 [{1 1}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 5} {2 4} {2 3} {3 3} {3 2}] 97  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{0 3} {0 4} {1 4} {1 3} {1 2}] 99  }]} OutOfBounds: []
2020/10/31 22:16:45 [10]: State: &{7 7 [{1 1}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 5} {2 5} {2 4} {2 3} {3 3}] 96  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{0 2} {0 3} {0 4} {1 4} {1 3}] 98  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:45 [11]: State: &{7 7 [{1 1}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 4} {3 5} {2 5} {2 4} {2 3}] 95  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 2} {0 2} {0 3} {0 4} {1 4}] 97  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:45 [12]: State: &{7 7 [{1 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 3} {3 4} {3 5} {2 5} {2 4}] 94  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{1 1} {1 2} {0 2} {0 3} {0 4} {0 4}] 100  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [13]: State: &{7 7 [{1 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 3} {3 3} {3 4} {3 5} {2 5}] 93  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{2 1} {1 1} {1 2} {0 2} {0 3} {0 4}] 99  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [14]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{1 3} {2 3} {3 3} {3 4} {3 5} {3 5}] 100  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{3 1} {2 1} {1 1} {1 2} {0 2} {0 3}] 98  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [15]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{1 4} {1 3} {2 3} {3 3} {3 4} {3 5}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{3 0} {3 1} {2 1} {1 1} {1 2} {0 2}] 97  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [16]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 4} {1 4} {1 3} {2 3} {3 3} {3 4}] 98  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{4 0} {3 0} {3 1} {2 1} {1 1} {1 2}] 96  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [17]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 4} {2 4} {1 4} {1 3} {2 3} {3 3}] 97  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{4 1} {4 0} {3 0} {3 1} {2 1} {1 1}] 95  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [18]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 3} {3 4} {2 4} {1 4} {1 3} {2 3}] 96  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{4 2} {4 1} {4 0} {3 0} {3 1} {2 1}] 94  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [19]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 3} {3 3} {3 4} {2 4} {1 4} {1 3}] 95  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 2} {4 2} {4 1} {4 0} {3 0} {3 1}] 93  }]} OutOfBounds: [{6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [20]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 2} {2 3} {3 3} {3 4} {2 4} {1 4}] 94  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 1} {5 2} {4 2} {4 1} {4 0} {3 0}] 92  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [21]: State: &{7 7 [{2 0}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 1} {2 2} {2 3} {3 3} {3 4} {2 4}] 93  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 0} {5 1} {5 2} {4 2} {4 1} {4 0}] 90  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:46 [22]: State: &{7 7 [{4 4}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{2 0} {2 1} {2 2} {2 3} {3 3} {3 4} {3 4}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{6 0} {5 0} {5 1} {5 2} {4 2} {4 1}] 88  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [23]: State: &{7 7 [{4 4} {4 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 0} {2 0} {2 1} {2 2} {2 3} {3 3} {3 4}] 97  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{6 1} {6 0} {5 0} {5 1} {5 2} {4 2}] 86  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [24]: State: &{7 7 [{4 4} {4 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 1} {3 0} {2 0} {2 1} {2 2} {2 3} {3 3}] 96  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{6 2} {6 1} {6 0} {5 0} {5 1} {5 2}] 84  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [25]: State: &{7 7 [{4 4} {4 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 2} {3 1} {3 0} {2 0} {2 1} {2 2} {2 3}] 95  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 2} {6 2} {6 1} {6 0} {5 0} {5 1}] 83  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [26]: State: &{7 7 [{4 4} {4 3}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{3 3} {3 2} {3 1} {3 0} {2 0} {2 1} {2 2}] 94  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 1} {5 2} {6 2} {6 1} {6 0} {5 0}] 82  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [27]: State: &{7 7 [{4 4}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{4 3} {3 3} {3 2} {3 1} {3 0} {2 0} {2 1} {2 1}] 100  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{4 1} {5 1} {5 2} {6 2} {6 1} {6 0}] 81  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [28]: State: &{7 7 [{3 6}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{4 4} {4 3} {3 3} {3 2} {3 1} {3 0} {2 0} {2 1} {2 1}] 100  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{4 0} {4 1} {5 1} {5 2} {6 2} {6 1}] 79  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:47 [29]: State: &{7 7 [{3 6}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{5 4} {4 4} {4 3} {3 3} {3 2} {3 1} {3 0} {2 0} {2 1}] 99  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 0} {4 0} {4 1} {5 1} {5 2} {6 2}] 77  }]} OutOfBounds: [{0 0} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:48 [30]: State: &{7 7 [{3 6} {4 6}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{5 3} {5 4} {4 4} {4 3} {3 3} {3 2} {3 1} {3 0} {2 0}] 98  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{6 0} {5 0} {4 0} {4 1} {5 1} {5 2}] 75  }]} OutOfBounds: [{0 0} {0 1} {0 2} {0 3} {0 4} {0 5} {0 6} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:48 [31]: State: &{7 7 [{3 6} {4 6}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{5 2} {5 3} {5 4} {4 4} {4 3} {3 3} {3 2} {3 1} {3 0}] 97  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{6 1} {6 0} {5 0} {4 0} {4 1} {5 1}] 73  }]} OutOfBounds: [{0 0} {0 1} {0 2} {0 3} {0 4} {0 5} {0 6} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:48 [32]: State: &{7 7 [{3 6} {4 6}] [{07ba7c7a-6533-4682-8769-fc2666b155c5 [{5 1} {5 2} {5 3} {5 4} {4 4} {4 3} {3 3} {3 2} {3 1}] 96  } {7b33dbd3-c9c5-461c-8d66-29ca715a9e43 [{5 1} {6 1} {6 0} {5 0} {4 0} {4 1}] 72 head-collision 07ba7c7a-6533-4682-8769-fc2666b155c5}]} OutOfBounds: [{0 0} {0 1} {0 2} {0 3} {0 4} {0 5} {0 6} {1 0} {2 0} {3 0} {4 0} {5 0} {6 0} {6 1} {6 2} {6 3} {6 4} {6 5} {6 6}]
2020/10/31 22:16:48 [DONE]: Game completed after 32 turns. Bob is the winner.
```
