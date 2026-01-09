# Python
from PIL import ImageTk, Image

# Tkinter
import tkinter as tk

# Internal
import frame_wrapper as fw

class CanvasWrapper(tk.Canvas):
    def __init__(self, parent, width=0, height=0, bg='', row=0, col=0, padx=0, pady=0, sticky='nsew'):
        # if width and height are not set, get parent width and height
        if width == 0:
            parent.update()
            width = parent.winfo_width()
        if height == 0:
            parent.update()
            width = parent.winfo_height()
        # if background color is not set, use default
        if bg == '':
            super().__init__(parent, width=width, height=height)
        else:
            super().__init__(parent, width=width, height=height, bg=bg)
        self.grid(row=row, column=col, padx=padx, pady=pady, sticky=sticky)
        # set instance vars
        self.parent = parent
        self.width = width
        self.height = height
        self.content_frame = None
        self.scrollbar = None

    def make_scrollable(self):
        self.scrollbar = tk.Scrollbar(self.parent, orient="vertical", command=self.yview)
        self.configure(yscrollcommand=self.scrollbar.set)
        self.scrollbar.grid(row=0, column=1, sticky="ns")
        self.bind_all("<MouseWheel>", self.on_mousewheel)

    def on_mousewheel(self, event):
        self.yview_scroll(int(-1 * (event.delta / 120)), "units")

    def create_content_frame_in_canvas(self):
        self.content_frame = fw.FrameWrapper(self)
        self.content_frame.grid_columnconfigure(0, weight=1)
        self.content_frame.grid_rowconfigure(0, weight=1)
        self.create_window((0, 0), window=self.content_frame, anchor="nw")
        return self.content_frame

    def display_an_image(self, image):
        self.clear_canvas()
        resized_img = image.resize((self.width, self.height), Image.Resampling.LANCZOS)
        # convert the image to a PhotoImage object
        photo = ImageTk.PhotoImage(resized_img)
        self.create_image(0, 0, anchor=tk.NW, image=photo)
        self.image = photo

    def draw_circle(self, x, y, color):
        radius = 3
        x1 = x - radius
        y1 = y - radius
        x2 = x + radius
        y2 = y + radius
        return self.create_oval(x1, y1, x2, y2, fill=color, outline=color)

    def draw_clubhead(self, x, y, color):
        x2 = x + 10
        y2 = y + 5
        return self.create_oval(x, y, x2, y2, fill=color, outline=color)

    def erase_circle(self, circle_id):
        self.delete(circle_id)

    def draw_line(self, x1, y1, x2, y2, color):
        return self.create_line(x1, y1, x2, y2, fill=color, width=3)

    def erase_line(self, line_id):
        self.delete(line_id)

    def clear_canvas(self):
        self.delete("all")
