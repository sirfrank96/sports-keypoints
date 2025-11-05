from concurrent import futures
import logging
import openpose

import grpc
import opencvandpose_pb2
import opencvandpose_pb2_grpc
import common_pb2
import common_pb2_grpc

#TODO: Implement all rpcs
class OpenCVAndPoseServiceServicer(opencvandpose_pb2_grpc.OpenCVAndPoseServiceServicer):
    def __init__(self):
        self.__init__

    def GetOpenPoseImage(self, request, context):
        print("GetOpenPoseImage grpc request")
        image = request.image.bytes
        processed_img = openpose.get_open_pose_image_from_bytes(image)
        print(f"Processed image. It's size is {len(processed_img)}")
        open_pose_image = common_pb2.Image(
            name=f"Processed image",
            bytes=processed_img
        )
        print("GetOpenPoseImage grpc request finished")
        return opencvandpose_pb2.GetOpenPoseImageResponse(
            image=open_pose_image
        )
    
    def GetOpenPoseData(self, request, context):
        return super().GetOpenPoseData(request, context)
    
    def GetOpenPoseHandImage(self, request, context):
        return super().GetOpenPoseHandImage(request, context)
    
    def GetOpenPoseHandData(self, request, context):
        return super().GetOpenPoseHandData(request, context)
    
    def GetOpenPoseAll(self, request, context):
        return super().GetOpenPoseAll(request, context)
    
    def GetOpenPoseImagesFromVideo(self, request_iterator, context):
        print("GetOpenPoseImagesFromVideo grpc request")
        img_idx = 0
        for get_open_pose_image_request in request_iterator:
            img_idx += 1
            image = get_open_pose_image_request.image.bytes
            processed_img = openpose.get_open_pose_image_from_bytes(image)
            print(f"Processed image #{img_idx}. It's size is {len(processed_img)}")
            open_pose_image = common_pb2.Image(
                name=f"Processed image #{img_idx}",
                bytes=processed_img
            )
            get_open_pose_image_response = opencvandpose_pb2.GetOpenPoseImageResponse(
                image=open_pose_image
            )
            yield get_open_pose_image_response
        print("GetOpenPoseImagesFromVideo grpc request finished")
    
    def GetOpenPoseDataFromVideo(self, request_iterator, context):
        return super().GetOpenPoseDataFromVideo(request_iterator, context)
    
    def GetOpenPoseHandImagesFromVideo(self, request_iterator, context):
        return super().GetOpenPoseHandImagesFromVideo(request_iterator, context)
    
    def GetOpenPoseHandDataFromVideo(self, request_iterator, context):
        return super().GetOpenPoseHandDataFromVideo(request_iterator, context)
    
    def GetOpenPoseAllFromVideo(self, request_iterator, context):
        return super().GetOpenPoseAllFromVideo(request_iterator, context)
    

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    opencvandpose_pb2_grpc.add_OpenCVAndPoseServiceServicer_to_server(
        OpenCVAndPoseServiceServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Waiting for opencv and openpose requests at port 50051")
    print("Waiting for sigint to stop services")
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()