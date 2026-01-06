// tests/newUser.test.fn.js
const request = require("supertest");
const API_URL = "http://localhost:42069";

async function createUser(agent, username, password, connType) {
  const res = await agent
    .post("/newUser")
    .send({ user_name: username, conn_type: connType, password: password })
    .set("Content-Type", "application/json");

  expect(res.statusCode).toBe(200);
  expect(res.body).toHaveProperty("id");
}

module.exports = createUser;

