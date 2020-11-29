import argparse
import grpc
import logging
import math
import os
import pandas as pd
import sklearn.metrics as metrics
import sys
import tensorflow as tf

from algorithm.proto import algorithm_io_pb2
from algorithm.proto import algorithm_io_pb2_grpc
from concurrent import futures
from data.dailyprices.proto import daily_prices_pb2
from datetime import datetime, timezone
from google.protobuf import timestamp_pb2
from registrar.registrar import Registrar
from supervisor.proto import supervisor_pb2


min_results = 15

"""
Map of precision to threshold, determined from Heads_or_Tails.ipynb.
"""
thresholds = dict({
    0.5: 0.47802734375,
    0.6: 0.58642578125,
    0.7: 0.73291015625,
    0.8: 0.79052734375,
    0.9: 0.83056640625})


def z_score_or_zero(x, mu, sigma):
    """
    Returns the z score unless sigma is prohibitively small, in which case
    returns 0.0. TODO: Move to a helper package.
    """
    if sigma == None or sigma < 1E-6:
        return 0.0
    return (x - mu) / sigma


def sqrt_or_zero(x):
    """
    Returns hte square root of x if x >= 0.0. TODO: Move to a helper package.
    """
    return math.sqrt(x) if x >= 0.0 else 0.0


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
        if ticker in prices:
            value += prices[ticker].close
    return value


def orders_for_target_units(algorithm_id, ticker, target_units, limit):
    """
    Makes orders for the target number of units batched as 0.9*target,
    0.9*(target - 0.9*target), and so on, and includes any number of single
    unit orders necessary to reach target.
    """
    orders = []
    placed = 0
    next_batch = int(target_units * 0.9)
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


class HeadsOrTails(algorithm_io_pb2_grpc.AlgorithmServicer):
    """
    Heads or Tails inference algorithm.
    """

    def skip_trading(self):
        """
        Optionally skip trading to speed things up. Won't query for portfolio
        and prices and will immediately call DoneTrading.
        """
        self.num_results += 1
        return self.num_results < min_results

    def trade_all_in_on_highest_probability(
            self, portfolio, prices, tickers, predictions):
        """
        This strategy attempts to go all in on the stock with highest
        probability of going up assuming this probability is greater than 0.5.
        Otherwise, buy SPY. Assumes all tickers are in prices.
        1. Determine desired ticker.
        2. Sell anything else that is held.
        3. Determine how much of desired ticker to buy.
        """
        desired_ticker = "SPY"
        max_prediction = 0.0
        for i in range(0, len(tickers)):
            if predictions[i] > max(thresholds[0.5], max_prediction):
                max_prediction = predictions[i]
                desired_ticker = tickers[i]

        to_sell = [ticker for ticker, units in portfolio.stocks.items()
                   if ticker != desired_ticker and units > 0.0]
        orders = sell_stocks_market_order(
            self.algorithm_id, to_sell, portfolio)

        price = prices[desired_ticker].close
        held = 0.0
        if desired_ticker in portfolio.stocks:
            held = portfolio.stocks[desired_ticker]
        target = portfolio_value(portfolio, prices) - held * price
        limit = 1.01*price
        target_units = int(target / limit)
        orders += orders_for_target_units(
            self.algorithm_id, desired_ticker, target_units, price*1.01)

        return orders

    def extract_features(self, ticker, daily_data):
        """
        Returns a feature vector that matches what is produced in
        heads_or_tails.go and compatible with the heads or tails model trained
        from Heads_or_Tails.ipynb.
        """
        if ticker not in daily_data.prices:
            return None
        if "SPY" not in daily_data.prices:
            return None
        prices = daily_data.prices[ticker]
        stats = daily_data.stats[ticker]
        market_prices = daily_data.prices["SPY"]
        market_stats = daily_data.stats["SPY"]

        price = prices.close
        market_price = market_prices.close

        z_vol_15 = z_score_or_zero(
            prices.volume, stats.moving_volume[15],
            sqrt_or_zero(stats.moving_volume_variance[15]))
        z_15 = z_score_or_zero(price, stats.moving_averages[15],
                               sqrt_or_zero(stats.moving_variance[15]))
        z_35 = z_score_or_zero(price, stats.moving_averages[35],
                               sqrt_or_zero(stats.moving_variance[35]))
        z_252 = z_score_or_zero(price, stats.moving_averages[252],
                                sqrt_or_zero(stats.moving_variance[252]))
        beta = stats.beta
        sigma_market_15 = sqrt_or_zero(market_stats.moving_variance[15])
        sigma_market_35 = sqrt_or_zero(market_stats.moving_variance[35])
        sigma_market_252 = sqrt_or_zero(market_stats.moving_variance[252])
        z_market_15 = z_score_or_zero(
            market_price, market_stats.moving_averages[15], sigma_market_15)
        z_market_35 = z_score_or_zero(
            market_price, market_stats.moving_averages[35], sigma_market_35)
        z_market_252 = z_score_or_zero(
            market_price, market_stats.moving_averages[252], sigma_market_252)

        return [0.0, z_vol_15, z_15, z_35, z_252, beta, sigma_market_15,
                sigma_market_35, sigma_market_252, z_market_15, z_market_35,
                z_market_252]

    def trade(self, portfolio, daily_data):
        """
        Tells the algorithm to trade. Mostly unimplemented.
        """
        tickers, prices = zip(
            *[(t, p) for t, p in daily_data.prices.items()])
        features = [self.extract_features(ticker, daily_data)
                    for ticker in tickers]
        predictions = self.model.predict(features)

        if self.strategy == 'all_in_on_highest_probability':
            return self.trade_all_in_on_highest_probability(
                portfolio, daily_data.prices, tickers, predictions)

        return []

    def test(self, test_data_path):
        """
        Evaluate the model on the test data and print a small summary.
        """
        hot_features = pd.read_csv(test_data_path).fillna(0.0)
        hot_labels = hot_features.pop('result').astype(int)
        hot_preds = self.model.predict(hot_features)

        fpr, tpr, _ = metrics.roc_curve(hot_labels, hot_preds)
        logging.info('Model loaded with %f AUC.' % metrics.auc(fpr, tpr))

    def __init__(self, id, model_dir, strategy, test_data_path):
        """
        Constructor for heads or tails model.
        """
        logging.basicConfig(
            format='%(asctime)s %(levelname)-8s %(message)s',
            level=logging.INFO,
            datefmt='%Y-%m-%d %H:%M:%S')

        self.algorithm_id = supervisor_pb2.AlgorithmId(algorithm_id=id)
        self.id = id
        self.registrar = Registrar()
        self.num_results = 0
        self.model = tf.keras.models.load_model(model_dir)
        self.strategy = strategy

        self.test(test_data_path)

    def Execute(self, input, context):
        """
        Implements the algorithm_io interface and tells the algorithm to
        execute any trades before calling DoneTrading.
        """
        timestamp = datetime.fromtimestamp(
            input.timestamp.seconds, timezone.utc).strftime('%Y-%m-%d')
        logging.info('Executing algorithm on %s' % timestamp)

        if not self.skip_trading():
            # Get prices, portfolio, and run the algorithm.
            portfolio = self.registrar.supervisor_stub.GetPortfolio(
                self.algorithm_id)
            daily_data = self.registrar.dailyprices_stub.Get(
                daily_prices_pb2.Request(timestamp=input.timestamp, version=0))
            orders = self.trade(portfolio, daily_data)

            # Submit the orders.
            for order in orders:
                self.registrar.supervisor_stub.PlaceOrder(order)

        # Indicate done trading.
        self.registrar.supervisor_stub.DoneTrading(self.algorithm_id)
        return algorithm_io_pb2.Output()

    algorithm_id = None
    id = None
    registrar = None
    num_results = 0
    model = None
    strategy = None


if __name__ == '__main__':
    # Parse flags
    parser = argparse.ArgumentParser(description='Heads or Tails algorithm.')
    parser.add_argument('--id', type=str, help='Algorithm name',
                        default='headsortails', required=False)
    parser.add_argument('--port', type=str, help='Serving port',
                        default='17506', required=False)
    parser.add_argument('--model_dir', type=str, help='Model directory',
                        default='algorithm/headsortails/train/model',
                        required=False)
    parser.add_argument('--strategy', type=str, help='Model strategy',
                        default='all_in_on_highest_probability',
                        required=False)
    parser.add_argument(
        '--test_data', type=str, help='Evaluate model on data.',
        default='algorithm/headsortails/train/data/2005_to_2015_data.csv',
        required=False)
    args = parser.parse_args()

    # Verify strategy
    strategy = args.strategy
    valid_strategies = set(['all_in_on_highest_probability'])
    if strategy not in valid_strategies:
        logging.error('Unknown strategy `%s`' % strategy)

    # Serve
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    algorithm_io_pb2_grpc.add_AlgorithmServicer_to_server(
        HeadsOrTails(args.id, args.model_dir, args.strategy, args.test_data),
        server)
    server.add_insecure_port('[::]:%s' % args.port)
    server.start()
    logging.info('Listening on `localhost:%s`' % args.port)
    server.wait_for_termination()
