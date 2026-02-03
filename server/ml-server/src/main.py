# Python
import logging

# GRPC
import grpc

# Internal
from common import cvclient
import training

# TODO: Make this start service that gets swing vidoes
# also start cv_client
def serve():
    setting_timeout_ms = 1000 * 60 * 3
    options = [('grpc.http2.settings_timeout', setting_timeout_ms)]
    channel = grpc.insecure_channel('localhost:50051', options=options)
    return cvclient.ComputerVisionClient(channel)

if __name__ == "__main__":
    logging.basicConfig()
    cv_client = cvclient.create_computer_vision_client()

    trainer = training.GolfMLTrainer(cv_client)
    trainer.train_classify_swing_model()

    #add cli so user can specify if want to train, infer, etc.
    #make this a UI loop where user can input stuff?
