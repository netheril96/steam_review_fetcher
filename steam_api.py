import requests
from typing import TypedDict
import time
from tqdm import tqdm


class App(TypedDict):
    appid: int
    name: str


def query_all_apps(session: requests.Session) -> list[App]:
    return session.get("http://api.steampowered.com/ISteamApps/GetAppList/v2").json()[
        "applist"
    ]["apps"]


def query_app_details(session: requests.Session, appid: int):
    return session.get(
        "https://store.steampowered.com/api/appdetails", params={"appids": str(appid)}
    ).json()


def query_app_reviews(session: requests.Session, appid: int):
    cursor = "*"

    def query():
        return session.get(
            f"https://store.steampowered.com/appreviews/{appid}",
            params={
                "json": "1",
                "filter": "recent",
                "cursor": cursor,
                "num_per_page": 100,
            },
        ).json()

    data = query()
    if not data["success"]:
        return []
    result = [data]
    with tqdm(total=data["query_summary"]["total_reviews"]) as t:
        while True:
            t.update(len(data["reviews"]))
            data = query()
            if not data["success"]:
                break
            result.append(data)
            if not data["reviews"]:
                break
            cursor = data["cursor"]
            if cursor == "*":
                break
            time.sleep(2)
    return result


if __name__ == "__main__":
    import random

    with requests.session() as s:
        print(query_app_reviews(s, 1041720))
