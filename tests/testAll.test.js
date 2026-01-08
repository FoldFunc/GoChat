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
const removeMessage = require("./tests/removeMessage.test.fn")
const queryMessageByBody = require("./tests/queryMessageByBody.test.fn")
const removeRoom = require("./tests/removeRoom.test.fn");
const sendFriendRequest = require("./tests/sendFriendRequest.test.fn")
const queryUserFriendReqestsFromUser = require("./tests/queryUserFriendReqestsFromUser.test.fn")
const acceptFriendRequest = require("./tests/acceptFriendRequest.test.fn")
describe("API Integration Tests", () => {
  const agent = createAgent();
  const username = "foldfunc";
  const password = "pass1";
  const usernamePriv = "foldfuncpriv";
  const passwordPriv = "pass1priv";
  const connType = true;
  const roomNamePublic = "foldfuncRoom";
  const roomNamePublicTest = "foldfuncRoomTest";
  const roomTypePublic = true;
  const roomNamePrivate = "foldfuncRoomPriv";
  const roomTypePrivate = false;
  let rooms = [];
  let user = {};
  let openRoom = {};
  let messageId = {};
  let requestId = {};
  const messageOpen = "Hello to an open room";
  const messageClose = "Hello to an close room";
  const messageOpenTest = "Hello to an open room test";

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
  });
  it("Adding a message to an open room as a test user", async () => {
    await sendMessageOpenRoom(agent, rooms[0], messageOpenTest); 
  });
  it("Query a message from an open room", async () => {
    messageId = await queryMessageByBody(agent, openRoom, messageOpenTest);
  })
  it("Removing a message from a room", async () => {
    await removeMessage(agent, messageId, openRoom);
  });
  it("Create test room for removal", async () => {
    await createRoom(agent, roomNamePublicTest, roomTypePublic);
  });
  it("Get test room id", async () => {
    testRoomId = await queryRoomByName(agent, roomNamePublicTest); 
  });
  it("Remove the test room", async () => {
    await removeRoom(agent, testRoomId);
  });
  it("Create a private user", async () => {
    await createUser(agent, usernamePriv, passwordPriv);
  });
  it("Send friend request from test user to foldfunc", async () => {
    user = await queryUser(agent, usernamePriv);
    await sendFriendRequest(agent, user)
  });
  it("logout test user so you can accept as priv user", async () => {
    await logoutUser(agent);
  });
  it("login priv user to accept the reqest", async () => {
    await loginUser(agent, usernamePriv, passwordPriv);
  });
  it("Find all reqeusts from the test user", async () => {
    user = await queryUser(agent, "test")
    // This is an []Object
    requestId = await queryUserFriendReqestsFromUser(agent, user)
  });
  it("Accept the chat friend request", async () => {
    await acceptFriendRequest(agent, requestId[0]);
  });
});

