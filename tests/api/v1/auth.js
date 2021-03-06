const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

describe('auth-endpoints', () => {
  const defaultOpts = {
    headers: {
      Authorization: `Bearer ${config.simpleAuthKey}`,
    },
    validateStatus: () => true,
  };

  let testUser = null;
  before(() => {
    return http.post('/api/v1/account', {
      username: 'authtest',
      password: 'test-pass',
      email: 'tps2@example.com',
    }, {
      headers: {
        Authorization: `SharedKey ${config.apiSharedKey}`,
      },
    }).then((resp) => {
      testUser = resp.data;
    });
  });

  describe('#simple', () => {
    it('Should allow login with username/password', () => {
      const payload = { username: 'authtest', password: 'test-pass' };
      return http.post('/api/v1/auth/simple', payload, defaultOpts)
        .then((resp) => {
          assert.equal(resp.status, 200);
          assert.equal(resp.data.id, testUser.id);
        });
    });

    it('Should allow login with email/password', () => {
      const payload = { username: 'tps2@example.com', password: 'test-pass' };
      return http.post('/api/v1/auth/simple', payload, defaultOpts)
        .then((resp) => {
          assert.equal(resp.status, 200);
          assert.equal(resp.data.id, testUser.id);
        });
    });

    it('Should return 403 if error logging in', () => {
      const payload = { username: 'authtest', password: 'bad-pass' };
      return http.post('/api/v1/auth/simple', payload, defaultOpts)
        .then((resp) => {
          assert.equal(resp.status, 403);
          assert.isUndefined(resp.data.id);
        });
    });
  });

  describe('#vouch', () => {
    it('Should return 401 if error logging in', () => {
      return http.get('/api/v1/auth/vouch', defaultOpts)
        .then((resp) => {
          assert.equal(resp.status, 401);
        });
    });
  });
});
