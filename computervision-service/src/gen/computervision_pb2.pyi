import common_pb2 as _common_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Mapping as _Mapping, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetPoseImageRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseImageResponse(_message.Message):
    __slots__ = ("success", "image")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    image: bytes
    def __init__(self, success: bool = ..., image: _Optional[bytes] = ...) -> None: ...

class GetPoseDataRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseDataResponse(_message.Message):
    __slots__ = ("success", "keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    keypoints: _common_pb2.Body25PoseKeypoints
    def __init__(self, success: bool = ..., keypoints: _Optional[_Union[_common_pb2.Body25PoseKeypoints, _Mapping]] = ...) -> None: ...

class GetPoseHandImageRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseHandImageResponse(_message.Message):
    __slots__ = ("success", "image")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    image: bytes
    def __init__(self, success: bool = ..., image: _Optional[bytes] = ...) -> None: ...

class GetPoseHandDataRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseHandDataResponse(_message.Message):
    __slots__ = ("success", "keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    keypoints: _common_pb2.Body25HandKeypoints
    def __init__(self, success: bool = ..., keypoints: _Optional[_Union[_common_pb2.Body25HandKeypoints, _Mapping]] = ...) -> None: ...

class GetPoseAllRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseAllResponse(_message.Message):
    __slots__ = ("success", "image", "pose_keypoints", "hand_keypoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    POSE_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    HAND_KEYPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    image: bytes
    pose_keypoints: _common_pb2.Body25PoseKeypoints
    hand_keypoints: _common_pb2.Body25HandKeypoints
    def __init__(self, success: bool = ..., image: _Optional[bytes] = ..., pose_keypoints: _Optional[_Union[_common_pb2.Body25PoseKeypoints, _Mapping]] = ..., hand_keypoints: _Optional[_Union[_common_pb2.Body25HandKeypoints, _Mapping]] = ...) -> None: ...
