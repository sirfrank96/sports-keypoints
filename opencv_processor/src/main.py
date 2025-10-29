import cv2 as cv
import os
import sys
pyopenpose_path = os.path.join(os.getcwd(), 'opencv_processor', 'src', 'openpose', 'python', 'openpose')
print("1: " + pyopenpose_path)
sys.path.append(pyopenpose_path)
import pyopenpose as openpose

assets_path = os.path.join(os.getcwd(), 'opencv_processor', 'assets')

#params = dict()
#params["model_folder"] = "openpose/models/"

#openposeWrapper = openpose.WrapperPython()
#openposeWrapper.configure(params)
#openposeWrapper.start()

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


#datum = op.Datum()
#datum.cvInputData = faceonimg
#openposeWrapper.emplaceAndPop([datum])

#print("Body keypoints: \n" + str(datum.poseKeypoints))

#cv.imshow("OpenPose Output", datum.cvOutputData)
#cv.waitKey(0)