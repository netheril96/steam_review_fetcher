import dateutil.parser
import steam_api
import os
import json
from tqdm import tqdm
import time
import requests
import traceback
import zstandard
import sqlite3
import enum
import tap
import dataclasses
import datetime
import dateutil


@dataclasses.dataclass
class AppInfoSummary:
    app_id: int
    release_date: datetime.date


class Mode(enum.IntFlag):
    D = 1
    R = 2


def main(mode: str):
    with requests.session() as session:
        all_apps = steam_api.query_all_apps(session)
        with open("apps.json", "w") as f:
            json.dump(all_apps, f, ensure_ascii=False)
        compressor = zstandard.ZstdCompressor(level=9)
        decompressor = zstandard.ZstdDecompressor()
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
            if Mode[mode] & Mode.D:
                print("Fetching all app details...")
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
                        details = steam_api.query_app_details(session, app["appid"])
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
            if Mode[mode] & Mode.R:
                summaries: list[AppInfoSummary] = []
                for app in tqdm(all_apps):
                    cursor.execute(
                        "select compressed_details from app_details where appid = ?",
                        (app["appid"],),
                    )
                    current_result = cursor.fetchone()
                    if not current_result:
                        continue
                    (data,) = current_result
                    details = json.loads(decompressor.decompress(data))
                    try:
                        data = details[str(app["appid"])]["data"]
                        if (
                            data["release_date"]["coming_soon"] is False
                            and data["release_date"]["date"].strip()
                        ):
                            summaries.append(
                                AppInfoSummary(
                                    app_id=app["appid"],
                                    release_date=dateutil.parser.parse(
                                        data["release_date"]["date"]
                                    ).date(),
                                )
                            )
                    except KeyError:
                        continue
                    except dateutil.parser.ParserError:
                        continue
                summaries.sort(key=lambda s: s.release_date, reverse=True)
                print("Fetching app reviews", summaries)
                for s in tqdm(summaries):
                    filename = os.path.join("app_reviews", str(s.app_id))
                    if os.path.exists(filename):
                        continue
                    reviews = steam_api.query_app_reviews(session, s.app_id)
                    with open(filename + ".tmp", "wb") as f:
                        f.write(
                            compressor.compress(
                                json.dumps(reviews, ensure_ascii=False).encode("utf-8")
                            )
                        )
                    os.replace(filename + ".tmp", filename)
                    time.sleep(0.88)


if __name__ == "__main__":
    tap.tapify(main)
