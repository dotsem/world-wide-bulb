import json
import os
import re
import sys
import urllib.request
import urllib.error

raw_url = os.environ.get("WWB_URL", "http://localhost:8080")
base_url = re.sub(r"/api/v1(?:/events)?/?$", "", raw_url).rstrip("/")

# region: get_state
def get_state():
    req = urllib.request.Request(f"{base_url}/api/v1/state")
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode("utf-8"))
        state_str = "ON" if data.get("state") else "OFF"
        print(f"State: {state_str} ({data.get('viewers', 0)} viewers)")
        return data
# endregion: get_state

# region: get_history
def get_history(limit=10):
    req = urllib.request.Request(f"{base_url}/api/v1/history?limit={limit}")
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode("utf-8"))
        items = data.get("items") or []
        print(f"History entries: {len(items)}")
        return data
# endregion: get_history

# region: post_toggle
def toggle_bulb():
    req = urllib.request.Request(
        f"{base_url}/api/v1/toggle",
        data=b"{}",
        headers={"Content-Type": "application/json"},
        method="POST"
    )
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode("utf-8"))
        state_str = "ON" if data.get("state") else "OFF"
        print(f"Toggled bulb, ID: {data.get('id')} (State: {state_str})")
        return data
# endregion: post_toggle

# region: post_reason
def submit_reason(toggle_id, reason):
    payload = json.dumps({"id": toggle_id, "reason": reason}).encode("utf-8")
    req = urllib.request.Request(
        f"{base_url}/api/v1/reason",
        data=payload,
        headers={"Content-Type": "application/json"},
        method="POST"
    )
    with urllib.request.urlopen(req) as resp:
        data = json.loads(resp.read().decode("utf-8"))
        print(f"Reason updated for {toggle_id}")
        return data
# endregion: post_reason

def main():
    print(f"Connecting to REST API at: {base_url}")
    get_state()
    get_history()
    toggle = toggle_bulb()
    toggle_id = toggle.get("id")
    if toggle_id:
        submit_reason(toggle_id, "Turning on the bulb for client demo")

if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"Error in REST client: {exc}", file=sys.stderr)
        sys.exit(1)
