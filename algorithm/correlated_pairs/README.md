# Correlated Pairs

## TODO

1. We should also consider rebalancing periodically.
2. We should build some tools for understanding the patterns. E.g. which pairs end up being good and bad in practice?
3. Better ML model, more examples, more features.
4. Rolling correlation may improve quality
5. Understand how to get a larger number of pairs (more "in position") and the impact that would have on quality.
6. (EASY) Don't enter a position for stocks that don't already have a significant holding.

## Services

1. Grafana runs on 3030 (`brew services grafana start`)
2. The db has a simple debug on 8080
3. Most things run on 175xx
4. Database `gravy` is a pure postgres with the `dailyprices` table etc
5. Database `gravy_timescale_output` is a timescaleDB output with tables for each run