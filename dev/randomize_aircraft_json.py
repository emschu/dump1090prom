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

import json
import os
import random
import secrets
import string
import sys

"""
This script generates random values for the following aircraft fields of a dump1090 aircraft.json - if available:
- hex
- flight
- lat
- lon

The randomization works in-place.

Call it like this: python randomize_aircraft_json.py ./../example/aircraft.json

"""

def generate_secure_random_hex(length=6):
    return secrets.token_hex(length // 2).lower()

def generate_secure_random_flight(length=6):
    return ''.join(random.choice(string.ascii_letters) for _ in range(length)).upper()

if len(sys.argv) < 2:
    print(f"Usage: python {os.path.basename(__file__)} <input_file_path>")
    sys.exit(1)

input_file_path = os.path.normcase(sys.argv[1])

if not os.path.exists(input_file_path) and os.path.isfile(input_file_path):
    print(f"File not found: {input_file_path}")
    sys.exit(1)

with open(input_file_path, "r", encoding="utf-8") as file:
    print(f"Processing file: {input_file_path}")

    jsonDict = json.load(file)
    for aircraft in jsonDict["aircraft"]:
        aircraft["hex"] = generate_secure_random_hex()
        if "flight" in aircraft:
            aircraft["flight"] = generate_secure_random_flight()
        if "lat" in aircraft:
            aircraft["lat"] = round(random.uniform(-90, 90), 6)
        if "lon" in aircraft:
            aircraft["lon"] = round(random.uniform(-180, 180), 6)

    with open(input_file_path, "w", encoding="utf-8") as writeFile:
        writeFile.write(json.dumps(jsonDict, indent=4))
    print("Updated file: ", input_file_path, "\n")

print("OK.")
exit(0)
