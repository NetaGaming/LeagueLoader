# LeagueLoader

**LeagueLoader** is a command line application designed to take a
list of summoner IDs from your database and then download their
KDA and CS.

## Installation From Release

* [Download the latest release](https://github.com/NetaGaming/LeagueLoader/releases/tag/v1.5)
  * You will need the *config.json.template*, *LeagueLoader*, and *tables.sql*
	  files
* Navigate to the folder where you saved these files
* Run the *tables.sql* file into MySQL

	```
	mysql -u yourusername -p yourpassword yourdatabase < tables.sql
	# or if you're already in MySQL
	source tables.sql
	```
* Update create a *config.json* file and update it to include your database
  credentials, Riot API key, and your rate limits

## Usage

After setting up your configuration and importing the database tables, all you
need to do is run **LeagueLoader**:

```sh
./LeagueLoader
```

Keep in mind that **LeagueLoader** looks for the configuration file in the
directory where you called the application. For example, a cronjob will
require `cd`'ing to the folder with the configuration file first.

```sh
*/60 * * * * cd /opt/leagueloader 1> /dev/null 2>> /tmp/leagueloader.err && ./leagueloader 1> /dev/null 2>> /tmp/ll.err
```

## Attribution

**LeagueLoader isn't endorsed by Riot Games and doesn't reflect the views or opinions of Riot Games or anyone officially involved in producing or managing League of Legends. League of Legends and Riot Games are trademarks or registered trademarks of Riot Games, Inc. League of Legends © Riot Games, Inc.**
