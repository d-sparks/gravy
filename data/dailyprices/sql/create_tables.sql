CREATE TABLE dailyprices (
  ticker VARCHAR(10) NOT NULL,
  exchange VARCHAR(10) NOT NULL,
  open FLOAT (8) NOT NULL,
  close FLOAT (8) NOT NULL,
	low FLOAT (8) NOT NULL,
	high FLOAT (8) NOT NULL,
	volume FLOAT (8) NOT NULL,
	date DATE NOT NULL
);

CREATE TABLE tradingdates (
  date DATE NOT NULL,
  nyse BOOL,
  nyse_etf BOOL,
  nasdaq BOOL,
  nasdaq_etf BOOL,
  nysemkt BOOL,
  nysemkt_etf BOOL
);