process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";

const request = require("supertest");
const https = require("https");

const agent = request.agent("https://foldfunc-server-production.up.railway.app");
agent.agent = new https.Agent({ rejectUnauthorized: false });

module.exports = () => agent;

