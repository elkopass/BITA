/*
Package tumble provides a trading strategy based on order book.

All necessary configuration can be provided by environment variables,
check out TradeConfig for the exact values to be passed.

The trade-bot tracks the "order book". If there are more lots in the
purchase orders than in the lots for sale a certain number of times,
then the robot buys the instrument at the market price, otherwise it sells,
immediately placing the order in the opposite direction, but with a
certain percentage of profit.

Tumble (OrderBook-based) strategy is ready-to-use in a Sandbox environment.
TradeBot will automatically create a sandbox account and do the same things as
in a real market (except it's just a sandbox and all money here is virtual).
*/
package tumble
