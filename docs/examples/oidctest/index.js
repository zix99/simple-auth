#!/usr/bin/env node
const express = require('express');
const axios = require('axios');
const config = require('./config');

const app = express();

// Simple page. You can either put the values in a form, or redirect back to your site, and then to simple-auth
const pageHome = `
<html>
  <head><title>OIDC Test</title></head>
  <body>
    <h1>OIDC Login</h1>
    <form action="${config.oauthGrantEndpoint}" method="GET">
      <input type="hidden" name="client_id" value="${config.clientId}">
      <input type="hidden" name="response_type" value="code">
      <input type="text" name="scope" value="email name">
      <input type="hidden" name="redirect_uri" value="${config.myUrl}/auth-callback">
      <button type="submit">Login</button>
    </form>
  </body>
</html>
`;
app.get('/', (req, res) => {
  res.send(pageHome);
});

function tradeCodeForAccessToken(code) {
  return axios.post(config.tokenEndpoint, {
    client_id: config.clientId,
    client_secret: config.clientSecret,
    grant_type: 'authorization_code',
    redirect_uri: `${config.myUrl}/auth-callback`,
    code,
  }).then((resp) => resp.data);
}

app.get('/auth-callback', (req, res) => {
  const { code } = req.query;

  tradeCodeForAccessToken(code)
    .then((token) => {
      res.send(token);
    })
    .catch((err) => {
      console.log(err);
      res.redirect('/');
    });
});

app.listen(config.port, () => {
  console.log(`Listening on http://0.0.0.0:${config.port}`);
});
