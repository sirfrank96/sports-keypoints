import user_pb2
import user_pb2_grpc

#client stub functions for user grpc service
#TODO: Email
class UserClient():
    def __init__(self, channel):
        self.stub = user_pb2_grpc.UserServiceStub(channel)

    def create_user(self, username, password):
        request = user_pb2.CreateUserRequest(user_name=username, password=password, email="blah@gmail.com")
        return self.stub.createUser(request)
    
    def register_user(self, username, password):
        request = user_pb2.RegisterUserRequest(user_name=username, password=password)
        return self.stub.registerUser(request)

    def read_user(self, session_token):
        request = user_pb2.ReadUserRequest(session_token=session_token)
        return self.stub.readUser(request)

    def update_user(self, session_token, username, password):
        request = user_pb2.UpdateUserRequest(session_token=session_token, user_name=username, password=password, email="blahmod@gmail.com")
        return self.stub.updateUser(request)

    def delete_user(self, session_token):
        request = user_pb2.DeleteUserRequest(session_token=session_token)
        return self.stub.deleteUser(request)
    