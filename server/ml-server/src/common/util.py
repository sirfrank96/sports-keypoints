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

def downsample_frames(frames, frames_dim):
    num_frames = len(frames)
    if num_frames <= frames_dim:
        return frames
    # downsample frames if num_frames is greater than required
    frame_rate = num_frames / frames_dim
    frame_rate = truncate_float_to_tenths(frame_rate)
    downsampled_frames = []
    last_idx = -1
    for i in range(0, num_frames):
        if last_idx >= frames_dim - 1:
            break
        idx_target = int(i/frame_rate)
        if idx_target > last_idx:
            downsampled_frames.append(frames[i])
            last_idx = idx_target
    return downsampled_frames

# returns x_train, x_test/val, y_train, y_test/val
def get_random_training_and_validation_sets(inputs, outputs, test_size=0.2):
    return train_test_split(inputs, outputs, test_size=test_size, shuffle=True)

def truncate_float_to_tenths(f):
    return math.trunc(f * 10) / 10.0

def load_all_body_datapoints_from_file():
    if os.path.exists(constants.BODY_DATAPOINTS_PATH):
        return np.load(constants.BODY_DATAPOINTS_PATH)
    else:
        return create_new_all_body_datapoints()

def create_new_all_body_datapoints():
    return np.empty((0, constants.FRAME_DIMENSION, constants.BODY_DATAPOINTS_DIMENSIONS))
    
def load_faceon_indices_from_file():
    if os.path.exists(constants.FACEON_TRUE_VALUES_PATH):
        return np.load(constants.FACEON_TRUE_VALUES_PATH)
    else:
        return create_new_faceon_indices()
    
def create_new_faceon_indices():
    return np.empty((0, constants.FRAME_DIMENSION))

def load_isswing_indices_from_file():
    if os.path.exists(constants.ISSWING_TRUE_VALUES_PATH):
        return np.load(constants.ISSWING_TRUE_VALUES_PATH)
    else:
        return create_new_isswing_indices()
    
def create_new_isswing_indices():
    return np.empty((0, constants.FRAME_DIMENSION))
    
def save_all_body_datapoints(datapoints):
    np.save(constants.BODY_DATAPOINTS_PATH, datapoints)

def save_faceon_indices(faceon_indices):
    np.save(constants.FACEON_TRUE_VALUES_PATH, faceon_indices)

def save_isswing_indices(isswing_indices):
    np.save(constants.ISSWING_TRUE_VALUES_PATH, isswing_indices)

def add_unlabeled_video(video_name, idx, unlabeled_dict):
    unlabeled_dict[video_name] = idx

def remove_unlabeled_video(video_name, unlabeled_dict):
    del unlabeled_dict[video_name] 
    
def load_unlabeled_dict():
    if os.path.exists(constants.UNLABELED_VIDEOS_DICT_PATH):
        with open(constants.UNLABELED_VIDEOS_DICT_PATH, 'r') as f:
            return json.load(f)
    else:
        return create_new_unlabeled_dict()
    
def create_new_unlabeled_dict():
    return {}

def save_unlabeled_dict(unlabeled_dict):
    with open(constants.UNLABELED_VIDEOS_DICT_PATH, 'w') as f:
        json.dump(unlabeled_dict, f)

def get_saved_video_path(video_name):
    return os.path.join(constants.VIDEOS_DIR_PATH, video_name)

def grab_saved_video(video_path):
    return cv2.VideoCapture(video_path)
