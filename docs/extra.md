# Дополнительно

## Консольная утилита trade-utils

Небольшая утилита trade-utils позволяет:

- получить список доступных для управления ботом
аккаунтов и их портфель (модуль `-mode accounts`)

- получить список доступных для торговли прямо сейчас 
инструментов (модуль `-mode figi`)

- подвести отчёт по совершённым операциям за последние 
сутки (модуль `-mode operations`).

### Сборка и запуск

```bash
$ go build -v -o trade-utils ./cmd/trade-utils/

$ ./trade-utils 
Usage: trade-utils -mode [accounts|figi|operations]
  -mode string
        running module
```
