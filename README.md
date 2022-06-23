# BattlesnakeOfficial/rules

[![codecov](https://codecov.io/gh/BattlesnakeOfficial/rules/branch/master/graph/badge.svg)](https://codecov.io/gh/BattlesnakeOfficial/rules)

[Battlesnake](https://play.battlesnake.com) rules and game logic, implemented as a Go module. This code is used in production at [play.battlesnake.com](https://play.battlesnake.com). Issues and contributions welcome!


## CLI for Running Battlesnake Games Locally

This repo provides a simple CLI tool to run games locally against your dev environment.

### Installation

Download precompiled binaries here: <br>
[https://github.com/BattlesnakeOfficial/rules/releases](https://github.com/BattlesnakeOfficial/rules/releases)

Install as a Go package. Requires Go 1.18 or higher. [[Download](https://golang.org/dl/)]
```
go install github.com/BattlesnakeOfficial/rules/cli/battlesnake@latest
```

Compile from source. Also requires Go 1.18 or higher.
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

### Can I run games locally?

Yes! [See the included CLI](cli/README.md).

### How is this different from the old Battlesnake engine?

The [old game engine](https://github.com/battlesnakeio/engine) was re-written in early 2020 to handle a higher volume of concurrent games. As part of that rebuild we moved the game logic into a separate Go module that gets compiled into the production engine.

This provides two benefits: it makes it much simpler/easier to build new game modes, and it allows the community to get more involved in game development (without the maintenance overhead of the entire game engine).

### Feedback

* **Do you have an issue or suggestions for this repository?** Head over to our [Feedback Repository](https://play.battlesnake.com/feedback) today and let us know!
