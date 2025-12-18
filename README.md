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
that stores users, images, and data about those images. It also is the server where the user will request the actual keypoints for different sports. The go-server uses gRPC (<https://grpc.io/docs/what-is-grpc/introduction/)
to implement computervision APIs, user APIs, and user requests for keypoints APIs. See the protos folder for protobuf definitions and services for the listed APIs.

The clients/python-golf-client is an example client that uses Tkinter GUI library to make it easy to create a user, login, select images, calibrate images, and calculate golf keypoints.

## Getting Started

The following are instructions on how to start the sports-keypoints services and connect the golf client application.

### Prerequisites

* Git: <https://git-scm.com/install/>
* Docker: <https://www.docker.com/get-started/>
* Python 3.10 (this is for the client gui application, see README in clients/python-golf-client for more information): <https://www.python.org/downloads/> 
* Optional: Machine with NVIDIA GPUs with 10+ GB of VRAM + Drivers + NVIDIA Container Toolkit (this is for the computer vision service, see README in computervision-service for more information)

### Usage

1. Clone the repo:
  `git clone https://github.com/sirfrank96/sports-keypoints.git`
2. Create .env file in root directory of the sports-keypoints repo and add environment variables:
  * If you do not have NVIDIA GPUs:
    ```
    #.env file
    PROCESSING_TYPE=cpu
    ```
  * If you do have NVIDIA GPUs and want to use them:
    ```
    #.env file
    PROCESSING_TYPE=gpu
    NUM_GPUS=all
    ```
3. Spin up the sports-keypoints service containers with docker compose:
   `docker compose up --build`
4. Navigate to client application:
   `cd clients/python-golf-client`
5. Create a virtual environment:
   `C:path\to\python310\python.exe -m venv python310_venv`
6. Activate virtual environment:
  * If Windows Command Prompt:
      `python310_venv\Scripts\activate.bat`
  * If Windows Powershell:
      `python310_venv\Scripts\Activate.ps1`
  * If Unix Shell (eg. bash):
      `source python310_venv/bin/activate`
7. Install requirements:
    `python -m pip install -r requirements.txt`
8. Navigate to main script:
    `cd src`
9. Run client application:
    `python main.py`
10. See README in clients/python-golf-client for more details on how to use the client application

### Alternatives To Docker

If you do no not want to run Docker (overhead too large or running too slowly), you can run sports-keypoints without docker.

Prerequisites: MongoDb (<https://www.mongodb.com/docs/manual/installation/>), Python 3.8 (<https://www.python.org/downloads/release/python-380/>), go 1.25.3 (<https://go.dev/doc/install>), Python 3.10 (<https://www.python.org/downloads/release/python-3100/>)

You will also need to follow the official openpose docs (<https://github.com/CMU-Perceptual-Computing-Lab/openpose/blob/master/doc/installation/0_index.md>) for installation prerequisites and how to build and compile for python api usage.

1. Start MongoDb:
   `C:path\to\mongo\mongdb.exe` (Mine was C:\Program Files\MongoDB\Server\8.2\bin\mongodb.exe)
2. Start computervision-service:
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
4. Start go backend server
   `cd server\go-server`
   `go run main.go`
6. Start client application (see above in Usage for instructions)
