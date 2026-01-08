async function addUserOpenRoom(agent, roomId) {
  const res = await agent
    .post("/addToOpenRoom")
    .send({ room_id: roomId})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message");
}

module.exports = addUserOpenRoom;
