const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

describe('/api/v1/account/check', () => {
  it('Should allow checking of username', () => {
    return http.post('/api/v1/account/check', { username: 'sloth' }, { headers: { Authorization: `SharedKey ${config.apiSharedKey}` } })
      .then((resp) => {
        assert.equal(resp.data.exists, false);
      });
  });
});
