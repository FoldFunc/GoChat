import http.client
import json
from urllib.parse import urlencode

HOST = "localhost"
PORT = 42069
HEADERS = {"Content-Type": "application/json"}


def new_user(user_name: str) -> int:
    conn = http.client.HTTPConnection(HOST, PORT)

    payload = {
        "user_name": user_name
    }

    conn.request(
        "POST",
        "/newUser",
        body=json.dumps(payload),
        headers=HEADERS
    )

    response = conn.getresponse()
    data = response.read().decode()
    conn.close()

    if response.status != 200:
        raise RuntimeError(f"newUser failed: {response.status} {data}")

    # {"id": 123}
    return json.loads(data)["id"]


def new_room(user_id: int, room_name: str) -> int:
    conn = http.client.HTTPConnection(HOST, PORT)

    payload = {
        "room_name": room_name,
        "user_id": user_id
    }

    conn.request(
        "POST",
        "/newRoom",
        body=json.dumps(payload),
        headers=HEADERS
    )

    response = conn.getresponse()
    data = response.read().decode()
    conn.close()

    if response.status != 200:
        raise RuntimeError(f"newRoom failed: {response.status} {data}")

    return json.loads(data)["id"]


def get_room(user_id: int, room_name: str) -> int:
    conn = http.client.HTTPConnection(HOST, PORT)

    params = urlencode({
        "id": user_id,
        "name": room_name
    })

    path = f"/getRoom?{params}"

    conn.request("GET", path)

    response = conn.getresponse()
    data = response.read().decode()
    conn.close()

    if response.status != 200:
        raise RuntimeError(f"getRoom failed: {response.status} {data}")

    return json.loads(data)["id"]


# -------------------------
# Test flow
# -------------------------

print("Creating user...")
user_id = new_user("user1")
print("User ID:", user_id)

print("\nCreating room...")
room_id = new_room(user_id, "my_room")
print("Room ID:", room_id)

print("\nFetching room...")
fetched_room_id = get_room(user_id, "my_room")
print("Fetched Room ID:", fetched_room_id)

assert room_id == fetched_room_id
print("\n✅ All tests passed")

