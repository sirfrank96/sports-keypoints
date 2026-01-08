# Python
from PIL import Image
import io
import datetime

# Tkinter
from tkinter import filedialog

# GRPC
from google.protobuf.timestamp_pb2 import Timestamp


def get_image_from_filesystem():
    # Open the file dialog and get the file path
    filepath = filedialog.askopenfilename(
        filetypes=[("Image Files", "*.png;*.jpg;*.jpeg;*.gif")]
    )
    # If a file was selected, return Image object
    if filepath:
        img = Image.open(filepath)
        return img
    else:
        return None

def get_image_bytes(image):
    buffer = io.BytesIO()
    image.save(buffer, format='JPEG')
    return buffer.getvalue()

def get_curr_grpc_timestamp():
    now_utc = datetime.datetime.now(datetime.timezone.utc)
    grpc_timestamp = Timestamp()
    grpc_timestamp.FromDatetime(now_utc) 
    return grpc_timestamp
 