package config

type TradeBotConfig struct {
	Env    string
	Token  string
	ApiURL string `default:"invest-public-api.tinkoff.ru:443"`
}
