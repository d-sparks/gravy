import argparse
import collections
import csv
import gravy.portfolio_helpers as portfolio_helpers
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
from gravy.math_helpers import *
from registrar.registrar import Registrar
from supervisor.proto import supervisor_pb2


"""
Don't trade until we've seen `min_results` examples, so we can build confidence
in the correlations.
"""
min_results = 100
"""
Hold pair trades for `period` days.
"""
period = 7


class Individual:
    """
    Represents the fundamentals of an asset in a given state.
    """

    def __init__(self, ticker, daily_data):
        prices = daily_data.prices[ticker]
        price = prices.close
        stats = daily_data.stats[ticker]
        market_prices = daily_data.prices["SPY"]
        market_price = market_prices.close
        market_stats = daily_data.stats["SPY"]
        self.price = price
        self.market_price = market_price
        self.z_vol_15 = z_score_or_zero(
            prices.volume, stats.moving_volume[15],
            sqrt_or_zero(stats.moving_volume_variance[15]))
        self.z_15 = z_score_or_zero(price, stats.moving_averages[15],
                                    sqrt_or_zero(stats.moving_variance[15]))
        self.z_35 = z_score_or_zero(price, stats.moving_averages[35],
                                    sqrt_or_zero(stats.moving_variance[35]))
        self.z_252 = z_score_or_zero(price, stats.moving_averages[252],
                                     sqrt_or_zero(stats.moving_variance[252]))
        self.beta = stats.beta
        self.sigma_market_15 = sqrt_or_zero(market_stats.moving_variance[15])
        self.sigma_market_35 = sqrt_or_zero(market_stats.moving_variance[35])
        self.sigma_market_252 = sqrt_or_zero(market_stats.moving_variance[252])
        self.z_market_15 = z_score_or_zero(
            market_price, market_stats.moving_averages[15],
            self.sigma_market_15)
        self.z_market_35 = z_score_or_zero(
            market_price, market_stats.moving_averages[35],
            self.sigma_market_35)
        self.z_market_252 = z_score_or_zero(
            market_price, market_stats.moving_averages[252],
            self.sigma_market_252)

    def csv_headers(self):
        """
        Returns csv headers (static).
        """
        return ["price", "market_price", "z_vol_15", "z_15", "z_35", "z_252",
                "beta", "sigma_market_15", "sigma_market_35",
                "sigma_market_252", "z_market_15", "z_market_35",
                "z_market_252"]

    def to_vector(self):
        """
        Returns self as a vector of floats (length 13).
        """
        return [self.price, self.market_price, self.z_vol_15, self.z_15,
                self.z_35, self.z_252, self.beta, self.sigma_market_15,
                self.sigma_market_35, self.sigma_market_252, self.z_market_15,
                self.z_market_35, self.z_market_252]


class Pair:
    """
    Pair of Individuals and their correlation.
    """

    def __init__(self, first, second, correlation, daily_data):
        self.first_ticker = first
        self.first_fundamentals = Individual(first, daily_data)
        self.second_ticker = second
        self.second_fundamentals = Individual(second, daily_data)
        self.correlation = correlation

    def to_vector(self, header=False):
        """
        Returns the concatenation of the two vectors corresponding to the two
        Individuals. (First, then second. Length 26.)
        """
        return (self.first_fundamentals.to_vector() +
                self.second_fundamentals.to_vector() + [self.correlation])


class CorrelatedPairs(algorithm_io_pb2_grpc.AlgorithmServicer):
    """
    Correlated Pairs algorithm.
    """

    def skip_trading(self):
        """
        Optionally skip trading to speed things up. Won't query for portfolio
        and prices and will immediately call DoneTrading.
        """
        return False

    def get_candidate_pairs(self, portfolio, daily_data):
        """
        Gets all candidate pairs from the daily data.
        """
        output = []
        for pair_stats in daily_data.pair_stats:
            if (pair_stats.first not in daily_data.prices or
                    pair_stats.second not in daily_data.prices):
                continue
            output.append(
                Pair(pair_stats.first, pair_stats.second,
                     pair_stats.correlation, daily_data))
        return output

    def filter_candidate_pairs(self, candidate_pairs):
        """
        Filters for candidate pairs that (1) aren't already in a position, and
        (2) pass the quality criteria.
        TODO: Implement.
        """
        return candidate_pairs

    def ground_truth(self, candidate_pair, daily_data):
        """
        Calculates and returns the ground truth label for the candidate pair
        given the daily_data after `period` days.
        NOTE: It would probably be better to store the actual estimate of the
        buy price. We imagine looking at a pair and deciding whether to open
        a position tomorrow, we should probably be comparing the prices of
        (day + 1) to those of (day + period + 1).
        """
        first = candidate_pair.first_ticker
        second = candidate_pair.second_ticker
        prices = daily_data.prices
        if first not in prices or second not in prices:
            return 0.0
        first_sale_price = (prices[first].open + prices[first].close) / 2.0
        first_buy_price = candidate_pair.first_fundamentals.price
        first_perf = (first_sale_price - first_buy_price) / first_buy_price
        second_sale_price = (prices[second].open + prices[second].close) / 2.0
        second_buy_price = candidate_pair.second_fundamentals.price
        second_perf = (second_sale_price - second_buy_price) / second_buy_price
        return first_perf - second_perf

    def write_training_examples(self, candidate_pairs, daily_data):
        """
        Calculate labels for each candidate pair and write them to csv. If you
        change this, make sure to change the headers in __init__.
        """
        for candidate_pair in candidate_pairs:
            label = self.ground_truth(candidate_pair, daily_data)
            self.output_csv.writerow(candidate_pair.to_vector() + [label])

    def maybe_process_training_data(self, candidate_pairs, daily_data):
        """
        Writes training data to the output csv file and records examples.
        TODO: Implement.
        """
        if not self.export_training_data:
            return
        if len(self.examples) >= period:
            self.write_training_examples(self.examples[0], daily_data)
        self.examples.append(candidate_pairs)

    def trade(self, portfolio, daily_data):
        """
        Tells the algorithm to trade. Mostly unimplemented.
        """
        candidate_pairs = self.get_candidate_pairs(portfolio, daily_data)
        self.maybe_process_training_data(candidate_pairs, daily_data)
        candidate_pairs = self.filter_candidate_pairs(candidate_pairs)
        # Close expiring pairs.
        # TODO: Implement this.
        # Create new long / short pairs.
        # TODO: Implement this.
        return []

    def __init__(self, id, export_training_data, training_data_path):
        """
        Constructor for correlated pairs model.
        """
        logging.basicConfig(
            format='%(asctime)s %(levelname)-8s %(message)s',
            level=logging.INFO,
            datefmt='%Y-%m-%d %H:%M:%S')
        # For core algorithm functionality.
        self.algorithm_id = supervisor_pb2.AlgorithmId(algorithm_id=id)
        self.id = id
        self.assets_in_pairs = set()
        self.pairs = collections.deque(maxlen=period)
        self.registrar = Registrar()
        # For training.
        self.export_training_data = export_training_data
        if export_training_data:
            self.training_data_path = training_data_path
            self.examples = collections.deque(maxlen=period)
            self.output_file = open(self.training_data_path, 'w')
            self.output_csv = csv.writer(self.output_file)
            # Write CSV headers.
            ind_headers = Individual.csv_headers(None)
            self.output_csv.writerow(
                ["first.%s" % s for s in ind_headers] +
                ["second.%s" % s for s in ind_headers] +
                ["correlation", "label"])

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


if __name__ == '__main__':
    # Parse flags
    parser = argparse.ArgumentParser(description='Correlated Pairsalgorithm.')
    parser.add_argument('--id', type=str, help='Algorithm name',
                        default='headsortails', required=False)
    parser.add_argument('--port', type=str, help='Serving port',
                        default='17507', required=False)
    parser.add_argument('--export_training_data', type=bool,
                        help='Export training data', default=False,
                        required=False)
    parser.add_argument(
        '--training_data_path', type=str, help='Training data path',
        default='algorithm/correlated_pairs/train/data/data.csv',
        required=False)
    args = parser.parse_args()
    # Serve
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    algorithm_io_pb2_grpc.add_AlgorithmServicer_to_server(
        CorrelatedPairs(args.id, args.export_training_data,
                        args.training_data_path),
        server)
    server.add_insecure_port('[::]:%s' % args.port)
    server.start()
    logging.info('Listening on `localhost:%s`' % args.port)
    server.wait_for_termination()
