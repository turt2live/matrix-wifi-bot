# matrix-wifi-bot

[![#wifibot:t2bot.io](https://img.shields.io/badge/matrix-%23wifibot:t2bot.io-brightgreen.svg)](https://matrix.to/#/#wifibot:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/matrix-wifi-bot.svg?branch=master)](https://travis-ci.org/turt2live/matrix-wifi-bot)

A bot that tracks what wifi networks it sees. Works well on a Raspberry Pi.

# Installing

Assuming Go 1.9 and `dep` are already installed:
```bash
# Get it
git clone https://github.com/turt2live/matrix-wifi-bot
cd matrix-wifi-bot

# Grab the dependencies
dep ensure

# Build it
go install

# Configure it (edit wifi-bot.yaml to meet your needs)
cp config.sample.yaml wifi-bot.yaml

# Run it
$GOPATH/bin/matrix-wifi-bot
```
