# dump1090prom
# Copyright (C) 2025 emschu[aet]mailbox.org
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public
# License along with this program.
# If not, see <https://www.gnu.org/licenses/>.

import pandas as pd
import requests
from io import StringIO
from markdown_it.rules_inline.backticks import regex

# script to fetch airline list from wikipedia and safe as CSV

links = "https://en.wikipedia.org/wiki/List_of_airline_codes"
headers = {"User-Agent": "Mozilla/5.0", "Connection": "close"}
response = requests.get(links, headers=headers)
if response.status_code != 200:
    raise Exception("Could not fetch data from URL: {links}")

tables = pd.read_html(StringIO(response.text), header=0)
df = pd.DataFrame(
    columns=["IATA", "ICAO", "Airline", "Call sign", "Country/Region", "Comments"]
)
df = pd.concat([df, tables[0]], ignore_index=True)

# clean up
df = df.loc[~df["Comments"].str.contains(r"^[dD]{1}efunct.*", regex=True, na=False)]
df = df.loc[df["Comments"] != "Became Lufthansa"]

with open("../data/wikipedia-airlines.csv", "w+") as f:
    f.write(df.to_csv(sep=";", escapechar="\\", index=False))

print("OK.")
exit(0)
