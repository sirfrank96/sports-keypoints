# Python
from PIL import Image
from functools import partial
from io import BytesIO
from enum import Enum
import grpc

# Tkinter
import tkinter as tk
from tkinter import messagebox, simpledialog

# Internal
from gen import golfkeypoints_pb2, common_pb2
import util
import canvas_wrapper as cw
import frame_wrapper as fw

class MainAppPage(fw.FrameWrapper):
    # class variables
    body_pose_field_descriptors = common_pb2.Body25PoseDatapoints.DESCRIPTOR.fields

    def __init__(self, parent, controller, user_client, golfkeypoints_client, session_token):
        super().__init__(parent)
        # configure frame grid
        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(0, weight=1)
        # set up canvasses and frames to put stuff on
        self.whole_canvas = cw.CanvasWrapper(self)
        self.whole_canvas.make_scrollable()
        self.content_frame = self.whole_canvas.create_content_frame_in_canvas()
        # images are 1080x2400, this content_canvas is 1/4 of the size
        self.content_canvas = cw.CanvasWrapper(self.content_frame, width=270, height=600, bg='white', row=12, col=1, padx=5, pady=5, sticky="w")
        # add label and buttons to allow user access initial golf_keypoints and user apis
        self.main_page_label = self.content_frame.add_label(text="This is the main page", row=0, col=0, padx=5, pady=5)
        self.select_new_input_image_button = self.content_frame.add_button(text="Select New Input Image", command=self.select_new_input_image, row=1, col=0, padx=5, pady=5)
        self.show_existing_images_button = self.content_frame.add_button(text="Show Existing Input Images", command=self.show_existing_input_images, row=2, col=0, padx=5, pady=5)
        self.show_user_button = self.content_frame.add_button(text="Show User Information", command=self.read_user, row=3, col=0, padx=5, pady=5)
        self.update_user_button = self.content_frame.add_button(text="Update User Information", command=self.update_user, row=4, col=0, padx=5, pady=5)
        self.delete_user_button = self.content_frame.add_button(text="Delete User", command=self.delete_user, row=5, col=0, padx=5, pady=5)
        # set instance vars
        self.parent = parent
        self.controller = controller
        # pass in clients for api access
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.session_token = session_token
        # initialize vars that will be set during input image upload process
        self.image_type = golfkeypoints_pb2.ImageType.IMAGE_TYPE_UNSPECIFIED 
        self.curr_input_image_id = ""
        self.curr_input_image = None
        # initialize vars that will be set during calibration process
        self.golf_specific_datapoints = golfkeypoints_pb2.GolfSpecificDatapoints(golf_ball=None, club_butt=None, club_head=None, shoulder_tilt=None)
        self.axes_calibration_image = None
        self.vanishing_point_calibration_image = None
        self.horizontal_axis = None
        self.vertical_axis = None
        self.first_line_at_target = None
        self.second_line_at_target = None
        self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_HEEL_LINE
        # initialize identify modes to none (will change when calibrating image)
        self.identify_mode = self.IdentifyMode.NONE
        self.identify_line_mode = self.IdentifyLineMode.NONE
        # initialize body datapoints vars so user can update it later
        self.body_datapoints = None

    # enum used to differentiate when user is clicking on image to identify golf ball, club butt, club head
    class IdentifyMode(Enum):
        NONE = 1
        GOLFBALL = 2
        CLUBBUTT = 3
        CLUBHEAD = 4

    # enum used to differentiate when user is drawing lines for calibration
    class IdentifyLineMode(Enum):
        NONE = 1
        HORAXIS = 2
        VERTAXIS = 3
        LINEATTARGET1 = 4
        LINEATTARGET2 = 5
    
    def select_new_input_image(self):
        img = util.get_image_from_filesystem()
        if img is not None:
            try:
                bytes = util.get_image_bytes(img)
                messagebox.showinfo("Bytes", f"Length of image bytes is: {len(bytes)}, length of image raw bytes: {len(img.tobytes())}")
                # get whether image is face on or dtl
                faceon_response = messagebox.askquestion("FaceOn or DTL", "Is this image Face On? (Yes for Face On, No for DTL)")
                self.image_type = golfkeypoints_pb2.ImageType.FACE_ON if faceon_response == "yes" else golfkeypoints_pb2.ImageType.DTL
                # get description of input image
                description_response = simpledialog.askstring("Input Image Description", prompt="Enter description for input image (eg. Driver DTL: feel pressure shift earlier)")
                # get timestamp
                timestamp = util.get_curr_grpc_timestamp()
                response = self.golfkeypoints_client.upload_input_image(session_token=self.session_token, image_type= self.image_type, image=bytes, description=description_response, timestamp=timestamp)
                messagebox.showinfo("Successfully Upload Input", f"response is {response}")
                self.curr_input_image_id = response.input_image_id
                self.curr_input_image = img
                self.display_input_image(self.curr_input_image)
            except grpc.RpcError as e:
                messagebox.showerror("Upload input image failed", f"Could not upload input image: {e.code()}: {e.details()}")
            
    def show_existing_input_images(self):
        self.content_canvas.clear_canvas()
        # get list of all input images
        try:
            response = self.golfkeypoints_client.list_input_images_for_user(session_token=self.session_token)
            for i, input_image_id in enumerate(response.input_image_ids):
                # for each input image, read and get the image bytes + information
                response = self.read_input_image(input_image_id)
                if response is not None:
                    curr_button = tk.Button(self.content_canvas, text=f"{response.timestamp.ToDatetime()}: {response.description}", command=partial(self.on_press_input_image_list_button, input_image_id, response))
                    self.content_canvas.create_window(100, 30+(i*50), window=curr_button)
        except grpc.RpcError as e:
            messagebox.showerror("List Images Failed", f"Could not get a list of images: {e.code()}: {e.details()}")

    def on_press_input_image_list_button(self, input_image_id, read_input_image_response):
        self.curr_input_image_id = input_image_id
        buffer = BytesIO(read_input_image_response.image)
        img = Image.open(buffer)
        self.curr_input_image = img
        self.image_type = read_input_image_response.image_type
        self.display_input_image(img)
    
    def read_input_image(self, input_image_id):
        try:
            response = self.golfkeypoints_client.read_input_image(session_token=self.session_token, input_image_id=input_image_id)
            messagebox.showinfo("Show Image", f"Response length of image: {len(response.image)}, ImageType: {response.image_type}, Calibrated: {response.calibrated}, Description: {response.description}, Timestamp: {response.timestamp.ToDatetime()}")
            return response
        except grpc.RpcError as e:
            messagebox.showerror("Show Image Failed", f"Could not get image: {e.code()}: {e.details()}")
            return None

    def read_user(self):
        try:
            response = self.user_client.read_user(self.session_token)
            messagebox.showinfo("Show User", f"User info: {response}")
        except grpc.RpcError as e:
            messagebox.showerror("Show User", f"Show user failed: {e.code()}: {e.details()}")

    def update_user(self):
        try:
            new_username = simpledialog.askstring("Input New Username", prompt="Enter new username, if same leave empty")
            new_password = simpledialog.askstring("Input New Password", prompt="Enter new password, if same leave empty", show="*")
            new_email = simpledialog.askstring("Input New Email", prompt="Enter new email, if same leave empty")
            response = self.user_client.update_user(session_token=self.session_token, username=new_username, password=new_password, email=new_email)
            messagebox.showinfo("Update User", f"Updated user info: {response}")
        except grpc.RpcError as e:
            messagebox.showerror("Update User", f"Update user failed: {e.code()}: {e.details()}")
    
    def delete_user(self):
        try:
            response = self.user_client.delete_user(self.session_token)
            messagebox.showinfo("Delete User", f"Successfully deleted user {response}")
        except grpc.RpcError as e:
            messagebox.showerror("Delete User", f"Delete user failed: {e.code()}: {e.details()}")

    def display_input_image(self, image):
        self.content_canvas.display_an_image(image)
        # create buttons around calibration, calculation and other golf keypoints apis
        self.calibrate_button = self.content_frame.add_button(text="Calibrate Image", command=self.calibrate_image, row=0, col=2, padx=5, pady=5)
        self.calculate_button = self.content_frame.add_button(text="Calculate Golf Keypoints", command=self.start_calculate_process, row=1, col=2, padx=5, pady=5)
        self.read_keypoints_button = self.content_frame.add_button(text="Show Golf Keypoints for Input Image", command=self.read_golf_keypoints, row=2, col=2, padx=5, pady=5)
        self.delete_input_image_button = self.content_frame.add_button(text="Delete Input Image", command=self.delete_input_image, row=3, col=2, padx=5, pady=5)
        self.delete_keypoints_button = self.content_frame.add_button(text="Delete Golf Keypoints for Input Image", command=self.delete_golf_keypoints, row=4, col=2, padx=5, pady=5)

    def modify_feet_line_method(self):
        response = messagebox.askquestion("Modify Feet Line Method", "Do you want to change the feet line method to toe line?")
        if response == "yes":
            self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_TOE_LINE
    
    def calibrate_image(self):
        self.modify_feet_line_method()
        additional_imgs_needed_response = messagebox.askquestion("Calibrate With Additional Images", "Do you want to calibrate with additional images? (If so, you will select images from filesystem. If not, you will click points on the current input image)")
        if additional_imgs_needed_response == "yes":
            self.calibrate_input_image()
        else:
            # start process for user drawing manual lines (see on_draw_line_on_input_image for rest of logic)
            self.identify_line_mode = self.IdentifyLineMode.HORAXIS
            messagebox.showinfo("Horizontal Axis Identify", "Please click and drag a line for the horizontal axis (ie. parallel to the ground)")
            self.content_canvas.bind("<ButtonPress-1>", self.on_draw_line_on_input_image)

    def calibrate_input_image(self):
        self.get_axes_calibration_image()
        if self.image_type == golfkeypoints_pb2.ImageType.DTL:
            self.get_vanishing_point_calibration_image()
        try:
            response = self.golfkeypoints_client.calibrate_input_image(session_token=self.session_token, input_image_id=self.curr_input_image_id, calibration_type=golfkeypoints_pb2.CalibrationType.FULL_CALIBRATION, feet_line_method=self.feet_line_method, calibration_image_axes=self.axes_calibration_image, calibration_image_vanishing_point=self.vanishing_point_calibration_image)
            messagebox.showinfo("Calibrate Input Image", f"Calibrate input image successful: {response}, calculate golf keypoints next")
            self.calibrate_button.config(state=tk.DISABLED)
        except grpc.RpcError as e:
            messagebox.showerror("Calibrate Input Image", f"Calibrate input image failed: {e.code()}: {e.details()}") 
    
    def calibrate_input_image_manual(self):
        try:
            response = self.golfkeypoints_client.calibrate_input_image_manual(session_token=self.session_token, input_image_id=self.curr_input_image_id, calibration_type=golfkeypoints_pb2.CalibrationType.FULL_CALIBRATION, feet_line_method=self.feet_line_method, horizontal_axis=self.horizontal_axis, vertical_axis=self.vertical_axis, first_line_at_target=self.first_line_at_target, second_line_at_target=self.second_line_at_target)
            messagebox.showinfo("Calibrate Input Image Manual", f"Calibrate input image manual successful: {response}, calculate golf keypoints next")
            self.calibrate_button.config(state=tk.DISABLED)
        except grpc.RpcError as e:
            messagebox.showerror("Calibrate Input Image Manual", f"Calibrate input image manual failed: {e.code()}: {e.details()}")

    def get_axes_calibration_image(self):
        messagebox.showinfo("Get Axes Calibration Image", "Please select your axes calibration image from your filesystem")
        img = util.get_image_from_filesystem()
        if img is not None:
            bytes = util.get_image_bytes(img)
            self.axes_calibration_image = bytes
            messagebox.showinfo("Axes Calibration Image", "Successfully set axes calibration image")
        else:
            messagebox.showerror("Axes Calibration Image", "Could not get axes calibration image")

    def get_vanishing_point_calibration_image(self):
        messagebox.showinfo("Get Vanishing Point Calibration Image", "Please select your vanishing point calibration image from your filesystem")
        img = util.get_image_from_filesystem()
        if img is not None:
            bytes = util.get_image_bytes(img)
            self.vanishing_point_calibration_image = bytes
            messagebox.showinfo("Vanishing Point Calibration Image", "Successfully set vanishing point calibration image")
        else:
            messagebox.showerror("Vanishing Point Calibration Image", "Could not get vanishing point calibration image")

    def start_calculate_process(self):
        if self.image_type == golfkeypoints_pb2.ImageType.DTL:
            self.input_shoulder_tilt()
        # start process of user clicking where golf equipment datapoints are
        self.identify_mode = self.IdentifyMode.GOLFBALL
        self.content_canvas.bind("<Button-1>", self.on_click_on_input_image)
        messagebox.showinfo("Golf Ball Identify", "Please click on the input image where the golf ball is")

    def input_shoulder_tilt(self):
        shoulder_tilt = simpledialog.askfloat("Shoulder Tilt", prompt="What is the shoulder tilt?")
        self.golf_specific_datapoints.shoulder_tilt.CopyFrom(common_pb2.Double(data=shoulder_tilt, warning= ""))   

    def calculate_golf_keypoints(self):
        try:
            response = self.golfkeypoints_client.calculate_golf_keypoints(session_token=self.session_token, input_image_id=self.curr_input_image_id, golf_specific_datapoints=self.golf_specific_datapoints)
            messagebox.showinfo("Calculate Golf Keypoints", f"Calculate Golf Keypoints successful")
            if response.output_image is not None:
                self.process_golf_keypoints(response.output_image, response.golf_keypoints)
                self.calculate_button.config(state=tk.DISABLED)
        except grpc.RpcError as e:
            messagebox.showerror("Calculate Golf Keypoints", f"Calculate Golf Keypoints failed: {e.code()}: {e.details()}") 

    def read_golf_keypoints(self):
        try:
            response = self.golfkeypoints_client.read_golf_keypoints(self.session_token, self.curr_input_image_id)
            messagebox.showinfo("Read Golf Keypoints", f"Read Golf Keypoints successful")
            if response.output_image is not None:
                self.process_golf_keypoints(response.output_image, response.golf_keypoints)
        except grpc.RpcError as e:
            messagebox.showerror("Read Golf Keypoints", f"Read golf keypoints failed: {e.code()}: {e.details()}")

    def process_golf_keypoints(self, output_image, golf_keypoints):
        self.body_datapoints = golf_keypoints.body_datapoints
        self.golf_specific_datapoints = golf_keypoints.golf_specific_datapoints
        buffer = BytesIO(output_image)
        img = Image.open(buffer)
        self.content_canvas.display_an_image(img)
        self.display_golf_keypoints_text(golf_keypoints)

    def close_popup(self, popup):
        self.whole_canvas.make_scrollable()
        popup.destroy()

    def add_tk_top_level_close_handler(self, popup):
        popup.protocol("WM_DELETE_WINDOW", partial(self.close_popup, popup))

    def display_golf_keypoints_text(self, golf_keypoints):
        # create popup window to show golf keypoints
        popup = tk.Toplevel(self)
        self.add_tk_top_level_close_handler(popup)
        popup.wm_title("Golf Keypoints Window")
        popup_canvas = cw.CanvasWrapper(popup, width=270, height=600)
        popup_content_frame = popup_canvas.create_content_frame_in_canvas()
        popup_content_frame.add_scrolled_text(text=f"{golf_keypoints}", row=0, col=0, padx=5, pady=5)
        popup_content_frame.add_button(text="Done Looking at Golf Keypoints", command=partial(self.done_looking_at_golf_keypoints, popup), row=1, col=0, padx=5, pady=5)
    
    def done_looking_at_golf_keypoints(self, popup):
        self.close_popup(popup)
        incorrect = messagebox.askyesno("Body Datapoints Update", "Are there body datapoints that computervision identified incorrectly?")
        if incorrect:
            self.select_body_datapoints_to_update(self.body_datapoints)

    def select_body_datapoints_to_update(self, body_datapoints):
        # create popup window to show body datapoints that can be updated
        popup = tk.Toplevel(self)
        self.add_tk_top_level_close_handler(popup)
        popup.wm_title("Body Datapoints Window")
        popup_canvas = cw.CanvasWrapper(popup, width=270, height=600)
        popup_canvas.make_scrollable()
        popup_content_frame = popup_canvas.create_content_frame_in_canvas()
        # iterate over body datapoints and create a button for each
        idx = 0
        for field in MainAppPage.body_pose_field_descriptors:
            name = field.name
            body_datapoint_value = getattr(body_datapoints, name)
            popup_content_frame.add_button(text=f"Modify {name}: {body_datapoint_value}", command=partial(self.update_body_datapoint, name), row=idx, col=0, padx=5, pady=5)
            idx += 1
        popup_content_frame.add_button(text="Done Updating Body Datapoints", command=partial(self.update_body_datapoints, popup), row=idx, col=0, padx=5, pady=5)
        return 
    
    def update_body_datapoint(self, field_name):
        x = simpledialog.askfloat("New Value ", prompt=f"What is the new x value for {field_name}")
        y = simpledialog.askfloat("New Value ", prompt=f"What is the new y value for {field_name}")
        if x != None and y != None:
            field = getattr(self.body_datapoints, field_name)
            setattr(field, "x", x)
            setattr(field, "y", y)
            setattr(field, "confidence", 1.0)

    def update_body_datapoints(self, popup):
        self.close_popup(popup)
        try:
            response = self.golfkeypoints_client.update_body_datapoints(self.session_token, self.curr_input_image_id, self.body_datapoints)
            self.process_golf_keypoints(response.updated_output_image, response.updated_golf_keypoints)
        except grpc.RpcError as e:
            messagebox.showerror("Update Body Datapoints", f"Update Body Datapoints failed: {e.code()}: {e.details()}")
        
    def delete_input_image(self):
        try:
            response = self.golfkeypoints_client.delete_input_image(self.session_token, self.curr_input_image_id)
            messagebox.showinfo("Delete Input Image", f"Successfully deleted input image {response}")
            # clear canvas
            self.content_canvas.clear_canvas()
        except grpc.RpcError as e:
            messagebox.showerror("Delete Input Image", f"Delete input image failed: {e.code()}: {e.details()}")

    def delete_golf_keypoints(self):
        try:
            response = self.golfkeypoints_client.delete_golf_keypoints(self.session_token, self.curr_input_image_id)
            messagebox.showinfo("Delete Golf Keypoints", f"Successfully deleted golf keypoints {response}")
            # go back to original input image
            self.content_canvas.display_an_image(self.curr_input_image)
        except grpc.RpcError as e:
            messagebox.showerror("Delete Golf Keypoints", f"Delete golf keypoints failed: {e.code()}: {e.details()}")

    # this function is the logic for drawing calibration lines manually
    # each time a user clicks the image, the first point of the line is saved
    # then when the user releases, the second point is saved a line is drawn between the two points
    # each time a line is successfully drawn, the identify_line_mode is switched to the next mode (eg. from IdentifyLineMode.HORAXIS -> IdentifyLineMode.VERTAXIS)
    # once the correct lines are drawn (axes for faceon, axes and vanishing point for dtl), this function will call the calibrate_input_image_manual function
    def on_draw_line_on_input_image(self, event):
        x = event.x
        y = event.y
        # because canvas is 1/4 the size of a 1080x2400 image, we have to scale the image back
        scaled_x = x*4
        scaled_y = y*4
        match self.identify_line_mode:
            case self.IdentifyLineMode.HORAXIS:
                if event.type == tk.EventType.ButtonPress:
                    self.horizontal_axis = common_pb2.Line()
                    self.horizontal_axis.first_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.content_canvas.bind("<ButtonRelease-1>", self.on_draw_line_on_input_image)
                elif event.type == tk.EventType.ButtonRelease:
                    self.horizontal_axis.second_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    first_x = self.horizontal_axis.first_point_on_line.x / 4
                    first_y = self.horizontal_axis.first_point_on_line.y / 4
                    line_id = self.content_canvas.draw_line(first_x, first_y, x, y, "red")
                    self.content_canvas.unbind("<ButtonRelease-1")
                    ok = messagebox.askokcancel("Line Drawn", "Is this the correct horizontal axis?")
                    if ok:
                        self.identify_line_mode = self.IdentifyLineMode.VERTAXIS
                        messagebox.showinfo("Vertical Axis Identify", "Please click and drag a line for the vertical axis (ie. center of frame, perpendicular to the horizontal axis)")
                    else:
                        self.content_canvas.erase_line(line_id)
                    self.content_canvas.bind("<ButtonPress-1>", self.on_draw_line_on_input_image)
            case self.IdentifyLineMode.VERTAXIS:
                if event.type == tk.EventType.ButtonPress:
                    self.vertical_axis = common_pb2.Line()
                    self.vertical_axis.first_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.content_canvas.bind("<ButtonRelease-1>", self.on_draw_line_on_input_image)
                elif event.type == tk.EventType.ButtonRelease:
                    self.vertical_axis.second_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    first_x = self.vertical_axis.first_point_on_line.x / 4
                    first_y = self.vertical_axis.first_point_on_line.y / 4
                    line_id = self.content_canvas.draw_line(first_x, first_y, x, y, "blue")
                    self.content_canvas.unbind("<ButtonRelease-1")
                    ok = messagebox.askokcancel("Line Drawn", "Is this the correct vertical axis?")
                    if ok:
                        if self.image_type == golfkeypoints_pb2.ImageType.DTL:
                            self.identify_line_mode = self.IdentifyLineMode.LINEATTARGET1
                            messagebox.showinfo("Line At Target 1 Identify", "Please click and drag a line for the first line at the target (ie. line on the ground pointing at the target)")
                        else:
                            self.identify_line_mode = self.IdentifyLineMode.NONE
                            self.calibrate_input_image_manual()
                            return
                    else:
                        self.content_canvas.erase_line(line_id)
                    self.content_canvas.bind("<ButtonPress-1>", self.on_draw_line_on_input_image)
            case self.IdentifyLineMode.LINEATTARGET1:
                if event.type == tk.EventType.ButtonPress:
                    self.first_line_at_target = common_pb2.Line()
                    self.first_line_at_target.first_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.content_canvas.bind("<ButtonRelease-1>", self.on_draw_line_on_input_image)
                elif event.type == tk.EventType.ButtonRelease:
                    self.first_line_at_target.second_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    first_x = self.first_line_at_target.first_point_on_line.x / 4
                    first_y = self.first_line_at_target.first_point_on_line.y / 4
                    line_id = self.content_canvas.draw_line(first_x, first_y, x, y, "green")
                    self.content_canvas.unbind("<ButtonRelease-1")
                    ok = messagebox.askokcancel("Line Drawn", "Is this the correct first line at target axis?")
                    if ok:
                        self.identify_line_mode = self.IdentifyLineMode.LINEATTARGET2
                        messagebox.showinfo("Line At Target 2 Identify", "Please click and drag a line for the second line at the target (ie. another line on the ground pointing at the target)")
                    else:
                        self.content_canvas.erase_line(line_id)
                    self.content_canvas.bind("<ButtonPress-1>", self.on_draw_line_on_input_image)
            case self.IdentifyLineMode.LINEATTARGET2:
                if event.type == tk.EventType.ButtonPress:
                    self.second_line_at_target = common_pb2.Line()
                    self.second_line_at_target.first_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.content_canvas.bind("<ButtonRelease-1>", self.on_draw_line_on_input_image)
                elif event.type == tk.EventType.ButtonRelease:
                    self.second_line_at_target.second_point_on_line.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    first_x = self.second_line_at_target.first_point_on_line.x / 4
                    first_y = self.second_line_at_target.first_point_on_line.y / 4
                    line_id = self.content_canvas.draw_line(first_x, first_y, x, y, "red")
                    self.content_canvas.unbind("<ButtonRelease-1")
                    ok = messagebox.askokcancel("Line Drawn", "Is this the correct second line at target axis?")
                    if ok:
                        self.identify_line_mode = self.IdentifyLineMode.NONE
                        self.calibrate_input_image_manual()
                    else:
                        self.content_canvas.erase_line(line_id)
                        self.content_canvas.bind("<ButtonPress-1>", self.on_draw_line_on_input_image)
        
    # this function is the logic for identifying golf ball and golf club points
    # each time a user presses the image for the corresponding prompt (eg. "Please click on the input image where the golf ball is"), 
    # the coordinates are saved to the golf_specific_datapoints object, then the identify_mode is switched to the next (eg. IdentifyMode.CLUBBUTT)
    def on_click_on_input_image(self, event):
        x = event.x
        y = event.y
        scaled_x = x*4
        scaled_y = y*4
        messagebox.showinfo("Clicked", f"Clicked at x: {x}, y: {y}, scaled x: {scaled_x}, scaled y: {scaled_y}")
        match self.identify_mode:
            case self.IdentifyMode.GOLFBALL:
                circle_id = self.content_canvas.draw_circle(x, y, "red")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the golf ball?")
                if ok:
                    self.golf_specific_datapoints.golf_ball.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.identify_mode = self.IdentifyMode.CLUBBUTT
                    messagebox.showinfo("Club Butt Identify", "Please click on the input image where the club butt is")
                else:
                    self.content_canvas.erase_circle(circle_id)
            case self.IdentifyMode.CLUBBUTT:
                circle_id = self.content_canvas.draw_circle(x, y, "blue")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the butt end of the club?")
                if ok:
                    self.golf_specific_datapoints.club_butt.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.identify_mode = self.IdentifyMode.CLUBHEAD
                    messagebox.showinfo("Club Head Identify", "Please click on the input image where the club head is")
                else:
                    self.content_canvas.erase_circle(circle_id)
            case self.IdentifyMode.CLUBHEAD:
                circle_id = self.content_canvas.draw_clubhead(x, y, "green")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the clubhead")
                if ok:
                    self.golf_specific_datapoints.club_head.CopyFrom(common_pb2.Datapoint(x=scaled_x, y=scaled_y, confidence=1.0))
                    self.identify_mode = self.IdentifyMode.NONE
                    self.content_canvas.unbind("<Button-1>")
                    self.calculate_golf_keypoints()
                else:
                    self.content_canvas.erase_circle(circle_id)
            case self.IdentifyMode.NONE:
                messagebox.showerror("Clicked", "Please click one of calibration buttons")
    