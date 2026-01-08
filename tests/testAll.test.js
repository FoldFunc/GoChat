// testAll.test.js
const createAgent = require("./testAgent");
const createUser = require("./tests/newUser.test.fn");
const loginUser = require("./tests/loginUser.test.fn");
const logoutUser = require("./tests/logoutUser.test.fn");
const createRoom = require("./tests/createRoom.test.fn")
const sendMessageOpenRoom = require("./tests/sendMessageOpenRoom.test.fn")
const sendMessageCloseRoom = require("./tests/sendMessageCloseRoom.test.fn")
const queryRooms = require("./tests/queryRooms.test.fn")
const queryUser = require("./tests/queryUser.test.fn")
const addUserCloseRoom = require("./tests/addUserCloseRoom.test.fn")
const queryRoomByName = require("./tests/queryRoomByName.test.fn")
const addUserOpenRoom = require("./tests/addUserOpenRoom.test.fn")

describe("API Integration Tests", () => {
  const agent = createAgent();
  const username = "foldfunc";
  const password = "pass1";
  const connType = true;
  const roomNamePublic = "foldfuncRoom";
  const roomTypePublic = true;
  const roomNamePrivate = "foldfuncRoomPriv";
  const roomTypePrivate = false;
  let rooms = [];
  let user = {};
  let openRoom = {};
  const messageOpen = "Hello to an open room";
  const messageClose = "Hello to an close room";
  it("Create user", async () => {
    await createUser(agent, username, password, connType);
  });

  it("Login", async () => {
    await loginUser(agent, username, password);
  });
  it("Create a room public room", async () => {
    await createRoom(agent, roomNamePublic, roomTypePublic);
  })
  it("Create a room private room", async () => {
    await createRoom(agent, roomNamePrivate, roomTypePrivate);
  })
  it("Query all of the user rooms", async () => {
    rooms = await queryRooms(agent);
  })
  it("Send a message to an open room", async () => {
    await sendMessageOpenRoom(agent, rooms[0], messageOpen); 
  })
  it("Send a message to an close room", async () => {
    await sendMessageCloseRoom(agent, rooms[1], messageClose);
  })
  it("Create test user", async () => {
    await createUser(agent, "test", "test", connType);
  });
  it("Query a test user by name", async () => {
    user = await queryUser(agent, "test");
  });
  it("Add user to a close room", async () => {
    await addUserCloseRoom(agent, rooms[1], user)
  });
  it("Logout", async () => {
    await logoutUser(agent);
  });
  it("Login test user", async () => {
    await loginUser(agent, "test", "test");
  });
  it("Query open room by name", async () => {
    openRoom = await queryRoomByName(agent, roomNamePublic);
  })
  it("Add user to a open room", async () => {
    await addUserOpenRoom(agent, openRoom);
  })
  // From here you are a test user beacouse the stupid agent can't have more than
  // one cookie at a time.
});

