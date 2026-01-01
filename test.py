import http.client
import json


def getRoom(name_value):
    conn = http.client.HTTPConnection("localhost", 42069)
    path = f"/getRoom?name={name_value}"
    conn.request("GET", path)

    response = conn.getresponse()
    data = response.read().decode()
    print("Raw response: ", data)
    try:
        message = json.loads(data)
        print("Parsed JSON: ", message)
    except:
        print("Response is not valid JSON")
    conn.close()
def newRoom(idd, name):
    conn = http.client.HTTPConnection("localhost", 42069)
    payload = {"room_id": idd, "room_name": name}
    json_data = json.dumps(payload)
    headers = {"Content-Type": "application/json"}
    conn.request("POST", "/newRoom", body=json_data, headers=headers)
    response = conn.getresponse()
    resp_data = response.read().decode()
    print("Status: ", response.status)
    print("Response: ", resp_data)
    conn.close()
print("Creating new room: ")
newRoom(0, "my_room")
print("Quering for the \"new_room\"")
getRoom("my_room")
