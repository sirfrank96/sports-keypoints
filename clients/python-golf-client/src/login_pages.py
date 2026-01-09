# Python
import grpc

# Tkinter
import tkinter as tk
from tkinter import messagebox

# Internal
import main_page
import frame_wrapper as fw


class InitialPage(fw.FrameWrapper):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        super().__init__(parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        # add buttons
        self.create_user_button = self.add_button(text="Create a new user", command=self.go_to_create_user_page, row=0, col=0, padx=10, pady=10)
        self.login_button = self.add_button(text="Login with existing user", command=self.go_to_login_page, row=1, col=0, padx=10, pady=10)

    def go_to_create_user_page(self):
        create_user_page = CreateUserPage(self.parent, self.controller, self.user_client, self.golfkeypoints_client)
        self.controller.show_frame(create_user_page)
    
    def go_to_login_page(self):
        login_page = LoginPage(self.parent, self.controller, self.user_client, self.golfkeypoints_client)
        self.controller.show_frame(login_page)        

class CreateUserPage(fw.FrameWrapper):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        super().__init__(parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        # add labels, entries, and button
        self.username_label = self.add_label(text="Username:", row=0, col=0, padx=5, pady=5)
        self.username_entry = self.add_entry(row=0, col=1, padx=5, pady=5)
        self.password_label = self.add_label(text="Password:", row=1, col=0, padx=5, pady=5)
        self.password_entry = self.add_entry(row=1, col=1, padx=5, pady=5, show="*") # Mask the password
        self.email_label = self.add_label(text="Email:", row=2, col=0, padx=5, pady=5)  
        self.email_entry = self.add_entry(row=2, col=1, padx=5, pady=5)
        self.login_button = self.add_button(text="Create New User", command=self.create_user, row=3, col=0, padx=10, pady=10)
        
    def create_user(self):
        username = self.username_entry.get()
        password = self.password_entry.get()
        email = self.email_entry.get()
        # replace with your actual login logic (e.g., database check)
        if username != "" and password != "" and email != "":
            try:
                response = self.user_client.create_user(username, password, email)
                messagebox.showinfo("CreateUser Response", f"Response: {response}")
                messagebox.showinfo("Create User Successful", "Welcome!")
                login_page = LoginPage(self.parent, self.controller, user_client=self.user_client, golfkeypoints_client=self.golfkeypoints_client)
                self.controller.show_frame(login_page)
            except grpc.RpcError as e:
                messagebox.showerror("Create User Failed", f"Invalid username: {e.code()}: {e.details()}")
        else:
            messagebox.showerror("Create User Failed", "Must input something for username, password, and email")

            

class LoginPage(fw.FrameWrapper):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        super().__init__(parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        # add labels, entries, and button
        self.username_label = self.add_label(text="Username:", row=0, col=0, padx=5, pady=5)
        self.username_entry = self.add_entry(row=0, col=1, padx=5, pady=5)
        self.password_label = self.add_label(text="Password:", row=1, col=0, padx=5, pady=5)
        self.password_entry = self.add_entry(row=1, col=1, padx=5, pady=5, show="*") # Mask the password
        self.login_button = self.add_button(text="Login", command=self.login, row=2, col=0, padx=10, pady=10)
    
    def login(self):
        username = self.username_entry.get()
        password = self.password_entry.get()
        # replace with your actual login logic (e.g., database check)
        if username != "" and password != "":
            try: 
                response = self.user_client.register_user(username, password)
                session_token = response.session_token
                messagebox.showinfo("Login Successful", "Welcome!")
                main_app_page = main_page.MainAppPage(self.parent, self.controller, user_client=self.user_client, golfkeypoints_client=self.golfkeypoints_client, session_token=session_token)
                self.controller.show_frame(main_app_page)
            except grpc.RpcError as e:
                messagebox.showerror("Login User Failed", f"{e.code()}: {e.details()}")
        else:
            messagebox.showerror("Login Failed", "Must input something for username and password")
