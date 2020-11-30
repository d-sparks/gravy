import grpc
import sys

from data.dailyprices.proto import daily_prices_pb2_grpc
from data.dailyprices.proto import daily_prices_pb2
from datetime import datetime
from google.protobuf import timestamp_pb2
from supervisor.proto import supervisor_pb2_grpc
from supervisor.proto import supervisor_pb2


# Location of google/protobuf/timestamp.proto
sys.path.insert(1, '/usr/local/include')


# TODO: These are duplicated/hardcoded here and in registrar.go but should
#       really be parameters/flags.
supervisor_url = "localhost:17500"
dailyprices_url = "localhost:17501"


class Registrar:
    """
    Registrar is a helper class for connecting to the various backends.
    Currently only used in algorithms, thus only connects to dailyprices and
    supervisor.
    """

    def __init__(self):
        dailyprices_channel = grpc.insecure_channel(dailyprices_url)
        self.dailyprices_stub = daily_prices_pb2_grpc.DataStub(
            dailyprices_channel)

        supervisor_channel = grpc.insecure_channel(supervisor_url)
        self.supervisor_stub = supervisor_pb2_grpc.SupervisorStub(
            supervisor_channel)

    dailyprices_stub = None
    supervisor_stub = None
