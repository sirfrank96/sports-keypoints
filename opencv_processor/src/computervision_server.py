from concurrent import futures
import logging
import openpose

import grpc
import computervision_pb2
import computervision_pb2_grpc


class ComputerVisionServicer(computervision_pb2_grpc.ComputerVisionServicer):
    def __init__(self):
        self.__init__

    def GetOpenPoseDTLImage(self, request_iterator, context):
        print("DTL grpc request")
        image = bytes()
        for new_image in request_iterator:
            image += new_image.image
        processedImg = openpose.getOpenPoseImageFromBytes(image)
        print("DTL grpc request processed")
        dtlImage = computervision_pb2.DTLImage(
            name="DTL Processed",
            image=processedImg
        )
        yield dtlImage
    
    def GetOpenPoseFaceOnImage(self, request_iterator, context):
        print("FaceOn grpc request")
        image = bytes()
        for new_image in request_iterator:
            image += new_image.image
        processedImg = openpose.getOpenPoseImageFromBytes(image)
        print("FaceOn grpc request processed")
        print(f"ProcessedImg size is {len(processedImg)}")
        faceOnImage = computervision_pb2.FaceOnImage(
            name="FaceOn Processed",
            image=processedImg
        )
        yield faceOnImage

    #def GetOpenPoseDTLImage(self, request, context):
    #    print("DTL grpc request")
    #    processedImg = openpose.getOpenPoseImageFromBytes(request.image)
    #    print("DTL grpc request returning")
    #    return computervision_pb2.DTLImage(
    #        name="DTL Processed",
    #        image=processedImg
    #    )
    
    #def GetOpenPoseFaceOnImage(self, request, context):
    #    print("Face on grpc request")
    #    processedImg = openpose.getOpenPoseImageFromBytes(request.image)
    #    print("Face on grpc request returning")
    #    return computervision_pb2.FaceOnImage(
    #        name="FaceOn Processed",
    #        image=processedImg
    #    )


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    computervision_pb2_grpc.add_ComputerVisionServicer_to_server(
        ComputerVisionServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Waiting for requests at port 50051")
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()