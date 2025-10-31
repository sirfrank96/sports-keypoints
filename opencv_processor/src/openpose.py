import cv2 as cv
import os
import traceback
import sys
import numpy as np
from pathlib import Path

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
    if hasattr(e, 'args'):
        print(f"Args: {e.args}")
    traceback.print_exc()


params = dict()
params["model_folder"] = curr_dir / r"..\3rdparty\openpose\models"



#Converts img bytes to Matlike?
def getOpenPoseImageFromBytes(img_bytes):
    print("in getopenposeimagefrombytes")
    nparr = np.frombuffer(img_bytes, np.uint8)
    img = cv.imdecode(nparr, cv.IMREAD_COLOR)
    return getOpenPoseImage(img)

#Grabs img from path and converts to Matlike?
def getOpenPoseImageFromPath(imgPath):
    img = cv.imread(imgPath, cv.IMREAD_COLOR)
    return getOpenPoseImage(img)
    
#Input is Matlike?
def getOpenPoseImage(img):
    print("in getOpenPoseImage")
    openposeWrapper = op.WrapperPython()
    openposeWrapper.configure(params)
    openposeWrapper.start()

    datum = op.Datum()
    datum.cvInputData = img
    openposeWrapper.emplaceAndPop(op.VectorDatum([datum]))

    #print("Body keypoints: \n" + str(datum.poseKeypoints))

    #cv.imshow("OpenPose Output", datum.cvOutputData)
    #cv.waitKey(0)

    img_encode = cv.imencode('.jpg', datum.cvOutputData)
    data_encode = np.array(img_encode[1])
    print(f"np array size is {data_encode.size}, len is {len(data_encode)}")
    byte_encode = data_encode.tobytes()
    print(f"size of byte is {len(byte_encode)}")

    print("encoded")
    
    return byte_encode