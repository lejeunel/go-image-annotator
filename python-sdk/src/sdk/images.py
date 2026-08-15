import typer
from .client import create_client, Client
from typing import Any

app = typer.Typer(help="list and download images")


class Images:
    def __init__(self, client: Client):
        self.client = client

    def get(self, image_id: str, collection: str):
        return self.client._request(
            "GET",
            f"/images/{collection}/{image_id}",
        )

    def list_images(
        self,
        filter: str | None = None,
        order: str | None = None,
        page: int = 1,
        pagesize: int = 10,
    ) -> dict:
        params: dict[str, Any] = {
            "page": page,
            "page_size": pagesize,
        }
        if filter is not None:
            params["filter"] = filter
        if order is not None:
            params["order"] = order

        return self.client._request(
            "GET",
            "/images",
            params=params,
        )


@app.command(help="query a page of image meta-data")
def list(filters: str = "", ordering: str = "", page: int = 1, pagesize: int = 10):
    cli = Images(create_client())
    print(
        cli.list_images(
            filter=filters,
            order=ordering,
            page=page,
            pagesize=pagesize,
        )
    )


@app.command(help="get an image by its id and collection")
def get(collection: str, id: str):
    cli = Images(create_client())
    print(cli.get(collection, id))
