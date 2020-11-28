import grpc
import logging

from algorithm.proto import algorithm_io_pb2
from algorithm.proto import algorithm_io_pb2_grpc
from concurrent import futures
from data.dailyprices.proto import daily_prices_pb2
from datetime import datetime, timezone
from google.protobuf import timestamp_pb2
from registrar.registrar import Registrar
from supervisor.proto import supervisor_pb2


class HeadsOrTails(algorithm_io_pb2_grpc.AlgorithmServicer):
    """
    Heads or Tails inference algorithm.
    """

    def skip_trading(self):
        return False

    def trade(self, portfolio, prices):
        return

    def __init__(self, id):
        logging.basicConfig(
            format='%(asctime)s %(levelname)-8s %(message)s',
            level=logging.INFO,
            datefmt='%Y-%m-%d %H:%M:%S')

        self.algorithm_id = supervisor_pb2.AlgorithmId(algorithm_id=id)
        self.id = id
        self.registrar = Registrar()

    def Execute(self, input, context):
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
