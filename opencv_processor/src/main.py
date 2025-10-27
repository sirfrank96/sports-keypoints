import cv2 as cv
import os

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