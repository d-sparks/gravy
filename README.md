# gravy

![](http://www.artyfactory.com/art_appreciation/animals_in_art/pablo_picasso/picasso_bull_plate_5.jpg)

## Getting started

#### To create a database from scratch:

1. Download the kaggle data and put it in `kaggle/data/`.
2. Install postgres and have it running on localhost.
3. `createdb gravy`
4. `cat kaggle/create_tables.sql | psql gravy`
5. `go run cmd/gravy/*.go` (takes > 5 hours)
6. `cat kaggle/create_indexes.sql | psql gravy`

#### To restore the database from a dump:

1. Download the kaggle database dump and put it in `kaggle/data`.
2. Install postgres and have it running on localhost.
3. `psql gravy < kaggle/data/pg_dump_output`. (Haven't tried this yet.)

#### To ask what are the five best days for GOOGL stock.

```
echo "select date, open, close, (close-open)/open as perf from dailyprices
      where ticker = 'GOOGL'
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
