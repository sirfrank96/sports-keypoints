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
import keras_models
import plot
import scraper
import util
import preprocess

import cv2 

def serve():
    setting_timeout_ms = 1000 * 60 * 3
    options = [('grpc.http2.settings_timeout', setting_timeout_ms)]
    channel = grpc.insecure_channel('localhost:50051', options=options)
    return cv.ComputerVisionClient(channel)

def get_body_datapoints_from_videos(num_downloaded_videos):
    # pull in videos
    all_datapoints = np.zeros((num_downloaded_videos, util.frames_dim_for_pose_estimation, 25*3))
    for i in range(0, num_downloaded_videos):
        vid = cv2.VideoCapture(os.path.join(curr_dir, "data", "videos", f'{i+1}.mp4'))
        # parse video into frames
        num_frames = 0
        frames = []
        while True:
            ok, image = vid.read()
            if ok:
                num_frames += 1
                frames.append(util.frame_to_bytes(image))
            else:
                break
        print(f'number of frames {num_frames}')
        vid.release()
        frames, _ = util.downsample_frames(frames, util.frames_dim_for_pose_estimation)
        print(f'number of frames after downsample {len(frames)}')
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

# TODO: Randomize train and validation data
def auto_encode_codeflow(num_downloaded_videos):
    all_datapoints = get_body_datapoints_from_videos(num_downloaded_videos)
    # plot datapoints
    for i in range(50, 100):
        plot.plot_datapoints_as_skeleton(all_datapoints[0][i])
    # put as input to autoencoder
    inputs, outputs = util.standardize_inputs_and_outputs(all_datapoints, all_datapoints)
    x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
    keras_models.autoencode(x_train, y_train, x_val, y_val) 

def encode_and_classify_codeflow(num_downloaded_videos, faceon_indices):
    all_datapoints = get_body_datapoints_from_videos(num_downloaded_videos)
    inputs, outputs = util.standardize_inputs_and_outputs(all_datapoints, faceon_indices)
    x_train, y_train, x_val, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
    keras_models.encode_and_classify(x_train, y_train, x_val, y_val, "models/classify_faceon_model_batch16_epochs50_early_stop.keras")

def output_new_faceon_model_from_existing_data():
    datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    all_datapoints = np.load(datapoints_file_path)
    print(f"size if all_datapoints is {all_datapoints.shape}")
    faceon_indices_file_path = os.path.join(curr_dir, "data", "faceon_classifier", 'faceon_indices.npy')
    faceon_indices = np.load(faceon_indices_file_path)
    print(f"size if faceon indices is {faceon_indices.shape}")
    inputs, outputs = util.standardize_inputs_and_outputs(all_datapoints, faceon_indices)
    x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
    keras_models.encode_and_classify(x_train, y_train, x_val, y_val, "models/classify_faceon_model_batch16_epochs50_early_stop.keras")

def output_new__isswing_model_from_existing_data():
    datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    all_datapoints = np.load(datapoints_file_path)
    print(f"size if all_datapoints is {all_datapoints.shape}")
    isswing_indices_file_path = os.path.join(curr_dir, "data", "isswing_classifier", 'is_swing_indices.npy')
    isswing_indices = np.load(isswing_indices_file_path)
    print(f"size if faceon indices is {isswing_indices.shape}, isswing_indices for 15th swing {isswing_indices[14]}")
    x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(all_datapoints, isswing_indices)
    keras_models.encode_and_classify(x_train, y_train, x_val, y_val, "models/classify_swing_model_openposedata_batch16_epochs50_early_stop.keras")

def train_classify_swing_model():
    # build x and y datasets
    #datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    #all_datapoints = np.load(datapoints_file_path)
    #print(f"size if all_datapoints is {all_datapoints.shape}")
    map_file = os.path.join(curr_dir, "data", "isswing_classifier", "video_start_end_frames_map.json")
    video_start_end_frames_map = preprocess.load_start_end_frames_map(map_file)
    all_datapoints = []
    vals = []
    for vid_id, start_end in video_start_end_frames_map.items():
        vid = cv2.VideoCapture(os.path.join(curr_dir, "data", "videos", f'{vid_id}.mp4'))
        # parse video into frames
        frame_idx = 0
        is_swing_arr = []
        frames = []
        while True:
            ok, image = vid.read()
            if ok:
                frames.append(util.frame_to_bytes(image))
                if frame_idx >= start_end["s"] and frame_idx <= start_end["e"]:
                    is_swing_arr.append(1)
                else:
                    is_swing_arr.append(0)
            else:
                break
            frame_idx += 1
        print(f'video id {vid_id} number of frames {frame_idx}')
        vid.release()
        frames, is_swing_arr = util.downsample_frames(frames, util.frames_dim_for_pose_estimation, is_swing_arr)
        print(f'number of frames after downsample {len(frames)}, length of is_swing_arr is {len(is_swing_arr)}')
        # query computervision-service for body datapoints
        datapoints_swing_i = None
        try:
            datapoints_swing_i = cv_client.get_pose_data_from_video(frames)
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # put bodydatapoints into tensors (num swings x num frames x num datapoints x datapoint x,y,conf) (add extra dimension for more num people)
        all_datapoints.append(datapoints_swing_i)
        vals.append(is_swing_arr)
    inputs, outputs = util.standardize_inputs_and_outputs(all_datapoints, vals)
    print(f"after pad size of videos arr is {len(inputs)}x{len(inputs[1])}x{len(inputs[1][0])}, size of vals arr is {len(outputs)}x{len(outputs[1])}")
    # save datapoints to data/bodydatapoints folder
    datapoints_file_path = os.path.join(curr_dir, "data", "bodydatapoints", 'body_datapoints.npy')
    np.save(datapoints_file_path, inputs)
    is_swing_file_path = os.path.join(curr_dir, "data", "isswing_classifier", "is_swing_indices.npy")
    np.save(is_swing_file_path, outputs)
    x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
    keras_models.encode_and_classify(x_train, y_train, x_val, y_val, "models/classify_swing_model_openposedata_batch16_epochs50_early_stop.keras")

if __name__ == "__main__":
    logging.basicConfig()
    cv_client = serve()

    #num_downloaded_videos, faceon_indices = scraper.download_faceon_and_dtl_videos()
    #auto_encode_codeflow(num_downloaded_videos)
    #faceon_indices_file_path = os.path.join(curr_dir, "data", "faceon_classifier", 'faceon_indices.npy')
    #np.save(faceon_indices_file_path, faceon_indices)
    #encode_and_classify_codeflow(num_downloaded_videos, faceon_indices)
    #output_new_isswing_model_from_existing_data()

    #preprocess.label_data_to_classify_swing()
    #train_classify_swing_model()
    output_new__isswing_model_from_existing_data()
    
