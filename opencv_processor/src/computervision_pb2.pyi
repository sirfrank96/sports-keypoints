from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class DTLImage(_message.Message):
    __slots__ = ("name", "image")
    NAME_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    name: str
    image: bytes
    def __init__(self, name: _Optional[str] = ..., image: _Optional[bytes] = ...) -> None: ...

class FaceOnImage(_message.Message):
    __slots__ = ("name", "image")
    NAME_FIELD_NUMBER: _ClassVar[int]
    IMAGE_FIELD_NUMBER: _ClassVar[int]
    name: str
    image: bytes
    def __init__(self, name: _Optional[str] = ..., image: _Optional[bytes] = ...) -> None: ...
