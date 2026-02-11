# Python
import cv2
import os
import argparse
import numpy as np
import uuid

# Internal
from src.common import util
from src.common import constants
from src.common import plot
from src.common import labeldata

# TODO: Make this labeling API accessible (save full video, and do the downsampling/extracting bodydata/labeling all at once??)
# TODO: Store labeled data in postgresql, for now just store as json file
class Labeler():
    def __init__(self):
        pass

    def label_unlabeled_videos(self):
        try: 
            # pull in data structures
            all_body_datapoints = util.load_all_body_datapoints_from_file()
            unlabeled_dict = util.load_unlabeled_dict()
            labeled_golf_swing_dataset = util.load_labeled_golf_swing_dataset()
            discarded_videos = []

            # plot datapoints for every 50 swing videos (every 10 frames)
            for i in range(0, len(all_body_datapoints), 50):
                for j in range(0, constants.FRAME_DIMENSION, 10):
                    print(f"BLAH: {all_body_datapoints[16]}")
                    plot.plot_datapoints_as_skeleton(all_body_datapoints[i][j])

            return
            
            # iterate over all unlabeled videos and label
            items = list(unlabeled_dict.items())
            for video_name, idx in items:
                self.label_video(video_name, idx, unlabeled_dict, labeled_golf_swing_dataset, discarded_videos)
            # save data structures
            self.save_data_structures(unlabeled_dict, all_body_datapoints, labeled_golf_swing_dataset, discarded_videos)
        except KeyboardInterrupt:
            #self.save_data_structures(unlabeled_dict, all_body_datapoints, labeled_golf_swing_dataset, discarded_videos)
            return

    def save_data_structures(self, unlabeled_dict, all_body_datapoints, labeled_golf_swing_dataset, discarded_videos):
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
        # shift labeled golf swing dataset indices
        new_videos_dict = {}
        new_idx = 0
        for video_name, video in sorted(labeled_golf_swing_dataset.videos.items(), key=lambda x: x[1]): # TODO: fix lambda
            old_idx = video.body_datapoints_idx
            if old_idx not in discarded_videos:
                video.body_datapoints_idx = new_idx
                new_videos_dict[video_name] = video
                new_idx += 1
        labeled_golf_swing_dataset.videos = new_videos_dict
        # save structs
        util.save_all_body_datapoints(all_body_datapoints)
        util.save_unlabeled_dict(unlabeled_dict)
        util.save_labeled_golf_swing_dataset(labeled_golf_swing_dataset)

    def label_video(self, video_name, idx, unlabeled_dict, labeled_golf_swing_dataset, discarded_videos):
        # grab actual video with video name
        video_path = util.get_saved_video_path(video_name)
        vid = util.grab_saved_video(video_path)
        # label process: check if video is ok, then prompt user to label faceon, isswing, etc.
        if self.verify_good_video(video_path, idx, vid):
            player_id = self.label_player_data(video_path, labeled_golf_swing_dataset.players)
            labeled_video = self.label_video_data(video_path, labeled_golf_swing_dataset.videos, player_id)
            self.label_frame_data(vid, labeled_video)
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

    def label_player_data(self, video_path, players_dict) -> str:
        os.startfile(video_path)
        name = self.prompt_user_for_information("Enter name of player: ", str)
        if name not in players_dict:
            players_dict[name] = labeldata.LabeledPlayer(
                id = str(uuid.uuid4()),
                level = self.prompt_user_for_information("What is the player's golf level? (...): ", str),
                handicap = self.prompt_user_for_information("What is the player's handicap? ", float),
                height = self.prompt_user_for_information("What is the player's height (cm)? ", float),
                weight = self.prompt_user_for_information("What is the player's weight (kg)? ", float),
                sex = self.prompt_user_for_information("What is the player's sex? (...) ", str),
                golf_handedness = self.prompt_user_for_information("What is the player's golf handedness? (...) ", str),
                real_handedness = self.prompt_user_for_information("What is the player's real handedness? (...) ", str),
                description = self.prompt_user_for_information("Additional information/description about the player. ", str),
            )
        return players_dict[name].id
        

    def label_video_data(self, video_path, videos_dict, player_id) -> labeldata.LabeledVideo:
        os.startfile(video_path)
        name = self.prompt_user_for_information("What is the name of the video: ", str)
        labeled_video = labeldata.LabeledVideo(
            player_id = player_id,
            feel = self.prompt_user_for_information("What was the player feeling in this swing? ", str),
            intended_shot = self.prompt_user_for_information("What was the intended shot? ", str),
            result = self.prompt_user_for_information("What was the result of the shot? ", str),
            shot_type = self.prompt_user_for_information("What was the shot type? (...) ", str),
            shot_shape = self.prompt_user_for_information("What was the shot shape? (...) ", str),
            orientation = self.prompt_user_for_information("What is the orientation? (...) ", str),
            club = self.prompt_user_for_information("What club was used? (...) ", str),
            ball_data = labeldata.BallData(
                ball_speed = self.prompt_user_for_information("What was the ball speed?", float),
                spin_rate = self.prompt_user_for_information("What was the spin rate?", float),
                apex_height = self.prompt_user_for_information("What was the apex height?", float),
                launch_angle = self.prompt_user_for_information("What was the launch angle?", float),
                spin_axis = self.prompt_user_for_information("What was the spin axis?", float),
                smash = self.prompt_user_for_information("What was the smash?", float),
                curve = self.prompt_user_for_information("What was the curve?", float),
                carry = self.prompt_user_for_information("What was the carry distance?", float),
                total_distance = self.prompt_user_for_information("What was the total distance?", float),
            ),
            club_data = labeldata.ClubData(
                club_speed = self.prompt_user_for_information("What was the club speed?", float),
                angle_of_attack = self.prompt_user_for_information("What was the angle of attack?", float),
                swing_direction = self.prompt_user_for_information("What was the swing direction?", float),
                face_to_path = self.prompt_user_for_information("What was the face to path?", float),
                dynamic_loft = self.prompt_user_for_information("What was the dynamic loft?", float),
            ),
            description = self.prompt_user_for_information("Additional information/description about the video. ", str),
        )
        videos_dict[name] = labeled_video
        return labeled_video
    
    def label_frame_data(self, vid, labeled_video):
        # get frames from video
        frames = []
        while True:
            ok, image = vid.read()
            if ok:
                frames.append(image)
            else:
                break
        if len(frames) != constants.FRAME_DIMENSION:
            print(f"Could not extract {constants.FRAME_DIMENSION} frames from video, got {len(frames)}, trying again...")
            self.label_frame_data(vid)
        # get start, end, and swing phases frames
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
            



    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    


    def label_isswing(self, video, labeled_golf_swing_dataset):
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

    def prompt_user_for_information(self, prompt, expected_type):
        # check that input is correct type
        user_input = input(prompt)
        try:
            expected_type(user_input)
        except ValueError:
            print(f"Incorrect type")
            self.prompt_user_for_information(prompt, expected_type)
        # verify with user that entered info is correct
        is_ok = input(f"You entered \"{user_input}\". Is this correct? 'y' for Yes, 'n' for No")
        if is_ok.lower() == "y":
            return user_input
        else:
            self.prompt_user_for_information(prompt, expected_type)
            
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
