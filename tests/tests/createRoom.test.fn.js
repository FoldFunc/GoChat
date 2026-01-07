async function createRoom(agent, roomName, roomType) {
  const res = await agent
    .post("/newRoom")
    .send({ room_name: roomName, room_type: roomType})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("id");
}

module.exports = createRoom;

