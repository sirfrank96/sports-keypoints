# Sports-Keypoints

Application that allows a user to upload images and get data about keypoints relevant to the sport.
(As of now, sports-keypoints only has functionality for golf (specifically golf setup keypoints).)

## Description

The goal of the project was to make a simplified Sportsbox.ai application that would track golf setup points and save them for comparison. For example,
if you wanted to track your stance width and your knee bend, you could compare numbers from an image you uploaded a couple months ago to one from today.

The project uses Openpose (<https://github.com/CMU-Perceptual-Computing-Lab/openpose.git>) as the main engine. The Openpose Body25 model provides pose
estimation for 25 body parts (see protos/common.proto for a full list). The computervision-service uses the Openpose python API to run the Body25 model and provide
the other services pose estimation.

The server/go-server is the main manager for the application. It manages access to the computervision-service to grab pose data. It manages the MongoDb client
that stores users, images, and data about those images. It also is the server where the user will request the actual keypoints for different sports. The go-server uses gRPC (<https://grpc.io/docs/what-is-grpc/introduction/>)
to implement computervision APIs, user APIs, and user requests for keypoints APIs. See the protos folder for protobuf definitions and services for the listed APIs.

The clients/python-golf-client is an example client that uses the Tkinter GUI library to make it easy to create a user, login, select images, calibrate images, and calculate golf keypoints.

## Getting Started

The following are instructions on how to start the sports-keypoints services and connect the golf client application.

### Prerequisites

* Git: <https://git-scm.com/install/>
* Docker: <https://www.docker.com/get-started/>
* Python 3.10 (this is for the client gui application, see README in clients/python-golf-client for more information): <https://www.python.org/downloads/> 
* Optional: Machine with NVIDIA GPUs with 5+ GB of VRAM + Drivers + NVIDIA Container Toolkit. (this is for the computer vision service, see README in computervision-service for more information)

### Usage

1. Clone the repo:<br>
   `git clone https://github.com/sirfrank96/sports-keypoints.git`
2. Navigate to root directory of repo:<br>
   `cd sports-keypoints`
3. Sync Openpose and Openpose's 3rd party submodules:<br>
   `git submodule update --init --recursive`
4. Create .env file in root directory of the sports-keypoints repo:<br>
  * If Windows:<br>
      `echo > .env`
  * If Unix Shell (eg. bash or mac zsh):<br>
      `touch .env`
5. Add environment variables to .env file:<br>
  * If you have NVIDIA GPUs and want to use them:<br>
    ```
    PROCESSING_TYPE=gpu
    NUM_GPUS=all
    ```
  * If you just want to use CPU:<br>
    ```
    PROCESSING_TYPE=cpu
    ```

6a. _(If using GPU then pull archived NVIDIA image:<br>
   `docker pull nvcr.io/nvidia/cuda:10.1-cudnn7-devel-ubuntu18.04`<br> 
   If you have issues, you may need to create an NVIDIA NGC Catalog account: <https://catalog.ngc.nvidia.com/> and then `docker login nvcr.io` before pulling the image.)_

6. Spin up the sports-keypoints service containers with docker compose:<br>
   `docker compose up --build`
7. While that builds, in another window navigate to the client application:<br>
   `cd sports-keypoints/clients/python-golf-client`
8. Create a virtual environment:<br>
   `C:path\to\python310\python.exe -m venv python310_venv` (Mine was C:Users\UserA\AppData\Local\Programs\Python\Python310\python.exe on Windows)
9. Activate virtual environment:<br>
  * If Windows Command Prompt:<br>
      `python310_venv\Scripts\activate.bat`
  * If Windows Powershell:<br>
      `python310_venv\Scripts\Activate.ps1`
  * If Unix Shell (eg. bash or mac zsh):<br>
      `source python310_venv/bin/activate`
10. Install requirements:<br>
    `python -m pip install -r requirements.txt`
11. Navigate to main script:<br>
    `cd src`
12. Once the backend services are running, run the client application:<br>
    `python main.py`
13. See README in clients/python-golf-client for more details on how to use the client application

### Alternatives To Docker

If you do no not want to run Docker (overhead too large or running too slowly), you can run sports-keypoints without docker.

Prerequisites: 
* MongoDb (<https://www.mongodb.com/docs/manual/installation/>)
* Python 3.7 (<https://www.python.org/downloads/release/python-370/>)
* go 1.25.3 (<https://go.dev/doc/install>)
* Python 3.10 (<https://www.python.org/downloads/release/python-3100/>)

You will also need to follow the official openpose docs (<https://github.com/CMU-Perceptual-Computing-Lab/openpose/blob/master/doc/installation/0_index.md>) for installation prerequisites and how to build and compile for python api usage.

1. Start MongoDb:<br>
   `C:path\to\mongo\mongdb.exe` (Mine was C:\Program Files\MongoDB\Server\8.2\bin\mongodb.exe on Windows)
2. Start computervision service:<br>
   ```
   cd computervision-service
   C:path\to\python38\python.exe -m venv python38_venv
   python38_venv\Scripts\activate.bat
   python -m pip install --upgrade pip
   python -m pip install -r requirements.txt
   set PYTHONPATH=C:path\to\sports-keypoints\computervision-service\3rdparty\openpose\build_windows\python\openpose\Release
   cd src
   python server.py
   ```
4. Start go backend server:<br>
   ```
   cd server\go-server
   go run main.go
   ```
6. Start client application (see above in Usage for instructions)

## Future Todos
* Get keypoints for any point in golf swing, not just setup
* Additional keypoints (eg. elbow bend, arm/hand positions throughout swing)
* Allow upload video and save individual frames with processed data
* Keypoints for other sports
