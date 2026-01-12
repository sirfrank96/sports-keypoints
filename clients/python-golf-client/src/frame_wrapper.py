# Tkinter
import tkinter as tk
import tkinter.scrolledtext as scrolledtext

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
    
    def add_text(self, text, row, col, padx, pady, sticky="nsew"):
        text_widget = tk.Text(self, wrap=tk.WORD)
        text_widget.grid(row=row, column=col, padx=padx, pady=pady, sticky=sticky)
        text_widget.insert("1.0", text)

    def add_scrolled_text(self, text, row, col, padx, pady, sticky="nsew"):
        text_widget = scrolledtext.ScrolledText(self, wrap=tk.WORD)
        text_widget.grid(row=row, column=col, padx=padx, pady=pady, sticky=sticky)
        text_widget.insert("1.0", text)
        text_widget.bind("<MouseWheel>", lambda event: self.stop_scroll_propogation(event, text_widget))

    def on_mousewheel(self, event, scrolled_text):
        scrolled_text.yview_scroll(int(-1 * (event.delta / 120)), "units")

    def stop_scroll_propogation(self, event, scrolled_text):
        self.on_mousewheel(event, scrolled_text)
        return "break"
