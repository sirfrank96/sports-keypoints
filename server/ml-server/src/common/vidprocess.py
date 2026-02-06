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
    body_datapoints: NDArray[np.float32]
    original_num_frames: int
    downsampled_num_frames: int
    downsample_rate: float

class VideoProcessor():
    def __init__(self, cv_client):
        self.cv_client = cv_client

    def process_video(self, video, video_name, num_frames) -> ProcessedVideo:
        # set up vars for outputting downsampled video
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')
        width = int(video.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(video.get(cv2.CAP_PROP_FRAME_HEIGHT))
        original_fps = video.get(cv2.CAP_PROP_FPS)
        total_frames = int(video.get(cv2.CAP_PROP_FRAME_COUNT))
        downsample_rate = total_frames/num_frames # TODO: TRUNCATE?
        new_fps = original_fps/downsample_rate
        output = cv2.VideoWriter(os.path.join(constants.VIDEOS_DIR_PATH, f"{video_name}"), fourcc, new_fps, (width, height))
        # parse video into frames and downsample
        frames_idx = 0
        last_idx = -1
        frames = []
        while True:
            if last_idx >= num_frames - 1:
                break
            ok, image = video.read()
            if ok:
                idx_target = int(frames_idx/downsample_rate)
                frames_idx += 1
                if idx_target > last_idx:
                    output.write(image)
                    frames.append(util.frame_to_bytes(image))
                    last_idx = idx_target
            else:
                break
        video.release()
        output.release()
        print(f"video is {video_name} original number of frames {frames_idx}, from constant: {total_frames}, downsampled num frames {len(frames)}")
        # query computervision-service for body datapoints
        datapoints = None
        try:
            datapoints = self.cv_client.get_pose_data_from_video(frames, constants.FRAME_DIMENSION)
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # plot datapoints for the swing video (every 10 frames)
        #for i in range(0, constants.FRAME_DIMENSION, 10):
            #plot.plot_datapoints_as_skeleton(datapoints[i])
        return ProcessedVideo(
            body_datapoints=datapoints,
            original_num_frames=frames_idx,
            downsampled_num_frames=len(frames),
            downsample_rate=frames_idx/len(frames),
        )
