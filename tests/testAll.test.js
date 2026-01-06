// testAll.test.js
const createAgent = require("./testAgent");
const createUser = require("./tests/newUser.test.fn");
const loginUser = require("./tests/loginUser.test.fn");
const logoutUser = require("./tests/logoutUser.test.fn");

describe("API Integration Tests", () => {
  const agent = createAgent();
  const username = "foldfunc";
  const password = "pass1";
  const connType = true;
  it("Create user", async () => {
    await createUser(agent, username, password, connType);
  });

  it("Login", async () => {
    await loginUser(agent, username, password);
  });

  it("Logout", async () => {
    await logoutUser(agent);
  });
});

