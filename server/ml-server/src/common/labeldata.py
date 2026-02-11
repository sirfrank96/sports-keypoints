from enum import Enum
from pydantic import BaseModel
import json
from typing import List

class Orientation(Enum):
    FACEON = 1
    DTL = 2
    NEITHER = 3

class GolfLevel(Enum):
    TOPPRO = 1
    PRO = 2
    TOPAM = 3
    GOODAM = 4
    OKAM = 5
    BADAM = 6

class Sex(Enum):
    UNKNOWN = 1
    NONBINARY = 2
    MALE = 3
    FEMALE = 4

class Handedness(Enum):
    LEFT = 1
    RIGHT = 2
    AMBIDEXTROUS = 3

class ShotType(Enum):
    FULLSHOT = 1
    PARTIALSHOT = 2
    SHORTGAME = 3
    PUTTING = 4
    OTHER = 5

class ShotShape(Enum):
    NA = 1
    STRAIGHT = 2
    DRAW = 3
    FADE = 4

class ShotHeight(Enum):
    NA = 1
    LOW = 2
    MID = 3
    HIGH = 4
    OTHER = 5

class Club(Enum):
    DRIVER = 1
    WOOD = 2
    HYBRID = 3
    LONGIRON = 4
    MIDIRON = 5
    SHORTIRON = 6
    WEDGE = 7
    PUTTER = 8

class SwingPhase(Enum):
    P1 = 1 # Setup
    P2 = 2 # Shaft parallel backswing
    P3 = 3 # Lead arm parallel backswing
    P4 = 4 # Top of backswing
    P5 = 5 # Lead arm parallel downswing
    P6 = 6 # Shaft parallel downswing
    P7 = 7 # Impact
    P8 = 8 # Shaft parallel follow through
    P9 = 9 # Trail arm parallel  follow through
    P10 = 10 # Finish

class BallData(BaseModel): # Launch monitor balldata
    ball_speed: float
    spin_rate: float
    apex_height: float
    launch_angle: float
    spin_axis: float
    smash: float
    curve: float
    carry_distance: float
    total_distance: float

class ClubData(BaseModel): # Launch monitor clubdata
    club_speed: float
    aoa: float
    swing_direction: float
    face_to_path: float
    dynamic_loft: float

class LabeledFrame(BaseModel):
    frame_idx: int
    is_swing: bool
    swing_phase: SwingPhase
    
class LabeledVideo(BaseModel):
    video_name: str
    body_datapoints_idx: int
    player_id: str
    intended_shot: str
    description: str
    feel_description: str
    result: str
    shot_type: ShotType
    shot_shape: ShotShape
    shot_height: ShotHeight
    orientation: Orientation
    club: Club
    ball_data: BallData
    club_data: ClubData
    labeled_frames: List[LabeledFrame]

class LabeledPlayer(BaseModel):
    player_id: str
    description: str
    golf_level: GolfLevel
    handicap: float
    height: float
    weight: float
    sex: Sex
    golf_handedness: Handedness
    real_handedness: Handedness

class LabeledGolfSwingDataset(BaseModel):
    players: dict[str, LabeledPlayer]
    videos: dict[str, LabeledVideo]
