# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

from . import daily_prices_pb2 as daily__prices__pb2


class DataStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Get = channel.unary_unary(
                '/dailyprices.Data/Get',
                request_serializer=daily__prices__pb2.Request.SerializeToString,
                response_deserializer=daily__prices__pb2.DailyData.FromString,
                )
        self.NewSession = channel.unary_unary(
                '/dailyprices.Data/NewSession',
                request_serializer=daily__prices__pb2.NewSessionRequest.SerializeToString,
                response_deserializer=daily__prices__pb2.NewSessionResponse.FromString,
                )
        self.TradingDatesInRange = channel.unary_unary(
                '/dailyprices.Data/TradingDatesInRange',
                request_serializer=daily__prices__pb2.Range.SerializeToString,
                response_deserializer=daily__prices__pb2.TradingDates.FromString,
                )


class DataServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Get(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def NewSession(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def TradingDatesInRange(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_DataServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Get': grpc.unary_unary_rpc_method_handler(
                    servicer.Get,
                    request_deserializer=daily__prices__pb2.Request.FromString,
                    response_serializer=daily__prices__pb2.DailyData.SerializeToString,
            ),
            'NewSession': grpc.unary_unary_rpc_method_handler(
                    servicer.NewSession,
                    request_deserializer=daily__prices__pb2.NewSessionRequest.FromString,
                    response_serializer=daily__prices__pb2.NewSessionResponse.SerializeToString,
            ),
            'TradingDatesInRange': grpc.unary_unary_rpc_method_handler(
                    servicer.TradingDatesInRange,
                    request_deserializer=daily__prices__pb2.Range.FromString,
                    response_serializer=daily__prices__pb2.TradingDates.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'dailyprices.Data', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Data(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Get(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/dailyprices.Data/Get',
            daily__prices__pb2.Request.SerializeToString,
            daily__prices__pb2.DailyData.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def NewSession(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/dailyprices.Data/NewSession',
            daily__prices__pb2.NewSessionRequest.SerializeToString,
            daily__prices__pb2.NewSessionResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def TradingDatesInRange(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/dailyprices.Data/TradingDatesInRange',
            daily__prices__pb2.Range.SerializeToString,
            daily__prices__pb2.TradingDates.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
