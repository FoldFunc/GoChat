async function queryUser(agent, name) {
  const res = await agent
    .post("/getIdByName")
    .send({search_name: name});
  expect(res.statusCode).toBe(200);
  return res.body;
}
module.exports = queryUser;
