# BITA

[![CI-CD](https://github.com/elkopass/BITA/actions/workflows/main.yml/badge.svg)](https://github.com/elkopass/BITA/actions/workflows/main.yml)

Trading **B**ot based on [**T**inkoff **I**nvest **A**PI](https://github.com/Tinkoff/investAPI)

![TradeBot logo](logo.png)

## Installation

Clone this repository first, and we are ready to go.

```shell script
$ git clone https://github.com/elkopass/BITA
```

This bot is fully configurable with environment variables 
using [envconfig](https://github.com/kelseyhightower/envconfig) library.

You can find full list of them in our 
[.env-example](https://github.com/elkopass/BITA/blob/main/cmd/trade-bot/.env-example) file.

Right before we started, obtain your API token in 
[settings](https://www.tinkoff.ru/invest/settings/) and set it as follows: `TRADEBOT_TOKEN=<your_api_token>`.

### Using Docker

Running trade-bot in a Docker-container is a preferable way.

```shell script
$ cp cmd/trade-bot/.env-example cmd/trade-bot/.env
$ vim cmd/trade-bot/.env # make sure to reconfigure it with your own data!
$ docker-compose up --build
```

### Manual

Let's set and export env variables first:

```shell script
$ cp cmd/trade-bot/.env-example cmd/trade-bot/.env
$ vim cmd/trade-bot/.env # make sure to reconfigure it with your own data!
$ export $(grep -v '^#' cmd/trade-bot/.env | xargs)
```

After that, you can build binary and run it using go v1.16+ as follows:
```shell script
$ go build -a -o trade-bot ./cmd/trade-bot
$ ./trade-bot
``` 

## License

This project is released under the MIT license. 
See [LICENSE](https://github.com/elkopass/BITA/blob/main/LICENSE) for details.
