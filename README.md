# League Loader

Pulls summoner and game data into a database for tracking and later consumption.

## What works

* downloads and updates summoner information (level and summoner name)
* downloads most recent games for each summoner
* merges only games that haven't been merged yet
* downloads statistics for each player in the game that is also in your database

## TODO

* General

    * download champion data (name and images)
    * download item data (name and images)
    * download spell data (name and images)
    * accept command line arguments (so that static data can be downloaded less often)

* Neta

    * track teams
    * separate games that were specific to League play

## Disclaimer

**This product is not endorsed, certified or otherwise approved in any way by Riot Games, Inc. or any of its affiliates.**
