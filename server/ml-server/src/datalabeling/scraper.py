# Python
import os
import re
import shutil
import threading
import argparse
import yt_dlp
import cv2
import grpc

# Internal
from src.common import constants
from src.common import util
from src.common import vidprocess
from src.common import cvclient


class Scraper():
    def __init__(self, cv_client):
        self.mutex = threading.Lock()
        #self.downloaded_videos = 0
        self.vid_processor = vidprocess.VideoProcessor(cv_client)
        #face_on_indices = [] # 1s at indices with face on videos, 0 otherwise

    def finished_hook(self, d):
        with self.mutex:
            if d['status'] == "finished":
                # increment number of downloaded videos
                #self.downloaded_videos += 1
                print("Finished downloading video")

    def download_shorts_from_creator(self, channel_url):
        shorts_url = f"{channel_url}/shorts"
        ydl_opts = {
            # Output template to save files in the specified directory with a specific name format
            'outtmpl': os.path.join(constants.VIDEOS_DIR_PATH, "temp", f'%(title)s.%(ext)s'),
            # Ensure only a single file format is downloaded and merged (best video + best audio)
            'format': 'bestvideo+bestaudio/best',
            # Merge the formats into an mp4 container
            'merge_output_format': 'mp4',
            # Print debug info to stderr (optional, useful for troubleshooting)
            #'verbose': True,
            # Select N items from the playlist
            #'playlist_items': f'1:{40}', 
            'playlist_items': f'1:{2}',
            # Add progress hook
            'progress_hooks': [self.finished_hook],
            #'postprocessor_hooks': [self.finished_hook],
        }
        try:
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                # The URL points to the shorts section of the channel
                ydl.download([shorts_url])
            print(f"Successfully downloaded shorts from {channel_url}")
        except Exception as e:
            print(f"An error occurred: {e}")
        self.postprocess_videos()
        #return self.downloaded_videos//2
        #clean_video_names()

    def postprocess_videos(self):
        temp_folder = os.path.join(constants.VIDEOS_DIR_PATH, "temp")
        fourcc = cv2.VideoWriter_fourcc(*'mp4v')
        # downsample and save video
        for file in os.listdir(temp_folder):
            print(f"file is {os.path.join(temp_folder, file)}")
            vid = cv2.VideoCapture(os.path.join(temp_folder, file))
            width = int(vid.get(cv2.CAP_PROP_FRAME_WIDTH))
            height = int(vid.get(cv2.CAP_PROP_FRAME_HEIGHT))
            original_fps = vid.get(cv2.CAP_PROP_FPS)
            processed_video = self.vid_processor.get_body_datapoints_from_video(vid, constants.FRAME_DIMENSION)
            print(f"processed video is {processed_video.downsample_rate}")
            new_fps = original_fps/processed_video.downsample_rate
            output = cv2.VideoWriter(os.path.join(constants.VIDEOS_DIR_PATH, f"{file.title}.mp4"), fourcc, new_fps, (width, height))
            for frame in processed_video.frames:
                output.write(util.bytes_to_frame(frame))
            output.release()
            os.remove(os.path.join(temp_folder, file))
            # extract body datapoints and save

    #def download_faceon_and_dtl_videos(self):
        #download_shorts_from_creator("https://www.youtube.com/@EKGolf")
        #print(f"downloaded {downloaded_videos//2} videos!, faceon indices is {face_on_indices[0::2]}")
        #return downloaded_videos//2, face_on_indices[0::2] # divide each of these in half, because progress_hook fires twice for each video download

    #def clean_video_names(self):
        #for original_filename in os.listdir(constants.VIDEOS_DIR_PATH):
            # strip leading zeros in filename
            #parts = original_filename.split('+', 1)
            #parts = parts[1].split('.', 1)
            #num_str = parts[0]
            #extension = parts[1]
            # Remove leading zeros from the number part
            #cleaned_num = str(int(num_str)) # Converts '01' to '1'
            #new_filename = f"{cleaned_num}.{extension}"
            #old_path = os.path.join(constants.VIDEOS_DIR_PATH, original_filename)
            #new_path = os.path.join(constants.VIDEOS_DIR_PATH, new_filename)
            #shutil.move(old_path, new_path)

if __name__ == "__main__":
    # create arg parser
    parser = argparse.ArgumentParser(
        description="Script that allows scraping youtube shorts from creators"
    )
    parser.add_argument(
        "--creator_url",
        type=str,
        required=True,
    )
    args = parser.parse_args()
    # initialize cv client for body datapoints requests
    cv_client = cvclient.create_computer_vision_client()
    # download youtube videos from requested creator url
    scraper = Scraper(cv_client)
    num_videos = scraper.download_shorts_from_creator(args.creator_url)
    print(f"num videos is {num_videos}")
