# Daily prices data

#### To create a database from scratch:

1. Install postgres and have it running on localhost. (This should happen by default if using docker compose)
1. Download and unzip [raw data](https://drive.google.com/file/d/1NzlJ81GiWpdtiWrnbMYRL4Z5I3dZR_C6/view?usp=drivesdk) and move the `data` directory to `gravy/data/dailyprices/raw`
1. `cat data/dailyprcies/sql/create_tables.sql | psql gravy`
1. `cat data/dailyprcies/sql/create_indexes.sql | psql gravy`
1. `go run cmd/data/dailyprices/pipeline/main.go` (takes > 5 hours)

Also include the S&P500 by downloading the historical prices from [Yahoo! Finance](https://finance.yahoo.com/quote/%5EGSPC/history/).

1. Put `^GSPC.csv` into `data/dailyprices/GSPC/raw`.
2. Run `go run cmd/data/dailyprices/GSPC/pipeline/main.go`

#### To restore the database from a dump:

1. Install postgres and have it running on localhost.
2. Download the `pg_dump_output` and put it in `data/dailyprices/raw`.
2. `psql gravy < data/dailyprices/raw/pg_dump_output`. (Haven't tried this yet.)
