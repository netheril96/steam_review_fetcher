import steam_api
import os
import json
from tqdm import tqdm
import time
import requests
import traceback
import zstandard
import sqlite3


def main():
    with requests.session() as s:
        all_apps = steam_api.query_all_apps(s)
        with open("apps.json", "w") as f:
            json.dump(all_apps, f, ensure_ascii=False)
        print("Fetching all app details...")
        compressor = zstandard.ZstdCompressor(level=9)
        with sqlite3.connect("raw.db") as conn:
            cursor = conn.cursor()
            cursor.executescript(
                """
                    create table if not exists
                        app_details
                        (appid integer primary key, compressed_details blob);
                    pragma journal_mode = wal;
                """
            )
            conn.commit()
            insert_count = 0
            for app in tqdm(all_apps):
                if not app.get("appid") or not app.get("name"):
                    continue
                cursor.execute(
                    "select count(1) from app_details where appid = ?",
                    (app["appid"],),
                )
                (count,) = cursor.fetchone()
                if count:
                    continue  # Already saved
                try:
                    details = steam_api.query_app_details(s, app["appid"])
                    data = compressor.compress(
                        json.dumps(details, ensure_ascii=False).encode("utf-8")
                    )
                except KeyboardInterrupt:
                    break
                except:
                    traceback.print_exc()
                    data = None
                cursor.execute(
                    "insert into app_details (appid, compressed_details) values (?, ?)",
                    (app["appid"], data),
                )
                conn.commit()
                insert_count += 1
                if insert_count % 200 == 0:
                    cursor.executescript(
                        """
                            pragma optimize;
                        """
                    )
                time.sleep(0.93)


if __name__ == "__main__":
    main()
