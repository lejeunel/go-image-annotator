import typer
from .images import app as images_app
from .users import app as users_app

app = typer.Typer()
app.add_typer(images_app, name="images")
app.add_typer(users_app, name="users")


def main():
    app()


if __name__ == "__main__":
    app()
