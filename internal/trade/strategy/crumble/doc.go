/*
Package crumble provides strategy based on moving average for a trade-bot.
See https://www.investopedia.com/terms/m/movingaverage.asp

All necessary configuration can be provided by environment variables,
check out TradeConfig for the exact values to be passed.

Strategy is pretty straightforward:
	1. For each Figi provided in global config TradeBot will create
       an independent TradeWorker.
	2. Each TradeWorker is going to perform some action in a loop:
	2.1. If trading for Figi is not available, if will proceed to sleep further.
	2.2. If it has an order on market, it will check it's status:
		 if order is fulfilled, it will go to the next stage
		 or sleep otherwise.
	2.3. The trading algorithm builds two MA, on a large interval
		 (TradeConfig.LongWindow) and a small one (TradeConfig.ShortWindow).
		 At the moment when the long exceeds the short, the robot
		 sells, in the opposite case, it buys. Also interval is configured
		 by TradeConfig.CandlesIntervalHours to query historic candles
		 from sdk.MarketDataService.
	2.4. If TradeWorker receives an interrupt signal, it will check a SellOnExit value
		 in global config. If it's 'true', bot will try to create a sell order based on
		 current market price. In other way it will just gracefully exit.

Сrumble (MA-based) strategy is ready-to-use in a Sandbox environment.
TradeBot will automatically create a sandbox account and do the same things
as in a real market (except it's just a sandbox and all money here is virtual).
*/
package crumble
