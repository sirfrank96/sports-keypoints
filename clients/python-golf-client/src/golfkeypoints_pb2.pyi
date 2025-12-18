import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ImageType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    IMAGE_TYPE_UNSPECIFIED: _ClassVar[ImageType]
    FACE_ON: _ClassVar[ImageType]
    DTL: _ClassVar[ImageType]

class CalibrationType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    NO_CALIBRATION: _ClassVar[CalibrationType]
    AXES_CALIBRATION_ONLY: _ClassVar[CalibrationType]
    AXES_AND_VANISHING_POINT_CALIBRATION: _ClassVar[CalibrationType]
    FULL_CALIBRATION: _ClassVar[CalibrationType]

class FeetLineMethod(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    FEET_LINE_METHOD_UNSPECIFIED: _ClassVar[FeetLineMethod]
    USE_HEEL_LINE: _ClassVar[FeetLineMethod]
    USE_TOE_LINE: _ClassVar[FeetLineMethod]
IMAGE_TYPE_UNSPECIFIED: ImageType
FACE_ON: ImageType
DTL: ImageType
NO_CALIBRATION: CalibrationType
AXES_CALIBRATION_ONLY: CalibrationType
AXES_AND_VANISHING_POINT_CALIBRATION: CalibrationType
FULL_CALIBRATION: CalibrationType
FEET_LINE_METHOD_UNSPECIFIED: FeetLineMethod
USE_HEEL_LINE: FeetLineMethod
USE_TOE_LINE: FeetLineMethod

class UploadInputImageRequest(_message.Message):
    __slots__ = ("session_token", "image_type", "image")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    IMAGE_TYPE_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    image_type: ImageType
    image: bytes
    def __init__(self, session_token: _Optional[str] = ..., image_type: _Optional[_Union[ImageType, str]] = ..., image: _Optional[bytes] = ...) -> None: ...

class UploadInputImageResponse(_message.Message):
    __slots__ = ("success", "input_image_id")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    success: bool
    input_image_id: str
    def __init__(self, success: bool = ..., input_image_id: _Optional[str] = ...) -> None: ...

class ListInputImagesForUserRequest(_message.Message):
    __slots__ = ("session_token",)
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    def __init__(self, session_token: _Optional[str] = ...) -> None: ...

class ListInputImagesForUserResponse(_message.Message):
    __slots__ = ("success", "input_image_ids")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_IDS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    input_image_ids: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, success: bool = ..., input_image_ids: _Optional[_Iterable[str]] = ...) -> None: ...

class ReadInputImageRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ...) -> None: ...

class ReadInputImageResponse(_message.Message):
    __slots__ = ("success", "image_type", "image", "calibration_type", "feet_line_method")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_TYPE_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    CALIBRATION_TYPE_FIELD_NUMBER: _ClassVar[int]
    FEET_LINE_METHOD_FIELD_NUMBER: _ClassVar[int]
    success: bool
    image_type: ImageType
    image: bytes
    calibration_type: CalibrationType
    feet_line_method: FeetLineMethod
    def __init__(self, success: bool = ..., image_type: _Optional[_Union[ImageType, str]] = ..., image: _Optional[bytes] = ..., calibration_type: _Optional[_Union[CalibrationType, str]] = ..., feet_line_method: _Optional[_Union[FeetLineMethod, str]] = ...) -> None: ...

class DeleteInputImageRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ...) -> None: ...

class DeleteInputImageResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class CalibrateInputImageRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id", "calibration_type", "feet_line_method", "calibration_image_axes", "calibration_image_vanishing_point", "golf_ball", "club_butt", "club_head", "shoulder_tilt")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    CALIBRATION_TYPE_FIELD_NUMBER: _ClassVar[int]
    FEET_LINE_METHOD_FIELD_NUMBER: _ClassVar[int]
    CALIBRATION_IMAGE_AXES_FIELD_NUMBER: _ClassVar[int]
    CALIBRATION_IMAGE_VANISHING_POINT_FIELD_NUMBER: _ClassVar[int]
    GOLF_BALL_FIELD_NUMBER: _ClassVar[int]
    CLUB_BUTT_FIELD_NUMBER: _ClassVar[int]
    CLUB_HEAD_FIELD_NUMBER: _ClassVar[int]
    SHOULDER_TILT_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    calibration_type: CalibrationType
    feet_line_method: FeetLineMethod
    calibration_image_axes: bytes
    calibration_image_vanishing_point: bytes
    golf_ball: _common_pb2.Keypoint
    club_butt: _common_pb2.Keypoint
    club_head: _common_pb2.Keypoint
    shoulder_tilt: _common_pb2.Double
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ..., calibration_type: _Optional[_Union[CalibrationType, str]] = ..., feet_line_method: _Optional[_Union[FeetLineMethod, str]] = ..., calibration_image_axes: _Optional[bytes] = ..., calibration_image_vanishing_point: _Optional[bytes] = ..., golf_ball: _Optional[_Union[_common_pb2.Keypoint, _Mapping]] = ..., club_butt: _Optional[_Union[_common_pb2.Keypoint, _Mapping]] = ..., club_head: _Optional[_Union[_common_pb2.Keypoint, _Mapping]] = ..., shoulder_tilt: _Optional[_Union[_common_pb2.Double, _Mapping]] = ...) -> None: ...

class CalibrateInputImageResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class CalculateGolfKeypointsRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ...) -> None: ...

class CalculateGolfKeypointsResponse(_message.Message):
    __slots__ = ("success", "output_image", "golf_keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    OUTPUT_IMAGE_FIELD_NUMBER: _ClassVar[int]
    GOLF_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    output_image: bytes
    golf_keypoints: GolfKeypoints
    def __init__(self, success: bool = ..., output_image: _Optional[bytes] = ..., golf_keypoints: _Optional[_Union[GolfKeypoints, _Mapping]] = ...) -> None: ...

class ReadGolfKeypointsRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ...) -> None: ...

class ReadGolfKeypointsResponse(_message.Message):
    __slots__ = ("success", "output_image", "golf_keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    OUTPUT_IMAGE_FIELD_NUMBER: _ClassVar[int]
    GOLF_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    output_image: bytes
    golf_keypoints: GolfKeypoints
    def __init__(self, success: bool = ..., output_image: _Optional[bytes] = ..., golf_keypoints: _Optional[_Union[GolfKeypoints, _Mapping]] = ...) -> None: ...

class UpdateBodyKeypointsRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id", "updated_body_keypoints")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    UPDATED_BODY_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    updated_body_keypoints: _common_pb2.Body25PoseKeypoints
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ..., updated_body_keypoints: _Optional[_Union[_common_pb2.Body25PoseKeypoints, _Mapping]] = ...) -> None: ...

class UpdateBodyKeypointsResponse(_message.Message):
    __slots__ = ("success", "updated_golf_keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    UPDATED_GOLF_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    updated_golf_keypoints: GolfKeypoints
    def __init__(self, success: bool = ..., updated_golf_keypoints: _Optional[_Union[GolfKeypoints, _Mapping]] = ...) -> None: ...

class DeleteGolfKeypointsRequest(_message.Message):
    __slots__ = ("session_token", "input_image_id")
    SESSION_TOKEN_FIELD_NUMBER: _ClassVar[int]
    INPUT_IMAGE_ID_FIELD_NUMBER: _ClassVar[int]
    session_token: str
    input_image_id: str
    def __init__(self, session_token: _Optional[str] = ..., input_image_id: _Optional[str] = ...) -> None: ...

class DeleteGolfKeypointsResponse(_message.Message):
    __slots__ = ("success",)
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    def __init__(self, success: bool = ...) -> None: ...

class GolfKeypoints(_message.Message):
    __slots__ = ("dtl_golf_setup_points", "faceon_golf_setup_points", "body_keypoints")
    DTL_GOLF_SETUP_POINTS_FIELD_NUMBER: _ClassVar[int]
    FACEON_GOLF_SETUP_POINTS_FIELD_NUMBER: _ClassVar[int]
    BODY_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    dtl_golf_setup_points: DTLGolfSetupPoints
    faceon_golf_setup_points: FaceOnGolfSetupPoints
    body_keypoints: _common_pb2.Body25PoseKeypoints
    def __init__(self, dtl_golf_setup_points: _Optional[_Union[DTLGolfSetupPoints, _Mapping]] = ..., faceon_golf_setup_points: _Optional[_Union[FaceOnGolfSetupPoints, _Mapping]] = ..., body_keypoints: _Optional[_Union[_common_pb2.Body25PoseKeypoints, _Mapping]] = ...) -> None: ...

class DTLGolfSetupPoints(_message.Message):
    __slots__ = ("spine_angle", "feet_alignment", "heel_alignment", "toe_alignment", "shoulder_alignment", "waist_alignment", "knee_bend", "distance_from_ball", "ulnar_deviation")
    SPINE_ANGLE_FIELD_NUMBER: _ClassVar[int]
    FEET_ALIGNMENT_FIELD_NUMBER: _ClassVar[int]
    HEEL_ALIGNMENT_FIELD_NUMBER: _ClassVar[int]
    TOE_ALIGNMENT_FIELD_NUMBER: _ClassVar[int]
    SHOULDER_ALIGNMENT_FIELD_NUMBER: _ClassVar[int]
    WAIST_ALIGNMENT_FIELD_NUMBER: _ClassVar[int]
    KNEE_BEND_FIELD_NUMBER: _ClassVar[int]
    DISTANCE_FROM_BALL_FIELD_NUMBER: _ClassVar[int]
    ULNAR_DEVIATION_FIELD_NUMBER: _ClassVar[int]
    spine_angle: _common_pb2.Double
    feet_alignment: _common_pb2.Double
    heel_alignment: _common_pb2.Double
    toe_alignment: _common_pb2.Double
    shoulder_alignment: _common_pb2.Double
    waist_alignment: _common_pb2.Double
    knee_bend: _common_pb2.Double
    distance_from_ball: _common_pb2.Double
    ulnar_deviation: _common_pb2.Double
    def __init__(self, spine_angle: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., feet_alignment: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., heel_alignment: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., toe_alignment: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., shoulder_alignment: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., waist_alignment: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., knee_bend: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., distance_from_ball: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., ulnar_deviation: _Optional[_Union[_common_pb2.Double, _Mapping]] = ...) -> None: ...

class FaceOnGolfSetupPoints(_message.Message):
    __slots__ = ("side_bend", "l_foot_flare", "r_foot_flare", "stance_width", "shoulder_tilt", "waist_tilt", "shaft_lean", "ball_position", "head_position", "chest_position", "mid_hip_position")
    SIDE_BEND_FIELD_NUMBER: _ClassVar[int]
    L_FOOT_FLARE_FIELD_NUMBER: _ClassVar[int]
    R_FOOT_FLARE_FIELD_NUMBER: _ClassVar[int]
    STANCE_WIDTH_FIELD_NUMBER: _ClassVar[int]
    SHOULDER_TILT_FIELD_NUMBER: _ClassVar[int]
    WAIST_TILT_FIELD_NUMBER: _ClassVar[int]
    SHAFT_LEAN_FIELD_NUMBER: _ClassVar[int]
    BALL_POSITION_FIELD_NUMBER: _ClassVar[int]
    HEAD_POSITION_FIELD_NUMBER: _ClassVar[int]
    CHEST_POSITION_FIELD_NUMBER: _ClassVar[int]
    MID_HIP_POSITION_FIELD_NUMBER: _ClassVar[int]
    side_bend: _common_pb2.Double
    l_foot_flare: _common_pb2.Double
    r_foot_flare: _common_pb2.Double
    stance_width: _common_pb2.Double
    shoulder_tilt: _common_pb2.Double
    waist_tilt: _common_pb2.Double
    shaft_lean: _common_pb2.Double
    ball_position: _common_pb2.Double
    head_position: _common_pb2.Double
    chest_position: _common_pb2.Double
    mid_hip_position: _common_pb2.Double
    def __init__(self, side_bend: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., l_foot_flare: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., r_foot_flare: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., stance_width: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., shoulder_tilt: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., waist_tilt: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., shaft_lean: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., ball_position: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., head_position: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., chest_position: _Optional[_Union[_common_pb2.Double, _Mapping]] = ..., mid_hip_position: _Optional[_Union[_common_pb2.Double, _Mapping]] = ...) -> None: ...
