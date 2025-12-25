# Computer Vision Service

Wrapper for Computer Vision Python libraries. Provides a simple interface for a client to request and receive computer vision data for an image.

## Description

The openpose.py provides a wrapper for the CMU OpenPose Python library (<https://github.com/CMU-Perceptual-Computing-Lab/openpose>). For now, the computervision-service 
only provides pose estimation. The APIs are accessible via gRPC and defined in sports-keypoints/protos/computervision.proto.<br>

Because pose estimation is a compute intensive task, the OpenPose library allows users to build the library to utilize GPUs. In particular, the library was tested with 
NVIDIA GPUs (AMD GPUs require OpenCL and may not work) (<https://github.com/CMU-Perceptual-Computing-Lab/openpose/blob/master/doc/installation/1_prerequisites.md>). 
This wrapper allows the user to run with GPU or with CPU only, depending on the user's system.

The gRPC ComputerVisionService provides APIs for submitting images and then receiving BODY25 model data about the image. This comes in 2 forms:
1. The actual image wih all of the BODY25 pose estimation points marked and lines drawn between to show the "skeleton"
2. A full gRPC message (Body25PoseKeypoints) where each field corresponds to a single BODY25 keypoint (eg. nose, left shoulder, etc.)

## Getting Started

The following are instructions on how to start the computervision-service.

### Prerequisites

* Docker: <https://www.docker.com/get-started/>
* Optional: If you want to use GPU:
*   NVIDIA GPU with Drivers >= 418.XX (<https://docs.nvidia.com/datacenter/tesla/pdf/NVIDIA_CUDA_Drivers_Support.pdf>) for CUDA 10.1
*   GPU should have 5+ GB of VRAM (run `nvidia-smi` to see available memory)
*   If you do not have Docker Desktop, you probably have to install the NVIDIA Container Toolkit to access your GPUs from the docker container (<https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/latest/install-guide.html>)  

### Usage

1. Depending on if you are using GPU or CPU only, create a Docker image from the appropriate Dockerfile
* GPU: `docker build -f Dockerfile.gpu -t computervision-gpu-image .`
* CPU-only: `docker build -f Dockerfile.cpu -t computervision-cpu-image .`
2. Spin up the Docker container from the image that you just created
* GPU: `docker run -p 50051:50051 --gpus all computervision-gpu-image:latest`
* CPU-only: `docker run -p 50051:50051 computervision-cpu-image:latest`

### Alternative To Docker

If you do no not want to run Docker (overhead too large, running too slowly, or want to more easily test locally), you can run the computervision-service without docker.
However, this requires a lot of work up front to build and compile the OpenPose Python library.

Prerequisites:
* CMake
* Python 3.7 (<https://www.python.org/downloads/release/python-370/>)

1. Create and start python virtual environment:<br>
  ```
  C:path\to\python38\python.exe -m venv python38_venv
  python38_venv\Scripts\activate.bat
  ```
2. Install dependencies:<br>
  ```
  python -m pip install --upgrade pip
  python -m pip install -r requirements.txt
  ```
3. Follow OpenPose installation instructions (<https://github.com/CMU-Perceptual-Computing-Lab/openpose/blob/master/doc/installation/0_index.md>) to build for and enable the Python API.
4. Allow Python to find created library binaries. Find where the binary is (.so file for Linux systems (eg. pyopenpose.cpython-36m-x86_64-linux-gnu.so), .pyd file for windows (eg. pyopenpose.cp37-win_amd64.pyd)),
  and set the PYTHONPATH to it. In Windows this looked like this:<br>
`set PYTHONPATH=C:path\to\sports-keypoints\computervision-service\3rdparty\openpose\build_windows\python\openpose\Release`
5. Navigate to where wrapper scripts are:<br>
`cd src`
6. Run server.py
`python server.py`

## Future Todos

* Implement APIs for hand estimation in order to be able to points for grip keypoints
* Add OpenCV wrapper APIs for object detection (golf ball, golf club, etc.)
* 3D reconstruction from multiple 2 images (ie. face on and dtl to form a 3D rendering)
* Use a more modern pose estimation library, so don't have to rely on legacy dependencies (eg. Caffe)
