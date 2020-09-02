# BattlesnakeOfficial/rules

[![codecov](https://codecov.io/gh/BattlesnakeOfficial/rules/branch/master/graph/badge.svg)](https://codecov.io/gh/BattlesnakeOfficial/rules)

[Battlesnake](https://play.battlesnake.com) rules and game logic.


## FAQ

**The Y-Axis appears to be implemented incorrectly?**

This is because the game rules implement an inverted Y-Axis. Older versions of the Battlesnake API operated this way, and several highly competitive Battlesnakes still rely on this behaviour and we'd still like to upport them. The current game engine accounts for this by translating the Y-Axis (or not) based on which version of the API each Battlesnake implements. More info [here](https://docs.battlesnake.com/guides/migrating-to-api-version-1) and [here](https://github.com/BattlesnakeOfficial/rules/issues/18).

In the future we might switch this to make the rules easier to develop? But until we drop support for the older API version it doesn't make sense to make that change.
