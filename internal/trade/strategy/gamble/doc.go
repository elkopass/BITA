/*
Package gamble provides naive strategy for a trade-bot.

All necessary configuration can be provided by environment variables,
check out TradeConfig for the exact values to be passed.

Strategy is pretty straightforward:
	1. For each Figi provided in global config TradeBot will create
       an independent TradeWorker.
	2. Each TradeWorker is going to perform some action in a loop:
	2.1. If trading for Figi is not available, if will proceed to sleep further.
	2.2. If it has an order on market, it will check it's status:
		 if order is fulfilled, it will go to the next stage
		 (TradeWorker.sellFlag will change it's value) or sleep otherwise.
	2.3. If bot has no instrument purchased, it will try to buy one:
	2.3.1. Two stock trends will be calculated first. Long trend will be taken
		   on TradeConfig.LongTrendIntervalSeconds period and a short one on
		   TradeConfig.ShortTrendIntervalSeconds. Trend itself is a coef. K in
		   f(x) = Kx + b, where f(x) approximates historical candles
           by a linear function.
	2.3.2. If short and long trend value are grater than specified thresholds
           (TradeConfig.ShortTrendToTrade and TradeConfig.LongTrendToTrade
		   correspondingly), bot will create an order to buy an instrument
		   or proceed to sleep in other way.
	2.4. If bot has an instrument, it will try to sell it.
		 Firstly, the current order book will be requested:
	2.4.1. If (close price / TradeWorker.orderPrice) is greater than
		   TradeConfig.TakeProfitCoef, bot will create an order to take profit.
	2.4.2. If (close price / TradeWorker.orderPrice) is below
		   TradeConfig.StopLossCoef, bot will create an order to stop further loss.
	2.4.3. Or it will sleep till an asset's price stays still.
	2.4.4. If order is not fulfilled longer than TradeConfig.SecondsToCancelOrder,
		   order will be cancelled.
	2.5. If TradeWorker receives an interrupt signal, it will check a SellOnExit value
		 in global config. If it's 'true', bot will try to create a sell order based on
		 current market price. In other way it will just gracefully exit.

Gamble (naive) strategy is ready-to-use in a Sandbox environment.
TradeBot will automatically create a sandbox account and do the same things
as in a real market (except it's just a sandbox and all money here is virtual).
 */
package gamble
