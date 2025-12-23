import tkinter as tk
from tkinter import messagebox, filedialog, simpledialog
from PIL import ImageTk, Image
from functools import partial
import io
from io import BytesIO
import grpc
from enum import Enum
from google.protobuf.timestamp_pb2 import Timestamp
import datetime

import golfkeypoints_pb2
import golf_keypoints_client as gc
import common_pb2


class MainAppPage(tk.Frame):
    def __init__(self, parent, controller, user_client, golfkeypoints_client, session_token):
        tk.Frame.__init__(self, parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.session_token = session_token
        self.parent = parent
        self.controller = controller 

        self.grid_columnconfigure(0, weight=1)
        self.grid_rowconfigure(0, weight=1)
        # create canvas inside the frame (self)
        self.whole_canvas = tk.Canvas(self)
        self.whole_canvas.grid(row=0, column=0, sticky="nsew")
        # create scrollbar in frame (self) that controls canvas 
        self.scrollbar = tk.Scrollbar(self, orient="vertical", command=self.whole_canvas.yview)
        self.whole_canvas.configure(yscrollcommand=self.scrollbar.set)
        self.scrollbar.grid(row=0, column=1, sticky="ns")
        # create frame inside the canvas for content inside the canvas
        self.content_frame = tk.Frame(self.whole_canvas)
        #self.content_frame.grid_columnconfigure(0, weight=1)
        #self.content_frame.grid_rowconfigure(0, weight=1)
        self.whole_canvas.create_window((0, 0), window=self.content_frame, anchor="nw")
        self.whole_canvas.bind_all("<MouseWheel>", partial(self.on_mousewheel, self.whole_canvas))

        # add labels and buttons
        self.main_app_label = tk.Label(self.content_frame, text="This is the main app")
        self.main_app_label.grid(row=0, column=0, padx=5, pady=5, sticky="w")

        self.select_new_input_image_button = tk.Button(self.content_frame, text="Select New Input Image", command=self.select_new_input_image)
        self.select_new_input_image_button.grid(row=1, column=0, padx=5, pady=5, sticky="w")

        self.show_input_images_button = tk.Button(self.content_frame, text="Show Previous Input Images", command=self.show_previous_input_images)
        self.show_input_images_button.grid(row=2, column=0, padx=5, pady=5, sticky="w")

        self.read_user_button = tk.Button(self.content_frame, text="Show User", command=self.read_user)
        self.read_user_button.grid(row=3, column=0, padx=5, pady=5, sticky="w")
        
        self.update_user_button = tk.Button(self.content_frame, text="Update User", command=self.update_user)
        self.update_user_button.grid(row=4, column=0, padx=5, pady=5, sticky="w")

        self.delete_user_button = tk.Button(self.content_frame, text="Delete User", command=self.delete_user)
        self.delete_user_button.grid(row=5, column=0, padx=5, pady=5, sticky="w")

        # 1/4 size of image -> canvas for images
        self.canvas = tk.Canvas(self.content_frame, width=270, height=600, bg='white')
        self.canvas.grid(row=12, column=1, padx=5, pady=5, sticky="w")

        self.identify_mode = self.IdentifyMode.NONE
        self.golf_ball = None
        self.club_butt = None
        self.club_head = None
        self.image_type = golfkeypoints_pb2.ImageType.IMAGE_TYPE_UNSPECIFIED 
        self.axes_calibration_image = None
        self.vanishing_point_calibration_image = None
        self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_HEEL_LINE
        self.shoulder_tilt = common_pb2.Double(data=0, warning="no shoulder tilt")

        self.curr_input_image_id = ""
        self.curr_input_image = None
        self.body_keypoints = None

    class IdentifyMode(Enum):
        NONE = 1
        GOLFBALL = 2
        CLUBBUTT = 3
        CLUBHEAD = 4
    
    def get_image_from_filesystem(self):
        # Open the file dialog and get the file path
        filepath = filedialog.askopenfilename(
            filetypes=[("Image Files", "*.png;*.jpg;*.jpeg;*.gif")]
        )
        # If a file was selected, call the display_image function
        if filepath:
            img = Image.open(filepath)
            return img
        else:
            return None
    
    def select_new_input_image(self):
        img = self.get_image_from_filesystem()
        if img is not None:
            try:
                bytes = self.get_image_bytes(img)
                messagebox.showinfo("Bytes", f"Length of image bytes is: {len(bytes)}, length of image raw bytes: {len(img.tobytes())}")
                # get whether image is face on or dtl
                faceon_response = messagebox.askquestion("FaceOn or DTL", "Is this image Face On? (Yes for Face On, No for DTL)")
                self.image_type = golfkeypoints_pb2.ImageType.FACE_ON if faceon_response == "yes" else golfkeypoints_pb2.ImageType.DTL
                # get description of input image
                description_response = simpledialog.askstring("Input Image Description", prompt="Enter description for input image (eg. Driver DTL: feel pressure shift earlier)")
                # get timestamp
                now_utc = datetime.datetime.now(datetime.timezone.utc)
                grpc_timestamp = Timestamp()
                grpc_timestamp.FromDatetime(now_utc)
                response = self.golfkeypoints_client.upload_input_image(session_token=self.session_token, image_type= self.image_type, image=bytes, description=description_response, timestamp=grpc_timestamp)
                messagebox.showinfo("Successfully Upload Input", f"response is {response}")
                self.curr_input_image_id = response.input_image_id
                self.curr_input_image = img
                self.display_input_image(self.curr_input_image)
            except grpc.RpcError as e:
                messagebox.showerror("Upload input image failed", f"Could not upload input image: {e.code()}: {e.details()}")
            
    def show_previous_input_images(self):
        self.clear_canvas()
        # get list of all input images
        try:
            response = self.golfkeypoints_client.list_input_images_for_user(session_token=self.session_token)
            for i, input_image_id in enumerate(response.input_image_ids):
                # for each input image, read and get the image bytes + information
                response = self.read_input_image(input_image_id)
                if response is not None:
                    buffer = BytesIO(response.image)
                    img = Image.open(buffer)
                    self.curr_input_image_id = input_image_id
                    self.curr_input_image = img
                    curr_button = tk.Button(self.canvas, text=f"{response.timestamp.ToDatetime()}: {response.description}", command=partial(self.display_input_image, self.curr_input_image))
                    self.canvas.create_window(100, 30+(i*50), window=curr_button)
        except grpc.RpcError as e:
            messagebox.showerror("List Images Failed", f"Could not get a list of images: {e.code()}: {e.details()}")
    
    def read_input_image(self, input_image_id):
        try:
            response = self.golfkeypoints_client.read_input_image(session_token=self.session_token, input_image_id=input_image_id)
            messagebox.showinfo("Show Image", f"Response length of image: {len(response.image)}")
            return response
        except grpc.RpcError as e:
            messagebox.showerror("Show Image Failed", f"Could not get image: {e.code()}: {e.details()}")
            return None
        
    def display_image(self, image):
        self.clear_canvas()
        # 1/4 the size to display on canvas
        resized_img = image.resize((270, 600), Image.Resampling.LANCZOS)
        # Convert the image to a PhotoImage object
        photo = ImageTk.PhotoImage(resized_img)
        self.canvas.create_image(0, 0, anchor=tk.NW, image=photo)
        self.canvas.image = photo

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
        self.display_image(image)

        self.axes_button = tk.Button(self.content_frame, text="Get Axes Calibration Image", command=self.get_axes_calibration_image)
        self.axes_button.grid(row=0, column=2, padx=5, pady=5, sticky="w")

        self.vanishing_point_button = tk.Button(self.content_frame, text="Get Vanishing Point Calibration Image", command=self.get_vanishing_point_calibration_image)
        self.vanishing_point_button.grid(row=1, column=2, padx=5, pady=5, sticky="w")

        self.identify_golf_ball_button = tk.Button(self.content_frame, text="Identify Golf Ball", command=self.identify_golf_ball)
        self.identify_golf_ball_button.grid(row=2, column=2, padx=5, pady=5, sticky="w")

        self.identify_club_butt_button = tk.Button(self.content_frame, text="Identify Club Butt", command=self.identify_club_butt)
        self.identify_club_butt_button.grid(row=3, column=2, padx=5, pady=5, sticky="w")

        self.identify_club_head_button = tk.Button(self.content_frame, text="Identify Club Head", command=self.identify_club_head)
        self.identify_club_head_button.grid(row=4, column=2, padx=5, pady=5, sticky="w")

        self.feet_line_button = tk.Button(self.content_frame, text="Modify Feet Line Method", command=self.modify_feet_line_method)
        self.feet_line_button.grid(row=5, column=2, padx=5, pady=5, sticky="w")

        self.shoulder_tilt_button = tk.Button(self.content_frame, text="Input Shoulder Tilt For DTL", command=self.input_shoulder_tilt)
        self.shoulder_tilt_button.grid(row=6, column=2, padx=5, pady=5, sticky="w")

        self.calibrate_button = tk.Button(self.content_frame, text="Calibrate Image", command=partial(self.calibrate_image))
        self.calibrate_button.grid(row=7, column=2, padx=5, pady=5, sticky="w")

        self.calculate_button = tk.Button(self.content_frame, text="Calculate Golf Keypoints", command=partial(self.calculate_golf_keypoints))
        self.calculate_button.grid(row=8, column=2, padx=5, pady=5, sticky="w")

        self.read_keypoints_button = tk.Button(self.content_frame, text="Read Keypoints for Input Image", command=partial(self.read_golf_keypoints))
        self.read_keypoints_button.grid(row=9, column=2, padx=5, pady=5, sticky="w")

        self.delete_input_img_button = tk.Button(self.content_frame, text="Delete Input Image", command=partial(self.delete_input_image))
        self.delete_input_img_button.grid(row=10, column=2, padx=5, pady=5, sticky="w")

        self.delete_keypoints_button = tk.Button(self.content_frame, text="Delete Keypoints for Input Image", command=partial(self.delete_golf_keypoints))
        self.delete_keypoints_button.grid(row=11, column=2, padx=5, pady=5, sticky="w")

    def get_axes_calibration_image(self):
        img = self.get_image_from_filesystem()
        if img is not None:
            bytes = self.get_image_bytes(img)
            self.axes_calibration_image = bytes
            messagebox.showinfo("Axes Calibration Image", "Successfully set axes calibration image")
        else:
            messagebox.showerror("Axes Calibration Image", "Could not get axes calibration image")

    def get_vanishing_point_calibration_image(self):
        if self.image_type == golfkeypoints_pb2.ImageType.FACE_ON:
            messagebox.showerror("Vanishing Point Calibration Image", "Vanishing point calibration is only used for DTL images")
        else:
            img = self.get_image_from_filesystem()
            if img is not None:
                bytes = self.get_image_bytes(img)
                self.vanishing_point_calibration_image = bytes
                messagebox.showinfo("Vanishing Point Calibration Image", "Successfully set vanishing point calibration image")
            else:
                messagebox.showerror("Vanishing Point Calibration Image", "Could not get vanishing point calibration image")

    def identify_golf_ball(self):
        self.identify_mode = self.IdentifyMode.GOLFBALL
        self.canvas.bind("<Button-1>", self.on_click_on_input_image)
        messagebox.showinfo("Golf Ball Identify", "Please identify golf ball")

    def identify_club_butt(self):
        self.identify_mode = self.IdentifyMode.CLUBBUTT
        self.canvas.bind("<Button-1>", self.on_click_on_input_image)
        messagebox.showinfo("Club Butt Identify", "Please identify club butt")

    def identify_club_head(self):
        self.identify_mode = self.IdentifyMode.CLUBHEAD
        self.canvas.bind("<Button-1>", self.on_click_on_input_image)
        messagebox.showinfo("Club Head Identify", "Please identify club head")

    def modify_feet_line_method(self):
        response = messagebox.askquestion("Modify Feet Line Method", "Do you want to change the feet line method to toe line?")
        if response == "yes":
            self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_TOE_LINE

    def input_shoulder_tilt(self):
        self.shoulder_tilt.data = simpledialog.askfloat("Shoulder Tilt", prompt="What is the shoulder tilt?")
        self.shoulder_tilt.warning = ""

    def determine_calibration_type(self):
        if self.image_type == golfkeypoints_pb2.FACE_ON:
            if self.axes_calibration_image is not None:
                return golfkeypoints_pb2.CalibrationType.FULL_CALIBRATION
            else:
                return golfkeypoints_pb2.CalibrationType.NO_CALIBRATION
        else:
            if self.axes_calibration_image is None:
                return golfkeypoints_pb2.CalibrationType.NO_CALIBRATION
            elif self.axes_calibration_image is not None and self.vanishing_point_calibration_image is not None:
                return golfkeypoints_pb2.CalibrationType.FULL_CALIBRATION
            elif self.vanishing_point_calibration_image is None:
                return golfkeypoints_pb2.CalibrationType.AXES_CALIBRATION_ONLY
    
    def calibrate_image(self):
        try:
            response = self.golfkeypoints_client.calibrate_input_image(session_token=self.session_token, input_image_id=self.curr_input_image_id, calibration_type=self.determine_calibration_type(), feet_line_method=self.feet_line_method, calibration_image_axes=self.axes_calibration_image, calibration_image_vanishing_point=self.vanishing_point_calibration_image, golf_ball=self.golf_ball, club_butt=self.club_butt, club_head=self.club_head, shoulder_tilt=self.shoulder_tilt)
            messagebox.showinfo("Calibrate Input Image", f"Calibrate input image successful: {response}, calculate golf keypoints next")
        except grpc.RpcError as e:
            messagebox.showerror("Calibrate Input Image", f"Calibrate input image failed: {e.code()}: {e.details()}")    

    def calculate_golf_keypoints(self):
        try:
            response = self.golfkeypoints_client.calculate_golf_keypoints(session_token=self.session_token, input_image_id=self.curr_input_image_id)
            messagebox.showinfo("Calculate Golf Keypoints", f"Calculate Golf Keypoints successful")
            if response.output_image is not None:
                self.process_golf_keypoints(response.output_image, response.golf_keypoints)
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
        buffer = BytesIO(output_image)
        img = Image.open(buffer)
        self.display_image(img)
        messagebox.showinfo("Golf Keypoints", f"{golf_keypoints}")
        self.body_keypoints = golf_keypoints.body_keypoints
        incorrect = messagebox.askyesno("Body Keypoints Update", "Are there body keypoints that computervision identified incorrectly?")
        if incorrect:
            self.select_body_keypoints_to_update(self.body_keypoints)

    body_pose_field_descriptors = common_pb2.Body25PoseKeypoints.DESCRIPTOR.fields

    def on_mousewheel(self, canvas, event):
        canvas.yview_scroll(int(-1 * (event.delta / 120)), "units")

    def select_body_keypoints_to_update(self, body_keypoints):
        # create popup window
        popup = tk.Toplevel(self)
        popup.wm_title("Body Keypoints Window")
        # create canvas inside popup window
        popup_canvas = tk.Canvas(popup)
        popup_canvas.grid(row=0, column=0, sticky="nsew")
        # create scrollbar in popup window that controls canvas 
        scrollbar = tk.Scrollbar(popup, orient="vertical", command=popup_canvas.yview)
        popup_canvas.configure(yscrollcommand=scrollbar.set)
        scrollbar.grid(row=0, column=1, sticky="ns")
        # create frame for content inside the canvas
        content_frame = tk.Frame(popup_canvas)
        popup_canvas.create_window((0, 0), window=content_frame, anchor="nw")
        popup_canvas.bind_all("<MouseWheel>", partial(self.on_mousewheel, popup_canvas))
        # create buttons for each body keypoint for selections
        idx = 0
        for field in self.body_pose_field_descriptors:
            name = field.name
            body_keypoint_value = getattr(body_keypoints, name)
            button = tk.Button(content_frame, text=f"Modify {name}: {body_keypoint_value}", command=partial(self.update_body_keypoint, name))
            button.grid(row=idx, column=0, padx=5, pady=5, sticky="w")
            idx += 1
        done_button = tk.Button(content_frame, text="Done Updating Body Keypoints", command=partial(self.update_body_keypoints, popup))
        done_button.grid(row=idx, column=0, padx=5, pady=5, sticky="w")
        return 
    
    def update_body_keypoint(self, field_name):
        x = simpledialog.askfloat("New Value ", prompt=f"What is the new x value for {field_name}")
        y = simpledialog.askfloat("New Value ", prompt=f"What is the new y value for {field_name}")
        field = getattr(self.body_keypoints, field_name)
        setattr(field, "x", x)
        setattr(field, "y", y)
        setattr(field, "confidence", 1.0)

    def update_body_keypoints(self, popup):
        popup.destroy
        # bind mousewheel back to whole canvas
        self.whole_canvas.bind_all("<MouseWheel>", partial(self.on_mousewheel, self.whole_canvas))
        try:
            response = self.golfkeypoints_client.update_body_keypoints(self.session_token, self.curr_input_image_id, self.body_keypoints)
            messagebox.showinfo("Update Body Keypoints", f"Update Body Keypoints successful: {response}")
        except grpc.RpcError as e:
            messagebox.showerror("Update Body Keypoints", f"Update Body keypoints failed: {e.code()}: {e.details()}")
        

    def delete_input_image(self):
        try:
            response = self.golfkeypoints_client.delete_input_image(self.session_token, self.curr_input_image_id)
            messagebox.showinfo("Delete Input Image", f"Successfully deleted input image {response}")
            # clear canvas
            self.clear_canvas()
        except grpc.RpcError as e:
            messagebox.showerror("Delete Input Image", f"Delete input image failed: {e.code()}: {e.details()}")

    def delete_golf_keypoints(self):
        try:
            response = self.golfkeypoints_client.delete_golf_keypoints(self.session_token, self.curr_input_image_id)
            messagebox.showinfo("Delete Golf Keypoints", f"Successfully deleted golf keypoints {response}")
            # go back to original input image
            self.display_image(self.curr_input_image)
        except grpc.RpcError as e:
            messagebox.showerror("Delete Golf Keypoints", f"Delete golf keypoints failed: {e.code()}: {e.details()}")

    def draw_circle(self, x, y, color):
        radius = 3
        x1 = x - radius
        y1 = y - radius
        x2 = x + radius
        y2 = y + radius
        return self.canvas.create_oval(x1, y1, x2, y2, fill=color, outline=color)


    def draw_clubhead(self, x, y, color):
        x2 = x + 10
        y2 = y + 5
        return self.canvas.create_oval(x, y, x2, y2, fill=color, outline=color)

    def erase_circle(self, circle_id):
        self.canvas.delete(circle_id)


    def on_click_on_input_image(self, event):
        x = event.x
        y = event.y
        scaled_x = x*4
        scaled_y = y*4
        messagebox.showinfo("Clicked", f"Clicked at x: {x}, y: {y}, scaled x: {scaled_x}, scaled y: {scaled_y}")
        match self.identify_mode:
            case self.IdentifyMode.GOLFBALL:
                circle_id = self.draw_circle(x, y, "red")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the golf ball?")
                if ok:
                    self.golf_ball = common_pb2.Keypoint(x=scaled_x, y=scaled_y, confidence=1.0)
                    self.identify_mode = self.IdentifyMode.NONE
                    self.canvas.unbind("<Button-1>")
                else:
                    self.erase_circle(circle_id)
            case self.IdentifyMode.CLUBBUTT:
                circle_id = self.draw_circle(x, y, "blue")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the butt end of the club?")
                if ok:
                    self.club_butt = common_pb2.Keypoint(x=scaled_x, y=scaled_y, confidence=1.0)
                    self.identify_mode = self.IdentifyMode.NONE
                    self.canvas.unbind("<Button-1>")
                else:
                    self.erase_circle(circle_id)
            case self.IdentifyMode.CLUBHEAD:
                circle_id = self.draw_clubhead(x, y, "green")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the clubhead")
                if ok:
                    self.club_head = common_pb2.Keypoint(x=scaled_x, y=scaled_y, confidence=1.0)
                    self.identify_mode = self.IdentifyMode.NONE
                    self.canvas.unbind("<Button-1>")
                else:
                    self.erase_circle(circle_id)
            case self.IdentifyMode.NONE:
                messagebox.showerror("Clicked", "Please click one of calibration buttons")


    def get_image_bytes(self, image):
        buffer = io.BytesIO()
        image.save(buffer, format='JPEG')
        return buffer.getvalue()  
    
    def clear_canvas(self):
        self.canvas.delete("all")
    