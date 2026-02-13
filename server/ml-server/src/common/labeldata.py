# Python
from enum import Enum
from pydantic import BaseModel, ConfigDict
from typing import List
import numpy as np

# Internal
from src.common import constants

class Orientation(Enum):
    UNKNOWN = 1
    FACEON = 2
    DTL = 3
    NEITHER = 4

class GolfLevel(Enum):
    UNKNOWN = 1
    TOPPRO = 2
    PRO = 3
    TOPAM = 4
    GOODAM = 5
    OKAM = 6
    BADAM = 7

class Sex(Enum):
    UNKNOWN = 1
    NONBINARY = 2
    MALE = 3
    FEMALE = 4

class Handedness(Enum):
    UNKNOWN = 1
    LEFT = 2
    RIGHT = 3
    AMBIDEXTROUS = 4

class Club(Enum):
    UNKNOWN = 1
    DRIVER = 2
    WOOD = 3
    HYBRID = 4
    LONGIRON = 5
    MIDIRON = 6
    SHORTIRON = 7
    WEDGE = 8
    PUTTER = 9

class PSystemPositions(BaseModel):
    P1: int = -1 # Setup
    P2: int = -1 # Shaft parallel backswing
    P3: int = -1 # Lead arm parallel backswing
    P4: int = -1 # Top of backswing
    P5: int = -1 # Lead arm parallel downswing
    P6: int = -1 # Shaft parallel downswing
    P7: int = -1 # Impact
    P8: int = -1 # Shaft parallel follow through
    P9: int = -1 # Trail arm parallel  follow through
    P10: int = -1 # Finish

class SwingPhase(Enum):
    NOT_SWING = 1
    SETUP = 2
    TAKEAWAY = 3
    MID_BACKSWING = 4
    LATE_BACKSWING = 5
    TOP_OF_SWING = 6
    EARLY_DOWNSWING = 7
    MID_DOWNSWING = 8
    LATE_DOWNSWING = 9
    IMPACT = 10
    EARLY_FOLLOW_THROUGH = 11
    MID_FOLLOW_THROUGH = 12
    LATE_FOLLOW_THROUGH = 13
    FINISH = 14

class BallData(BaseModel): # Launch monitor balldata
    ball_speed: float = np.nan
    spin_rate: float = np.nan
    apex_height: float = np.nan
    launch_angle: float = np.nan
    spin_axis: float = np.nan
    smash: float = np.nan
    curve: float = np.nan
    carry_distance: float = np.nan
    total_distance: float = np.nan

class ClubData(BaseModel): # Launch monitor clubdata
    club_speed: float = np.nan
    angle_of_attack: float = np.nan
    swing_direction: float = np.nan
    face_to_path: float = np.nan
    dynamic_loft: float = np.nan

class LabeledFrame(BaseModel):
    frame_idx: int
    swing_phase: SwingPhase = SwingPhase.NOT_SWING
    
class LabeledVideo(BaseModel):
    id: str
    body_datapoints_idx: int
    player_id: str
    intended_shot: str = ""
    description: str = ""
    feel_description: str = ""
    result: str = ""
    orientation: Orientation = Orientation.UNKNOWN
    club: Club = Club.UNKNOWN
    ball_data: BallData = BallData()
    club_data: ClubData = ClubData()
    labeled_frames: List[LabeledFrame] = [LabeledFrame(frame_idx=idx) for idx in range(constants.FRAME_DIMENSION)]

class LabeledPlayer(BaseModel):
    id: str
    description: str = ""
    golf_level: GolfLevel = GolfLevel.UNKNOWN
    handicap: float = np.nan
    height: float = np.nan
    weight: float = np.nan
    sex: Sex = Sex.UNKNOWN
    golf_handedness: Handedness = Handedness.UNKNOWN
    real_handedness: Handedness = Handedness.UNKNOWN

class LabeledGolfSwingDataset(BaseModel):
    model_config = ConfigDict(use_enum_values=True)
    players: dict[str, LabeledPlayer]
    videos: dict[str, LabeledVideo]
