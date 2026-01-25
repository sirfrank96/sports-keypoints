# Python
from pathlib import Path
import os
import sys
import json
import cv2
import numpy as np

# Internal
import keras_models
import util

# TODO: Split into Training and Inference modules with shared util functions

def preprocess_data():
    # pull start of swing classifier
    # run through and pull out frames that are needed
    # normalize number of frames

    # send each frame to get body_keypoints analyzed

    # pull faceon vs not face on classifier
    # divide dataset into face on and not face on

    # pull dtl vs not classifier
    # remove non dtls
    return

curr_dir = Path(__file__).parent.resolve()

def load_start_end_frames_map(map_file):
    # load map if saved previously
    video_start_end_frames_map = {}
    try:
        with open(map_file, 'r') as f:
            video_start_end_frames_map = json.load(f)
    except FileNotFoundError:
        print("No map, will create one")
    return video_start_end_frames_map

# TODO: Skip 35th video (rory at night hitting a bunch of shots), and 36th (normal and slomo video in 1)
def label_data_to_classify_swing():
    # load map if saved previously
    map_file = os.path.join(curr_dir, "data", "frames", "video_start_end_frames_map.json")
    video_start_end_frames_map = load_start_end_frames_map(map_file)
    # manually label start frame and end frames for 40 videos
    while True:
        # pull in video
        video_input = input("Which video do you want to identify? Type 'q' to exit. ")
        if video_input.lower() == "q":
            break
        video_num = 0
        try:
            video_num = int(video_input)
            if video_num < 1 or video_num > 40: # TODO: pass in num videos available?
                print("please enter a valid video number")
                video_num = 0
                continue
        except ValueError:
            print("please enter a valid number")
            continue
        vid = cv2.VideoCapture(os.path.join(curr_dir, "data", "videos", f'{video_num}.mp4'))
        frames = []
        while True:
            ok, image = vid.read()
            if ok:
                frames.append(image)
            else:
                break
        print(f'number of frames {len(frames)}')
        vid.release()
        frame_idx = 0
        # prompt user to identify start and end frames for the video
        while True:
            # show image
            resized_img = cv2.resize(frames[frame_idx], (400, 800))
            cv2.imshow(f"current frame: {frame_idx+1}", resized_img)
            cv2.waitKey(0)
            frame_input = input("If this is the start frame, press 's'. If this is the end frame, press 'e'. Type a number to go to a specific frame. Or us 'a'/'d' (wasd) to navigate. Type 'q' to exit. ")
            if frame_input.lower() == "q":
                break
            # TODO: Make a class/struct instead of "s" and "e"
            elif frame_input.lower() == "s":
                video_start_end_frames_map.setdefault(video_num, {})["s"] = frame_idx
            elif frame_input.lower() == "e":
                video_start_end_frames_map.setdefault(video_num, {})["e"] = frame_idx
            elif frame_input.lower() == "a":
                if frame_idx < 0:
                    print("already at first frame, cannot go lower")
                    continue
                frame_idx -= 1
            elif frame_input.lower() == "d":
                if frame_idx >= len(frames):
                    print("already at last frame, cannot go higher")
                    continue
                frame_idx += 1
            else:
                try:
                    prev_frame_idx = frame_idx
                    frame_idx = int(frame_input) - 1
                    if frame_idx < 0 or frame_idx >= len(frames):
                        print("please enter a valid frame number")
                        frame_idx = prev_frame_idx
                        continue
                except ValueError:
                    print("please enter a valid number")
                    continue

    # store data back in json file
    with open(map_file, 'w') as f:
        json.dump(video_start_end_frames_map, f)
    return