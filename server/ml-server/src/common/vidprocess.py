# Python
import os
import cv2
import numpy as np
import grpc
from dataclasses import dataclass
from numpy.typing import NDArray
import subprocess
import threading

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

    # TODO: Checking that number of frames is correct
    def process_video(self, video, video_name) -> ProcessedVideo:
        # parse video into frames and downsample
        frames_idx = 0
        frames = []
        while True:
            ok, image = video.read()
            if ok:
                frames.append(util.frame_to_bytes(image))
                frames_idx += 1
            else:
                break
        video.release()
        print(f"video is {video_name} number of frames {frames_idx}")
        # query computervision-service for body datapoints
        datapoints = None
        try:
            datapoints = self.cv_client.get_pose_data_from_video(frames, constants.FRAME_DIMENSION)
            print(f"datapoints are {datapoints}")
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        return ProcessedVideo(
            body_datapoints=datapoints,
            original_num_frames=frames_idx,
            downsampled_num_frames=len(frames),
            downsample_rate=frames_idx/len(frames),
        )

    def process_video_and_output_new_video(self, video, video_name, num_frames) -> ProcessedVideo:
        # set up vars for outputting downsampled video
        width = int(video.get(cv2.CAP_PROP_FRAME_WIDTH))
        height = int(video.get(cv2.CAP_PROP_FRAME_HEIGHT))
        width = width - (width % 2) # width and height should be multiples of 2
        height = height - (height % 2)
        original_fps = video.get(cv2.CAP_PROP_FPS)
        total_frames = int(video.get(cv2.CAP_PROP_FRAME_COUNT))
        downsample_rate = total_frames/num_frames
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
            '-pix_fmt', 'yuv420p',
            '-crf', '23',
            '-preset', 'medium',
            '-movflags', '+faststart',
            output_dir
        ], stdin=subprocess.PIPE, stderr=subprocess.PIPE, stdout=subprocess.PIPE)
        #], stdin=subprocess.PIPE, stderr=subprocess.DEVNULL, stdout=subprocess.DEVNULL)
        # Read stderr in background to prevent deadlock
        stderr_output = []
        def read_stderr():
            for line in process.stderr:
                stderr_output.append(line)
        
        stderr_thread = threading.Thread(target=read_stderr, daemon=True)
        stderr_thread.start()
        # parse video into frames and downsample
        frames_idx = 0
        last_idx = -1
        frames = []
        frames_written_to_ffmpeg = 0
        try: 
            while True:
                if last_idx >= num_frames - 1:
                    break
                ok, image = video.read()
                if ok:
                    idx_target = int(frames_idx/downsample_rate)
                    frames_idx += 1
                    if idx_target > last_idx:
                        raw_bytes = image.tobytes()
                        expected_size = width * height * 3
                        if len(raw_bytes) != expected_size:
                            raise ValueError(f"Frame size mismatch: expected {expected_size}, got {len(raw_bytes)}")
                        process.stdin.write(raw_bytes)
                        frames_written_to_ffmpeg += 1
                        #process.stdin.flush()
                        frames.append(util.frame_to_bytes(image))
                        last_idx = idx_target
                else:
                    break
            print(f"Finished reading. Frames written to ffmpeg: {frames_written_to_ffmpeg}")
            video.release()
            process.stdin.close()
            #process.wait()
            process.wait(timeout=60)
            stderr_thread.join(timeout=5)
            print(f"FFmpeg finished with return code: {process.returncode}")
            stderr_text = b''.join(stderr_output).decode('utf-8', errors='ignore')
            print(f"FFmpeg stderr:\n{stderr_text}")
            if process.returncode != 0:
                raise RuntimeError(f"FFmpeg failed with return code {process.returncode}")
            # Verify output file exists and has size
            if os.path.exists(output_dir):
                file_size = os.path.getsize(output_dir)
                print(f"Output file size: {file_size} bytes")
                if file_size < 1000:
                    raise RuntimeError(f"Output file too small ({file_size} bytes), likely corrupted")
            else:
                raise RuntimeError(f"Output file not created: {output_dir}")
        except Exception as e:
            print(f"Exception: {e}")
            process.kill()
            raise
        print(f"video is {video_name} original number of frames {frames_idx}, from constant: {total_frames}, downsampled num frames {len(frames)}")
        # query computervision-service for body datapoints
        datapoints = None
        try:
            datapoints = self.cv_client.get_pose_data_from_video(frames, constants.FRAME_DIMENSION)
            print(f"datapoints are {datapoints}")
        except grpc.RpcError as e:
            print(f'grpc error {e}')
        return ProcessedVideo(
            body_datapoints=datapoints,
            original_num_frames=frames_idx,
            downsampled_num_frames=len(frames),
            downsample_rate=frames_idx/len(frames),
        )
