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

    def test_many_stocks_some_already_held(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "usd": 664.0,
            "stocks": {
                "MSFT": 7.0,
                "NVDA": 14.0,
                "GOOG": 7.0,
            }
        }), supervisor_pb2.Portfolio())
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
            "MSFT": 7.0,
            "NVDA": 1.0,
        }
        expected_limit = {
            "APPL": 151.5,
            "FB": 7.07,
            "FORD": 1.01,
            "GM": 0.505,
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
        self.assertEqual(len(orders), 6)
        self.assertEqual(volume, expected_volume)
        self.assertEqual(limit, expected_limit)
        self.assertLess(657.0, total_limit)
        self.assertLess(total_limit, 664.0)

    def test_partial(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "usd": 100,
            "stocks": {
                "APPL": 1.0,
            }
        }), supervisor_pb2.Portfolio())
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 25.0},
                "GOOG": {"close": 25.0},
                "FB": {"close": 2.5},
                "APPL": {"close": 33.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "FB": 13.0,
            "GOOG": 1.0,
        }
        expected_limit = {
            "FB": 2.525,
            "GOOG": 25.25,
        }
        # Create orders to be tested
        orders = portfolio_helpers.invest_approximately_uniformly_in_targets(
            test_algorithm_id, portfolio, daily_data, ["GOOG", "FB"], 80.0)
        volume = dict()
        limit = dict()
        total_limit = 0.0
        for order in orders:
            if order.ticker not in volume:
                volume[order.ticker] = 0.0
            volume[order.ticker] += order.volume
            limit[order.ticker] = order.limit
            total_limit += order.limit * order.volume
        # Assertions
        self.assertEqual(len(orders), 2)
        self.assertEqual(volume, expected_volume)
        self.assertEqual(limit, expected_limit)
        self.assertLess(58, total_limit)
        self.assertLess(total_limit, 80.0)


class TestOrdersSortedDescending(unittest.TestCase):

    def test_happy_path(self):
        orders = [
            supervisor_pb2.Order(volume=5.0, limit=5.0),
            supervisor_pb2.Order(volume=6.0, limit=8.0),
            supervisor_pb2.Order(volume=7.0, limit=9.0),
            supervisor_pb2.Order(volume=8.0, limit=2.0),
            supervisor_pb2.Order(volume=-8.0, limit=2.0),
            supervisor_pb2.Order(volume=-9.0, limit=3.0),
        ]
        sorted_orders = portfolio_helpers.orders_sorted_descending(orders)
        for i in range(1, len(sorted_orders)):
            prev_order = sorted_orders[i-1]
            order = sorted_orders[i]
            # Sells before buys.
            self.assertFalse(prev_order.volume > 0.0 and order.volume <= 0.0)
            if prev_order.volume <= 0.0:
                continue
            # Buys in descending order.
            self.assertLess(order.limit * order.volume,
                            prev_order.limit * prev_order.volume)


class TestSellOverweightTargetStocks(unittest.TestCase):

    def test_simple(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "stocks": {
                "MSFT": 6.0,
                "GOOG": 7.0,
                "FB": 6.0,
                "APPL": 7.0,
            },
            "usd": -50.0,
        }), supervisor_pb2.Portfolio())
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 10.0},
                "GOOG": {"close": 10.0},
                "FB": {"close": 11.0},
                "APPL": {"close": 8.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "APPL": -1.0,
            "FB": -1.0,
            "GOOG": -2.0,
            "MSFT": -1.0,
        }
        # Create orders to be tested
        targets = ["GOOG", "FB", "APPL", "MSFT"]
        orders = portfolio_helpers.sell_overweight_target_stocks(
            test_algorithm_id, portfolio, daily_data, targets)
        # Assertions
        for order in orders:
            self.assertEqual(order.volume, expected_volume[order.ticker])
            self.assertEqual(order.stop, 0.99 *
                             daily_data.prices[order.ticker].close)
        self.assertEqual(len(orders), 4)

    def test_restricted(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "stocks": {
                "MSFT": 6.0,
                "GOOG": 7.0,
                "FB": 6.0,
                "APPL": 7.0,
            },
            "usd": -50.0,
        }), supervisor_pb2.Portfolio())
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 10.0},
                "GOOG": {"close": 10.0},
                "FB": {"close": 11.0},
                "APPL": {"close": 8.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "APPL": -1.0,
            "FB": -1.0,
        }
        # Create orders to be tested
        targets = ["FB", "APPL"]
        orders = portfolio_helpers.sell_overweight_target_stocks(
            test_algorithm_id, portfolio, daily_data, targets)
        # Assertions
        for order in orders:
            self.assertEqual(order.volume, expected_volume[order.ticker])
            self.assertEqual(order.stop, 0.99 *
                             daily_data.prices[order.ticker].close)
        self.assertEqual(len(orders), 2)

    def test_restricted_inverted(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "stocks": {
                "MSFT": 6.0,
                "GOOG": 7.0,
                "FB": 6.0,
                "APPL": 7.0,
            },
            "usd": -50.0,
        }), supervisor_pb2.Portfolio())
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 10.0},
                "GOOG": {"close": 10.0},
                "FB": {"close": 11.0},
                "APPL": {"close": 8.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "APPL": -1.0,
            "FB": -1.0,
        }
        # Create orders to be tested
        targets = ["GOOG", "MSFT", "BOGUS"]
        orders = portfolio_helpers.sell_overweight_target_stocks(
            test_algorithm_id, portfolio, daily_data, targets,
            invert_targets=True)
        # Assertions
        for order in orders:
            self.assertEqual(order.volume, expected_volume[order.ticker])
            self.assertEqual(order.stop, 0.99 *
                             daily_data.prices[order.ticker].close)
        self.assertEqual(len(orders), 2)

    def test_underweight(self):
        # Inputs
        portfolio = json_format.Parse(json.dumps({
            "stocks": {
                "MSFT": 0.0,
                "GOOG": 7.0,
                "FB": 6.0,
                "APPL": 7.0,
            },
            "usd": 10.0,
        }), supervisor_pb2.Portfolio())
        daily_data = json_format.Parse(json.dumps({
            "prices": {
                "MSFT": {"close": 10.0},
                "GOOG": {"close": 10.0},
                "FB": {"close": 11.0},
                "APPL": {"close": 8.0},
            },
        }), daily_prices_pb2.DailyData())
        # Expected outputs
        expected_volume = {
            "APPL": -1.0,
            "FB": -1.0,
            "GOOG": -2.0,
        }
        # Create orders to be tested
        targets = ["GOOG", "FB", "APPL", "MSFT"]
        orders = portfolio_helpers.sell_overweight_target_stocks(
            test_algorithm_id, portfolio, daily_data, targets)
        # Assertions
        for order in orders:
            self.assertEqual(order.volume, expected_volume[order.ticker])
            self.assertEqual(order.stop, 0.99 *
                             daily_data.prices[order.ticker].close)
        self.assertEqual(len(orders), 3)


if __name__ == '__main__':
    unittest.main()
