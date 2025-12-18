from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Double(_message.Message):
    __slots__ = ("data", "warning")
    DATA_FIELD_NUMBER: _ClassVar[int]
    WARNING_FIELD_NUMBER: _ClassVar[int]
    data: float
    warning: str
    def __init__(self, data: _Optional[float] = ..., warning: _Optional[str] = ...) -> None: ...

class Keypoint(_message.Message):
    __slots__ = ("x", "y", "confidence")
    X_FIELD_NUMBER: _ClassVar[int]
    Y_FIELD_NUMBER: _ClassVar[int]
    CONFIDENCE_FIELD_NUMBER: _ClassVar[int]
    x: float
    y: float
    confidence: float
    def __init__(self, x: _Optional[float] = ..., y: _Optional[float] = ..., confidence: _Optional[float] = ...) -> None: ...

class Body25PoseKeypoints(_message.Message):
    __slots__ = ("nose", "neck", "r_shoulder", "r_elbow", "r_wrist", "l_shoulder", "l_elbow", "l_wrist", "midhip", "r_hip", "r_knee", "r_ankle", "l_hip", "l_knee", "l_ankle", "r_eye", "l_eye", "r_ear", "l_ear", "l_big_toe", "l_small_toe", "l_heel", "r_big_toe", "r_small_toe", "r_heel")
    NOSE_FIELD_NUMBER: _ClassVar[int]
    NECK_FIELD_NUMBER: _ClassVar[int]
    R_SHOULDER_FIELD_NUMBER: _ClassVar[int]
    R_ELBOW_FIELD_NUMBER: _ClassVar[int]
    R_WRIST_FIELD_NUMBER: _ClassVar[int]
    L_SHOULDER_FIELD_NUMBER: _ClassVar[int]
    L_ELBOW_FIELD_NUMBER: _ClassVar[int]
    L_WRIST_FIELD_NUMBER: _ClassVar[int]
    MIDHIP_FIELD_NUMBER: _ClassVar[int]
    R_HIP_FIELD_NUMBER: _ClassVar[int]
    R_KNEE_FIELD_NUMBER: _ClassVar[int]
    R_ANKLE_FIELD_NUMBER: _ClassVar[int]
    L_HIP_FIELD_NUMBER: _ClassVar[int]
    L_KNEE_FIELD_NUMBER: _ClassVar[int]
    L_ANKLE_FIELD_NUMBER: _ClassVar[int]
    R_EYE_FIELD_NUMBER: _ClassVar[int]
    L_EYE_FIELD_NUMBER: _ClassVar[int]
    R_EAR_FIELD_NUMBER: _ClassVar[int]
    L_EAR_FIELD_NUMBER: _ClassVar[int]
    L_BIG_TOE_FIELD_NUMBER: _ClassVar[int]
    L_SMALL_TOE_FIELD_NUMBER: _ClassVar[int]
    L_HEEL_FIELD_NUMBER: _ClassVar[int]
    R_BIG_TOE_FIELD_NUMBER: _ClassVar[int]
    R_SMALL_TOE_FIELD_NUMBER: _ClassVar[int]
    R_HEEL_FIELD_NUMBER: _ClassVar[int]
    nose: Keypoint
    neck: Keypoint
    r_shoulder: Keypoint
    r_elbow: Keypoint
    r_wrist: Keypoint
    l_shoulder: Keypoint
    l_elbow: Keypoint
    l_wrist: Keypoint
    midhip: Keypoint
    r_hip: Keypoint
    r_knee: Keypoint
    r_ankle: Keypoint
    l_hip: Keypoint
    l_knee: Keypoint
    l_ankle: Keypoint
    r_eye: Keypoint
    l_eye: Keypoint
    r_ear: Keypoint
    l_ear: Keypoint
    l_big_toe: Keypoint
    l_small_toe: Keypoint
    l_heel: Keypoint
    r_big_toe: Keypoint
    r_small_toe: Keypoint
    r_heel: Keypoint
    def __init__(self, nose: _Optional[_Union[Keypoint, _Mapping]] = ..., neck: _Optional[_Union[Keypoint, _Mapping]] = ..., r_shoulder: _Optional[_Union[Keypoint, _Mapping]] = ..., r_elbow: _Optional[_Union[Keypoint, _Mapping]] = ..., r_wrist: _Optional[_Union[Keypoint, _Mapping]] = ..., l_shoulder: _Optional[_Union[Keypoint, _Mapping]] = ..., l_elbow: _Optional[_Union[Keypoint, _Mapping]] = ..., l_wrist: _Optional[_Union[Keypoint, _Mapping]] = ..., midhip: _Optional[_Union[Keypoint, _Mapping]] = ..., r_hip: _Optional[_Union[Keypoint, _Mapping]] = ..., r_knee: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ankle: _Optional[_Union[Keypoint, _Mapping]] = ..., l_hip: _Optional[_Union[Keypoint, _Mapping]] = ..., l_knee: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ankle: _Optional[_Union[Keypoint, _Mapping]] = ..., r_eye: _Optional[_Union[Keypoint, _Mapping]] = ..., l_eye: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ear: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ear: _Optional[_Union[Keypoint, _Mapping]] = ..., l_big_toe: _Optional[_Union[Keypoint, _Mapping]] = ..., l_small_toe: _Optional[_Union[Keypoint, _Mapping]] = ..., l_heel: _Optional[_Union[Keypoint, _Mapping]] = ..., r_big_toe: _Optional[_Union[Keypoint, _Mapping]] = ..., r_small_toe: _Optional[_Union[Keypoint, _Mapping]] = ..., r_heel: _Optional[_Union[Keypoint, _Mapping]] = ...) -> None: ...

class Body25HandKeypoints(_message.Message):
    __slots__ = ("l_wrist", "l_thumb1", "l_thumb2", "l_thumb3", "l_thumb", "l_index1", "l_index2", "l_index3", "l_index", "l_middle1", "l_middle2", "l_middle3", "l_middle", "l_ring1", "l_ring2", "l_ring3", "l_ring", "l_pinky1", "l_pinky2", "l_pinky3", "l_pinky", "r_wrist", "r_thumb1", "r_thumb2", "r_thumb3", "r_thumb", "r_index1", "r_index2", "r_index3", "r_index", "r_middle1", "r_middle2", "r_middle3", "r_middle", "r_ring1", "r_ring2", "r_ring3", "r_ring", "r_pinky1", "r_pinky2", "r_pinky3", "r_pinky")
    L_WRIST_FIELD_NUMBER: _ClassVar[int]
    L_THUMB1_FIELD_NUMBER: _ClassVar[int]
    L_THUMB2_FIELD_NUMBER: _ClassVar[int]
    L_THUMB3_FIELD_NUMBER: _ClassVar[int]
    L_THUMB_FIELD_NUMBER: _ClassVar[int]
    L_INDEX1_FIELD_NUMBER: _ClassVar[int]
    L_INDEX2_FIELD_NUMBER: _ClassVar[int]
    L_INDEX3_FIELD_NUMBER: _ClassVar[int]
    L_INDEX_FIELD_NUMBER: _ClassVar[int]
    L_MIDDLE1_FIELD_NUMBER: _ClassVar[int]
    L_MIDDLE2_FIELD_NUMBER: _ClassVar[int]
    L_MIDDLE3_FIELD_NUMBER: _ClassVar[int]
    L_MIDDLE_FIELD_NUMBER: _ClassVar[int]
    L_RING1_FIELD_NUMBER: _ClassVar[int]
    L_RING2_FIELD_NUMBER: _ClassVar[int]
    L_RING3_FIELD_NUMBER: _ClassVar[int]
    L_RING_FIELD_NUMBER: _ClassVar[int]
    L_PINKY1_FIELD_NUMBER: _ClassVar[int]
    L_PINKY2_FIELD_NUMBER: _ClassVar[int]
    L_PINKY3_FIELD_NUMBER: _ClassVar[int]
    L_PINKY_FIELD_NUMBER: _ClassVar[int]
    R_WRIST_FIELD_NUMBER: _ClassVar[int]
    R_THUMB1_FIELD_NUMBER: _ClassVar[int]
    R_THUMB2_FIELD_NUMBER: _ClassVar[int]
    R_THUMB3_FIELD_NUMBER: _ClassVar[int]
    R_THUMB_FIELD_NUMBER: _ClassVar[int]
    R_INDEX1_FIELD_NUMBER: _ClassVar[int]
    R_INDEX2_FIELD_NUMBER: _ClassVar[int]
    R_INDEX3_FIELD_NUMBER: _ClassVar[int]
    R_INDEX_FIELD_NUMBER: _ClassVar[int]
    R_MIDDLE1_FIELD_NUMBER: _ClassVar[int]
    R_MIDDLE2_FIELD_NUMBER: _ClassVar[int]
    R_MIDDLE3_FIELD_NUMBER: _ClassVar[int]
    R_MIDDLE_FIELD_NUMBER: _ClassVar[int]
    R_RING1_FIELD_NUMBER: _ClassVar[int]
    R_RING2_FIELD_NUMBER: _ClassVar[int]
    R_RING3_FIELD_NUMBER: _ClassVar[int]
    R_RING_FIELD_NUMBER: _ClassVar[int]
    R_PINKY1_FIELD_NUMBER: _ClassVar[int]
    R_PINKY2_FIELD_NUMBER: _ClassVar[int]
    R_PINKY3_FIELD_NUMBER: _ClassVar[int]
    R_PINKY_FIELD_NUMBER: _ClassVar[int]
    l_wrist: Keypoint
    l_thumb1: Keypoint
    l_thumb2: Keypoint
    l_thumb3: Keypoint
    l_thumb: Keypoint
    l_index1: Keypoint
    l_index2: Keypoint
    l_index3: Keypoint
    l_index: Keypoint
    l_middle1: Keypoint
    l_middle2: Keypoint
    l_middle3: Keypoint
    l_middle: Keypoint
    l_ring1: Keypoint
    l_ring2: Keypoint
    l_ring3: Keypoint
    l_ring: Keypoint
    l_pinky1: Keypoint
    l_pinky2: Keypoint
    l_pinky3: Keypoint
    l_pinky: Keypoint
    r_wrist: Keypoint
    r_thumb1: Keypoint
    r_thumb2: Keypoint
    r_thumb3: Keypoint
    r_thumb: Keypoint
    r_index1: Keypoint
    r_index2: Keypoint
    r_index3: Keypoint
    r_index: Keypoint
    r_middle1: Keypoint
    r_middle2: Keypoint
    r_middle3: Keypoint
    r_middle: Keypoint
    r_ring1: Keypoint
    r_ring2: Keypoint
    r_ring3: Keypoint
    r_ring: Keypoint
    r_pinky1: Keypoint
    r_pinky2: Keypoint
    r_pinky3: Keypoint
    r_pinky: Keypoint
    def __init__(self, l_wrist: _Optional[_Union[Keypoint, _Mapping]] = ..., l_thumb1: _Optional[_Union[Keypoint, _Mapping]] = ..., l_thumb2: _Optional[_Union[Keypoint, _Mapping]] = ..., l_thumb3: _Optional[_Union[Keypoint, _Mapping]] = ..., l_thumb: _Optional[_Union[Keypoint, _Mapping]] = ..., l_index1: _Optional[_Union[Keypoint, _Mapping]] = ..., l_index2: _Optional[_Union[Keypoint, _Mapping]] = ..., l_index3: _Optional[_Union[Keypoint, _Mapping]] = ..., l_index: _Optional[_Union[Keypoint, _Mapping]] = ..., l_middle1: _Optional[_Union[Keypoint, _Mapping]] = ..., l_middle2: _Optional[_Union[Keypoint, _Mapping]] = ..., l_middle3: _Optional[_Union[Keypoint, _Mapping]] = ..., l_middle: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ring1: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ring2: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ring3: _Optional[_Union[Keypoint, _Mapping]] = ..., l_ring: _Optional[_Union[Keypoint, _Mapping]] = ..., l_pinky1: _Optional[_Union[Keypoint, _Mapping]] = ..., l_pinky2: _Optional[_Union[Keypoint, _Mapping]] = ..., l_pinky3: _Optional[_Union[Keypoint, _Mapping]] = ..., l_pinky: _Optional[_Union[Keypoint, _Mapping]] = ..., r_wrist: _Optional[_Union[Keypoint, _Mapping]] = ..., r_thumb1: _Optional[_Union[Keypoint, _Mapping]] = ..., r_thumb2: _Optional[_Union[Keypoint, _Mapping]] = ..., r_thumb3: _Optional[_Union[Keypoint, _Mapping]] = ..., r_thumb: _Optional[_Union[Keypoint, _Mapping]] = ..., r_index1: _Optional[_Union[Keypoint, _Mapping]] = ..., r_index2: _Optional[_Union[Keypoint, _Mapping]] = ..., r_index3: _Optional[_Union[Keypoint, _Mapping]] = ..., r_index: _Optional[_Union[Keypoint, _Mapping]] = ..., r_middle1: _Optional[_Union[Keypoint, _Mapping]] = ..., r_middle2: _Optional[_Union[Keypoint, _Mapping]] = ..., r_middle3: _Optional[_Union[Keypoint, _Mapping]] = ..., r_middle: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ring1: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ring2: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ring3: _Optional[_Union[Keypoint, _Mapping]] = ..., r_ring: _Optional[_Union[Keypoint, _Mapping]] = ..., r_pinky1: _Optional[_Union[Keypoint, _Mapping]] = ..., r_pinky2: _Optional[_Union[Keypoint, _Mapping]] = ..., r_pinky3: _Optional[_Union[Keypoint, _Mapping]] = ..., r_pinky: _Optional[_Union[Keypoint, _Mapping]] = ...) -> None: ...
