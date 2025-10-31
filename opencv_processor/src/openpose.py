import cv2 as cv
import os
import traceback
import sys
import numpy as np
from pathlib import Path

# Import pyopenpose, make sure sys can find paths necessary dlls and binaries
curr_dir = Path(__file__).parent.resolve()
sys.path.append(curr_dir / r"..\3rdparty\openpose\build_windows\python\openpose\Release")
os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\x64\Release")
os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\bin")
try:
    import pyopenpose as op
except ImportError as e:
    print("ERROR: {e}")
    print(f"Type of error: {type(e)}")
    if hasattr(e, 'name'):
        print(f"Module name: {e.name}")
    if hasattr(e, 'path'):
        print(f"Path: {e.path}")
    if hasattr(e, 'msg'):
        print(f"Msg: {e.msg}")
    traceback.print_exc()

# Location of OpenPose models
params = dict()
params["model_folder"] = curr_dir / r"..\3rdparty\openpose\models"


# Helper function: converts bytes object to Matlike
def getOpenPoseImageFromBytes(img_bytes):
    nparr = np.frombuffer(img_bytes, np.uint8)
    img = cv.imdecode(nparr, cv.IMREAD_COLOR)
    return getOpenPoseImage(img)

# Helper function: Grabs img from path and converts to Matlike
def getOpenPoseImageFromPath(imgPath):
    img = cv.imread(imgPath, cv.IMREAD_COLOR)
    return getOpenPoseImage(img)
    
# Input is Matlike, Feed into OpenPose, Return processed image as bytes
def getOpenPoseImage(img):
    # Get OpenPose python wrapper
    openposeWrapper = op.WrapperPython()
    openposeWrapper.configure(params)
    openposeWrapper.start()

    # Process input image to get OpenPose image with skeleton lines
    datum = op.Datum()
    datum.cvInputData = img
    openposeWrapper.emplaceAndPop(op.VectorDatum([datum]))

    # Encode data as jpg, then convert to bytes object
    img_encode = cv.imencode('.jpg', datum.cvOutputData)
    data_encode = np.array(img_encode[1])
    byte_encode = data_encode.tobytes()
    
    return byte_encode