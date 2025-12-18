import golfkeypoints_pb2
import golfkeypoints_pb2_grpc

# client stub functions for golfkeypoints grpc service
class GolfKeypointsClient():
    def __init__(self, channel):
        self.stub = golfkeypoints_pb2_grpc.GolfKeypointsServiceStub(channel)

    def upload_input_image(self, session_token, image_type, image):
        request = golfkeypoints_pb2.UploadInputImageRequest(session_token=session_token, image_type=image_type, image=image)
        return self.stub.UploadInputImage(request)
    
    def list_input_images_for_user(self, session_token):
        request = golfkeypoints_pb2.ListInputImagesForUserRequest(session_token=session_token)
        return self.stub.ListInputImagesForUser(request)
    
    def read_input_image(self, session_token, input_image_id):
        request = golfkeypoints_pb2.ReadInputImageRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.ReadInputImage(request)
    
    def delete_input_image(self, session_token, input_image_id):
        request = golfkeypoints_pb2.DeleteInputImageRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.DeleteInputImage(request)
    
    def calibrate_input_image(self, session_token, input_image_id, calibration_type, feet_line_method, calibration_image_axes, calibration_image_vanishing_point, golf_ball, club_butt, club_head, shoulder_tilt):
        request = golfkeypoints_pb2.CalibrateInputImageRequest(session_token=session_token, input_image_id=input_image_id, calibration_type=calibration_type, feet_line_method=feet_line_method, calibration_image_axes=calibration_image_axes, calibration_image_vanishing_point=calibration_image_vanishing_point, golf_ball=golf_ball, club_butt=club_butt, club_head=club_head, shoulder_tilt=shoulder_tilt)
        return self.stub.CalibrateInputImage(request)
    
    def calculate_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.CalculateGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.CalculateGolfKeypoints(request)
    
    def read_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.ReadGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.ReadGolfKeypoints(request)
    
    def update_body_keypoints(self, session_token, input_image_id, updated_body_keypoints):
        request = golfkeypoints_pb2.UpdateBodyKeypointsRequest(session_token=session_token, input_image_id=input_image_id, updated_body_keypoints=updated_body_keypoints)
        return self.stub.UpdateBodyKeypoints(request)
    
    def delete_golf_keypoints(self, session_token, input_image_id):
        request = golfkeypoints_pb2.DeleteGolfKeypointsRequest(session_token=session_token, input_image_id=input_image_id)
        return self.stub.DeleteGolfKeypoints(request)
    