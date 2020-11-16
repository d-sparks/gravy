CREATE INDEX date_ix ON dailyprices (date);
CREATE INDEX ticker_ix ON dailyprices (ticker);
CREATE INDEX date_ticker ON dailyprices (date, ticker);
CREATE UNIQUE INDEX date_exchange_ticker_pkey ON dailyprices (date, exchange, ticker);

CREATE UNIQUE INDEX trading_date_ix ON tradingdates (date);