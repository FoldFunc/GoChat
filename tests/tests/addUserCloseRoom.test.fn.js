async function addUserCloseRoom(agent, roomId, user) {
  const res = await agent
    .post("/addToCloseRoom")
    .send({ room_id: Number(roomId.Id), user_id: Number(user.id)})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message");
}

module.exports = addUserCloseRoom;
