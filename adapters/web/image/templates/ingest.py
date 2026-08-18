from sdk.client import create_client
from sdk.images import Images
from pathlib import Path

PATH_TO_IMAGE = "/home/user/images/test.png"
COLLECTION = "my-collection"

if __name__ == "__main__":
    # Configure HTTP client. This looks for your
    # API token and the service URL in your env variables
    client = create_client()

    # instantiate service
    images = Images(client)

    # ingest
    r = images.ingest(Path(PATH_TO_IMAGE), COLLECTION)
