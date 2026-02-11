# Python
import os
import threading
import argparse
import yt_dlp
import sys

# Internal
from src.common import constants

class Scraper():
    def __init__(self):
        self.mutex = threading.Lock()

    def finished_hook(self, d):
        with self.mutex:
            if d['status'] == "finished":
                print(f"Finished downloading video {d["filename"]}")

    def download_shorts_from_creator(self, channel_url, start_idx):
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
            # Select items from the playlist starting at start_idx
            'playlist_items': f'{start_idx}:{start_idx + 1}',
            # Add progress hook
            'progress_hooks': [self.finished_hook],
        }
        try:
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                # The URL points to the shorts section of the channel
                ydl.download([shorts_url])
            print(f"Successfully downloaded shorts from {channel_url}")
        except Exception as e:
            print(f"An error occurred: {e}")

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
    parser.add_argument(
        "--start_idx",
        type=int,
        required=False,
    )
    args = parser.parse_args()
    # download youtube videos from requested creator url
    scraper = Scraper()
    start_idx = 0
    if args.start_idx is not None:
        start_idx = args.start_idx
    num_videos = scraper.download_shorts_from_creator(args.creator_url, start_idx)
