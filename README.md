# BattlesnakeOfficial/rules

[![codecov](https://codecov.io/gh/BattlesnakeOfficial/rules/branch/master/graph/badge.svg)](https://codecov.io/gh/BattlesnakeOfficial/rules)

[Battlesnake](https://play.battlesnake.com) rules and game logic, implemented as a Go module. Issues and contributions welcome!


## FAQ

### How is this different from the old Battlesnake engine?

The [old game engine](https://github.com/battlesnakeio/engine) was re-written in early 2020 to handle a higher volume of concurrent games. As a result, the majority of engine code was concerned with scalable game execution and not the actual game rules. As part of that refactor we moved the game logic implementation into a separate go module that gets compiled into the production engine. 

This provides two benefits: it makes it much simpler/easier to build new game modes, and it allows the community to get more involved in game development (without the maintenance overhead of the entire game engine).


### Can I run games locally?

This was functionality that an older version of the game logic provided, including a release binary. We'd like to recreate that behaviour and include it in this repo, but for now it doesn't exist. See [Issue #20](https://github.com/BattlesnakeOfficial/rules/issues/20) for more info.


### The Y-Axis appears to be implemented incorrectly?

This is because the game rules implement an inverted Y-Axis. Older versions of the Battlesnake API operated this way, and several highly competitive Battlesnakes still rely on this behaviour and we'd still like to upport them. The current game engine accounts for this by translating the Y-Axis (or not) based on which version of the API each Battlesnake implements. More info [here](https://docs.battlesnake.com/guides/migrating-to-api-version-1) and [here](https://github.com/BattlesnakeOfficial/rules/issues/18).

In the future we might switch this to make the rules easier to develop? But until we drop support for the older API version it doesn't make sense to make that change.

## CLI Build Instructions

From the root folder of the project, run:

```
go build -o battlesnake cli/main.go
```

This will create the `rules` command, that you can use. It currently has one verb: `play`.

```
Use the CLI to configure and play a game of Battlesnake against 
multiple snakes, with multiple rulesets.

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

For more details, see the [CLI-Specific README](cli/README.md)
