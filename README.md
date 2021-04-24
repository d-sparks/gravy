# gravy

![](https://i.pinimg.com/originals/fe/68/55/fe68554a8fc5edad57e7e19a6bb51ec5.jpg)

## Overview

Gravy consists of several microservices that operate over grpc:

1. Data sources (e.g. `data/dailyprices`)
2. Supervisor (`supervisor`)
3. Algorithms (e.g. `algorithms/buyandhold`)

With all of these running, one can begin a backtest via `cmd/begin_backtest/main.go`. Typically one will persist the data sources (which generally cache their outputs in memory) and make a "study" (`studies/`), which is a shell script which configures the supervisor and algorithms to run against the persisting data sources.

Backtests currently output various debug logs (usually to a temp directory) and also a TimescaleDB output. Visualize the output with Grafana.

## Technologies / dependencies

1. gRPC
2. proto3
3. Golang
4. Python 3
5. PostGRES
6. TimescaleDB
7. Grafana
8. Jupyter / colab
9. Typical Python libraries (pandas, scipy, numpy, Tensorflow, keras, sklearn, matplotlib, etc)

## Examples

### To run gravy:

1. Follow the instructions in `data/dailyprices` and `data/assetids` to get the basic `dailyprices`, `gravy_timescale_output`, and `assetids` dbs
2. Run a persisting data source with `go run cmd/data/dailyprices/main.go`
3. Run a study, such as `sh ./studies/correlated_pairs.sh`
4. This should start populating a new table `timescaleout${TIMESTAMP}` in the `gravy_timescale_output` db
5. Run grafana, e.g. `brew services start grafana`
6. Go to `localhost:3000` and bring the `${TIMESTAMP}` to visualize the results
7. (there are also some files output in `/tmp/fizzybuzzy/` including a log of individual buy/sell orders)

<img width="1282" alt="grafana_example" src="https://user-images.githubusercontent.com/7853117/115964692-ebc43280-a4e2-11eb-84c1-8deaaecb65d3.png">

### To ask what are the five best days for GOOG stock.

```
echo "select date, open, close, (close-open)/open as perf from dailyprices
      where ticker = 'GOOG'
      order by perf desc limit 5
;" | psql gravy
```

### To visualize prices data, use Colab:

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
