async function sendFriendRequest(agent, userId, message) {
  console.log("userId: ", userId)
  const res = await agent
    .post("/sendUserRequest")
    .send({ send_id: Number(userId.id), message: message})
    .set("Content-Type", "application/json");
  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("message");
  expect(res.body.message).toBe("User request sent");
}

module.exports = sendFriendRequest;
