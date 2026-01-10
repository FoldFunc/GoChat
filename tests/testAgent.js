// testAgent.js
const request = require("supertest");

const API_URL = "http://localhost:8080";

module.exports = () => request.agent(API_URL);

