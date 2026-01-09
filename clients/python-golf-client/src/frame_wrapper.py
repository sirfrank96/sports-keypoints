# Tkinter
import tkinter as tk

class FrameWrapper(tk.Frame):
    def __init__(self, parent):
        super().__init__(parent)

    def add_label(self, text, row, col, padx, pady, sticky="w"):
        label = tk.Label(self, text=text)
        label.grid(row=row, column=col, padx=padx, pady=pady, sticky=sticky)
        return label

    def add_button(self, text, command, row, col, padx, pady, sticky="w"):
        button = tk.Button(self, text=text, command=command)
        button.grid(row=row, column=col, padx=padx, pady=pady, sticky=sticky)
        return button
    
    def add_entry(self, row, col, padx, pady, show=""):
        entry = tk.Entry(self)
        if show != "":
            entry = tk.Entry(self, show=show)
        entry.grid(row=row, column=col, padx=padx, pady=pady)
        return entry
