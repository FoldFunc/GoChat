async function queryRooms(agent) {
  const res = await agent.get("/queryUserRooms");
  expect(res.statusCode).toBe(200);
  expect(Array.isArray(res.body)).toBe(true);
    if (res.body.length > 0) {
    expect(res.body[0]).toHaveProperty("Id");
    expect(res.body[0]).toHaveProperty("UserId");
    expect(res.body[0]).toHaveProperty("Name");
    expect(res.body[0]).toHaveProperty("Type");
  }
  return res.body
}
module.exports = queryRooms;
