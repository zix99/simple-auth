#!/usr/bin/env node
const express = require('express');
const cookieParser = require('cookie-parser');
const jwt = require('jsonwebtoken');

const PORT = process.env.PORT || 8080;
const AUTHURL = process.env.AUTHURL;
const JWTKEY = process.env.JWTKEY;

const app = express();

app.use(cookieParser());

// Simplistic auth middleware
app.use((req, res, next) => {
  const authCookie = req.cookies.auth;
  if (!authCookie) {
    // You could redirect here..
    return res.redirect(AUTHURL);
  }

  return jwt.verify(authCookie, JWTKEY, (err, decoded) => {
    if (err) {
      return res.status(401).send('Invalid token');
    }
    req.auth = decoded;
    return next();
  });
});

// Only can get if passes auth middleware
app.get('/', (req, res) => {
  res.send(`Hello!<br>
  Your auth cookie is: ${req.cookies.auth}<br>
  Your token decodes to: ${JSON.stringify(req.auth)}<br>
  <br>
  <a href="${AUTHURL}/#/manage">Click here to manage your account</a>`);
});

app.listen(PORT, () => {
  console.log(`Listening on http://0.0.0.0:${PORT}`);
});
