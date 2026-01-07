// testAll.test.js
const createAgent = require("./testAgent");
const createUser = require("./tests/newUser.test.fn");
const loginUser = require("./tests/loginUser.test.fn");
const logoutUser = require("./tests/logoutUser.test.fn");
const createRoom = require("./tests/createRoom.test.fn")
const sendMessageOpenRoom = require("./tests/sendMessageOpenRoom.test.fn")
const queryRooms = require("./tests/queryRooms.test.fn")

describe("API Integration Tests", () => {
  const agent = createAgent();
  const username = "foldfunc";
  const password = "pass1";
  const connType = true;
  const roomName = "foldfuncRoom";
  const roomType = true;
  let rooms = [];
  const messageOpen = "Hello to an open room";
  it("Create user", async () => {
    await createUser(agent, username, password, connType);
  });

  it("Login", async () => {
    await loginUser(agent, username, password);
  });
  it("Create a room", async () => {
    await createRoom(agent, roomName, roomType);
  })
  it("Query all of the user rooms", async () => {
    rooms = await queryRooms(agent);
  })
  it("Send a message to an open room", async () => {
    await sendMessageOpenRoom(agent, rooms[0], messageOpen); 
  })
  it("Logout", async () => {
    await logoutUser(agent);
  });
});

