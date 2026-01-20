# Python
import numpy as np
import sys

# Internal
from gen import computervision_pb2, computervision_pb2_grpc

import plot

# client stub functions for computervision grpc service
class ComputerVisionClient():
    def __init__(self, channel):
        self.stub = computervision_pb2_grpc.ComputerVisionServiceStub(channel)
        self.min_x = sys.float_info.max
        self.max_x = sys.float_info.min
        self.min_y = sys.float_info.max
        self.max_y = sys.float_info.min

    def generate_requests_from_frames(self, frames):
        for frame in frames:
            request = computervision_pb2.GetPoseDataRequest(image=frame)
            yield request

    def get_pose_data_from_video(self, frames):
        response_iterator = self.stub.GetPoseDataFromVideo(self.generate_requests_from_frames(frames))
        datapoints_for_swing = np.zeros((150, 25*3))
        frame_idx = 0
        for response in response_iterator:
            np_arr_dp = self.convert_body_25_datapoints_to_np_arr(response.datapoints)
            np_arr_normalized = self.normalize_datapoints(np_arr_dp)
            plot.convert_data_to_body25_datapoints(np_arr_normalized)
            datapoints_for_swing[frame_idx] = np_arr_normalized
            frame_idx += 1
        if frame_idx < 150:
            for i in range(frame_idx, 150):
                datapoints_for_swing[i] = datapoints_for_swing[frame_idx-1]
        return datapoints_for_swing
    

    def convert_body_25_datapoints_to_np_arr(self, datapoints):
        dps_arr = []
        for descriptor in datapoints.DESCRIPTOR.fields:
            datapoint = getattr(datapoints, descriptor.name)
            for fd in datapoint.DESCRIPTOR.fields:
                attr = getattr(datapoint, fd.name)
                dps_arr.append(attr)
                if attr != 0:
                    if fd.name == "x":
                        if attr < self.min_x:
                            self.min_x = attr
                        if attr > self.max_x:
                            self.max_x = attr
                    elif fd.name == "y":
                        if attr < self.min_y:
                            self.min_y = attr
                        if attr > self.max_y:
                            self.max_y = attr
        return np.array(dps_arr)

    # crop frames to midhip to neck length*3 by midhip to neck length*3
    # center around neck (all points are offset from that)???? Do i need to do this???
    # normalize lengths to between 0 and 1 (xval- minx / (maxx - minx)) same for y
    def normalize_datapoints(self, np_arr_dp):
        #print(f'neck: {self.neck_x}, {self.neck_y}, min: {self.min_x}, {self.min_y}, max: {self.max_x}, {self.max_y}')
        normalized_arr = np.zeros((25*3))
        for i in range(0, 25*3, 3):
            x = np_arr_dp[i]
            if x != 0:
                normalized_x = (x - self.min_x) / (self.max_x - self.min_x)
            else:
                normalized_x = 0
            y = np_arr_dp[i+1]
            if y != 0:
                normalized_y = (y - self.min_y) / (self.max_y - self.min_y) 
            else:
                normalized_y = 0
            normalized_arr[i] = normalized_x
            normalized_arr[i+1] = normalized_y
            normalized_arr[i+2] = np_arr_dp[i+2]
        return normalized_arr
