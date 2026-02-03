# Python
import os
import cv2
import numpy as np
import grpc
from dataclasses import dataclass
from numpy.typing import NDArray

# Internal
from . import util
from . import constants
from . import plot

@dataclass
class ProcessedVideo:
    frames: list[bytes]
    body_datapoints: NDArray[np.float32]
    original_num_frames: int
    downsampled_num_frames: int
    downsample_rate: float

class VideoProcessor():
    def __init__(self, cv_client):
        self.cv_client = cv_client

    def get_body_datapoints_from_video(self, video, num_frames) -> ProcessedVideo:
        # parse video into frames
        frames_idx = 0
        frames = []
        while True:
            ok, image = video.read()
            if ok:
                frames_idx += 1
                frames.append(util.frame_to_bytes(image))
            else:
                break
        video.release()
        # downsample to consistent number of frames
        print(f'number of frames in the video {frames_idx}')
        frames, _ = util.downsample_frames(frames, num_frames)
        print(f'number of frames after downsample {len(frames)}')
        # query computervision-service for body datapoints
        datapoints = None
        try:
            datapoints = self.cv_client.get_pose_data_from_video(frames, constants.FRAME_DIMENSION)
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # plot datapoints for the swing video (every 10 frames)
        for i in range(0, constants.FRAME_DIMENSION, 10):
            plot.plot_datapoints_as_skeleton(datapoints[i])
        return ProcessedVideo(
            frames=frames,
            body_datapoints=datapoints,
            original_num_frames=frames_idx,
            downsampled_num_frames=len(frames),
            downsample_rate=frames_idx/len(frames),
        )

    def get_body_datapoints_from_all_videos(self, num_videos, num_frames):
        # check if body_datapoints.npy already exists
        if os.path.exists(constants.BODY_DATAPOINTS_PATH):
            all_datapoints = np.load(constants.BODY_DATAPOINTS_PATH)
            return all_datapoints
        # pull in videos
        all_datapoints = np.zeros((num_videos, num_frames, constants.BODY_DATAPOINTS_DIMENSIONS))
        for i in range(0, num_videos):
            vid = cv2.VideoCapture(os.path.join(constants.VIDEOS_DIR_PATH, f'{i+1}.mp4'))
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
            vid.release()
            # downsample to consistent number of frames
            print(f'number of frames in the video {num_frames}')
            frames, _ = util.downsample_frames(frames, num_frames)
            print(f'number of frames after downsample {len(frames)}')
            # query computervision-service for body datapoints
            datapoints_swing_i = None
            try:
                datapoints_swing_i = self.cv_client.get_pose_data_from_video(frames)
            except grpc.RpcError as e:
                print(f'grpc error {e}')
            # put bodydatapoints into tensors (num swings x num frames x num datapoints x datapoint x,y,conf) (add extra dimension for more num people?)
            all_datapoints[i] = datapoints_swing_i
        # save datapoints to file
        np.save(constants.BODY_DATAPOINTS_PATH, all_datapoints)
        return all_datapoints

    def get_body_datapoints_and_isswing_indices_from_videos(self, num_videos, num_frames):
        # check if body_datapoints.npy and isswing_indices.npy already exist
        all_datapoints = []
        isswing_indices = []
        if os.path.exists(constants.BODY_DATAPOINTS_PATH):
            all_datapoints = np.load(constants.BODY_DATAPOINTS_PATH)
        if os.path.exists(constants.ISSWING_TRUE_VALUES_PATH):
            isswing_indices = np.load(constants.ISSWING_TRUE_VALUES_PATH)
        if len(all_datapoints) == 0 and len(isswing_indices) == 0:
            return all_datapoints, isswing_indices
        # pull in videos and map of start end frames to downsample and return
        video_start_end_frames_map = util.load_start_end_frames_map(constants.START_END_FRAMES_MAP_PATH)
        for vid_id, start_end in video_start_end_frames_map.items():
            vid = cv2.VideoCapture(os.path.join(constants.VIDEOS_DIR_PATH, f'{vid_id}.mp4'))
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
            frames, is_swing_arr = util.downsample_frames(frames, util.frames_dim_for_swing_classification, is_swing_arr)
            print(f'number of frames after downsample {len(frames)}, length of is_swing_arr is {len(is_swing_arr)}')
            # query computervision-service for body datapoints
            datapoints_swing_i = None
            try:
                datapoints_swing_i = self.cv_client.get_pose_data_from_video(frames, util.frames_dim_for_swing_classification)
                print(f"datapoints for swing {frame_idx} is {datapoints_swing_i}")
            except grpc.RpcError as e:
                print(f'grpc error {e}')
            # put bodydatapoints into tensors (num swings x num frames x num datapoints x datapoint x,y,conf) (add extra dimension for more num people)
            all_datapoints.append(datapoints_swing_i)
            # put is_swing_arr as np array into vals
            is_swing_np_arr = np.zeros(constants.FRAME_DIMENSION)
            is_swing_np_arr[:len(is_swing_arr)] = is_swing_arr
            isswing_indices.append(is_swing_np_arr)
        # save datapoints and indices to files
        np.save(constants.BODY_DATAPOINTS_PATH, all_datapoints)
        np.save(constants.ISSWING_TRUE_VALUES_PATH, isswing_indices)
        return all_datapoints, isswing_indices
