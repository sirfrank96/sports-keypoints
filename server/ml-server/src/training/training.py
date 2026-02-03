# Python
import numpy as np
import argparse
from enum import Enum

# Internal
from . import keras_models
from src.common import util
from src.common import constants

# TODO: Auto retrain when hit certain number of videos
class GolfMLTrainer():
    def __init__(self):
        pass

    def create_golf_swing_feature_extractor(self):
        # get body datapoints from frames
        all_datapoints = util.load_all_body_datapoints_from_file()
        # put as input to created autoencoder and pull out the encoder
        inputs, outputs = all_datapoints, all_datapoints
        x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
        autoencoder = keras_models.create_sequential_autoencoder(x_train, y_train, x_val, y_val, constants.FRAME_DIMENSION, constants.BODY_DATAPOINTS_DIMENSIONS) 
        feature_extractor = keras_models.get_encoding_layer_from_autoencoder(autoencoder)
        # save feature extractor
        feature_extractor.save(constants.FEATURE_EXTRACTOR_MODEL_PATH)

    def create_golf_swing_faceon_classifier(self):
        # get body datapoints from frames
        all_datapoints = util.load_all_body_datapoints_from_file()
        # grab true values on whether videos are face on or not
        faceon_indices = util.load_faceon_indices_from_file()
        # create sequential classifier for whether video/frame is face on or not
        inputs, outputs = all_datapoints, faceon_indices
        x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
        faceon_classifier = keras_models.create_sequential_classifier(x_train, y_train, x_val, y_val, constants.FRAME_DIMENSION, constants.BODY_DATAPOINTS_DIMENSIONS)
        faceon_classifier.save(constants.CLASSIFY_FACEON_MODEL_PATH)

    def create_golf_swing_isswing_classifier(self):
        # get body datapoints from frames
        all_datapoints = util.load_all_body_datapoints_from_file()
        # grab true values of whether each frame is part of the swing or not
        isswing_indices = util.load_isswing_indices_from_file()
        # create sequential classifier for whether frame is swing or not
        inputs, outputs = all_datapoints, isswing_indices
        x_train, x_val, y_train, y_val = util.get_random_training_and_validation_sets(inputs, outputs)
        isswing_classifier = keras_models.create_sequential_classifier(x_train, y_train, x_val, y_val, constants.FRAME_DIMENSION, constants.BODY_DATAPOINTS_DIMENSIONS)
        isswing_classifier.save(constants.CLASSIFY_ISSWING_MODEL_PATH)


class Model(Enum):
    FEATURE_EXTRACTOR = "feature_extractor"
    FACEON_CLASSIFIER = "faceon_classifier"
    ISSWING_CLASSIFIER = "isswing_classifier"

def model_type(value: str) -> Model:
    for model in Model:
        if value == model.value:
            return model
    raise argparse.ArgumentTypeError(f"Invalid model: {value}")


if __name__ == "__main__":
    # create arg parser
    parser = argparse.ArgumentParser(
        description="Script that allows manual training"
    )
    parser.add_argument(
        "--model",
        type=model_type,
        required=True,
    )
    args = parser.parse_args()
    # train model based on arg
    golf_ml_trainer = GolfMLTrainer()
    if args.model == Model.FEATURE_EXTRACTOR:
        golf_ml_trainer.create_golf_swing_feature_extractor()
    elif args.model == Model.FACEON_CLASSIFIER:
        golf_ml_trainer.create_golf_swing_faceon_classifier()
    elif args.model == Model.ISSWING_CLASSIFIER:
        golf_ml_trainer.create_golf_swing_isswing_classifier()
