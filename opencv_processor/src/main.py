import cv2 as cv
import os
import traceback
import sys
from pathlib import Path

cwd = Path.cwd()

sys.path.append(cwd / r"opencv_processor\3rdparty\openpose\build_windows\python\openpose\Release")
os.add_dll_directory(cwd / r"opencv_processor\3rdparty\openpose\build_windows\x64\Release")
os.add_dll_directory(cwd / r"opencv_processor\3rdparty\openpose\build_windows\bin")
try:
    import pyopenpose as openpose
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

assets_path = os.path.join(os.getcwd(), 'opencv_processor', 'assets')

faceonimg = cv.imread(os.path.join(assets_path, 'faceon.jpg'))
faceonimg = cv.resize(faceonimg, (0,0), fx=0.25, fy=0.25)
cv.imshow('Face On', faceonimg)
cv.waitKey(0)
cv.destroyAllWindows()

dtlimg = cv.imread(os.path.join(assets_path, 'dtl.jpg'))
dtlimg = cv.resize(dtlimg, (0,0), fx=0.25, fy=0.25)
cv.imshow('DTL', dtlimg)
cv.waitKey(0)
cv.destroyAllWindows()

params = dict()
params["model_folder"] = cwd / r"opencv_processor\3rdparty\openpose\models"

openposeWrapper = openpose.WrapperPython()
openposeWrapper.configure(params)
openposeWrapper.start()

datum = openpose.Datum()
datum.cvInputData = faceonimg
openposeWrapper.emplaceAndPop(openpose.VectorDatum([datum]))

print("Body keypoints: \n" + str(datum.poseKeypoints))

cv.imshow("OpenPose Output", datum.cvOutputData)
cv.waitKey(0)