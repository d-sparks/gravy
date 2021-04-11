# Correlated Pairs

## TODO

1. We should also consider rebalancing periodically.
2. There are division by 0 errors when some closing prices are 0.
3. We should build some tools for understanding the patterns. E.g. which pairs end up being good and bad in practice?
4. Better ML model, more examples, more features.

## How to measure strategies

Need somewhere to consolidate thoughts/documentation/decisions. (Here for now.)

A notion of "in position" holdings and "out of position" holdings

* `portfolio_value`
* `usd` (uninvested capital)
* `position_value`
* `oop_value`

We have `portfolio_value = usd + in_position_value + oop_value`.

We measure the overall distribution of out of position holdings.
Some strategies will have no out of position holdings, which is fine.

* `ooo_deviation_min`
* `ooo_deviation_max`
* `ooo_deviation_10p`
* `ooo_deviation_25p`
* `ooo_deviation_50p`
* `ooo_deviation_75p`
* `ooo_deviation_90p`

Of course basic metrics of the return value

* `alpha_15`
* `alpha_35`
* `alpha_252`
* `beta_15`
* `beta_35`
* `beta_252`

Total value of buys and sells

* `buys_value`
* `sells_value`

Some strategies will reason in terms of discrete "positions" which are opened and then closed.
We track all closing positions in a separate database for drilling down to specifics.
In the main time series, we track some aggregated metrics of all positions closing on that day.

* `num_opening_positions`
* `num_closing_positions`
* `closing_pos_return_min`
* `closing_pos_return_max`
* `closing_pos_return_mean`

## What is a position?

There is some question remaining about how to measure the performance of a position.

