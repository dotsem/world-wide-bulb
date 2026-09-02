import json
import os
import urllib.request

url = os.environ.get("WWB_URL", "https://wwb.dotsem.be/api/v1/events")
req = urllib.request.Request(url, headers={"Accept": "text/event-stream"})

print(f"Connecting to {url}...")

with urllib.request.urlopen(req) as stream:
    event_type = "message"
    for line in stream:
        decoded = line.decode("utf-8").strip()
        if decoded.startswith("event:"):
            event_type = decoded[6:].strip()
        elif decoded.startswith("data:"):
            raw_data = decoded[5:].strip()
            if not raw_data:
                continue
            data = json.loads(raw_data)
            if event_type == "state_changed":
                state = "ON" if data.get("state") else "OFF"
                print(f"Lamp state is now: {state} (ID: {data.get('id')})")
            elif event_type == "reason_updated":
                print(f"Reason: {data.get('reason')}")
