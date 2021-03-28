import math

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


def invest_approximately_uniformly_in_targets(algorithm_id, portfolio,
                                              daily_data, targets):
    """
    Attempt to approximately invest uniformly across all assets.
    TODO: Test this function! It is a port of the tested Go version, but there
          could be errors.
    """
    if len(targets) == 0:
        return []
    total_limit_of_orders = 0.0
    target = 0.99 * portfolio_value(portfolio,
                                    daily_data.prices) / len(targets)
    target_investments = dict()
    # Note: this will represent a total limit of
    #
    #   1.01 * \sum_stock floor(target / price) * price <=
    #   1.01 * \sum_stock target =
    #   1.01 * 0.99 * \sum_stocks portfolioValue / # stocks =
    #   1.01 * 0.99 * portfolioValue =
    #   0.9999 * portfolioValue
    #
    # Thus, the investment is safe.
    orders = []
    for ticker in daily_data.prices:
        prices = daily_data.prices[ticker]
        if ticker not in targets:
            continue
        volume = math.floor(target / prices.close)
        limit = 1.01 * prices.close
        orders.append(supervisor_pb2.Order(algorithm_id=algorithm_id,
                                           ticker=ticker, volume=volume,
                                           limit=limit))
        target_investments[ticker] = volume * limit
        total_limit_of_orders += volume * limit
    # Continue investing until we can't get closer to a uniform investment.
    while portfolio.usd - total_limit_of_orders > 0.0:
        next_ticker = None
        next_improvement = 0.0
        next_price = 0.0
        for ticker in daily_data.prices:
            prices = daily_data.prices[ticker]
            if ticker not in targets:
                continue
            current_target = target_investments[ticker]
            close_price = prices.close
            if close_price + total_limit_of_orders > portfolio.usd:
                continue
            hypothetical_delta = abs(close_price + current_target - target)
            current_delta = abs(current_target - target)
            improvement = current_delta - hypothetical_delta
            if improvement > next_improvement:
                next_improvement = improvement
                next_ticker = ticker
                next_price = close_price
        if next_improvement == 0.0:
            # No improvement can be made.
            break
        # Place an order for a next_ticker.
        limit = 1.01 * next_price
        orders.append(supervisor_pb2.Order(algorithm_id=algorithm_id,
                                           ticker=next_ticker,
                                           volume=1.0,
                                           limit=limit))
        target_investments[next_ticker] += limit
        total_limit_of_orders += limit
    return orders


def invest_approximately_uniformly(algorithm_id, portfolio, daily_data):
    """
    Invest approximately uniformly in *all* assets.
    """
    targets = set([ticker for ticker in daily_data.prices])
    return invest_approximately_uniformly_in_targets(
        algorithm_id, portfolio, daily_data, targets)
