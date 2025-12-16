import cv2 as cv
import os
import traceback
import sys
from pathlib import Path

curr_dir = Path(__file__).parent.resolve()
sys.path.append(curr_dir / r"..\3rdparty\openpose\build_windows\python\openpose\Release")
os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\x64\Release")
os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\bin")
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
    traceback.print_exc()


# Standard CV read and display image
assets_path = curr_dir / r"..\assets"

faceonimg = cv.imread(os.path.join(assets_path, 'faceon.jpg'))
faceonimg = cv.resize(faceonimg, (0,0), fx=0.5, fy=0.5)
cv.imshow('Face On', faceonimg)
cv.waitKey(0)
cv.destroyAllWindows()

dtlimg = cv.imread(os.path.join(assets_path, 'dtl.jpg'))
dtlimg = cv.resize(dtlimg, (0,0), fx=0.5, fy=0.5)
cv.imshow('DTL', dtlimg)
cv.waitKey(0)
cv.destroyAllWindows()


# Standard 2D openpose
params = dict()
params["model_folder"] = curr_dir / r"..\3rdparty\openpose\models"
params["hand"] = True

openposeWrapper = openpose.WrapperPython()
openposeWrapper.configure(params)
openposeWrapper.start()

datumFaceOn = openpose.Datum()
datumFaceOn.cvInputData = faceonimg
openposeWrapper.emplaceAndPop(openpose.VectorDatum([datumFaceOn]))

print("Body keypoints: \n" + str(datumFaceOn.poseKeypoints))
print("Hand keypoints: \n" + str(datumFaceOn.handKeypoints))

cv.imshow("OpenPose Output", datumFaceOn.cvOutputData)
cv.waitKey(0)

#TODO: 3D for multiple images (very difficult probably), update render threshold, write json?