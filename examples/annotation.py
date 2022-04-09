# Copyright 2020 Wearless Tech Inc All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


import grpc
import video_streaming_pb2_grpc, video_streaming_pb2
import argparse
import time

def annotate(stub, device_name, event_type):
    """ Sending annotation to Chrysalis Cloud """


    annotation_request = video_streaming_pb2.AnnotateRequest()
    annotation_request.device_name = device_name
    annotation_request.type = event_type
    annotation_request.start_timestamp = int(round(time.time() * 1000))
    annotation_request.end_timestamp = int(round(time.time() * 1000))
    try:
        resp = stub.Annotate(annotation_request)
        print(resp)
    except grpc.RpcError as rpc_error_call:
        print("start proxy failed with", rpc_error_call.code(), rpc_error_call.details())


if __name__ == "__main__":
    # Initialize parser 
    parser = argp