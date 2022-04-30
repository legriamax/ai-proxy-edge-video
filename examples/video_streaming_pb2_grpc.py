# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import video_streaming_pb2 as video__streaming__pb2


class ImageStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.VideoLatestImage = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/VideoLatestImage',
                request_serializer=video__streaming__pb2.VideoFrameRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.VideoFrame.FromString,
                )
        self.VideoBufferedImage = channel.unary_stream(
                '/chrys.cloud.videostreaming.v1beta1.Image/VideoBufferedImage',
                request_serializer=video__streaming__pb2.VideoFrameBufferedRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.VideoFrame.FromString,
                )
        self.VideoProbe = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/VideoProbe',
                request_serializer=video__streaming__pb2.VideoProbeRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.VideoProbeResponse.FromString,
                )
        self.ListStreams = channel.unary_stream(
                '/chrys.cloud.videostreaming.v1beta1.Image/ListStreams',
                request_serializer=video__streaming__pb2.ListStreamRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.ListStream.FromString,
                )
        self.Annotate = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/Annotate',
                request_serializer=video__streaming__pb2.AnnotateRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.AnnotateResponse.FromString,
                )
        self.Proxy = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/Proxy',
                request_serializer=video__streaming__pb2.ProxyRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.ProxyResponse.FromString,
                )
        self.Storage = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/Storage',
                request_serializer=video__streaming__pb2.StorageRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.StorageResponse.FromString,
                )
        self.SystemTime = channel.unary_unary(
                '/chrys.cloud.videostreaming.v1beta1.Image/SystemTime',
                request_serializer=video__streaming__pb2.SystemTimeRequest.SerializeToString,
                response_deserializer=video__streaming__pb2.SystemTimeResponse.FromString,
                )


class ImageServicer(object):
    """Missing associated documentation comment in .proto file."""

    def VideoLatestImage(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def VideoBufferedImage(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def VideoProbe(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListStreams(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Annotate(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Proxy(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def Storage(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def SystemTime(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ImageServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'VideoLatestImage': grpc.unary_unary_rpc_method_handler(
                    servicer.VideoLatestImage,
                    request_deserializer=video__streaming__pb2.VideoFrameRequest.FromString,
                    response_serializer=video__streaming__pb2.VideoFrame.SerializeToString,
            ),
            'VideoBufferedImage': grpc.unary_stream_rpc_method_handler(
                    servicer.VideoBufferedImage,
                    request_deserializer=video__streaming__pb2.VideoFrameBufferedRequest.FromString,
                    response_serializer=video__streaming__pb2.VideoFrame.SerializeToString,
            ),
            'VideoProbe': grpc.unary_unary_rpc_method_handler(
                    servicer.VideoProbe,
                    request_deserializer=video__streaming__pb2.VideoProbeRequest.FromString,
                    response_serializer=video__streaming__pb2.VideoProbeResponse.SerializeToString,
            ),
            'ListStreams': grpc.unary_stream_rpc_method_handler(
                    servicer.ListStreams,
                    request_deserializer=video__streaming__pb2.ListStreamRequest.FromString,
                    response_serializer=video__streaming__pb2.ListStream.SerializeToString,
            ),
            'Annotate': grpc.unary_unary_rpc_method_handler(
                    servicer.Annotate,
                    request_deserializer=video__streaming__pb2.AnnotateRequest.FromString,
                    response_serializer=video__streaming__pb2.AnnotateResponse.SerializeToString,
            ),
            'Proxy': grpc.unary_unary_rpc_method_handler(
                    servicer.Proxy,
                    request_deserializer=video__streaming__pb2.ProxyRequest.FromString,
                    response_serializer=video__streaming__pb2.ProxyResponse.SerializeToString,
            ),
            'Storage': grpc.unary_unary_rpc_method_handler(
                    servicer.Storage,
                    request_deserializer=video__streaming__pb2.StorageRequest.FromString,
                    response_serializer=video__streaming__pb2.StorageResponse.SerializeToString,
            ),
            'SystemTime': grpc.unary_unary_rpc_method_handler(
                    servicer.SystemTime,
                    request_deserializer=video__streaming__pb2.SystemTimeRequest.FromString,
                    response_serializer=video__streaming__pb2.SystemTimeResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'chrys.cloud.videostreaming.v1beta1.Image', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Image(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def VideoLatestImage(request,
            target,
            options=(),
            channel_credentials=None,
            call_creden