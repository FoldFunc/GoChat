// tests/loginUser.test.fn.js
async function loginUser(agent, username, password){
  const res = await agent
    .post("/login")
    .send({ user_name: username, user_password: password })
    .set("Content-Type", "application/json");

  expect(res.statusCode).toBe(200);
  expect(res.body.message).toBe("logged in");
  expect(res.headers['set-cookie']).toBeDefined();

  return res.headers['set-cookie'][0];
}

module.exports = loginUser;

