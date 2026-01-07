async function sendMessageCloseRoom(agent, roomId, message) {
  const res = await agent
    .post("/sendMessageCloseRoom")
    .send({ room_id: Number(roomId.Id), body: message})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message");
  expect(res.body).toHaveProperty("id");
}

module.exports = sendMessageCloseRoom;


