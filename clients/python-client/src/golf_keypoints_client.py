import golfkeypoints_pb2
import golfkeypoints_pb2_grpc

#client stub functions for golfkeypoints grpc service
class GolfKeypointsClient():
    def __init__(self, channel):
        self.stub = golfkeypoints_pb2_grpc.GolfKeypointsServiceStub(channel)

    def upload_input_image(self, session_token, image_type, image):
        request = golfkeypoints_pb2.UploadInputImageRequest(session_token=session_token, image_type=image_type, image=image)
        return self.stub.uploadInputImage(request)
    
    def list_input_images_for_user(self, session_token):
        request = golfkeypoints_pb2.ListInputImagesForUserRequest(session_token=session_token)
        return self.stub.listInputImagesForUser(request)
    
    def read_input_image(self, session_token, input_image_id):
        request = golfkeypoints_pb2.ReadInputImageRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.readInputImage(request)
    
    def delete_input_image(self, session_token, input_image_id):
        request = golfkeypoints_pb2.DeleteGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.deleteGolfKeypoints(request)
    
    def calibrate_input_image(self, session_token, input_image_id, calibration_type, feet_line_method, calibration_image_axes, calibration_image_vanishing_point, golf_ball, club_butt, club_head):
        request = golfkeypoints_pb2.CalibrateInputImageRequest(session_token=session_token, input_image_id=input_image_id, calibration_type=calibration_type, feet_line_method=feet_line_method, calibration_image_axes=calibration_image_axes, calibration_image_vanishing_point=calibration_image_vanishing_point, golf_ball=golf_ball, club_butt=club_butt, club_head=club_head)
        return self.stub.calibrateInputImage(request)
    
    def calculate_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.CalculateGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.calculateGolfKeypoints(request)
    
    def read_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.ReadGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.calculateGolfKeypoints(request)
    
    def update_body_keypoints(self, session_token, input_image_id, updated_body_keypoints):
        request = golfkeypoints_pb2.UpdateBodyKeypointsRequest(session_token=session_token, input_image_id=input_image_id, updated_body_keypoints=updated_body_keypoints)
        return self.stub.updateBodyKeypoints(request)
    
    def delete_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.DeleteGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.deleteGolfKeypoints(request)
    