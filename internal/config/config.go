package config

type TradeBotConfig struct {
	Token  string
	Env    string `default:"UNKNOWN"`
	ApiURL string `default:"invest-public-api.tinkoff.ru:443"`
}
