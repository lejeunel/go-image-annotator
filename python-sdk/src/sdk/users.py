import typer
from .client import create_client, Client

app = typer.Typer(help="manage users and inspect identity")


class Users:
    def __init__(self, client: Client):
        self.client = client

    def whoami(self):
        return self.client._request("GET", f"/whoami")


@app.command(help="get current users's identity")
def whoami():
    print(Users(create_client()).whoami())
