import cv2
import math

from keras.preprocessing import sequence
from sklearn.model_selection import train_test_split

def frame_to_bytes(frame):
    ok, buf = cv2.imencode('.jpg', frame)
    if ok:
        bytes = buf.tobytes() 
        return bytes
    else:
        return None 

frames_dim_for_pose_estimation = 150
#frames_dim_for_swing_classification = 300
frames_dim_for_swing_classification = 150

def downsample_frames(frames, frames_dim, vals=None):
    num_frames = len(frames)
    # let keras pad frames if too small
    if num_frames <= frames_dim:
        return frames, vals
    # downsample frames if num_frames is greater than required
    frame_rate = num_frames / frames_dim
    frame_rate = truncate_float_to_tenths(frame_rate)
    downsampled_frames = []
    downsampled_vals = None
    if vals is not None:
        downsampled_vals = []
    last_idx = -1
    for i in range(0, num_frames):
        if last_idx >= frames_dim - 1:
            break
        idx_target = int(i/frame_rate)
        if idx_target > last_idx:
            downsampled_frames.append(frames[i])
            if vals is not None:
                downsampled_vals.append(vals[i])
            last_idx = idx_target
    return downsampled_frames, downsampled_vals

def standardize_inputs_and_outputs(inputs, outputs):
    padded_inputs = sequence.pad_sequences(inputs, padding='post', value=0.0)
    padded_outputs = sequence.pad_sequences(outputs, padding='post', value=0.0)
    return padded_inputs, padded_outputs

# returns x_train, x_test/val, y_train, y_test/val
def get_random_training_and_validation_sets(inputs, outputs, test_size=0.2):
    #train_test_split(inputs, outputs, test_size=test_size, shuffle=True)
    return inputs[:30], inputs[30:], outputs[:30], outputs[30:]

def truncate_float_to_tenths(f):
    return math.trunc(f * 10) / 10.0
