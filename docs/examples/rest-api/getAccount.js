#!/usr/bin/env node
const axios = require('axios');

const sharedSecret = 'test';
const baseUrl = 'http://localhost:9002';

const headers = {
  Authorization: `SharedKey ${sharedSecret}`,
  'X-Account-UUID': 'e4fe94ae-0cfd-44e1-878d-b93a25e38fab',
};

axios.get(`${baseUrl}/api/v1/account`, { headers })
  .then((resp) => {
    console.dir(resp.data);
  }).catch((err) => {
    console.log(err.message);
    if (err.response) {
      console.dir(err.response.data);
    }
  });
