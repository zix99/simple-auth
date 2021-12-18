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
  const sharedKeyOpts = {
    headers: {
      Authorization: `SharedKey ${config.apiSharedKey}`,
    },
  };

  let testUser = null;
  before(() => {
    return http.post('/api/v1/account', {
      username: 'authtest',
      password: 'test-pass',
      email: 'tps2@example.com',
    }, sharedKeyOpts).then((resp) => {
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
    let sessionCookie = null;
    it('Should allow creation of a session', () => {
      const payload = { username: 'tps2@example.com', password: 'test-pass' };
      return http.post('/api/v1/auth/session', payload, sharedKeyOpts) // shared key gets around CSRF
        .then((resp) => {
          assert.equal(resp.status, 200);
          const cookies = resp.headers['set-cookie'][0];
          [sessionCookie] = cookies.split(';');
          assert.isNotEmpty(sessionCookie);
        });
    });

    it('Should return 401 if error logging in', () => {
      return http.get('/api/v1/auth/vouch', defaultOpts)
        .then((resp) => {
          assert.equal(resp.status, 401);
        });
    });

    it('Should allow forwarding vouch', () => {
      const opts = {
        maxRedirects: 0,
        validateStatus: () => true,
      };
      return http.get('/api/v1/auth/vouch?forward=1', opts)
        .then((resp) => {
          assert.equal(resp.status, 307);
        });
    });

    it('Should allow forwarding vouch with continue', () => {
      const opts = {
        maxRedirects: 0,
        validateStatus: () => true,
        params: {
          forward: '1',
          continue: 'http://asdf.com',
        },
      };
      return http.get('/api/v1/auth/vouch?forward=1', opts)
        .then((resp) => {
          assert.equal(resp.status, 307);
          assert.include(resp.headers.location, '?continue=http%3A%2F%2Fasdf.com');
        });
    });

    it('Should allow forwarding vouch with continue', () => {
      const opts = {
        maxRedirects: 0,
        validateStatus: () => true,
        params: {
          forward: '1',
        },
        headers: {
          'X-Forwarded-Host': 'asdf.com',
          'X-Forwarded-Proto': 'http',
          'X-Forwarded-Uri': '/abc',
        },
      };
      return http.get('/api/v1/auth/vouch?forward=1', opts)
        .then((resp) => {
          assert.equal(resp.status, 307);
          assert.include(resp.headers.location, '?continue=http%3A%2F%2Fasdf.com%2Fabc');
        });
    });

    it('Should allow vouching with a session cookie', () => {
      const headers = {
        cookie: sessionCookie,
      };
      return http.get('/api/v1/auth/vouch', { headers })
        .then((resp) => {
          assert.equal(resp.status, 200);
          assert.isNotEmpty(resp.headers['x-user-id']);
        });
    });

    it('Should allow forwarding vouch with a session cookie', () => {
      const headers = {
        cookie: sessionCookie,
      };
      return http.get('/api/v1/auth/vouch?forward=1', { headers })
        .then((resp) => {
          assert.equal(resp.status, 200);
        });
    });
  });
});
