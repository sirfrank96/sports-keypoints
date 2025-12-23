import cv2 as cv
import os
import traceback
import sys
import numpy as np
from pathlib import Path
import platform

# Import pyopenpose, make sure sys can find paths necessary dlls and binaries
curr_dir = Path(__file__).parent.resolve()
isWindows = False
try:
    system = platform.system()
    if system == "Windows":
        isWindows = True
        sys.path.append(curr_dir / r"..\3rdparty\openpose\build_windows\python\openpose\Release")
        os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\x64\Release")
        os.add_dll_directory(curr_dir / r"..\3rdparty\openpose\build_windows\bin")
        print(curr_dir)
    elif os == "Linux":
        sys.path.append('/usr/local/python/openpose')
except ImportError as e:
    print('Error getting platform system')

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

class OpenPoseManager():

    def __init__(self):
        super().__init__()
        # Location of OpenPose models
        params = dict()
        params["number_people_max"] = 1
        params["face"] = False
        params["hand"] = False
        if isWindows:
            params["model_folder"] = curr_dir / r"..\3rdparty\openpose\models"
        else:
            params["model_folder"] = curr_dir / r"../3rdparty/openpose/models"
        # Get OpenPose python wrapper
        self.openpose_wrapper = op.WrapperPython()
        self.openpose_wrapper.configure(params)
        self.openpose_wrapper.start()

    # Helper function: converts bytes object to Matlike
    def get_image_from_bytes(self, img_bytes):
        nparr = np.frombuffer(img_bytes, np.uint8)
        img = cv.imdecode(nparr, cv.IMREAD_COLOR)
        return img

    # Helper function: Grabs img from path and converts to Matlike
    def get_image_from_path(self, imgPath):
        img = cv.imread(imgPath, cv.IMREAD_COLOR)
        return img

    def run_open_pose(self, img):
        datum = op.Datum()
        datum.cvInputData = img
        self.openpose_wrapper.emplaceAndPop(op.VectorDatum([datum]))
        return datum

    # Input is Matlike, Feed into OpenPose, Return processed image as bytes
    def get_open_pose_image(self, datum, field_descriptors):
        # Add labels
        output_img = datum.cvOutputData
        keypoints = datum.poseKeypoints
        if keypoints.shape != (): # Check if keypoints are detected
            idx = 0
            curr_keypoints = keypoints[0]
            for keypoint in curr_keypoints:
                pos = (int(keypoint[0]), int(keypoint[1]))
                if pos[0] > 0 and pos[1] > 0: # Check if point is valid
                    cv.putText(output_img, field_descriptors[idx].name, pos, cv.FONT_HERSHEY_SIMPLEX, 1.0, (255, 255, 255), 1, cv.LINE_AA)
                idx += 1
        # Encode data as jpg, then convert to bytes object
        img_encode = cv.imencode('.jpg', output_img)
        data_encode = np.array(img_encode[1])
        byte_encode = data_encode.tobytes()   
        return byte_encode

    def get_open_pose_data(self, datum):
        return datum.poseKeypoints[0]        
