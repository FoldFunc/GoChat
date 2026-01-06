// testAgent.js
const request = require("supertest");

const API_URL = "http://localhost:42069";

module.exports = () => request.agent(API_URL);

