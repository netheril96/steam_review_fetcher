import requests
from typing import TypedDict


class App(TypedDict, total=False):
    appid: int
    name: str


def query_all_apps() -> list[App]:
    return requests.get("http://api.steampowered.com/ISteamApps/GetAppList/v2").json()[
        "applist"
    ]["apps"]


if __name__ == "__main__":
    print(query_all_apps())
