import sys
import json
import fitz  
import re
import os

def extract_text(file_path):
    if not os.path.exists(file_path):
        raise FileNotFoundError(file_path)

    text = ""
    if file_path.lower().endswith(".pdf"):
        with fitz.open(file_path) as doc:
            for page in doc:
                text += page.get_text()
    else:
        raise ValueError("Unsupported file type")

    return text

def extract_contact(text: str):
    name = text.split("\n")[0]
    emails = re.findall(r"[\w\.-]+@[\w\.-]+", text)
    phones = re.findall(r"\+?\d[\d\s\-]{8,}\d", text)
    return {
        "name": name.strip(),
        "email": emails,
        "number": phones,
        "location": "",
        "social": {}
    }

def parse_resume(file_path):
    text = extract_text(file_path)
    contact = extract_contact(text)

    result = {
        "raw": {"text": text[:5000]},
        "metadata": {"file_name": os.path.basename(file_path)},
        "sections": {
            "contact": {
                "type": "contact",
                "content": contact
            }
        }
    }
    print(json.dumps(result))  # Output JSON to stdout

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "file path required"}))
        sys.exit(1)

    file_path = sys.argv[1]
    parse_resume(file_path)
