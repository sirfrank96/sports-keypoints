from concurrent import futures
import logging
import openpose

import grpc
import opencvandpose_pb2
import opencvandpose_pb2_grpc
import common_pb2
import common_pb2_grpc

#TODO: Rename python module to opencv_service/computervision_service
# Implement rpcs
class OpenCVAndPoseServiceServicer(opencvandpose_pb2_grpc.OpenCVAndPoseServiceServicer):
    def __init__(self):
        self.__init__

    def GetOpenPoseImage(self, request_iterator, context):
        print("GetOpenPoseImage grpc request")
        image = bytes()
        for new_image in request_iterator:
            image += new_image.bytes
        processedImg = openpose.getOpenPoseImageFromBytes(image)
        print("GetOpenPoseImage grpc request processed")
        print(f"ProcessedImg size is {len(processedImg)}")
        openPoseImage = common_pb2.Image(
            name="OpenPose Image Processed",
            bytes=processedImg
        )
        yield openPoseImage

    def GetOpenPoseData(self, request_iterator, context):
        return super().GetOpenPoseData(request_iterator, context)
    
    def GetOpenPoseHandImage(self, request_iterator, context):
        return super().GetOpenPoseHandImage(request_iterator, context)
    
    def GetOpenPoseHandData(self, request_iterator, context):
        return super().GetOpenPoseHandData(request_iterator, context)
    
    def GetOpenPoseAllData(self, request_iterator, context):
        return super().GetOpenPoseAllData(request_iterator, context)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    opencvandpose_pb2_grpc.add_OpenCVAndPoseServiceServicer_to_server(
        OpenCVAndPoseServiceServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Waiting for requests at port 50051")
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()