import tkinter as tk
from tkinter import messagebox
import grpc
from functools import partial

import user_client as uc
import main_page


class InitialPage(tk.Frame):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        tk.Frame.__init__(self, parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        #Create new user
        self.create_user_button = tk.Button(self, text="Create a new user", command=self.go_to_create_user_page)
        self.create_user_button.grid(row=0, column=0, columnspan=2, pady=10)
        # Login Button
        self.login_button = tk.Button(self, text="Login with existing user", command=self.go_to_login_page)
        self.login_button.grid(row=1, column=0, columnspan=2, pady=10)

    def go_to_create_user_page(self):
        create_user_page = CreateUserPage(self.parent, self.controller, self.user_client, self.golfkeypoints_client)
        self.controller.show_frame(create_user_page)
    
    def go_to_login_page(self):
        login_page = LoginPage(self.parent, self.controller, self.user_client, self.golfkeypoints_client)
        self.controller.show_frame(login_page)        

class CreateUserPage(tk.Frame):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        tk.Frame.__init__(self, parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        #same as login, but button is createuser instead of login
        #messages are different: success -> created user, fail -> unable to create user
        # Username Label and Entry
        self.username_label = tk.Label(self, text="Username:")
        self.username_label.grid(row=0, column=0, padx=5, pady=5, sticky="w")
        #username_label.place(relx=0.5, rely=0.5, anchor=tk.CENTER)
        self.username_entry = tk.Entry(self)
        self.username_entry.grid(row=0, column=1, padx=5, pady=5)

        # Password Label and Entry
        self.password_label = tk.Label(self, text="Password:")
        self.password_label.grid(row=1, column=0, padx=5, pady=5, sticky="w")
        self.password_entry = tk.Entry(self, show="*") # Mask the password
        self.password_entry.grid(row=1, column=1, padx=5, pady=5)

        # Email Label and Entry
        self.email_label = tk.Label(self, text="Email:")
        self.email_label.grid(row=2, column=0, padx=5, pady=5, sticky="w")
        self.email_entry = tk.Entry(self)
        self.email_entry.grid(row=2, column=1, padx=5, pady=5)

        # Create User Button
        self.login_button = tk.Button(self, text="Create New User", command=self.create_user)
        self.login_button.grid(row=3, column=0, columnspan=2, pady=10)
    
    def create_user(self):
        username = self.username_entry.get()
        password = self.password_entry.get()
        email = self.email_entry.get()

        # Replace with your actual login logic (e.g., database check)
        if username != "" and password != "" and email != "":
            #API: CreateUser
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

            

class LoginPage(tk.Frame):
    def __init__(self, parent, controller, user_client, golfkeypoints_client):
        tk.Frame.__init__(self, parent)
        self.user_client = user_client
        self.golfkeypoints_client = golfkeypoints_client
        self.parent = parent
        self.controller = controller
        # Username Label and Entry
        self.username_label = tk.Label(self, text="Username:")
        self.username_label.grid(row=0, column=0, padx=5, pady=5, sticky="w")
        #username_label.place(relx=0.5, rely=0.5, anchor=tk.CENTER)
        self.username_entry = tk.Entry(self)
        self.username_entry.grid(row=0, column=1, padx=5, pady=5)

        # Password Label and Entry
        self.password_label = tk.Label(self, text="Password:")
        self.password_label.grid(row=1, column=0, padx=5, pady=5, sticky="w")
        self.password_entry = tk.Entry(self, show="*") # Mask the password
        self.password_entry.grid(row=1, column=1, padx=5, pady=5)

        # Login Button
        self.login_button = tk.Button(self, text="Login", command=self.login)
        self.login_button.grid(row=2, column=0, columnspan=2, pady=10)
    
    def login(self):
        username = self.username_entry.get()
        password = self.password_entry.get()

        # Replace with your actual login logic (e.g., database check)
        if username != "" and password != "":
            #API: RegisterUser
            try: 
                response = self.user_client.register_user(username, password)
                session_token = response.session_token
                messagebox.showinfo("Login User Response", f"Response: {response}")
                messagebox.showinfo("Login Successful", "Welcome!")
                main_app_page = main_page.MainAppPage(self.parent, self.controller, user_client=self.user_client, golfkeypoints_client=self.golfkeypoints_client, session_token=session_token)
                self.controller.show_frame(main_app_page)
            except grpc.RpcError as e:
                messagebox.showerror("Login User Failed", f"{e.code()}: {e.details()}")
        else:
            messagebox.showerror("Login Failed", "Must input something for username and password")
