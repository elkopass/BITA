# Базовая конфигурация

Торговый бот полностью конфигурируется через переменные окружения
с использованием библиотеки [envconfig](https://github.com/kelseyhightower/envconfig).

В файлах [.env-example](https://github.com/elkopass/BITA/blob/main/cmd/trade-bot/.env-example)
приведены шаблоны конфигов. Перед запуском trade-bot или trade-utils их 
необходимо модифицировать, передав туда все необходимые параметры, такие как 
токен API, ID аккаунта и т.д. 

Модифицированные значения требуется либо экспортировать перед запуском бота, либо
передать в контейнер в виде .env-файла (см. раздел "Инсталляция").

## Глобальные параметры

```bash
## (обязательный) массив FIGI, разделённых запятой
TRADEBOT_FIGI=<figi1>,<figi2>
## (обязательный) должна ли использоваться песочница
TRADEBOT_IS_SANDBOX=true
## (обязательный) токен API (https://www.tinkoff.ru/invest/settings/)
TRADEBOT_TOKEN=<your_api_token>
## (обязательный) ID аккаунта для торговли (не нужен для песочницы)
TRADEBOT_ACCOUNT_ID=<your_account_id>

## окружение для запуска (попадает в логи): DEV, TEST, PROD
# TRADEBOT_ENV=UNSPECIFIED
## торговая стратегия, доступны для выбора: gamble, crumble, tumble
# TRADEBOT_STRATEGY=gamble
## при значение true воркер будет продавать купленный инструмент 
## по рыночной цене при прерывании (сигнал SIGINT)
# TRADEBOT_SELL_ON_EXIT=false
## уровень логирования, один из: DEBUG, INFO, WARN, ERROR
# TRADEBOT_LOG_LEVEL=INFO
```

## Prometheus-экспортер

```bash
## должен ли бот экспортировать метрики
# METRICS_ENABLED=true
## адрес и порт для экспорта метрики
# METRICS_ADDR=:8080
## эндпоинт для получения метрик
# METRICS_ENDPOINT=/metrics
```

## Конфигурация стратегий

Переменные окружения для различных стратегий начинаются 
с их уникального префикса (например, `GAMBLE_STRATEGY_`).

Детали конфигураций каждой отдельной стратегии приведены 
в разделе "Торговые стратегии".
