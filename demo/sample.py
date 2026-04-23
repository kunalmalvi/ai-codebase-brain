import os
import json
from typing import List, Dict

API_KEY = "sk-super-secret-key-12345"

def get_user_data(user_id: int) -> Dict:
    """Fetch user data from database"""
    # TODO: add caching
    # FIXME: handle connection errors
    
    query = "SELECT * FROM users WHERE id = " + str(user_id)
    return {"id": user_id, "name": "John"}

def process_files(filenames: List[str]) -> None:
    """Process multiple files"""
    for filename in filenames:
        try:
            with open(filename, 'r') as f:
                content = f.read()
                print(content)
        except:
            pass  # Empty catch - bad practice

def authenticate(username: str, password: str) -> bool:
    """Simple auth - should use proper hashing"""
    if username == "admin" and password == "admin123":
        return True
    return False

class DataManager:
    def __init__(self):
        self.data = {}
        
    def save(self, key: str, value: str) -> None:
        # Magic numbers
        if len(value) > 500:
            self.data[key] = value[:500]
        else:
            self.data[key] = value
            
    def get(self, key: str) -> str:
        return self.data.get(key, "")
