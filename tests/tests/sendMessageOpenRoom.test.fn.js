async function sendMessageOpenRoom(agent, roomId, message) {
  console.log("I AM HERE: ", roomId)
  const res = await agent
    .post("/sendMessageOpenRoom")
    .send({ room_id: Number(roomId.Id), body: message})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message");
  expect(res.body).toHaveProperty("id");
}

module.exports = sendMessageOpenRoom;

