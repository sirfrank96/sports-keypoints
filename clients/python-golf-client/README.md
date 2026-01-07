# Python Golf Client

Client application that uses the golf keypoints client and connects to the sports-keypoints/server/go-server. 

## Description

This application implements the gRPC client for GolfKeypointsService and UserService in sports-keypoints/protos. Using the Python Tkinter GUI library, it allows the user
to register, login, and make changes to their user profile. It will then allow the user to select images to upload, calibrate, and the calculate golf keypoints.

The application takes either Face On or DTL (down the line) golf setup images (.jpg/jpeg files). Golf setup images will probably just be screenshots you take in a swing video right before you start your swing.

In addition to the setup image, you must calibrate your image. Because the golf setup images are 2D images in 3D space, so information is lost. So, in order to get keypoints such as alignments, bends, etc. that are relative to your 3D space, you must provide additional information. You can either draw lines that will give horizontal/vertical axes and vanishing point lines, or submit calibration images. 

## Getting Started

The following are instructions on how to start the python golf client application.

### Prerequisites

* Python 3.10 (<https://www.python.org/downloads/release/python-3100/>)

### Usage

1. Make sure an instance of the go-server is running (see README in the root directory of this repo to get started).
2. Create a virtual environment:<br>
   `C:path\to\python310\python.exe -m venv python310_venv` (Mine was C:\Users\UserA\AppData\Local\Programs\Python\Python310\python.exe on Windows)
3. Activate virtual environment:<br>
  * If Windows Command Prompt:<br>
      `python310_venv\Scripts\activate.bat`
  * If Windows Powershell:<br>
      `python310_venv\Scripts\Activate.ps1`
  * If Unix Shell (eg. bash or mac zsh):<br>
      `source python310_venv/bin/activate`
4. Install requirements:<br>
    `python -m pip install -r requirements.txt`
5. Navigate to main script:<br>
    `cd src`
6. Once the backend services are running, run the client application:<br>
    `python main.py`
7. Once the application launches, register a user by typing in a username, password, and email
8. Login using the same username and password
9. To, upload an image for the first time, press the `Select New Input Image` button. This will launch another window that will allow you to browse your filesystem for images.<br>
_(If you are uploading both a Face On and DTL image, it is recommended to upload the Face On first, because DTL can use additional information provided from the Face On results)_
10. Once you select the image you want, a popup will appear asking you if this image is Face On or DTL (Down the line). Press `Yes` for Face On, press `No` for DTL. <br>    
<img width="445" height="187" alt="dtl-or-faceon-screenshot" src="https://github.com/user-attachments/assets/17dc6ddb-eefd-4960-9fc6-283d42cadad3" />

11. Next input a description for the input image. Recommended to add what club is being used, faceon or dtl, and any other useful identifying descriptions. For example: <br>
<img width="528" height="142" alt="input-img-description-screenshot" src="https://github.com/user-attachments/assets/be2aeb71-f02a-474a-9886-4e654879f4b3" />

12. Once the image displays in the center, you must enter additional information to calibrate it. Go down the list of buttons on the right side and click and follow
instructions for each. Once each is done, the button will be disabled and grayed out.
13. Once necessary buttons are completed, press the `Calibrate Image` button. You will see a prompt asking how you want to calibrate the image. Press `Yes` if you have calibration images saved (see below on how to take those images). Press `No` if you want to draw axes and vanishing point lines yourself to calibrate (drawing is probably more accurate).<br>
<img width="511" height="197" alt="press-calibrate-button-screenshot" src="https://github.com/user-attachments/assets/d241208b-e1a6-40ef-86a4-e851246649a8" />

14. Once a successful response comes back for image calibration, press the `Calculate Golf Keypoints` button.
15. Take a look at the data about setup keypoints and pose estimation points.
16. Once you close the the window with that data, you will get a prompt that asks you if any of the pose estimation points are inaccurate. Click `Yes` if you would
like to modify those and recalculate, press `No` otherwise.
<img width="477" height="182" alt="body-keypoints-incorrect-prompt-screenshot" src="https://github.com/user-attachments/assets/d9582759-db65-4ff1-9654-7d4fa16e539c" />

18. If you pressed `Yes`, you will see another window that lists all of the pose estimation points and their coordinates in the image. Click all buttons for body points
that you would like to modify. (The coordinate plane origin is the top left corner of the image. As you go right, the x value increases. As you go down, the y value increases).
19. Once you have modified all pose estimation points, scroll to the bottom of that window, and press the `Done Updating Body Keypoints` button to recalculate.
20. Take a look at the new data about setup keypoints

#### Drawing Lines for Calibration

When prompted to draw a line, you will click on the image where you want the line to start and release where you want the line to end. The app will show the line you drew and ask if it is correct. If so, continue. If not, redraw the line.

Both Face On and DTL images requires you to draw the horizontal and vertical axes.
* The horizontal axis is any line running from the left to the right side of the image, parallel to the ground
* The vertical axis is any line running from the top to the bottom of the image, perpendicular to the ground

DTL images require an additional two "vanishing point" lines that identify where the camera is pointing/the target. The app needs two lines, in order to calculate the vanishing point. These lines could be:
* A line from the ball to the target (should look pretty vertical)
* A line along your feet line (if your feet are pointing at the target).
* The side edges of a mat if you are hitting off of mats.
Remember that parallel lines converge in the the distance, so these lines should look like they are converging to each other as they get closer to the top of the image.

Here is an example image of the four lines that should be drawn:<br>
<img width="823" height="363" alt="example-draw-lines-screenshot" src="https://github.com/user-attachments/assets/01ac0f9c-29dc-42b5-a30d-0bf95f3490ff" />

#### Taking Images for Calibration

Calibration images should be taken from the same camera position as the golf setup image.<br>

Both Face On and DTL images require an "Axes Calibration Image". This image provides vertical and horizontal axes, relative to the ground (ie. the horizontal axis is 
parallel to the ground, and the vertical axis is perpendicular to the ground). To take this image, do the following:<br>
* From where the ball is, stand facing the camera, as centered horizontally in the frame as possible.
* Stand upright and straight with feet wider than shoulder width perpendicular to the camera angle.
  * Face On: Align the camera perpendicular to the target line, Feet should run parallel to the target line.
  * DTL: Align the camera to point at the target, Feet should run perpendicular to the target line.
 
DTL images require an additional "Vanishing Point Calibration Image). This image provides information to identify where the camera is pointing/the target 
(ie. the vanishing point), in order to track alignments. To take this image, do the following:<br>
* Align the center of the camera frame to point at the target
* Take your stance (shoulder width or wider) so that your feet are off center from camera frame (ie. from left of the center of the frame, it will look like your feet are pointing to the right)
* Aim your feet parallel to target line (recommended to use an alignment stick) (parallel lines converge in the distance, so probably right at the target/1 foot away)
* Make sure either both your heels or both your toes are visible in image

Save these images in your filesystem where they can be easily accessed.

## Future Todos

1. Make the application look nicer (right now it is pretty bare bones)
2. Add zoom functionality to more easily identify where pose keypoints are
3. Allow user to click on points in the input image when updating pose estimation points, and update image to reflect new pose estimation points
4. Add a back button to go back and forth between frames
5. Clean up the code (pull out common functions, move code out of main_page.py and into separate modular files)
