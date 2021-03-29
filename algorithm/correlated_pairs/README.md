# Correlated Pairs

## TODO

1. The introduction of correlated pairs caused the dailyprices service to use a lot of memory. Need to figure that out.
2. For every pair established, both on the opening and closing, we can leak small amounts of capital. We should reinvest that capital into stocks that aren't in pairs and are underweighted in the portfolio.
3. We should also consider rebalancing periodically.
4. There are division by 0 errors when some closing prices are 0.
5. We should build some tools for understanding the patterns. E.g. which pairs end up being good and bad in practice?
6. Better ML model, more examples, more features.