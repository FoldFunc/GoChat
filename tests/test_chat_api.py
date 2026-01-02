import http.client
import json
import pytest
from typing import Dict, Any

HOST = "localhost"
PORT = 42069
HEADERS = {"Content-Type": "application/json"}


# -------------------------
# HTTP helper
# -------------------------

def post(path: str, payload: Dict[str, Any]) -> Dict[str, Any]:
    conn = http.client.HTTPConnection(HOST, PORT)
    conn.request("POST", path, json.dumps(payload), HEADERS)

    response = conn.getresponse()
    body = response.read().decode()
    conn.close()

    if response.status != 200:
        raise RuntimeError(f"{path} failed: {response.status} {body}")

    return json.loads(body) if body else {}


# -------------------------
# API helpers
# -------------------------

def new_user(user_name: str) -> int:
    return post("/newUser", {"user_name": user_name})["id"]


def new_room(user_id: int, room_name: str, is_public: bool) -> int:
    return post("/newRoom", {
        "room_name": room_name,
        "user_id": int(user_id),
        "room_type": is_public,
    })["id"]


def access_private_room(admin_id: int, user_id: int, room_id: int) -> None:
    post("/addToCloseRoom", {
        "room_id": int(room_id),
        "admin_id": int(admin_id),
        "user_id": int(user_id),
    })


def enter_public_room(user_id: int, room_id: int) -> None:
    post("/addToOpenRoom", {
        "room_id": int(room_id),
        "user_id": int(user_id),
    })


def send_message_private(user_id: int, room_id: int, message: str) -> int:
    return post("/sendMessageCloseRoom", {
        "room_id": int(room_id),
        "user_id": int(user_id),
        "message": message,
    })["id"]


def send_message_public(user_id: int, room_id: int, message: str) -> int:
    return post("/sendMessageOpenRoom", {
        "room_id": room_id,
        "user_id": user_id,
        "message": message,
    })["id"]


def remove_message(user_id: int, room_id: int, message_id: int) -> None:
    post("/removeMessage", {
        "room_id": int(room_id),
        "user_id": int(user_id),
        "mess_id": int(message_id),
    })


def remove_room(admin_id: int, room_id: int) -> None:
    post("/removeRoom", {
        "room_id": int(room_id),
        "user_id": int(admin_id),
    })


# -------------------------
# Tests
# -------------------------

def test_public_and_private_room_flow():
    # Create users
    admin_id = new_user("admin_user")
    test_user_id = new_user("test_user")

    assert isinstance(admin_id, int)
    assert isinstance(test_user_id, int)

    # Create rooms
    public_room_id = new_room(admin_id, "public_room", True)
    private_room_id = new_room(admin_id, "private_room", False)

    assert public_room_id > 0
    assert private_room_id > 0

    # Access rooms
    access_private_room(admin_id, test_user_id, private_room_id)
    enter_public_room(test_user_id, public_room_id)

    # Send messages
    public_msg_id = send_message_public(
        test_user_id,
        public_room_id,
        "Hello public room"
    )

    private_msg_id = send_message_private(
        test_user_id,
        private_room_id,
        "Hello private room"
    )

    assert int(public_msg_id) > 0
    assert int(private_msg_id) > 0

    # Remove messages
    remove_message(test_user_id, public_room_id, public_msg_id)
    remove_message(test_user_id, private_room_id, private_msg_id)

    # Cleanup
    remove_room(admin_id, public_room_id)
    remove_room(admin_id, private_room_id)

