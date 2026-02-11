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
from . import labeldata

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
    
def save_all_body_datapoints(datapoints):
    np.save(constants.BODY_DATAPOINTS_PATH, datapoints)

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

def load_labeled_golf_swing_dataset():
    if os.path.exists(constants.LABELED_GOLF_SWING_DATASET_PATH):
        with open(constants.LABELED_GOLF_SWING_DATASET_PATH, 'r') as f:
            return json.load(f)
    else:
        return create_new_labeled_golf_swing_dataset()
    
def create_new_labeled_golf_swing_dataset():
    return labeldata.LabeledGolfSwingDataset(players={}, videos={})

def save_labeled_golf_swing_dataset(labeled_golf_swing_dataset):
    with open(constants.LABELED_GOLF_SWING_DATASET_PATH, 'w') as f:
        json.dump(labeled_golf_swing_dataset, f)

def get_saved_video_path(video_name):
    return os.path.join(constants.VIDEOS_DIR_PATH, video_name)

def grab_saved_video(video_path):
    return cv2.VideoCapture(video_path)
