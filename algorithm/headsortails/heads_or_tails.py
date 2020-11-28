import gravy

import grpc
import logging
import math

from algorithm.proto import algorithm_io_pb2
from algorithm.proto import algorithm_io_pb2_grpc
from concurrent import futures
from data.dailyprices.proto import daily_prices_pb2
from datetime import datetime, timezone
from google.protobuf import timestamp_pb2
from registrar.registrar import Registrar
from supervisor.proto import supervisor_pb2


min_results = 15


def z_score_or_zero(x, mu, sigma):
    """
    Returns the z score unless sigma is prohibitively small, in which case
    returns 0.0. TODO: Move to a helper package.
    """
    if sigma < 1E-6:
        return 0.0
    return (x - mu) / sigma


class HeadsOrTails(algorithm_io_pb2_grpc.AlgorithmServicer):
    """
    Heads or Tails inference algorithm.
    """

    def skip_trading(self):
        """
        Optionally skip trading to speed things up. Won't query for portfolio
        and prices and will immediately call DoneTrading.
        """
        return False

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

        z_vol_15 = z_score_or_zero(prices.volume, stats.moving_volume[15],
                                   math.sqrt(stats.moving_volume_variance[15]))
        z_15 = z_score_or_zero(price, stats.moving_averages[15],
                               math.sqrt(stats.moving_variance[15]))
        z_35 = z_score_or_zero(price, stats.moving_averages[35],
                               math.sqrt(stats.moving_variance[35]))
        z_252 = z_score_or_zero(price, stats.moving_averages[252],
                                math.sqrt(stats.moving_variance[252]))
        beta = stats.beta
        sigma_market_15 = math.sqrt(market_stats.moving_variance[15])
        sigma_market_35 = math.sqrt(market_stats.moving_variance[35])
        sigma_market_252 = math.sqrt(market_stats.moving_variance[252])
        z_market_15 = z_score_or_zero(
            market_price, market_stats.moving_averages[15], sigma_market_15)
        z_market_35 = z_score_or_zero(
            market_price, market_stats.moving_averages[35], sigma_market_35)
        z_market_252 = z_score_or_zero(
            market_price, market_stats.moving_averages[252], sigma_market_252)

        return [z_vol_15, z_15, z_35, z_252, beta, sigma_market_15,
                sigma_market_35, sigma_market_252, z_market_15, z_market_35,
                z_market_252]

    def trade(self, portfolio, daily_data):
        """
        Tells the algorithm to trade. Mostly unimplemented.
        """
        self.num_results += 1
        if self.num_results < min_results:
            return

        for ticker in daily_data.prices.items():
            print(extract_features(ticker, daily_data))
            break

        return

    def __init__(self, id):
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

    def Execute(self, input, context):
        """
        Implements the algorithm_io interface and tells the algorithm to
        execute any trades before calling DoneTrading.
        """
        timestamp = datetime.fromtimestamp(
            input.timestamp.seconds, timezone.utc).strftime('%Y-%m-%d')
        logging.info('Executing algorithm on %s' % timestamp)

        if not self.skip_trading():
            portfolio = self.registrar.supervisor_stub.GetPortfolio(
                self.algorithm_id)
            daily_data = self.registrar.dailyprices_stub.Get(
                daily_prices_pb2.Request(timestamp=input.timestamp, version=0))
            self.trade(portfolio, daily_data)

        # Indicate done trading.
        self.registrar.supervisor_stub.DoneTrading(self.algorithm_id)
        return algorithm_io_pb2.Output()

    algorithm_id = None
    id = None
    registrar = None
    num_results = 0


if __name__ == '__main__':
    id = 'headsortails'
    port = '17506'
    # model_dir

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    algorithm_io_pb2_grpc.add_AlgorithmServicer_to_server(
        HeadsOrTails(id), server)
    server.add_insecure_port('[::]:%s' % port)
    server.start()
    logging.info('Listening on `localhost:%s`' % port)
    server.wait_for_termination()
