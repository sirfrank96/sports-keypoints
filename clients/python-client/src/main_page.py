import tkinter as tk
from tkinter import messagebox, filedialog
from PIL import ImageTk, Image
from functools import partial
import io
from io import BytesIO
import grpc

import golfkeypoints_pb2
import golf_keypoints_client as gc


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

        self.open_button = tk.Button(self, text="Select New Input Image", command=self.select_new_input_image)
        self.open_button.grid(row=1, column=0, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Show Previous Input Images", command=self.show_previous_input_images)
        self.open_button.grid(row=2, column=0, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Delete User", command=self.delete_user)
        self.open_button.grid(row=3, column=0, padx=5, pady=5, sticky="w")

        #1/4 size of image
        self.canvas = tk.Canvas(self, width=270, height=600, bg='white')
        self.canvas.grid(row=4, column=1, padx=5, pady=5, sticky="w")
        self.golf_ball = None
        self.club_butt = None
        self.club_head = None
        self.validClicks = 0
        self.session_token = session_token
        self.image_type = golfkeypoints_pb2.ImageType.IMAGE_TYPE_UNSPECIFIED 
        self.axes_calibration_image = None
        self.vanishing_point_calibration_image = None
        self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_HEEL_LINE
    
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
                bytes = self.getImageBytes(img)
                messagebox.showinfo("Bytes", f"Length of image bytes is: {len(bytes)}, length of image raw bytes: {len(img.tobytes())}")
                response = messagebox.askquestion("FaceOn or DTL", "Is this image Face On? (Yes for Face On, No for DTL)")
                self.image_type = golfkeypoints_pb2.ImageType.FACE_ON if response == "yes" else golfkeypoints_pb2.ImageType.DTL
                response = self.golfkeypoints_client.upload_input_image(session_token=self.session_token, image_type= self.image_type, image=bytes)
                messagebox.showinfo("Successfully Upload Input", f"response is {response}")
                self.display_input_image(img, response.input_image_id)
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
            self.display_input_image(img, input_image_id)
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

    def display_input_image(self, image, input_image_id):
        self.display_image(image)
        self.canvas.bind("<Button-1>", self.on_click)
        messagebox.showinfo("Please identify golf ball, butt end of golf club, and club head") 

        self.open_button = tk.Button(self, text="Get Axes Calibration Image", command=self.get_axes_calibration_image)
        self.open_button.grid(row=0, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Get Vanishing Point Calibration Image", command=self.get_vanishing_point_calibration_image)
        self.open_button.grid(row=1, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Modify Feet Line Method", command=self.modify_feet_line_method)
        self.open_button.grid(row=2, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Calibrate Image", command=partial(self.calibrate_image, input_image_id))
        self.open_button.grid(row=3, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Calculate Golf Keypoints", command=partial(self.calculate_golf_keypoints, input_image_id))
        self.open_button.grid(row=4, column=2, padx=5, pady=5, sticky="w")

        self.open_button = tk.Button(self, text="Delete Input Image", command=partial(self.delete_input_image, input_image_id))
        self.open_button.grid(row=5, column=2, padx=5, pady=5, sticky="w")

    def get_axes_calibration_image(self):
        img = self.get_image_from_filesystem()
        if img is not None:
            bytes = self.getImageBytes(img)
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
                bytes = self.getImageBytes(img)
                self.vanishing_point_calibration_image = bytes
                messagebox.showinfo("Vanishing Point Calibration Image", "Successfully set vanishing point calibration image")
            else:
                messagebox.showerror("Vanishing Point Calibration Image", "Could not get vanishing point calibration image")

    def modify_feet_line_method(self):
        response = messagebox.askquestion("Modify Feet Line Method", "Do you want to change the feet line method to toe line?")
        if response == "yes":
            self.feet_line_method = golfkeypoints_pb2.FeetLineMethod.USE_TOE_LINE

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
    
    def calibrate_image(self, input_image_id):
        try:
            response = self.golfkeypoints_client.calibrate_input_image(session_token=self.session_token, input_image_id=input_image_id, calibration_type=self.determine_calibration_type(), feet_line_method=self.feet_line_method, calibration_image_axes=self.axes_calibration_image, calibration_image_vanishing_point=self.vanishing_point_calibration_image, golf_ball=self.golf_ball, club_butt=self.club_butt, club_head=self.club_head)
            messagebox.showinfo("Calibrate Input Image", f"Calibrate input image successful: {response}, calculate golf keypoints next")
        except grpc.RpcError as e:
            messagebox.showerror("Calibrate Input Image", f"Calibrate input image failed: {e.code()}: {e.details()}")    

    def calculate_golf_keypoints(self, input_image_id):
        try:
            response = self.golfkeypoints_client.calculate_golf_keypoints(session_token=self.session_token, input_image_id=input_image_id)
            messagebox.showinfo("Calculate Golf Keypoints", f"Calculate Golf Keypoints successful")
            if response.output_image is not None:
                buffer = BytesIO(response.output_image)
                img = Image.open(buffer)
                self.display_image(img)
                messagebox.showinfo("Golf Keypoints", f"{response.golf_keypoints}")
        except grpc.RpcError as e:
            messagebox.showerror("Calculate Golf Keypoints", f"Calculate Golf Keypoints failed: {e.code()}: {e.details()}")  
    
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


    def on_click(self, event):
        x = event.x
        y = event.y
        scaled_x = x*4
        scaled_y = y*4
        messagebox.showinfo("Clicked", f"Clicked at x: {x}, y: {y}, scaled x: {scaled_x}, scaled y: {scaled_y}")
        match self.validClicks:
            case 0:
                circle_id = self.draw_circle(x, y, "red")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the golf ball?")
                if ok:
                    self.validClicks += 1
                else:
                    self.erase_circle(circle_id)
            case 1:
                circle_id = self.draw_circle(x, y, "blue")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the butt end of the club?")
                if ok:
                    self.validClicks += 1
                else:
                    self.erase_circle(circle_id)
            case 2:
                circle_id = self.draw_clubhead(x, y, "green")
                ok = messagebox.askokcancel("Clicked", "Is this the correct spot for the clubhead")
                if ok:
                    messagebox.showinfo("Got club and ball info", "Blah")
                    self.validClicks = 0
                    #API: Calibrate Input Image
                    #try:
                        #bytes = self.getImageBytes(self.img)
                        #messagebox.showinfo("Bytes", f"Length of image bytes is: {len(bytes)}, length of image raw bytes: {len(self.img.tobytes())}")
                        #response = self.golfkeypoints_client.upload_input_image(session_token=self.session_token, image_type=golfkeypoints_pb2.ImageType.DTL, image=bytes)
                        #messagebox.showinfo("Successfully Upload Input", f"response is {response}")
                    #except grpc.RpcError as e:
                        #messagebox.showerror("Upload input image failed", f"Could not upload input image: {e.code()}: {e.details()}") 
                else:
                    self.erase_circle(circle_id)

    def getImageBytes(self, image):
        buffer = io.BytesIO()
        image.save(buffer, format='JPEG')
        return buffer.getvalue()  