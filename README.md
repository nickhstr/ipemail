# ipemail
A simple program to notify when a user's public IP address has updated, using gmail to send messages.

## Installation
```sh
$ git clone https://github.com/nickhstr/ipemail.git
$ cd ipemail
$ go install
```

## Usage
Some environment variables need to be set, either as part of the execution of `ipemail`, or in a `.env` file.

### .env file
By default, `ipemail` will look for a `.env` file in the same directory as where the command is called, or a path to a `.env` file can be specified with the environment variable `ENV_FILE`.

Example:
```sh
$ ENV_FILE=/home/me/ipemail/ip_file.txt ipemail
```

And here's an example `.env` file:
```
EMAIL_FROM_ADDRESS=from_someone@gmail.com
EMAIL_FROM_USER=John Doe
EMAIL_FROM_PASSWORD=Password
EMAIL_TO_ADDRESS=to_someone@gmail.com
LAST_IP_DIR=/home/me/go/src/github.com/nickhstr/ipemail
LAST_IP_FILE=ip_file.txt
```
