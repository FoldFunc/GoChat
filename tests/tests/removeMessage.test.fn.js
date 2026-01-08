async function removeMessage(agent, messageId, roomId){
  const res = await agent
    .post("/removeMessage")
    .send({ room_id: roomId, mess_id: messageId})
    .set("Content-Type", "application/json");

  expect(res.statusCode).toBe(200);
  expect(res.body.message).toBe("message deleted");
}

module.exports = removeMessage;

