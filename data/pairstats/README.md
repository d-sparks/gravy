# Pair stats

One thing we'd like to track is the 252 day rolling correlation of every pair of stocks. Because for 10k+ stocks this can amount to significant data, we break down as follows. To store this data historically even at daily frequency is `10000 * 10000 * 252 * 15 = 378 billion` floats.

## Historical data

We can chunk the cartesian product Stocks x Stocks by segmenting Stocks = Stocks_1 + ... + Stocks_n and processing Stocks_i x Stocks_j for (i, j) in n x n. As mentioned, this may be a lot to store, so we can:

```
# Phase 1: Calculate number of data points.
mu = streaming mean
for each chunk Stocks_i x Stocks_j:
    track 252 day correlation of each (s, t) in Stocks_i x Stocks_j
    for each trading day:
        for each pair of stocks:
            if both listed:
                observe correlation and variance of correlation
                if correlation and variance have 252 observations:
                    increase count by 1
log count
```