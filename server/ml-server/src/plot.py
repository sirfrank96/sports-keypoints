# Matplotlib
import matplotlib.pyplot as plt

# Internal
from gen import common_pb2, common_pb2_grpc

# converts flattened 75 index array to grpc object 
def convert_data_to_body25_datapoints(data):
    body_25_datapoints = common_pb2.Body25PoseDatapoints()
    field_idx = 0
    for descriptor in body_25_datapoints.DESCRIPTOR.fields:
            data_start_idx = field_idx * 3
            x = data[data_start_idx]
            y = data[data_start_idx+1]
            confidence = data[data_start_idx+2]
            field = common_pb2.Datapoint(x=x, y=y, confidence=confidence)
            curr_field = getattr(body_25_datapoints, descriptor.name)
            curr_field.CopyFrom(field)
            field_idx += 1
    return body_25_datapoints

def get_all_x_coordinates(data):
    x = []
    for i in range(0, 25*3, 3):
        curr_x = data[i]
        x.append(curr_x)
    return x
     
def get_all_y_coordinates(data):
    y = []
    for i in range(0, 25*3, 3):
        curr_y = data[i+1]
        y.append(curr_y)
    return y

def draw_line_between_two_points(point1, point2):
    if point1.x == 0.0 and point1.y == 0.0:
        return
    if point2.x == 0.0 and point2.y == 0.0:
        return
    plt.plot([point1.x, point2.x], [point1.y, point2.y], linestyle='-', color='green')


def plot_datapoints_as_skeleton(data):
    body_25_datapoints = convert_data_to_body25_datapoints(data)
    x = get_all_x_coordinates(data)
    y = get_all_y_coordinates(data)
    plt.scatter(x, y, color='red', marker = 'o')

    # draw head area
    draw_line_between_two_points(body_25_datapoints.l_ear, body_25_datapoints.l_eye)
    draw_line_between_two_points(body_25_datapoints.r_ear, body_25_datapoints.r_eye)
    draw_line_between_two_points(body_25_datapoints.l_eye, body_25_datapoints.nose)
    draw_line_between_two_points(body_25_datapoints.r_eye, body_25_datapoints.nose)
    draw_line_between_two_points(body_25_datapoints.nose, body_25_datapoints.neck)
    # draw upper body
    draw_line_between_two_points(body_25_datapoints.l_shoulder, body_25_datapoints.neck)
    draw_line_between_two_points(body_25_datapoints.r_shoulder, body_25_datapoints.neck)
    draw_line_between_two_points(body_25_datapoints.l_shoulder, body_25_datapoints.l_elbow)
    draw_line_between_two_points(body_25_datapoints.l_elbow, body_25_datapoints.l_wrist)
    draw_line_between_two_points(body_25_datapoints.r_shoulder, body_25_datapoints.r_elbow)
    draw_line_between_two_points(body_25_datapoints.r_elbow, body_25_datapoints.r_wrist)
    draw_line_between_two_points(body_25_datapoints.midhip, body_25_datapoints.neck)
    # draw lower body
    draw_line_between_two_points(body_25_datapoints.l_hip, body_25_datapoints.midhip)
    draw_line_between_two_points(body_25_datapoints.r_hip, body_25_datapoints.midhip)
    draw_line_between_two_points(body_25_datapoints.l_hip, body_25_datapoints.l_knee)
    draw_line_between_two_points(body_25_datapoints.l_knee, body_25_datapoints.l_ankle)
    draw_line_between_two_points(body_25_datapoints.r_hip, body_25_datapoints.r_knee)
    draw_line_between_two_points(body_25_datapoints.r_knee, body_25_datapoints.r_ankle)
    # draw feet
    draw_line_between_two_points(body_25_datapoints.l_ankle, body_25_datapoints.l_heel)
    draw_line_between_two_points(body_25_datapoints.l_ankle, body_25_datapoints.l_small_toe)
    draw_line_between_two_points(body_25_datapoints.l_ankle, body_25_datapoints.l_big_toe)
    draw_line_between_two_points(body_25_datapoints.r_ankle, body_25_datapoints.r_heel)
    draw_line_between_two_points(body_25_datapoints.r_ankle, body_25_datapoints.r_small_toe)
    draw_line_between_two_points(body_25_datapoints.r_ankle, body_25_datapoints.r_big_toe)
    
    plt.gca().invert_yaxis()
    plt.show()
