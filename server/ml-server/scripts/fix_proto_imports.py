import re
from pathlib import Path

# Get the script's directory, then go up one level to ml-server root
script_dir = Path(__file__).parent
project_root = script_dir.parent
gen_dir = project_root / "src" / "gen"

# Fix both *_pb2.py, *_pb2_grpc.py, and *.pyi files
for proto_file in gen_dir.glob("*_pb2*.py*"):  # Catches .py and .pyi
    content = proto_file.read_text()
    
    # Change: import something_pb2 as something__pb2
    # To: from . import something_pb2 as something__pb2
    # Handles both "as something__pb2" and "as _something_pb2" (pyi style)
    content = re.sub(
        r'^import (\w+_pb2(?:_grpc)?) as ',
        r'from . import \1 as ',
        content,
        flags=re.MULTILINE
    )
    proto_file.write_text(content)

print("Fixed proto imports!")
