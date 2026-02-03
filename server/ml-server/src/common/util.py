# Python
import cv2
import math
import json
import os
import numpy as np

# Keras
from sklearn.model_selection import train_test_split

# Internal
from . import constants

def frame_to_bytes(frame: np.ndarray) -> bytes:
    ok, buf = cv2.imencode('.jpg', frame)
    if ok:
        b = buf.tobytes() 
        return b
    else:
        return None
    
def bytes_to_frame(frame_bytes: bytes) -> np.ndarray:
    nparr = np.frombuffer(frame_bytes, np.uint8)
    frame = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    return frame

def downsample_frames(frames, frames_dim, vals=None):
    num_frames = len(frames)
    if num_frames <= frames_dim:
        return frames, vals
    # downsample frames if num_frames is greater than required
    frame_rate = num_frames / frames_dim
    frame_rate = truncate_float_to_tenths(frame_rate)
    downsampled_frames = []
    downsampled_vals = None
    if vals is not None:
        downsampled_vals = []
    last_idx = -1
    for i in range(0, num_frames):
        if last_idx >= frames_dim - 1:
            break
        idx_target = int(i/frame_rate)
        if idx_target > last_idx:
            downsampled_frames.append(frames[i])
            if vals is not None:
                downsampled_vals.append(vals[i])
            last_idx = idx_target
    return downsampled_frames, downsampled_vals

# returns x_train, x_test/val, y_train, y_test/val
def get_random_training_and_validation_sets(inputs, outputs, test_size=0.2):
    return train_test_split(inputs, outputs, test_size=test_size, shuffle=True)

def truncate_float_to_tenths(f):
    return math.trunc(f * 10) / 10.0

def load_start_end_frames_map(map_file):
    # load map if saved previously
    video_start_end_frames_map = {}
    try:
        with open(map_file, 'r') as f:
            video_start_end_frames_map = json.load(f)
    except FileNotFoundError:
        print("No map, will create one")
    return video_start_end_frames_map

def load_all_body_datapoints_from_file():
    all_datapoints = None
    if os.path.exists(constants.BODY_DATAPOINTS_PATH):
        all_datapoints = np.load(constants.BODY_DATAPOINTS_PATH)
    return all_datapoints
    
def load_faceon_indices_from_file():
    faceon_indices = None
    if os.path.exists(constants.FACEON_TRUE_VALUES_PATH):
        faceon_indices = np.load(constants.FACEON_TRUE_VALUES_PATH)
    return faceon_indices

def load_isswing_indices_from_file():
    isswing_indices = None
    if os.path.exists(constants.ISSWING_TRUE_VALUES_PATH):
        isswing_indices = np.load(constants.ISSWING_TRUE_VALUES_PATH)
    return isswing_indices
