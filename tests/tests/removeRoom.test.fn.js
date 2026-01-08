async function removeRoom(agent, roomId){
  const res = await agent
    .post("/removeRoom")
    .send({ room_id: roomId})
    .set("Content-Type", "application/json");

  expect(res.statusCode).toBe(200);
  expect(res.body.message).toBe("room deleted");
}

module.exports = removeRoom;

