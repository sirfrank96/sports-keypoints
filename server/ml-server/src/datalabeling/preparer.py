# Python
import os
import numpy as np
import argparse

# Internal
from src.common import util
from src.common import constants
from src.common import vidprocess
from src.common import cvclient

class Preparer():
    def __init__(self, cv_client):
        self.vid_processor = vidprocess.VideoProcessor(cv_client)

    def process_videos_in_folder(self, folder_path, all_body_datapoints, unlabeled_dict, start_idx, delete_original=False):
        additional_body_datapoints = []
        # iterate over all files in folder
        for file in os.listdir(folder_path):
            if file == "temp":
                continue
            # grab video from temp folder and get information about it
            vid = util.grab_saved_video(os.path.join(folder_path, file))
            if delete_original:
                # process video: downsample and get body datapoints
                processed_video = self.vid_processor.process_video_and_output_new_video(vid, file, constants.FRAME_DIMENSION)
                # delete temp video after processing
                os.remove(os.path.join(folder_path, file))
            else:
                processed_video = self.vid_processor.process_video(vid, file)
            # append to additional_body_datapoints
            additional_body_datapoints.append(processed_video.body_datapoints)
            # add to unlabeled videos 
            util.add_unlabeled_video(file, start_idx, unlabeled_dict)
            start_idx += 1
        if len(additional_body_datapoints) > 0:
            additional_body_datapoints_np_arr = np.array(additional_body_datapoints)
            all_body_datapoints = np.concatenate((all_body_datapoints, additional_body_datapoints_np_arr), axis=0)
        return all_body_datapoints, unlabeled_dict

    def prepare_new_videos_for_labelling(self):
        # load data structures to store processed information
        unlabeled_dict = util.load_unlabeled_dict()
        all_body_datapoints = util.load_all_body_datapoints_from_file()
        print(f"Before start preparer: len of unlabeled dict is {len(unlabeled_dict)}, len of all_body_datapoints is {len(all_body_datapoints)}, len of videos dir {len(os.listdir(constants.VIDEOS_DIR_PATH))}, len of temp dir is {len(os.listdir(os.path.join(constants.VIDEOS_DIR_PATH, "temp")))}")
        # process all videos in temp folder
        start_idx = len(all_body_datapoints)
        temp_folder = os.path.join(constants.VIDEOS_DIR_PATH, "temp")
        all_body_datapoints, unlabeled_dict = self.process_videos_in_folder(temp_folder, all_body_datapoints, unlabeled_dict, start_idx, True)
        # save body datapoints and unlabeled videos to file
        print(f"After preparer: len of unlabeled dict is {len(unlabeled_dict)}, len of all_body_datapoints is {len(all_body_datapoints)}, len of videos dir {len(os.listdir(constants.VIDEOS_DIR_PATH))}, len of temp dir is {len(os.listdir(os.path.join(constants.VIDEOS_DIR_PATH, "temp")))}")
        util.save_all_body_datapoints(all_body_datapoints)
        util.save_unlabeled_dict(unlabeled_dict)

    def prepare_all_videos_for_labelling(self):
        # create new empty data structures for preparing all videos
        unlabeled_dict = util.create_new_unlabeled_dict()
        all_body_datapoints = util.create_new_all_body_datapoints()
        print(f"Before start preparer: len of unlabeled dict is {len(unlabeled_dict)}, len of all_body_datapoints is {len(all_body_datapoints)}, len of videos dir {len(os.listdir(constants.VIDEOS_DIR_PATH))}, len of temp dir is {len(os.listdir(os.path.join(constants.VIDEOS_DIR_PATH, "temp")))}")
        # reprocess all videos in videos folder
        all_body_datapoints, unlabeled_dict = self.process_videos_in_folder(constants.VIDEOS_DIR_PATH, all_body_datapoints, unlabeled_dict, 0, False)
        # process all videos in temp folder
        start_idx = len(all_body_datapoints)
        temp_folder = os.path.join(constants.VIDEOS_DIR_PATH, "temp")
        all_body_datapoints, unlabeled_dict = self.process_videos_in_folder(temp_folder, all_body_datapoints, unlabeled_dict, start_idx, True)
        # save body datapoints and unlabeled videos to file
        print(f"After preparer: len of unlabeled dict is {len(unlabeled_dict)}, len of all_body_datapoints is {len(all_body_datapoints)}, len of videos dir {len(os.listdir(constants.VIDEOS_DIR_PATH))}, len of temp dir is {len(os.listdir(os.path.join(constants.VIDEOS_DIR_PATH, "temp")))}")
        util.save_all_body_datapoints(all_body_datapoints)
        util.save_unlabeled_dict(unlabeled_dict)

if __name__ == "__main__":
    # create arg parser
    parser = argparse.ArgumentParser(
        description="Script that will prepare videos for labeling (downsample, get body datapoints, etc.)"
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--new_videos",
        action="store_true",
        required=False,
    )
    group.add_argument(
        "--all_videos",
        action="store_true",
        required=False,
    )
    args = parser.parse_args()
    # initialize cv client for body datapoints requests
    cv_client = cvclient.create_computer_vision_client()
    preparer = Preparer(cv_client)
    # pull in data structures
    if args.new_videos:
        preparer.prepare_new_videos_for_labelling()
    elif args.all_videos:
        preparer.prepare_all_videos_for_labelling()
