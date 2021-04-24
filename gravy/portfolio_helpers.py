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
                                              daily_data, targets,
                                              investment_limit=None,
                                              UP=1.01, DOWN=0.99):
    """
    Attempt to approximately invest uniformly across all assets.
    NOTE: Unlike the Go version, this function will work on a portfolio that is
          not all cash.
    """
    if len(targets) == 0:
        return []
    total_limit_of_orders = 0.0
    if investment_limit == None:
        investment_limit = portfolio.usd
    num_assets = len(set([ticker for ticker in daily_data.prices]))
    target = DOWN * portfolio_value(portfolio,
                                    daily_data.prices) / num_assets
    target_investments = dict()
    # Note: this will represent a total limit of
    #
    #   1.01 * \sum_stock floor(target / price) * price <=
    #   1.01 * \sum_stock target =
    #   1.01 * 0.99 * \sum_stocks portfolioValue / # stocks =
    #   1.01 * 0.99 * portfolioValue =
    #   0.9999 * portfolioValue
    #
    # Thus, the investment is safe as long as UP * DOWN <= 1.
    orders = []
    for ticker in daily_data.prices:
        target_investments[ticker] = 0.0
        prices = daily_data.prices[ticker]
        if ticker not in targets or prices.close < 1e-4:
            continue
        volume = math.floor(target / prices.close) - portfolio.stocks[ticker]
        if volume <= 0.0:
            continue
        limit = UP * prices.close
        if total_limit_of_orders + volume * limit >= investment_limit:
            continue
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
            if close_price + total_limit_of_orders >= investment_limit:
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
        limit = UP * next_price
        orders.append(supervisor_pb2.Order(algorithm_id=algorithm_id,
                                           ticker=next_ticker,
                                           volume=1.0,
                                           limit=limit))
        target_investments[next_ticker] += limit
        total_limit_of_orders += limit
    return orders


def invest_approximately_uniformly(algorithm_id, portfolio, daily_data,
                                   investment_limit=None, UP=1.01, DOWN=0.99):
    """
    Invest approximately uniformly in *all* assets.
    """
    targets = set([ticker for ticker in daily_data.prices])
    return invest_approximately_uniformly_in_targets(
        algorithm_id, portfolio, daily_data, targets, investment_limit)


def total_order_limit(orders):
    """
    Returns the total limit of a iterable of orders.
    """
    return sum([order.limit * order.volume for order in orders])


def remaining_limit(portfolio, orders):
    """
    Returns the amount of investment limit remaining after orders are accounted
    for.
    """
    return portfolio.usd - total_order_limit(orders)


def orders_sorted_descending(orders):
    """
    Sort orders such that all sell orders come before all buy orders, and sort
    the buys by volume * limit, descending.
    """
    sells = [order for order in orders if order.volume <= 0.0]
    buys = [order for order in orders if order.volume > 0.0]
    buys = sorted(buys, key=lambda order: -order.limit * order.volume)
    return sells + buys


def divide_or_zero(num, denom):
    """
    Divides the two quantities or returns 0 if denom is extremely small.
    """
    if denom < 1e-6:
        return 0.0
    return num / denom


def sell_overweight_target_stocks(algorithm_id, portfolio, daily_data, targets,
                                  stop_frac=0.99, invert_targets=False):
    """
    Creates sell orders for stocks that are above the target for uniformity.
    """
    if len(targets) == 0:
        return []
    num_assets = len(set([ticker for ticker in daily_data.prices]))
    target = portfolio_value(portfolio, daily_data.prices) / num_assets
    orders = []
    for ticker in portfolio.stocks:
        if (ticker not in daily_data.prices or
            (invert_targets and ticker in targets) or
                (not invert_targets and ticker not in targets)):
            continue
        price = daily_data.prices[ticker].close
        quantity = portfolio.stocks[ticker]
        weight = price * quantity
        # Only sell if over target by more than 10%.
        if (weight - target) / target >= 0.1:
            delta_units = round(target / price) - portfolio.stocks[ticker]
            orders.append(supervisor_pb2.Order(algorithm_id=algorithm_id,
                                               ticker=ticker,
                                               volume=delta_units,
                                               stop=stop_frac * price))
    return orders
