from . import common_pb2 as _common_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

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
    __slots__ = ("success", "datapoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATAPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    datapoints: _common_pb2.Body25PoseDatapoints
    def __init__(self, success: bool = ..., datapoints: _Optional[_Union[_common_pb2.Body25PoseDatapoints, _Mapping]] = ...) -> None: ...

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
    __slots__ = ("success", "datapoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    DATAPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    datapoints: _common_pb2.Body25HandDatapoints
    def __init__(self, success: bool = ..., datapoints: _Optional[_Union[_common_pb2.Body25HandDatapoints, _Mapping]] = ...) -> None: ...

class GetPoseAllRequest(_message.Message):
    __slots__ = ("image",)
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    image: bytes
    def __init__(self, image: _Optional[bytes] = ...) -> None: ...

class GetPoseAllResponse(_message.Message):
    __slots__ = ("success", "image", "pose_datapoints", "hand_datapoints")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    POSE_DATAPOINTS_FIELD_NUMBER: _ClassVar[int]
    HAND_DATAPOINTS_FIELD_NUMBER: _ClassVar[int]
    success: bool
    image: bytes
    pose_datapoints: _common_pb2.Body25PoseDatapoints
    hand_datapoints: _common_pb2.Body25HandDatapoints
    def __init__(self, success: bool = ..., image: _Optional[bytes] = ..., pose_datapoints: _Optional[_Union[_common_pb2.Body25PoseDatapoints, _Mapping]] = ..., hand_datapoints: _Optional[_Union[_common_pb2.Body25HandDatapoints, _Mapping]] = ...) -> None: ...
