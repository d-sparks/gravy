# Daily prices data

#### To create a database from scratch:

1. Download `historical_stock_prices.csv` and `historical_stocks.csv` and put them in `data/dailyprices/raw`
2. Install postgres and have it running on localhost.
3. `createdb gravy`
4. `cat data/dailyprcies/sql/create_tables.sql | psql gravy`
5. `go run cmd/data/dailyprices/pipeline/main.go` (takes > 5 hours)
6. `cat kaggle/create_indexes.sql | psql gravy`

#### To restore the database from a dump:

1. Install postgres and have it running on localhost.
2. Download the `pg_dump_output` and put it in `data/dailyprices/raw`.
2. `psql gravy < data/dailyprices/raw/pg_dump_output`. (Haven't tried this yet.)

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
