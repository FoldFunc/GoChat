async function queryRoomByName(agent, roomName) {
  const res = await agent 
    .post("/queryUserRoom")
    .send({room_name: roomName});
  expect(res.statusCode).toBe(200);
  return res.body
}
module.exports = queryRoomByName;
