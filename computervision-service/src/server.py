from concurrent import futures
import logging
import openpose

import grpc
import computervision_pb2
import computervision_pb2_grpc
import common_pb2

# TODO: Implement all rpcs
class ComputerVisionServiceServicer(computervision_pb2_grpc.ComputerVisionServiceServicer):

    body_pose_field_descriptors = common_pb2.Body25PoseKeypoints.DESCRIPTOR.fields

    def __init__(self):
        super().__init__()
        self.open_pose_mgr = openpose.OpenPoseManager()

    def GetPoseImage(self, request, context):
        print("GetPoseImage grpc request")
        image_bytes = request.image
        image = self.open_pose_mgr.get_image_from_bytes(image_bytes)
        datum = self.open_pose_mgr.run_open_pose(image)
        processed_img = self.open_pose_mgr.get_open_pose_image(datum, self.body_pose_field_descriptors)
        print(f"Processed image. It's size is {len(processed_img)}")
        print("GetPoseImage grpc request finished")
        return computervision_pb2.GetPoseImageResponse(
            success=True,
            image=processed_img
        )
    
    def processOpenPoseData(self, data):
        body_25_pose_keypoints = common_pb2.Body25PoseKeypoints()
        body_25_pose_keypoints_descriptor = body_25_pose_keypoints.DESCRIPTOR
        for field_descriptor in body_25_pose_keypoints_descriptor.fields:
            field_name = field_descriptor.name
            field_number = field_descriptor.number
            keypoint = common_pb2.Keypoint(
                x=data[field_number-1][0],
                y=data[field_number-1][1],
                confidence=data[field_number-1][2]
            )
            curr_field = getattr(body_25_pose_keypoints, field_name)
            curr_field.CopyFrom(keypoint)
        return body_25_pose_keypoints
    
    def GetPoseData(self, request, context):
        print("GetPoseData grpc request")
        image_bytes = request.image
        image = self.open_pose_mgr.get_image_from_bytes(image_bytes)
        datum = self.open_pose_mgr.run_open_pose(image)
        data = self.open_pose_mgr.get_open_pose_data(datum)
        print(f"Processed image. Data is {data}. Length of data is {len(data)}")
        body_25_pose_keypoints = self.processOpenPoseData(data)
        print(f"Converted data array to Keypoints {body_25_pose_keypoints}")
        print("GetPoseData grpc request finished")
        return computervision_pb2.GetPoseDataResponse(
            success=True,
            keypoints=body_25_pose_keypoints
        )
    
    def GetPoseHandImage(self, request, context):
        return super().GetOpenPoseHandImage(request, context)
    
    def GetPoseHandData(self, request, context):
        return super().GetOpenPoseHandData(request, context)
    
    def GetPoseAll(self, request, context):
        print("GetPoseAll grpc request")
        # run openpose
        image_bytes = request.image
        image = self.open_pose_mgr.get_image_from_bytes(image_bytes)
        datum = self.open_pose_mgr.run_open_pose(image)
        # get image
        processed_img = self.open_pose_mgr.get_open_pose_image(datum, self.body_pose_field_descriptors)
        # get data
        data = self.open_pose_mgr.get_open_pose_data(datum)
        body_25_pose_keypoints = self.processOpenPoseData(data)
        print(f"GetPoseall converted data array to Keypoints {body_25_pose_keypoints}")
        print("GetPoseAll grpc request finished")
        return computervision_pb2.GetPoseAllResponse(
            success=True,
            image=processed_img,
            pose_keypoints=body_25_pose_keypoints
        )
    
    def GetPoseImagesFromVideo(self, request_iterator, context):
        print("GetPoseImagesFromVideo grpc request")
        img_idx = 0
        for get_open_pose_image_request in request_iterator:
            img_idx += 1
            image_bytes = get_open_pose_image_request.image
            image = self.open_pose_mgr.get_image_from_bytes(image_bytes)
            processed_img = self.open_pose_mgr.get_open_pose_image(image, self.body_pose_field_descriptors)
            print(f"Processed image #{img_idx}. It's size is {len(processed_img)}")
            get_open_pose_image_response = computervision_pb2.GetPoseImageResponse(
                image=processed_img
            )
            yield get_open_pose_image_response
        print("GetPoseImagesFromVideo grpc request finished")
    
    def GetPoseDataFromVideo(self, request_iterator, context):
        return super().GetPoseDataFromVideo(request_iterator, context)
    
    def GetPoseHandImagesFromVideo(self, request_iterator, context):
        return super().GetPoseHandImagesFromVideo(request_iterator, context)
    
    def GetPoseHandDataFromVideo(self, request_iterator, context):
        return super().GetPoseHandDataFromVideo(request_iterator, context)
    
    def GetPoseAllFromVideo(self, request_iterator, context):
        return super().GetPoseAllFromVideo(request_iterator, context)
    

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    computervision_pb2_grpc.add_ComputerVisionServiceServicer_to_server(
        ComputerVisionServiceServicer(), server
    )
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Waiting for computervision requests at port 50051")
    print("Waiting for sigint to stop services")
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig()
    serve()
