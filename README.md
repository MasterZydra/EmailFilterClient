# Email Filter Client

The **Email Filter Client** is a tool designed to connect to multiple IMAP email accounts, process incoming emails, and filter messages based on a blacklist.  
It automates the task of identifying and moving unwanted emails to the trash, ensuring your inbox stays clean and organized.

## Table of Contents
- [Features](#features)
- [Configuration](#configuration)
  - [`config.json`](#configjson)
  - [`blacklist.json`](#blacklistjson)
- [Web frontend](#web-frontend)
  - [Security](#security)
  - [Routes](#routes)
- [Run application](#run-application)

## Features
- Connects to multiple IMAP email accounts.
- Moves emails to the trash based on a customizable blacklist.
- Maintains state to avoid reprocessing already handled emails.
- Configurable interval for periodic email account processing.

## Configuration
The tool requires two main configuration files: `config.json` and `blacklist.json`.

### `config.json`
This file contains the configuration for IMAP connections and the processing interval. Below is an example:

```json
{
  "interval": 5,
  "imapConnections": [
    {
      "host": "imap.example.com:993",
      "email": "user1@example.com",
      "password": "password123"
    },
    {
      "host": "imap.example.org:993",
      "email": "user2@example.com",
      "password": "password456"
    }
  ]
}
```

- `interval`: The time interval (in minutes) between processing cycles.
- `imapConnections`: A list of IMAP accounts to connect to, each with:
    - `host`: The IMAP server address and port.
    - `email`: The email address.
    - `password`: The password for the account.

### `blacklist.json`
This file contains a list of email addresses or domains to filter. Below is an example:

```json
{
  "from": [
    "spam@example.com",
    "ads@marketing.com",
    "@spam.domain",
    "unwanted.org"
  ]
}
```

- `from`: A list of email addresses or domains that should be filtered.

## Web frontend
The program also starts a webserver.

Run the program with the default port 8080:
```bash
$ go run webserver.go
```

Run the program with a custom port:
```bash
$ go run webserver.go -port=9090
```

### Security
To protect the web frontend a password for basic auth (username: mailadmin) can be passed:
```bash
$ go run main.go -basicAuthPassword=mySecretKey
```

### Routes
`/` - Shows a list of all available routes.  
`/log` - Returns the content of the log file.  
`/log/clear` - Clears the content of the log file.  
`/config` - Returns the content of the `config.json` file.  
`/config/update` - Updates the `config.json` file with new data.  
`/blacklist` - Returns the content of the `blacklist.json` file.  
`/blacklist/update` - Updates the `blacklist.json` file with new data.  

## Run application
### Docker
#### Build it yourself
```bash
$ docker build -t email-filter-client .
$ docker run --rm -d --env PORT=8080 --env BASIC_AUTH_PASSWORD=mySecurePassword -p 8080:8080 -v $(pwd)/config:/app/config -v $(pwd)/log:/app/log --name email-filter-client email-filter-client
```

#### Use image from container registry
```bash
$ docker run --rm -d --env PORT=8080 --env BASIC_AUTH_PASSWORD=mySecurePassword -p 8080:8080 -v $(pwd)/config:/app/config -v $(pwd)/log:/app/log --name email-filter-client ghcr.io/masterzydra/email-filter-client:latest
```

### Without docker
```bash
$ go build -o . ./...
$ ./EmailFilterClient -port=8081 -basicAuthPassword=ASecurePassword
```

