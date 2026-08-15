import requests
import os


class Client:
    def __init__(self, api_url: str, session: requests.Session):
        self.session = session
        self.api_url = api_url

    def _request(self, method: str, path: str, **kwargs) -> dict:
        url = f"{self.api_url.rstrip('/')}/{path.lstrip('/')}"
        response = self.session.request(method, url, **kwargs)
        try:
            response.raise_for_status()
        except requests.HTTPError as e:
            try:
                message = response.json()["error"]
            except (ValueError, KeyError):
                message = response.text

            raise requests.HTTPError(
                f"{e}: {message}",
                response=response,
            ) from e

        return response.json()

    def whoami(self) -> dict:
        return self._request("GET", "/whoami")


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
