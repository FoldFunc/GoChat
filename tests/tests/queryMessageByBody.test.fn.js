async function queryMessageByBody(agent, roomId, messageBody) {
  const res = await agent
    .post("/queryMessageFromRoom")
    .send({room_id: Number(roomId), message_body: messageBody});
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message_id");
  expect(Number.isInteger(res.body.message_id)).toBe(true);
  return res.body.message_id;
}
module.exports = queryMessageByBody;

