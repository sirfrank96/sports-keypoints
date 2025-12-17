import tkinter as tk
from tkinter import messagebox, filedialog, simpledialog
from PIL import ImageTk, Image
from functools import partial
import io
from io import BytesIO
import grpc
from enum import Enum

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

        self.main_app_label = tk.Label(self, text="This is the main app")
        self.main_app_label.grid(row=0, column=0, padx=5, pady=5, sticky="w")

        self.select_new_input_image_button = tk.Button(self, text="Select New Input Image", command=self.select_new_input_image)
        self.select_new_input_image_button.grid(row=1, column=0, padx=5, pady=5, sticky="w")

        self.show_input_images_button = tk.Button(self, text="Show Previous Input Images", command=self.show_previous_input_images)
        self.show_input_images_button.grid(row=2, column=0, padx=5, pady=5, sticky="w")

        self.read_user_button = tk.Button(self, text="Show User", command=self.read_user)
        self.read_user_button.grid(row=3, column=0, padx=5, pady=5, sticky="w")
        
        self.update_user_button = tk.Button(self, text="Update User", command=self.update_user)
        self.update_user_button.grid(row=4, column=0, padx=5, pady=5, sticky="w")

        self.delete_user_button = tk.Button(self, text="Delete User", command=self.delete_user)
        self.delete_user_button.grid(row=5, column=0, padx=5, pady=5, sticky="w")

        #1/4 size of image
        self.canvas = tk.Canvas(self, width=270, height=600, bg='white')
        self.canvas.grid(row=8, column=1, padx=5, pady=5, sticky="w")

        self.curr_input_image_id = ""
        self.identify_mode = self.IdentifyMode.NONE
        self.golf_ball = None
        self.club_butt = None
        self.club_head = None
        self.image_type = golfkeypoints_pb2.ImageType.IMAGE_TYPE_UNSPECIFIED 
        self.axes_calibration_image = None
        self.vanishing_point_calibration_image = None
        self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_HEEL_LINE
        self.shoulder_tilt = common_pb2.Double(data=0, warning="no shoulder tilt")

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
                response = messagebox.askquestion("FaceOn or DTL", "Is this image Face On? (Yes for Face On, No for DTL)")
                self.image_type = golfkeypoints_pb2.ImageType.FACE_ON if response == "yes" else golfkeypoints_pb2.ImageType.DTL
                response = self.golfkeypoints_client.upload_input_image(session_token=self.session_token, image_type= self.image_type, image=bytes)
                messagebox.showinfo("Successfully Upload Input", f"response is {response}")
                self.curr_input_image_id = response.input_image_id
                self.display_input_image(img)
            except grpc.RpcError as e:
                messagebox.showerror("Upload input image failed", f"Could not upload input image: {e.code()}: {e.details()}")
            
    def show_previous_input_images(self):
        self.canvas.delete("all")
        try:
            response = self.golfkeypoints_client.list_input_images_for_user(session_token=self.session_token)
            for i, input_image_id in enumerate(response.input_image_ids):
                curr_button = tk.Button(self.canvas, text=f"{input_image_id}", command=partial(self.read_input_image, input_image_id))
                self.canvas.create_window(100, 30+(i*50), window=curr_button)
        except grpc.RpcError as e:
            messagebox.showerror("List Images Failed", f"Could not get a list of images: {e.code()}: {e.details()}")
    
    def read_input_image(self, input_image_id):
        try:
            response = self.golfkeypoints_client.read_input_image(session_token=self.session_token, input_image_id=input_image_id)
            messagebox.showinfo("Show Image", f"Response length of image: {len(response.image)}")
            buffer = BytesIO(response.image)
            img = Image.open(buffer)
            self.curr_input_image_id = input_image_id
            self.display_input_image(img)
        except grpc.RpcError as e:
            messagebox.showerror("Show Image Failed", f"Could not get image: {e.code()}: {e.details()}")
    
    def display_image(self, image):
        self.canvas.delete("all")
        #1/4 the size to display on canvas
        resized_img = image.resize((270, 600), Image.Resampling.LANCZOS)
        # Convert the image to a PhotoImage object
        photo = ImageTk.PhotoImage(resized_img)
        self.canvas.create_image(0, 0, anchor=tk.NW, image=photo)
        self.canvas.image = photo

    def display_input_image(self, image):
        self.display_image(image)

        self.open_button = tk.Button(self, text="Get Axes Calibration Image", command=self.get_axes_calibration_image)
        self.open_button.grid(row=0, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Get Vanishing Point Calibration Image", command=self.get_vanishing_point_calibration_image)
        self.open_button.grid(row=1, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Identify Golf Ball", command=self.identify_golf_ball)
        self.open_button.grid(row=2, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Identify Club Butt", command=self.identify_club_butt)
        self.open_button.grid(row=3, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Identify Club Head", command=self.identify_club_head)
        self.open_button.grid(row=4, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Modify Feet Line Method", command=self.modify_feet_line_method)
        self.open_button.grid(row=5, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Input Shoulder Tilt For DTL", command=self.input_shoulder_tilt)
        self.open_button.grid(row=6, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Calibrate Image", command=partial(self.calibrate_image))
        self.open_button.grid(row=7, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Calculate Golf Keypoints", command=partial(self.calculate_golf_keypoints))
        self.open_button.grid(row=8, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Delete Input Image", command=partial(self.delete_input_image))
        self.open_button.grid(row=9, column=2, padx=5, pady=5, sticky="w")

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
                buffer = BytesIO(response.output_image)
                img = Image.open(buffer)
                self.display_image(img)
                messagebox.showinfo("Golf Keypoints", f"{response.golf_keypoints}")
                incorrect = messagebox.askyesno("Body Keypoints Update", "Are there body keypoints that computervision identified incorrectly?")
                if incorrect:
                    self.update_body_keypoints(response.golf_keypoints.body_keypoints)
        except grpc.RpcError as e:
            messagebox.showerror("Calculate Golf Keypoints", f"Calculate Golf Keypoints failed: {e.code()}: {e.details()}") 

    def update_body_keypoints(self, body_keypoints):
        #display a button for each body keypoint
        #user clicks on button
        #self.canvas.bind button1
        #add bodykeypoints identifymode
        return 

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

    def delete_input_image(self, input_image_id):
        try:
            response = self.golfkeypoints_client.delete_input_image(self.session_token, input_image_id)
            messagebox.showinfo("Delete Input Image", f"Successfully deleted input image {response}")
        except grpc.RpcError as e:
            messagebox.showerror("Delete Input Image", f"Delete input image failed: {e.code()}: {e.details()}")

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
    