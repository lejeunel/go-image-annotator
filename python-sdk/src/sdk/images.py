import typer
from .client import create_client, Client
from typing import Any, BinaryIO
from pathlib import Path
import json
import mimetypes

app = typer.Typer(help="list and download images")


class Images:
    def __init__(self, client: Client):
        self.client = client

    def get(self, collection: str, id: str):
        return self.client._request(
            "GET",
            f"/images/{collection}/{id}",
        )

    def list(
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

    def download(self, id: str) -> bytes:
        return self.client._request_bytes("GET", f"/raw/{id}")

    def ingest(self, path: Path, collection: str) -> dict:
        mimetype, _ = mimetypes.guess_type(path)
        if mimetype is None:
            raise ValueError(f"could not guess mimetype for file {path}")

        with open(path, "rb") as f:
            return self.client._request(
                "POST",
                "/images",
                files={
                    "metadata": (
                        None,
                        json.dumps({"collection": collection}),
                        "application/json",
                    ),
                    "image": (
                        None,
                        f,
                        mimetype,
                    ),
                },
            )


@app.command(help="query a page of image meta-data")
def list(filters: str = "", ordering: str = "", page: int = 1, pagesize: int = 10):
    print(
        Images(create_client()).list(
            filter=filters,
            order=ordering,
            page=page,
            pagesize=pagesize,
        )
    )


@app.command(help="get an image by its id and collection")
def get(collection: str, id: str):
    print(Images(create_client()).get(collection, id))


@app.command(help="download an image to a local path")
def download(id: str, out: Path):
    data = Images(create_client()).download(id)
    with open(out, "wb") as f:
        f.write(data)


@app.command(help="ingest an image into a collection")
def ingest(collection: str, path: Path):
    Images(create_client()).ingest(path, collection)
