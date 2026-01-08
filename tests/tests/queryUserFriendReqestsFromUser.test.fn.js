async function queryUserFriendReqestsFromUser(agent, user) {
  const res = await agent
    .post("/viewUserRequestsFromUser")
    .send({user_id: Number(user.Id)});
  expect(res.statusCode).toBe(200);
  return res.body;
}
module.exports = queryUserFriendReqestsFromUser;
