import steam_api
import os
import json
from tqdm import tqdm
import time
import requests


def main():
    os.chdir("D:/ML/steam/raw")
    with requests.session() as s:
        all_apps = steam_api.query_all_apps(s)
        with open("apps.json", "x") as f:
            json.dump(all_apps, f, ensure_ascii=False)
        print("Fetching all app details...")
        for app in tqdm(all_apps):
            if app.get("appid") and app.get("name"):
                details = steam_api.query_app_details(s, app["appid"])
                with open(f'app_details/{app["appid"]}.json', "x") as f:
                    json.dump(details, f, ensure_ascii=False)
                time.sleep(3.47)


if __name__ == "__main__":
    main()
