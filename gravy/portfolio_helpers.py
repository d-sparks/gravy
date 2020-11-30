from supervisor.proto import supervisor_pb2


def sell_stocks_with_stop(algorithm_id, tickers, portfolio, stop):
    """
    Creates orders to sell all tickers at the stop price determined by the
    given `stop` lambda. TODO: Move to a helper package.
    """
    orders = []
    for ticker in tickers:
        if ticker not in portfolio.stocks:
            continue
        units = portfolio.stocks[ticker]
        orders += [supervisor_pb2.Order(
            algorithm_id=algorithm_id, ticker=ticker,
            volume=-units, stop=stop(ticker))]
    return orders


def sell_stocks_market_order(algorithm_id, tickers, portfolio):
    """
    Sells all held tickers with market orders. TODO: Move to a helper package.
    """
    return sell_stocks_with_stop(
        algorithm_id, tickers, portfolio, lambda ticker: 0.0)


def portfolio_value(portfolio, prices):
    """
    Calculates the portfolio value at closing prices. TODO: Move to a helper
    package.
    """
    value = portfolio.usd
    for ticker, units in portfolio.stocks.items():
        value += units * prices[ticker].close
    return value


def orders_for_target_units(algorithm_id, ticker, target_units, limit):
    """
    Makes orders for the target number of units batched as 0.9*target,
    0.9*(target - 0.9*target), and so on, and includes any number of single
    unit orders necessary to reach target.
    """
    orders = []
    placed = 0
    next_batch = max(1, int(target_units * 0.9))
    while True:
        if placed + next_batch > target_units:
            if next_batch == 1:
                break
            next_batch = max(1, int((target_units-placed) * 0.9))
            continue
        orders += [supervisor_pb2.Order(
            algorithm_id=algorithm_id, ticker=ticker,
            volume=next_batch, limit=limit)]
        placed += next_batch
    return orders
