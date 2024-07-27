import steam_api
import os
import json
from tqdm import tqdm
import time
import requests
import traceback


def main():
    with requests.session() as s:
        all_apps = steam_api.query_all_apps(s)
        with open("apps.json", "w") as f:
            json.dump(all_apps, f, ensure_ascii=False)
        print("Fetching all app details...")
        for app in tqdm(all_apps):
            if app.get("appid") and app.get("name"):
                filename = f'app_details/{app["appid"]}.json'
                try:
                    with open(filename, "x") as f:
                        details = steam_api.query_app_details(s, app["appid"])
                        json.dump(details, f, ensure_ascii=False)
                    time.sleep(0.93)
                except FileExistsError:
                    continue
                except KeyboardInterrupt:
                    os.remove(filename)
                    raise
                except:
                    traceback.print_exc()
                    os.remove(filename)
                    continue


if __name__ == "__main__":
    main()
