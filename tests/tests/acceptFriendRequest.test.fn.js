async function acceptFriendRequest(agent, requestId) {
  const res = await agent
    .post("/acceptChatRequest")
    .send({request_id: requestId});
  expect(res.statusCode).toBe(200);
}
module.exports = acceptFriendRequest;
