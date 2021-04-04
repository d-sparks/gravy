import gravy.portfolio_helpers as portfolio_helpers
import json
import unittest
import google.protobuf.json_format as json_format

from data.dailyprices.proto import daily_prices_pb2
from supervisor.proto import supervisor_pb2


test_algorithm_id = supervisor_pb2.AlgorithmId(algorithm_id="test")


class TestInvestApproximatelyUniformly(unittest.TestCase):

    def test_many_stocks(self):
        # Inputs
        portfolio = supervisor_pb2.Portfolio(usd=1000.0)
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 10.0},
                "GOOG": {"close": 20.0},
                "FB": {"close": 7.0},
                "APPL": {"close": 150.0},
                "NVDA": {"close": 9.0},
                "GM": {"close": 0.50},
                "FORD": {"close": 1.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "APPL": 1.0,
            "FB": 20.0,
            "FORD": 141.0,
            "GM": 282.0,
            "GOOG": 7.0,
            "MSFT": 14.0,
            "NVDA": 15.0,
        }
        expected_limit = {
            "APPL": 151.5,
            "FB": 7.07,
            "FORD": 1.01,
            "GM": 0.505,
            "GOOG": 20.2,
            "MSFT": 10.1,
            "NVDA": 9.09,
        }
        # Create orders to be tested
        orders = portfolio_helpers.invest_approximately_uniformly(
            test_algorithm_id, portfolio, daily_data)
        volume = dict()
        limit = dict()
        total_limit = 0.0
        for order in orders:
            volume[order.ticker] = order.volume
            limit[order.ticker] = order.limit
            total_limit += order.limit * order.volume
        # Assertions
        self.assertEqual(len(orders), 7)
        self.assertEqual(volume, expected_volume)
        self.assertEqual(limit, expected_limit)
        self.assertLess(996.0, total_limit)


if __name__ == '__main__':
    unittest.main()
