# Email Filter Client

The **Email Filter Client** is a tool designed to connect to multiple IMAP email accounts, process incoming emails, and filter messages based on a blacklist.  
It automates the task of identifying and moving unwanted emails to the trash, ensuring your inbox stays clean and organized.

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