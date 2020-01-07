CREATE TABLE tickers (
	ticker VARCHAR(5) PRIMARY KEY,
	exchange VARCHAR(10),
	name TEXT,
	sector TEXT,
	industry TEXT
);

CREATE TABLE dailyprices (
	ticker VARCHAR(5) NOT NULL,
	open FLOAT(8) NOT NULL,
	close FLOAT(8) NOT NULL,
	adj_close FLOAT(8)  NOT NULL,
	low FLOAT(8) NOT NULL,
	high FLOAT(8) NOT NULL,
	volume FLOAT(8) NOT NULL,
	date DATE NOT NULL
);

CREATE TABLE tradingdates (
	date DATE NOT NULL
);
