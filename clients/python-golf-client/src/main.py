# start clients to golfkeypoints and user
# launch gui for user to create user, login, upload inputimages, mark golf ball and club,
# upload calibration images, display output, display data, update body keypoints

# Tkinter
import tkinter as tk

# Python
import logging
import sys
import os
from pathlib import Path

# GRPC
import grpc

# Internal
curr_dir = Path(__file__).parent.resolve()
sys.path.append(os.path.join(curr_dir, 'gen'))
import user_client as uc
import golf_keypoints_client as gc
import login_pages as login

class GolfKeypointsClientApp(tk.Tk):
    def __init__(self, user_client, golfkeypoints_client):
        super().__init__()
        self.title("Golf Setup and Keypoint Client")
        self.geometry('900x900')
        # create initial frame that everything will be put inside of
        container = tk.Frame(self)  
        container.pack(side = "top", fill = "both", expand = True) 
        container.grid_rowconfigure(0, weight = 1)
        container.grid_columnconfigure(0, weight = 1)
        # create initial page frame and display
        initial_page = login.InitialPage(container, self, user_client=user_client, golfkeypoints_client=golfkeypoints_client)
        self.show_frame(initial_page)
    
    # this function will display that frame provided
    # frames inside the initial container will be passed a reference to the initial container (controller)
    # frames will then call controller.show_frame(frame_name) to switch what is displayed
    def show_frame(self, frame):
        frame.grid(row = 0, column = 0, sticky ="nsew")
        frame.tkraise()         

session_token = ""

def serve():
    #options=[('grpc.max_send_message_length', 10000000), ('grpc.max_receive_message_length', 10000000)]
    setting_timeout_ms = 1000 * 60 * 3
    options = [('grpc.http2.settings_timeout', setting_timeout_ms)]
    channel = grpc.insecure_channel('localhost:50052', options=options)
    user_client = uc.UserClient(channel)
    golfkeypoints_client = gc.GolfKeypointsClient(channel)
    app = GolfKeypointsClientApp(user_client, golfkeypoints_client)
    app.mainloop()

if __name__ == "__main__":
    logging.basicConfig()
    serve()
