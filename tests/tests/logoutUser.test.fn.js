// tests/logoutUser.test.fn.js
async function logoutUser(agent) {
  const res = await agent.get("/logout");
  expect(res.statusCode).toBe(200);
  expect(res.body.message).toBe("logged out");
}

module.exports = logoutUser;

