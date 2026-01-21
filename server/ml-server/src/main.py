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
import scraper

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

frames_dim = 150

def standardize_num_frames(frames):
    num_frames = len(frames)
    # pad frames if num_frames is less than required
    if num_frames <= frames_dim:
        last_elem = frames[-1]
        padding = [last_elem] * (frames_dim - num_frames)
        frames.extend(padding)
        return frames
    # downsample frames if num_frames is greater than required
    est_frame_rate = num_frames // frames_dim
    frame_range = frames_dim * est_frame_rate
    buffer = (num_frames - frame_range) // 2 # buffer on beginning and end of video to "center" the video
    downsampled_frames = []
    num_downsampled_frames = 0
    for i in range(0+buffer, num_frames-buffer, est_frame_rate):
        if num_downsampled_frames == frames_dim:
            break
        downsampled_frames.append(frames[i])
        num_downsampled_frames += 1
    return downsampled_frames

def get_body_datapoints_from_videos(num_downloaded_videos):
    # pull in videos
    all_datapoints = np.zeros((num_downloaded_videos, frames_dim, 25*3))
    for i in range(0, num_downloaded_videos):
        vid = cv2.VideoCapture(os.path.join(curr_dir, "data", "videos", f'{i+1}.mp4'))
        # parse video into frames
        num_frames = 0
        frames = []
        while True:
            ok, image = vid.read()
            if ok:
                num_frames += 1
                frames.append(frame_to_bytes(image))
            else:
                break
        print(f'number of frames {num_frames}')
        vid.release()
        frames = standardize_num_frames(frames)
        print(f'number of frames after standardize {len(frames)}')
        # query computervision-service for body datapoints
        datapoints_swing_i = None
        try:
            datapoints_swing_i = cv_client.get_pose_data_from_video(frames)
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # put bodydatapoints into tensors (num swings x num frames x num datapoints x datapoint x,y,conf) (add extra dimension for more num people)
        all_datapoints[i] = datapoints_swing_i
    # save datapoints to data/bodydatapoints folder
    datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    np.save(datapoints_file_path, all_datapoints)
    return all_datapoints

def auto_encode_codeflow(num_downloaded_videos):
    all_datapoints = get_body_datapoints_from_videos(num_downloaded_videos)
    # plot datapoints
    for i in range(50, 100):
        plot.plot_datapoints_as_skeleton(all_datapoints[0][i])
    # put as input to autoencoder
    autoencoder.autoencode(all_datapoints[:4], all_datapoints[:4], all_datapoints[4:], all_datapoints[4:]) 

def encode_and_classify_codeflow(num_downloaded_videos, faceon_indices):
    all_datapoints = get_body_datapoints_from_videos(num_downloaded_videos)
    autoencoder.encode_and_classify(all_datapoints[:30], faceon_indices[:30], all_datapoints[30:,], faceon_indices[30:])

def output_new_model_from_existing_data():
    datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    all_datapoints = np.load(datapoints_file_path)
    print(f"size if all_datapoints is {all_datapoints.shape}")
    faceon_indices_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'faceon_indices.npy')
    faceon_indices = np.load(faceon_indices_file_path)
    print(f"size if faceon indices is {faceon_indices.shape}")
    autoencoder.encode_and_classify(all_datapoints[:30], faceon_indices[:30], all_datapoints[30:,], faceon_indices[30:])

if __name__ == "__main__":
    logging.basicConfig()
    cv_client = serve()

    #num_downloaded_videos, faceon_indices = scraper.download_faceon_and_dtl_videos()
    #auto_encode_codeflow(num_downloaded_videos)
    #faceon_indices_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'faceon_indices.npy')
    #np.save(faceon_indices_file_path, faceon_indices)
    #encode_and_classify_codeflow(num_downloaded_videos, faceon_indices)
    output_new_model_from_existing_data()
    
