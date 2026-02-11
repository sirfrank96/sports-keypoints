# Python
import os
import cv2
import numpy as np
import grpc
from dataclasses import dataclass
from numpy.typing import NDArray
import subprocess

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
        width = int(video.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(video.get(cv2.CAP_PROP_FRAME_HEIGHT))
        width = width - (width % 2)
        height = height - (height % 2)
        original_fps = video.get(cv2.CAP_PROP_FPS)
        total_frames = int(video.get(cv2.CAP_PROP_FRAME_COUNT))
        print(f"total frames {total_frames}, num_frames {num_frames}")
        downsample_rate = total_frames/num_frames # TODO: TRUNCATE?
        if downsample_rate < 1:
            downsample_rate = 1
        new_fps = original_fps/downsample_rate
        output_dir = os.path.join(constants.VIDEOS_DIR_PATH, f"{video_name}")
        process = subprocess.Popen([
            'ffmpeg', '-y',
            '-f', 'rawvideo',
            '-vcodec', 'rawvideo',
            '-s', f'{width}x{height}',
            '-pix_fmt', 'bgr24',
            '-r', str(new_fps),
            '-i', '-',
            '-c:v', 'libx264',
            '-crf', '23',
            '-preset', 'medium',
            '-movflags', '+faststart',
            output_dir
        #], stdin=subprocess.PIPE, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        ], stdin=subprocess.PIPE, stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL)
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
                    process.stdin.write(image.tobytes())
                    #process.stdin.flush()
                    frames.append(util.frame_to_bytes(image))
                    if last_idx == 0:
                        print(f"{util.frame_to_bytes(image)[10:15]}")
                    last_idx = idx_target
            else:
                break
        video.release()
        process.stdin.close()
        process.wait()

        # Wait for ffmpeg with timeout
        #try:
            #stdout, stderr = process.communicate(timeout=30)
            #print(f"FFmpeg return code: {process.returncode}")
            #print(f"FFmpeg stderr:\n{stderr.decode()}")
            #if process.returncode != 0:
                #raise RuntimeError(f"FFmpeg failed with code {process.returncode}")
        #except subprocess.TimeoutExpired:
            #print("FFmpeg timed out!")
            #process.kill()
            #stdout, stderr = process.communicate()
            #print(f"FFmpeg stderr:\n{stderr.decode()}")
            #raise RuntimeError("FFmpeg timed out")

        print(f"video is {video_name} original number of frames {frames_idx}, from constant: {total_frames}, downsampled num frames {len(frames)}")
        # query computervision-service for body datapoints
        datapoints = None
        try:
            datapoints = self.cv_client.get_pose_data_from_video(frames, constants.FRAME_DIMENSION)
            print(f"datapoints are {datapoints}")
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        # plot datapoints for the swing video (every 10 frames)
        for i in range(0, constants.FRAME_DIMENSION, 10):
            plot.plot_datapoints_as_skeleton(datapoints[i])
        return ProcessedVideo(
            body_datapoints=datapoints,
            original_num_frames=frames_idx,
            downsampled_num_frames=len(frames),
            downsample_rate=frames_idx/len(frames),
        )
