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
            # iterate over all unlabeled videos and label
            items = list(unlabeled_dict.items())
            print(f"Before labeling: len all_body_datapoints is {len(all_body_datapoints)}, unlabeled_dict {unlabeled_dict}, len labeled_golf_dataset {len(labeled_golf_swing_dataset)}, discarded: {discarded_videos}")
            for video_name, idx in items:
                print(f"video name {video_name}, idx is {idx}")
                self.label_video(video_name, idx, unlabeled_dict, labeled_golf_swing_dataset, discarded_videos)
            print(f"After labeling: len all_body_datapoints is {len(all_body_datapoints)}, unlabeled_dict {unlabeled_dict}, len labeled_golf_dataset {len(labeled_golf_swing_dataset)}, discarded: {discarded_videos}")
            # save data structures
            self.save_data_structures(unlabeled_dict, all_body_datapoints, labeled_golf_swing_dataset, discarded_videos)
        except KeyboardInterrupt:
            self.save_data_structures(unlabeled_dict, all_body_datapoints, labeled_golf_swing_dataset, discarded_videos)
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
        for video_name, video in sorted(labeled_golf_swing_dataset.videos.items(), key=lambda x: x[1].body_datapoints_idx):
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
        if self.verify_good_video(video_path, idx):
            player_id = self.label_player_data(video_path, labeled_golf_swing_dataset.players)
            labeled_video = self.label_video_data(video_path, idx, labeled_golf_swing_dataset.videos, player_id)
            self.label_frame_data(vid, labeled_video)
        else:
            vid.release()
            video_path = os.path.join(constants.VIDEOS_DIR_PATH, video_name)
            os.remove(video_path)
            discarded_videos.append(idx)
        vid.release()
        # remove from unlabeled dict
        util.remove_unlabeled_video(video_name, unlabeled_dict)

    def verify_good_video(self, video_path, idx) -> bool:
        # display video
        os.startfile(video_path)
        # prompt user to look at video and yes to keep, no to remove 
        user_input = input("Is this video good for the dataset? Type 'y' to continue, 'n' to remove. (ie. clearly faceon or dtl, not multiple people, etc.) ")
        # remove if not faceon or dtl or videos with multiple swings/people
        if user_input.lower() == "n":
            return False
        elif user_input.lower() == "y":
            return True
        else:
            print("please enter a valid answer ('y' or 'n')")
            return self.verify_good_video(video_path, idx)

    def label_player_data(self, video_path, players_dict) -> str:
        os.startfile(video_path)
        name = self.prompt_user_for_information("Enter name of player: ", str)
        if name not in players_dict:
            players_dict[name] = labeldata.LabeledPlayer(
                id = str(uuid.uuid4()),
                golf_level = self.prompt_user_for_information("What is the player's golf level? (TOPPRO, PRO, TOPAM, GOODAM, OKAM, BADAM): ", labeldata.GolfLevel, enum_type=True),
                handicap = self.prompt_user_for_information("What is the player's handicap? ", float),
                height = self.prompt_user_for_information("What is the player's height (cm)? ", float),
                weight = self.prompt_user_for_information("What is the player's weight (kg)? ", float),
                sex = self.prompt_user_for_information("What is the player's sex? (NONBINARY, MALE, FEMALE) ", labeldata.Sex, enum_type=True),
                golf_handedness = self.prompt_user_for_information("What is the player's golf handedness? (LEFT, RIGHT, AMBIDEXTROUS) ", labeldata.Handedness, enum_type=True),
                real_handedness = self.prompt_user_for_information("What is the player's real handedness? (LEFT, RIGHT, AMBIDEXTROUS) ", labeldata.Handedness, enum_type=True),
                description = self.prompt_user_for_information("Additional information/description about the player. ", str),
            )
        return players_dict[name].id
        

    def label_video_data(self, video_path, idx, videos_dict, player_id) -> labeldata.LabeledVideo:
        os.startfile(video_path)
        name = self.prompt_user_for_information("What is the name of the video: ", str)
        labeled_video = labeldata.LabeledVideo(
            id = str(uuid.uuid4()),
            player_id = player_id,
            body_datapoints_idx = idx,
            feel_description = self.prompt_user_for_information("What was the player feeling in this swing? ", str),
            intended_shot = self.prompt_user_for_information("What was the intended shot? ", str),
            result = self.prompt_user_for_information("What was the result of the shot? ", str),
            orientation = self.prompt_user_for_information("What is the orientation? (FACEON, DTL, NEITHER) ", labeldata.Orientation, enum_type=True),
            club = self.prompt_user_for_information("What club was used? (DRIVER, WOOD, HYBRID, LONGIRON, MIDIRON, SHORTIRON, WEDGE, PUTTER) ", labeldata.Club, enum_type=True),
            ball_data = labeldata.BallData(
                ball_speed = self.prompt_user_for_information("What was the ball speed (mph)? ", float),
                spin_rate = self.prompt_user_for_information("What was the spin rate? ", float),
                apex_height = self.prompt_user_for_information("What was the apex height (ft.)? ", float),
                launch_angle = self.prompt_user_for_information("What was the launch angle? ", float),
                spin_axis = self.prompt_user_for_information("What was the spin axis? ", float),
                smash = self.prompt_user_for_information("What was the smash? ", float),
                curve = self.prompt_user_for_information("What was the curve (yd.)? ", float),
                carry_distance = self.prompt_user_for_information("What was the carry distance (yd.)? ", float),
                total_distance = self.prompt_user_for_information("What was the total distance (yd.)? ", float),
            ),
            club_data = labeldata.ClubData(
                club_speed = self.prompt_user_for_information("What was the club speed (mph)? ", float),
                angle_of_attack = self.prompt_user_for_information("What was the angle of attack? ", float),
                swing_direction = self.prompt_user_for_information("What was the swing direction? ", float),
                face_to_path = self.prompt_user_for_information("What was the face to path? ", float),
                dynamic_loft = self.prompt_user_for_information("What was the dynamic loft? ", float),
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
            return self.label_frame_data(vid)
        # get swing phases frames
        psystem_positions = labeldata.PSystemPositions()
        frame_idx = 0
        while True:
            # show each frame
            resized_img = cv2.resize(frames[frame_idx], (400, 800))
            cv2.imshow(f"current frame: {frame_idx+1}", resized_img)
            cv2.waitKey(0)
            # prompt user for input about frame
            frame_input = input("If this frame is one of the P-System positions (P1-P10), type which one (eg. 'P8'). If not, type a number to go to a specific frame. Or use 'a'/'d' (wasd) to navigate. Type 'q' to exit: ")
            if frame_input.lower() == "q":
                break
            elif frame_input.lower() == "p1":
                psystem_positions.P1 = frame_idx
            elif frame_input.lower() == "p2":
                psystem_positions.P2 = frame_idx
            elif frame_input.lower() == "p3":
                psystem_positions.P3 = frame_idx
            elif frame_input.lower() == "p4":
                psystem_positions.P4 = frame_idx
            elif frame_input.lower() == "p5":
                psystem_positions.P5 = frame_idx
            elif frame_input.lower() == "p6":
                psystem_positions.P6 = frame_idx
            elif frame_input.lower() == "p7":
                psystem_positions.P7 = frame_idx
            elif frame_input.lower() == "p8":
                psystem_positions.P8 = frame_idx
            elif frame_input.lower() == "p9":
                psystem_positions.P9 = frame_idx
            elif frame_input.lower() == "p10":
                psystem_positions.P10 = frame_idx
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
                    print("please enter a valid number of P-System position")
                    continue
        # confirm with user that start and end frames look good
        user_input = input(f"Are these P-System positions correct {psystem_positions}? Type 'y' to continue, 'n' if not: ")            
        if user_input.lower() != "y":
            return self.label_frame_data(vid, labeled_video)
        # convert p system position to swing phases 
        labeled_frames = labeled_video.labeled_frames
        p1 = psystem_positions.P1
        if p1 > 0:
            setup_frame = labeled_frames[p1]
            setup_frame.swing_phase = labeldata.SwingPhase.SETUP
        p2 = psystem_positions.P2
        if p1 > 0 and p2 > 0:
            for i in range(p1+1, p2+1):
                takeaway_frame = labeled_frames[i]
                takeaway_frame.swing_phase = labeldata.SwingPhase.TAKEAWAY
        p3 = psystem_positions.P3
        if p2 > 0 and p3 > 0:
            for i in range(p2+1, p3+1):
                mid_backswing_frame = labeled_frames[i]
                mid_backswing_frame.swing_phase = labeldata.SwingPhase.MID_BACKSWING
        p4 = psystem_positions.P4
        if p3 > 0 and p4 > 0:
            for i in range(p3+1, p4+1):
                late_backswing_frame = labeled_frames[i]
                late_backswing_frame.swing_phase = labeldata.SwingPhase.LATE_BACKSWING
        if p4 > 0:
            top_of_swing_frame = labeled_frames[p4]
            top_of_swing_frame.swing_phase = labeldata.SwingPhase.TOP_OF_SWING
        p5 = psystem_positions.P5
        if p4 > 0 and p5 > 0:
            for i in range(p4+1, p5+1):
                early_downswing_frame = labeled_frames[i]
                early_downswing_frame.swing_phase = labeldata.SwingPhase.EARLY_DOWNSWING
        p6 = psystem_positions.P6
        if p5 > 0 and p6 > 0:
            for i in range(p5+1, p6+1):
                mid_downswing_frame = labeled_frames[i]
                mid_downswing_frame.swing_phase = labeldata.SwingPhase.MID_DOWNSWING
        p7 = psystem_positions.P7
        if p6 > 0 and p7 > 0:
            for i in range(p6+1, p7+1):
                late_downswing_frame = labeled_frames[i]
                late_downswing_frame.swing_phase = labeldata.SwingPhase.LATE_DOWNSWING
        if p7 > 0:
            impact_frame = labeled_frames[p7]
            impact_frame.swing_phase = labeldata.SwingPhase.IMPACT
        p8 = psystem_positions.P8
        if p7 > 0 and p8 > 0:
            for i in range(p7+1, p8+1):
                early_follow_through_frame = labeled_frames[i]
                early_follow_through_frame.swing_phase = labeldata.SwingPhase.EARLY_FOLLOW_THROUGH
        p9 = psystem_positions.P9
        if p8 > 0 and p9 > 0:
            for i in range(p8+1, p9+1):
                mid_follow_through_frame = labeled_frames[i]
                mid_follow_through_frame.swing_phase = labeldata.SwingPhase.MID_FOLLOW_THROUGH
        p10 = psystem_positions.P10
        if p9 > 0 and p10 > 0:
            for i in range(p9+1, p10+1):
                late_follow_through_frame = labeled_frames[i]
                late_follow_through_frame.swing_phase = labeldata.SwingPhase.LATE_FOLLOW_THROUGH
        if p10 > 0:
            finish_frame = labeled_frames[p10]
            finish_frame.swing_phase = labeldata.SwingPhase.FINISH

    def prompt_user_for_information(self, prompt, expected_type, enum_type=False):
        # check that input is correct type
        user_input = input(prompt)
        try:
            if enum_type:
                parsed_val = expected_type[user_input.upper()].value
            else:
                parsed_val = expected_type(user_input)
        except ValueError:
            print(f"Incorrect type")
            return self.prompt_user_for_information(prompt, expected_type, enum_type)
        except KeyError:
            print("Incorrect enum type")
            return self.prompt_user_for_information(prompt, expected_type, enum_type)
        # verify with user that entered info is correct
        is_ok = input(f"You entered \"{user_input}\". Is this correct? 'y' for Yes, 'n' for No: ")
        if is_ok.lower() == "y":
            return parsed_val
        else:
            return self.prompt_user_for_information(prompt, expected_type, enum_type)
            
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
