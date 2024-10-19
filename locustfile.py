# locustfile.py

from locust import HttpUser, task, between


class MyUser(HttpUser):
    wait_time = between(1, 5)

    @task
    def load_main_page(self):
        self.client.get("/")
