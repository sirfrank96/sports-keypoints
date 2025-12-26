# go-server

The go-server is the main hub for the logic for the sports-keypoints repo. The go-server connects the implementation for the keypoints-server, the code for storing data in MongoDB, and the implementation of the computervision client. It also calculates the actual golf specific (or sport specific) keypoints.

The standard control flow of a sports-keypoints API call is as follows:
1. User makes a sport (golf) keypoints API request via a client application (eg. CalibrateInputImage, CalculateGolfKeypoints, etc.)
2. The go-server receives the request:
* The gRPC unary interceptor (keypoints-server/unary_interceptor.go) receives the request, verifies the session cookie, and puts the userId into the context that will passed around for the remainder of the process
* The keypoints server (keypoints-server/golf_keypoints_server.go) verifies the information in the request (makes sure required fields are set, etc.), and then forwards the request to the controller
* Controller receives the request (controller/golf_keypoints_listener)
  * Controller makes sure that the user exists in MongoDB
  * Controller executes logic to read necessary information from MongoDB, make requests to the computervision server for pose estimation, and calculate sport keypoints give the pose estimation points
* Controller sends the response back to keypoints server, which sends the response back to the user

## Directories

The following is a list of the folders within the go-server and their functions.

* controller:<br>
The central point of the go-server. Contains instances of a database manager, computervision client, and handles requests from the keypoints-server. Also contains logic for the calculation of golf setup points.

* cv-client:<br>
Implements the ComputerVisionServiceClient gRPC APIs. Make requests to the computervision service for pose estimation points.

* db:<br>
Contains code for CRUD MongoDB operations for users, input images, and keypoints for each input image. Also contains the struct definitions that are serialized into bson objects for MongoDB storage. Database operations are protected by a mutex handled by the DbManager.

* keypoints-server:<br>
Implements the UserServiceServer and GolfKeypointsServiceServer gRPC APIs. Is the first point of entry for users wanting to get keypoints for their image. Handles verification of session cookies and verification of requests coming in. 

* sports-keypoints-proto:<br>
Contains GoLang gRPC generated files containing client and server code from .proto files in the protos directory in the root directory of the sports-keypoints repo.

* util:<br>
Provides utility functions for the go-server. This includes a custom warning interface and struct that allows APIs to continue even if there is some missing information for only a specific part of the request. It also includes structs and vector math to help easily calculate keypoints given coordinate pose keypoints. 

## Getting Started

The following are instructions on how to start the go-server.

### Prerequisites

* Docker: <https://www.docker.com/get-started/>

### Usage

1. Start the computervision-service (see the README in sports-keypoints/computervision-service for more information)

2. Start the MongoDB and go-server containers via docker compose:<br>
`docker compose up --build`

### Alternative to Docker

Prerequisites
* go 1.25.3 (<https://go.dev/doc/install>)
* MongoDB (<https://www.mongodb.com/docs/manual/installation/>)

1. Start MongoDB:<br>
`C:path\to\mongo\mongdb.exe` (Mine was C:\Program Files\MongoDB\Server\8.2\bin\mongodb.exe on Windows)<br>

2. Start go-server:
`go run main.go`

## Future Todos

- Handle video API requests
- Extend to other sports, like baseball
- Rate limit APIs
