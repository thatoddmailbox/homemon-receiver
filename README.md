# homemon-receiver
This program receives packets sent over UDP by [homemon-daemon](https://github.com/thatoddmailbox/homemon-daemon).

Note that this is designed to work together with [homemon-server](https://github.com/thatoddmailbox/homemon-server). To that end, it inserts data into homemon-server's database, and so you will need to run the migrations in that repository, even if you're only using this code.

## Configuration
Create a config.toml file, in the same directory you run the receiver from. This file should look something like:

```toml
Port = 4321
Token = "token"
DatabaseDSN = "username:password@tcp(localhost:3306)/database"
```

### Options
* `Port` - the port to listen for UDP packets on. Can be anything, as long as it matches homemon-daemon.
* `Token` - the token to use for authentication of messages. Must match homemon-daemon's `Token`. You can generate a random token in the correct format by running this program with the `-generate-token` flag.
* `DatabaseDSN` - The database DSN to use to connect to MySQL. You should fill this in with your database credentials. This should be the same database that homemon-server works with.