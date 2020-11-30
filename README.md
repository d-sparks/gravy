# gravy

![](https://i.pinimg.com/originals/fe/68/55/fe68554a8fc5edad57e7e19a6bb51ec5.jpg)

## Overview

Gravy consists of several microservices that operate over grpc:

1. Data sources (e.g. `data/dailyprices`)
2. Supervisor (`supervisor`)
3. Algorithms (e.g. `algorithms/buyandhold`)

In theory, one can leave the data sources and supervisor running in between backtesting sessions.

The supervisor is responsible for managing the backtest and will communicate with algorithms.

#### To run gravy:

1. Follow the instructions in `data/dailyprices` to get the basic `dailyprices` db.
2. `go run cmd/data/dailyprices/main.go`
3. Run a study, such as `./studies/heads_or_tails.sh`

This should create a few files in `/tmp/foo` that are the output of the backtest.

#### To ask what are the five best days for GOOG stock.

```
echo "select date, open, close, (close-open)/open as perf from dailyprices
      where ticker = 'GOOG'
      order by perf desc limit 5
;" | psql gravy
```

#### To visualize data, use Colab:

1. `pip3 install psycopg2`
2. `./jupyter-up.sh`
3. Go to `https://colab.research.google.com` and make a new Notebook.
4. On the upper right, use the "Connect" menu to "Connect to local runtime...".
5. Run the following code.

```
import psycopg2, matplotlib.pyplot as plt, pandas as pd, pandas.io.sql as psql

conn = psycopg2.connect("host=localhost port=5432 dbname=gravy")

MSFTVSGOOGL = """
SELECT M.close AS msft, G.close AS googl, M.date AS date
FROM dailyprices M INNER JOIN dailyprices G ON M.date = G.date
WHERE M.ticker = 'MSFT' AND G.ticker = 'GOOGL' AND M.date > '2006-01-02'
ORDER BY date;"""

df = pd.read_sql(MSFTVSGOOGL, conn)
df["msft_norm"] = df["msft"] / df["msft"][0]
df["googl_norm"] = df["googl"] / df["googl"][0]

df.plot("date", ["msft_norm", "googl_norm"], figsize=(14, 7))
plt.show()
```
