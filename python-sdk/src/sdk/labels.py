import typer
from .client import create_client, Client
from typing import Any

app = typer.Typer(help="list and inspect labels")


class Labels:
    def __init__(self, client: Client):
        self.client = client

    def get(self, name: str):
        return self.client._request(
            "GET",
            f"/labels/{name}",
        )

    def list(
        self,
        page: int = 1,
        pagesize: int = 10,
    ) -> dict:
        params: dict[str, Any] = {
            "page": page,
            "page_size": pagesize,
        }
        return self.client._request(
            "GET",
            "/labels",
            params=params,
        )


@app.command(help="fetch a page of labels")
def list(page: int = 1, pagesize: int = 10):
    print(
        Labels(create_client()).list(
            page=page,
            pagesize=pagesize,
        )
    )


@app.command(help="inspect a label by its name")
def get(name: str):
    print(Labels(create_client()).get(name))
