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