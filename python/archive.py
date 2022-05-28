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

import av
import time
import threading, queue
import os
import datetime
from random import randint

class StoreMP4VideoChunks(threading.Thread):

    def __init__(self, queue=None, path=None, device_id=None, video_stream=None, audio_stream=None):
        threading.Thread.__init__(self) 
        self.in_video_stream = video_stream
        self.in_audio_stream = audio_stream
        self.path = os.path.join(path, '') + device_id
        self.device_id = device_id
        self.q = queue
        if not os.path.exists(self.path):
            os.makedirs(self.path)
    
    def run(self):
        while True:
            try:
                archive_group = self.q.get(timeout=5) # 5s timeout
                # print(archive_group.packet_group, archive_group.start_timestamp)
                self.saveToMp4(archive_group.packet_group, archive_group.start_timestamp)
            except queue.Empty:
                continue

            self.q.task_done()
        pass

    def saveToMp4(self, packet_store, start_timestamp):
        minimum_dts = -