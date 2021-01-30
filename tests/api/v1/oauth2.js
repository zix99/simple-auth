const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

describe('oauth', () => {
  let testUser = null;
  let headers = {};
  before(() => {
    return http.post('/api/v1/account', {
      username: 'oauthtest',
      password: 'test-pass',
      email: 'oauth@example.com',
    }, {
      headers: {
        Authorization: `SharedKey ${config.apiSharedKey}`,
      },
    }).then((resp) => {
      testUser = resp.data;
      headers = {
        Authorization: `SharedKey ${config.apiSharedKey}`,
        'X-Account-UUID': testUser.id,
      };
    });
  });

  it('should describe a client', () => {
    return http.get('/api/v1/auth/oauth2/client/testid')
      .then((resp) => {
        assert.deepEqual(resp.data, {
          name: 'Test Client',
          author: 'sa',
          author_url: 'http://sa.com',
        });
      });
  });

  // The below code has to be executed in this order
  // These are state variables
  let code = null;
  let token = null;

  it('should not allow auto-granting before first grant', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/grant', {
      client_id: 'testid',
      response_type: 'code',
      scope: 'a',
      redirect_uri: 'http://example.com/redirect',
      state: 'statetoken',
      auto: true,
    }, { headers }));
  });

  it('should failed grant if invalid client', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/grant', {
      client_id: 'testid-bad',
      response_type: 'code',
      scope: 'a',
      redirect_uri: 'http://example.com/redirect',
      state: 'statetoken',
    }, { headers }));
  });

  it('should allow granting a token', () => {
    return http.post('/api/v1/auth/oauth2/grant', {
      client_id: 'testid',
      response_type: 'code',
      scope: 'a',
      redirect_uri: 'http://example.com/redirect',
      state: 'statetoken',
    }, { headers }).then((resp) => {
      assert.equal(resp.data.state, 'statetoken');
      assert.lengthOf(resp.data.code, 6);
      code = resp.data.code;
    });
  });

  it('should not trade token for bad client', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'authorization_code',
      code,
      redirect_uri: 'http://example.com/redirect',
      client_id: 'testid-bad',
      client_secret: 'client-secret',
    }));
  });

  it('should allow trading code for token', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'authorization_code',
      code,
      redirect_uri: 'http://example.com/redirect',
      client_id: 'testid',
      client_secret: 'client-secret',
    }).then((resp) => {
      token = resp.data;
      console.dir(token);
    });
  });

  it('should not allow trading the code twice', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'authorization_code',
      code,
      redirect_uri: 'http://example.com/redirect',
      client_id: 'testid',
      client_secret: 'client-secret',
    }));
  });

  it('should show the token in the client list of tokens', () => {
    return http.get('/api/v1/auth/oauth2', { headers })
      .then((resp) => {
        const clientIds = resp.data.tokens.map((x) => x.client_id);
        assert.include(clientIds, 'testid');

        const tokens = resp.data.tokens.map((x) => x.short_token);
        assert.include(tokens, token.access_token.substring(0, 5));
      });
  });

  it('should allow auto-granting when token already exists, and re-use token', () => {
    return http.post('/api/v1/auth/oauth2/grant', {
      client_id: 'testid',
      response_type: 'code',
      scope: 'a',
      redirect_uri: 'http://example.com/redirect',
      state: 'statetoken',
      auto: true,
    }, { headers }).then((resp) => {
      const autoCode = resp.data.code;
      return http.post('/api/v1/auth/oauth2/token', {
        grant_type: 'authorization_code',
        code: autoCode,
        redirect_uri: 'http://example.com/redirect',
        client_id: 'testid',
        client_secret: 'client-secret',
      });
    }).then((resp) => {
      assert.equal(resp.data.access_token, token.access_token);
    });
  });

  it('should allow inspecting token', () => {
    // TODO
  });

  it('should reject refresh token with bad secret', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'refresh_token',
      refresh_token: token.refresh_token,
      client_id: 'testid',
      client_secret: 'client-secret-bad',
    }));
  });

  it('should reject refresh token with bad clientid', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'refresh_token',
      refresh_token: token.refresh_token,
      client_id: 'testid-bad',
      client_secret: 'client-secret',
    }));
  });

  it('should allow trading refresh token for new token', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'refresh_token',
      refresh_token: token.refresh_token,
      client_id: 'testid',
      client_secret: 'client-secret',
    }).then((resp) => {
      assert.notEqual(resp.data.access_token, token.access_token);
      assert.isUndefined(resp.data.refresh_token);
    });
  });

  it('should allow granting token via credentials', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'a',
      client_id: 'testid',
      client_secret: 'client-secret',
    }).then((resp) => {
      assert.notEmpty(resp.data.access_token);
      assert.notEmpty(resp.data.refresh_token);
    });
  });

  it('should fail granting token for bad credentials', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass-bad',
      scope: 'a',
      client_id: 'testid',
      client_secret: 'client-secret',
    }));
  });

  it('should fail granting token for bad secret', () => {
    return assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'a',
      client_id: 'testid',
      client_secret: 'client-secret-bad',
    }));
  });
});
