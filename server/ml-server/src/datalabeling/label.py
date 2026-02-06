# Python
import cv2
import os
import argparse
import numpy as np

# Internal
from src.common import util
from src.common import constants
from src.common import plot

class Labeler():
    def __init__(self):
        pass

    def label_unlabeled_videos(self):
        try: 
            # pull in data structures
            all_body_datapoints = util.load_all_body_datapoints_from_file()

            # plot datapoints for every 50 swing videos (every 10 frames)
            for i in range(0, len(all_body_datapoints), 50):
                for j in range(0, constants.FRAME_DIMENSION, 10):
                    plot.plot_datapoints_as_skeleton(all_body_datapoints[i][j])

            unlabeled_dict = util.load_unlabeled_dict()
            faceon_indices = util.load_faceon_indices_from_file()
            isswing_indices = util.load_isswing_indices_from_file()
            discarded_videos = []
            # iterate over all unlabeled videos and label
            items = list(unlabeled_dict.items())
            for video_name, idx in items:
                self.label_video(video_name, idx, unlabeled_dict, faceon_indices, isswing_indices, discarded_videos)
            # save data structures
            self.save_data_structures(unlabeled_dict, all_body_datapoints, faceon_indices, isswing_indices, discarded_videos)
        except KeyboardInterrupt:
            self.save_data_structures(unlabeled_dict, all_body_datapoints, faceon_indices, isswing_indices, discarded_videos)

    def save_data_structures(self, unlabeled_dict, all_body_datapoints, faceon_indices, isswing_indices, discarded_videos):
        # remove all discarded videos from body datapoints
        all_body_datapoints = np.delete(all_body_datapoints, discarded_videos, axis=0)
        # shift unlabeled dict indices
        new_unlabeled_dict = {}
        new_idx = 0
        for video_name, old_idx in sorted(unlabeled_dict.items(), key=lambda x: x[1]):
            if old_idx not in discarded_videos:
                new_unlabeled_dict[video_name] = new_idx
                new_idx += 1
        unlabeled_dict = new_unlabeled_dict
        # save structs
        util.save_all_body_datapoints(all_body_datapoints)
        util.save_unlabeled_dict(unlabeled_dict)
        util.save_faceon_indices(faceon_indices)
        util.save_isswing_indices(isswing_indices)

    def label_video(self, video_name, idx, unlabeled_dict, faceon_indices, isswing_indices, discarded_videos):
        # grab actual video with video name
        video_path = util.get_saved_video_path(video_name)
        vid = util.grab_saved_video(video_path)
        # label process: check if video is ok, then prompt user to label faceon, isswing, etc.
        if self.verify_good_video(video_path, idx, vid):
            self.label_isfaceon(vid, faceon_indices)
            self.label_isswing(vid, isswing_indices)
        else:
            video_path = os.path.join(constants.VIDEOS_DIR_PATH, video_name)
            os.remove(video_path)
            discarded_videos = discarded_videos.append(idx)
        vid.release()
        # remove from unlabeled dict
        util.remove_unlabeled_video(video_name, unlabeled_dict)

    def verify_good_video(self, video_path, idx) -> bool:
        # display video
        os.startfile(video_path)
        # prompt user to look at video and yes to keep, no to remove 
        user_input = input("Is this video good for the dataset? Type 'y' to continue, 'n' to remove. (ie. clearly faceon or dtl, not multiple people, etc.)")
        # remove if not faceon or dtl or videos with multiple swings/people
        if user_input.lower() == "n":
            return False
        elif user_input.lower() == "y":
            return True
        else:
            print("please enter a valid answer ('y' or 'n')")
            self.verify_good_video(video_path, idx)

    def label_isfaceon(self, video, faceon_indices):
        # prompt user for whether video is dtl or faceon 
        user_input = input("Is this video face on? Type 'y' for face on, 'n' for dtl.")
        # update faceon_indices with array of all zeroes or ones of size number of frames
        if user_input.lower() == "y":
            faceon_indices = np.concatenate([faceon_indices, [np.ones(constants.FRAME_DIMENSION)]], axis=0)
        elif user_input.lower() == "n":
            faceon_indices = np.concatenate([faceon_indices, [np.zeros(constants.FRAME_DIMENSION)]], axis=0)
        else:
            print("please enter a valid answer ('y' or 'n')")
            self.label_isfaceon(video)

    def label_isswing(self, video, isswing_indices):
        frames = []
        while True:
            ok, image = video.read()
            if ok:
                frames.append(image)
            else:
                break
        # prompt user to identify start and end frames for the video
        start_frames = [constants.FRAME_DIMENSION]
        end_frames = [0]
        frame_idx = 0
        while True:
            # show each frame
            resized_img = cv2.resize(frames[frame_idx], (400, 800))
            cv2.imshow(f"current frame: {frame_idx+1}", resized_img)
            cv2.waitKey(0)
            # prompt user for input about frame
            frame_input = input("If this is the start frame, press 's'. If this is the end frame, press 'e'. Type a number to go to a specific frame. Or us 'a'/'d' (wasd) to navigate. Type 'q' to exit. ")
            if frame_input.lower() == "q":
                break
            elif frame_input.lower() == "s":
                start_frames.append(frame_idx)
            elif frame_input.lower() == "e":
                end_frames.append(frame_idx)
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
        # confirm with user that start and end frames look good
        user_input = input(f"Are these start frames: {start_frames} and end frames: {end_frames} correct? Type 'y' to continue, 'n' if not.")
        if user_input.lower() == "y":
            # if we don't have the same number of start and end frames, restart
            if len(start_frames) != len(end_frames):
                print(f"number of start frames {len(start_frames)} should equal number of end frames {len(end_frames)}")
                self.label_isswing(video)
            # build array of 1s and 0s, where 1s is "during a swing", 0 is not
            currswing_indices = np.zeros(constants.FRAME_DIMENSION)
            for i in range(0, len(start_frames)):
                start_frame, end_frame = start_frames[i], end_frames[i]
                # if start frame is bigger than end frame, restart
                if start_frame > end_frame:
                    print("start frame is greater than end frame, please retry")
                    self.label_isswing(video)
                else:
                    currswing_indices[start_frame:end_frame] = 1
            # append this videos isswing indices to the list of all videos isswing indices
            isswing_indices = np.concatenate([isswing_indices, [currswing_indices]], axis=0)
        elif user_input.lower() == "n":
            self.label_isswing(video)
        else:
            print("please enter a valid answer ('y' or 'n')")
            self.label_isswing(video)
            
if __name__ == "__main__":
    # create arg parser
    parser = argparse.ArgumentParser(
        description="Script that allows labelling stored videos"
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--unlabeled",
        action="store_true",
        required=False,
    )
    args = parser.parse_args()
    # label based on arg
    labeler = Labeler()
    if args.unlabeled:
        labeler.label_unlabeled_videos()
