import sys
from pathlib import Path
import os
import re
import shutil
import threading

import yt_dlp

curr_dir = Path(__file__).parent.resolve()
videos_folder = os.path.join(curr_dir, "data", "videos")

#def search_for_faceon(info):
#    print(f"blah: {info['title']}")
#    res = re.search('FO|Face On', info['title'], re.IGNORECASE)
#    if res == None:
#        return "dont download"
#    return None

def search_for_faceon(info):
    res = re.search('FO|Face On', info['filename'], re.IGNORECASE)
    return res != None

mutex = threading.Lock()
downloaded_videos = 0
face_on_indices = [] # 1s at indices with face on videos, 0 otherwise

def finished_hook(d):
    global downloaded_videos
    global face_on_indices
    with mutex:
        if d['status'] == "finished":
            if search_for_faceon(d):
                face_on_indices.append([1]*150)
            else:
                face_on_indices.append([0]*150)
            # increment number of downloaded videos
            downloaded_videos += 1

def download_shorts_from_creator(channel_url):
    shorts_url = f"{channel_url}/shorts"
    ydl_opts = {
        # Filter for videos that are face on
        #'match_filter': search_for_faceon,
        # Output template to save files in the specified directory with a specific name format
        'outtmpl': os.path.join(videos_folder, f'%(title)s+%(autonumber)s.%(ext)s'),
        # Ensure only a single file format is downloaded and merged (best video + best audio)
        'format': 'bestvideo+bestaudio/best',
        # Merge the formats into an mp4 container
        'merge_output_format': 'mp4',
        # Print debug info to stderr (optional, useful for troubleshooting)
        #'verbose': True,
        # Select N items from the playlist
        'playlist_items': f'1:{40}', 
        # Add progress hook
        'progress_hooks': [finished_hook],
    }
    try:
        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            # The URL points to the shorts section of the channel
            ydl.download([shorts_url])
        print(f"Successfully downloaded shorts from {channel_url}")
    except Exception as e:
        print(f"An error occurred: {e}")
    clean_video_names()

def download_faceon_and_dtl_videos():
    download_shorts_from_creator("https://www.youtube.com/@EKGolf")
    print(f"downloaded {downloaded_videos//2} videos!, faceon indices is {face_on_indices[0::2]}")
    return downloaded_videos//2, face_on_indices[0::2] # divide each of these in half, because progress_hook fires twice for each video download

def clean_video_names():
    for original_filename in os.listdir(videos_folder):
        # strip leading zeros in filename
        parts = original_filename.split('+', 1)
        parts = parts[1].split('.', 1)
        num_str = parts[0]
        extension = parts[1]
        # Remove leading zeros from the number part
        cleaned_num = str(int(num_str)) # Converts '01' to '1'
        new_filename = f"{cleaned_num}.{extension}"
        old_path = os.path.join(videos_folder, original_filename)
        new_path = os.path.join(videos_folder, new_filename)
        shutil.move(old_path, new_path)
