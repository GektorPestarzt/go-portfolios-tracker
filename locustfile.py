from time import sleep
from enum import Enum
import random
from secrets import token_hex

from locust import HttpUser, task, between
from locust.exception import RescheduleTask


class FamilyStatus(str, Enum):
    SINGLE = "SINGLE"
    MARRIED = "MARRIED"
    DIVORCED = "DIVORCED"
    WIDOWED = "WIDOWED"


class NormalUser(HttpUser):
    wait_time = between(0.2, 0.5)
    base_path = "/api/account/tinkoff"

    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.headers = {
            "Content-Type": "application/json",
            "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjU5OTUwNjIsInVzZXJuYW1lIjoicnVzbGFuIn0.4-b64kcZuHKK11lUWAhzBSZqkLyi4Wini3SVKkaDBQE"
        }
        self.profile_ids = []

    @task(1)
    def create_profile(self):
        payload = {
            "token": "t.PSA_VJvF2o48h7XaBtMnL52wNzxXgIN0In1Owgg9cHybqo4PVZJs6opJi-mb5wAAIOSRGZwKWVq8MHEx41BMQQ"
        }
        with self.client.post("/api/accounts/tinkoff", json=payload, headers=self.headers, catch_response=True) as response:
            if response.status_code == 201:
                profile_id = response.json().get("id")
                if profile_id:
                    self.profile_ids.append(profile_id)
            else:
                response.failure(f"Failed to create profile: {response.status_code}")

    @task(2)
    def update_profile(self):
        if self.profile_ids:
            profile_id = random.choice(self.profile_ids)
            payload = {
                "token": "token_number"
            }
            with self.client.put(f"/api/accounts/tinkoff/{profile_id}", json=payload, headers=self.headers, catch_response=True) as response:
                if response.status_code != 201:
                    response.failure(f"Failed to update profile: {response.status_code}")
        else:
            self.create_profile()