# BattlesnakeOfficial/rules

[![codecov](https://codecov.io/gh/BattlesnakeOfficial/rules/branch/master/graph/badge.svg)](https://codecov.io/gh/BattlesnakeOfficial/rules)

[Battlesnake](https://play.battlesnake.com) rules and game logic, implemented as a Go module. This code is used in production at [play.battlesnake.com](https://play.battlesnake.com). Issues and contributions welcome!


## CLI for Running Battlesnake Games Locally

This repo provides a simple CLI tool to run games locally against your dev environment.

### Installation

Download precompiled binaries here: <br>
[https://github.com/BattlesnakeOfficial/rules/releases](https://github.com/BattlesnakeOfficial/rules/releases)

Install as a Go package. Requires Go 1.15 or higher. [[Download](https://golang.org/dl/)]
```
go get github.com/BattlesnakeOfficial/rules/cli/battlesnake
```

Compile from source. Also requires Go 1.15 or higher.
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

For more details, see the [CLI README](cli/README.md).


## FAQ

### How is this different from the old Battlesnake engine?

The [old game engine](https://github.com/battlesnakeio/engine) was re-written in early 2020 to handle a higher volume of concurrent games. As part of that rebuild we moved the game logic into a separate Go module that gets compiled into the production engine.

This provides two benefits: it makes it much simpler/easier to build new game modes, and it allows the community to get more involved in game development (without the maintenance overhead of the entire game engine).


### Can I run games locally?

Yes! [See the included CLI](cli/README.md).


### The Y-Axis appears to be implemented incorrectly!?!?

This is because the game rules implement an inverted Y-Axis. Older versions of the Battlesnake API operated this way, and several highly competitive Battlesnakes still rely on this behaviour. The current game engine accounts for this by translating the Y-Axis (or not) based on which version of the API each Battlesnake implements. [More info here](https://docs.battlesnake.com/guides/migrating-to-api-version-1) and [here](https://github.com/BattlesnakeOfficial/rules/issues/18).

In the future we might switch this to make the rules easier to develop? But until we drop support for the older API version it doesn't make sense to make that change.
