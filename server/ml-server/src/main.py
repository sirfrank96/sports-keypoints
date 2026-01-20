# Python
import logging
from pathlib import Path
import os
import sys
import numpy as np

# GRPC
import grpc

# Internal
curr_dir = Path(__file__).parent.resolve()
sys.path.append(os.path.join(curr_dir, 'gen'))
import cvclient as cv
import autoencoder
import plot

import cv2

def frame_to_bytes(frame):
    ok, buf = cv2.imencode('.jpg', frame)
    if ok:
        bytes = buf.tobytes() 
        return bytes
    else:
        return None   

def serve():
    setting_timeout_ms = 1000 * 60 * 3
    options = [('grpc.http2.settings_timeout', setting_timeout_ms)]
    channel = grpc.insecure_channel('localhost:50051', options=options)
    return cv.ComputerVisionClient(channel)

if __name__ == "__main__":
    logging.basicConfig()
    cv_client = serve()
    # pull in videos
    num_videos = 5
    all_datapoints = np.zeros((num_videos, 150, 25*3))
    for i in range(num_videos):
        dtlvid = cv2.VideoCapture(f'E:\\Pictures\\Phone Pictures\\Swing Videos\\Dtl-ml\\dtl{i}.mp4')
        # parse video into frames
        num_frames = 0
        frames = []
        while True:
            ok, image = dtlvid.read()
            if ok:
                num_frames += 1
                frames.append(frame_to_bytes(image))
            else:
                break
        print(f'number of frames {num_frames}')
        dtlvid.release()
        frames = frames[0:150]
        # query computervision-service for body datapoints
        datapoints_swing_i = None
        try:
            datapoints_swing_i = cv_client.get_pose_data_from_video(frames)
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # put bodydatapoints into tensors (num swings x num frames x num datapoints x datapoint x,y,conf) (add extra dimension for more num people)
        all_datapoints[i] = datapoints_swing_i
    # put as input to autoencoder
    autoencoder.autoencode(all_datapoints[:4], all_datapoints[:4], all_datapoints[4:], all_datapoints[4:])
