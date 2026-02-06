# Python
from pathlib import Path
import os

curr_dir = Path(__file__).parent.resolve()

# ml-server src root directory
ROOT_DIR = os.path.join(curr_dir, '..')

# Dimensions constants
PREPROCESSED_VIDEOS = 38
NUM_VIDEOS = 38
FRAME_DIMENSION = 150
BODY_DATAPOINTS_DIMENSIONS = 25 * 3 # 25 datapoints, each with 3 features (x, y, confidence)

# Model paths
FEATURE_EXTRACTOR_MODEL_PATH = os.path.join(ROOT_DIR, "models", "feature_extractor.keras")
CLASSIFY_FACEON_MODEL_PATH = os.path.join(ROOT_DIR, "models", "classify_faceon_model.keras")
CLASSIFY_ISSWING_MODEL_PATH = os.path.join(ROOT_DIR, "models", "classify_swing_model.keras")
CLASSIFY_ISSWING_CNN_MODEL_PATH = os.path.join(ROOT_DIR, "models", "classify_swing_model_cnn.keras")

# Data paths
VIDEOS_DIR_PATH = os.path.join(ROOT_DIR, "data", "videos")
BODY_DATAPOINTS_PATH = os.path.join(ROOT_DIR, "data", "body_datapoints.npy")
FACEON_TRUE_VALUES_PATH = os.path.join(ROOT_DIR, "data", "faceon_indices.npy")
ISSWING_TRUE_VALUES_PATH = os.path.join(ROOT_DIR, "data", "isswing_indices.npy")
UNLABELED_VIDEOS_DICT_PATH = os.path.join(ROOT_DIR, "data", "unlabeled_videos_dict.json")
