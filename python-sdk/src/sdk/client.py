import requests
import os


class Client:
    def __init__(self, api_url: str, session: requests.Session):
        self.session = session
        self.api_url = api_url

    def _request(self, method: str, path: str, **kwargs) -> dict:
        url = f"{self.api_url.rstrip('/')}/{path.lstrip('/')}"
        response = self.session.request(method, url, **kwargs)
        self._raise_for_status(response)

        return response.json()

    def _request_bytes(self, method: str, path: str, **kwargs) -> bytes:

        url = f"{self.api_url.rstrip('/')}/{path.lstrip('/')}"

        response = self.session.request(method, url, **kwargs)
        self._raise_for_status(response)
        return response.content

    def _raise_for_status(self, r: requests.Response):
        try:
            r.raise_for_status()
        except requests.HTTPError as e:
            try:
                error = r.json()["error"]
            except (ValueError, KeyError):
                error = r.text

            raise requests.HTTPError(
                f"API request failed ({r.status_code} {r.reason}): {error}",
                response=r,
            ) from e


def create_session() -> requests.Session:
    session = requests.Session()
    token = os.environ["GOIA_API_TOKEN"]
    if token == "":
        raise ValueError("GOIA_API_TOKEN env variable has not been set")

    session.headers.update(
        {
            "Authorization": f"Bearer {token}",
        }
    )
    return session


def create_client() -> Client:
    api_url = os.environ["GOIA_API_URL"]
    if api_url == "":
        raise ValueError("GOIA_API_URL env variable has not been set")
    return Client(api_url, create_session())
