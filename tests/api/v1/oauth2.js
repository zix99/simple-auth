const { assert } = require('chai');
const jwt = require('jsonwebtoken');
const http = require('../../http');
const config = require('../../config');

function isRejectedWith(errCode, p) {
  return p.then(() => {
    assert.fail('Promise must be rejected');
  }).catch((err) => {
    if (!err.response) {
      assert.fail(err);
    } else {
      assert.equal(err.response.status, 400);
      assert.equal(err.response.data.error, errCode);
    }
  });
}

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

  it('should not allow granting a token with bad scopes', () => {
    return isRejectedWith('invalid_scope', http.post('/api/v1/auth/oauth2/grant', {
      client_id: 'testid',
      response_type: 'code',
      scope: 'a c',
      redirect_uri: 'http://example.com/redirect',
      state: 'statetoken',
    }, { headers }));
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

      assert.notEmpty(token.access_token);
      assert.notEmpty(token.refresh_token);
      assert.notEmpty(token.id_token);
    });
  });

  it('should now have a valid id token', () => {
    const decoded = jwt.verify(token.id_token, 'this-is-a-test-key');
    assert.notEmpty(decoded.sub);
    assert.notEmpty(decoded.aud);
    assert.equal(decoded.iss, 'simple-auth');
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

  it('should allow inspecting token', () => {
    return http.post('/api/v1/auth/oauth2/token_info', { token: token.access_token })
      .then((resp) => {
        assert.isTrue(resp.data.active);
        assert.equal(resp.data.token_type, 'access_token');
        assert.equal(resp.data.scope, 'a');
        assert.notEmpty(resp.data.sub);
        assert.isNumber(resp.data.exp);
        assert.isNumber(resp.data.iat);
        assert.equal(resp.data.client_id, 'testid');
        assert.equal(resp.data.aud, 'testid');
        assert.equal(resp.data.iss, 'simple-auth');
      });
  });

  it('Should return non-active when bad token', () => {
    return http.post('/api/v1/auth/oauth2/token_info', { token: 'fake' })
      .then((resp) => {
        assert.equal(200, resp.status);
        assert.isFalse(resp.data.active);
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

  it('Should successfully revoke all tokens', () => {
    return http.delete('/api/v1/auth/oauth2/token', { params: { client_id: 'testid' }, headers });
  });

  it('Should have no visible tokens after revoking', () => {
    return http.get('/api/v1/auth/oauth2', { headers })
      .then((resp) => {
        assert.equal(resp.data.tokens.length, 0);
      });
  });

  it('Should allow revoking single token', () => {
    let tk = null;
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'a',
      client_id: 'testid',
      client_secret: 'client-secret',
    }).then((resp) => {
      tk = resp.data;
      const params = {
        client_id: 'testid',
        token: tk.refresh_token,
      };
      return http.delete('/api/v1/auth/oauth2/token', { params, headers });
    }).then(() => http.get('/api/v1/auth/oauth2', { headers }))
      .then((resp) => {
        const tokens = resp.data.tokens.map((x) => x.short_token);
        assert.include(tokens, tk.access_token.substring(0, 5));
        assert.notInclude(tokens, tk.refresh_token.substring(0, 5));
      });
  });
});

describe('oauth2#credentials', () => {
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
      assert.notEmpty(resp.data.id_token);
    });
  });

  it('should fail granting via credentials when bad scope', () => {
    return isRejectedWith('invalid_scope', http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'a c',
      client_id: 'testid',
      client_secret: 'client-secret',
    }));
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

  it('should fail granting token for credentials with bad secret', () => {
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

describe('OAuth2#Single use token', () => {
  let token = null;
  let access = null;

  it('should allow granting token via credentials', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'name',
      client_id: 'singleissue',
      client_secret: 'si-secret',
    }).then((resp) => {
      assert.notEmpty(resp.data.access_token);
      assert.notEmpty(resp.data.refresh_token);
      assert.isUndefined(resp.data.id_token);
      token = resp.data;
    });
  });

  it('Should introspect refresh token', () => {
    return http.post('/api/v1/auth/oauth2/token_info', {
      token: token.refresh_token,
    }).then((resp) => {
      assert.isTrue(resp.data.active);
    });
  });

  it('should allow trading refresh token for new token', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'refresh_token',
      refresh_token: token.refresh_token,
      client_id: 'singleissue',
      client_secret: 'si-secret',
    }).then((resp) => {
      assert.notEqual(resp.data.access_token, token.access_token);
      assert.isUndefined(resp.data.refresh_token);
      access = resp.data;
      console.dir(access);
    });
  });

  it('Should introspect access token', () => {
    return http.post('/api/v1/auth/oauth2/token_info', {
      token: access.access_token,
    }).then((resp) => {
      assert.isTrue(resp.data.active);
    });
  });

  it('Should not accept access token after 1 second', (done) => {
    setTimeout(() => {
      http.post('/api/v1/auth/oauth2/token_info', {
        token: access.access_token,
      }).then((resp) => {
        if (resp.data.active) {
          done(new Error('Failed'));
        } else {
          done();
        }
      }).catch((err) => {
        done(err);
      });
    }, 1200);
  });

  it('Issuing 2 tokens will make the first invalid', () => {
    const tokenReq = {
      grant_type: 'refresh_token',
      refresh_token: token.refresh_token,
      client_id: 'singleissue',
      client_secret: 'si-secret',
    };
    let token1;
    let token2;

    return http.post('/api/v1/auth/oauth2/token', tokenReq)
      .then((resp) => {
        token1 = resp.data;
        return http.post('/api/v1/auth/oauth2/token', tokenReq);
      }).then((resp) => {
        token2 = resp.data;
      }).then(() => {
        assert.notEqual(token1.access_token, token2.access_token);
        return Promise.all([
          http.post('/api/v1/auth/oauth2/token_info', { token: token1.access_token })
            .then((ti) => assert.isFalse(ti.data.active)),
          http.post('/api/v1/auth/oauth2/token_info', { token: token2.access_token })
            .then((ti) => assert.isTrue(ti.data.active)),
        ]);
      });
  });

  it('Issuing a new refresh token will make the first invalid, along with last token', () => {
    return http.post('/api/v1/auth/oauth2/token', {
      grant_type: 'password',
      username: 'oauthtest',
      password: 'test-pass',
      scope: 'name',
      client_id: 'singleissue',
      client_secret: 'si-secret',
    }).then((resp) => {
      assert.notEmpty(resp.data.access_token);
      assert.notEmpty(resp.data.refresh_token);
      assert.isUndefined(resp.data.id_token);

      return Promise.all([
        assert.isRejected(http.post('/api/v1/auth/oauth2/token', {
          grant_type: 'refresh_token',
          refresh_token: token.refresh_token,
          client_id: 'singleissue',
          client_secret: 'si-secret',
        })),
        http.post('/api/v1/auth/oauth2/token_info', { token: token.refresh_token }).then((ti) => assert.isFalse(ti.data.active)),
        http.post('/api/v1/auth/oauth2/token_info', { token: token.access_token }).then((ti) => assert.isFalse(ti.data.active)),
        http.post('/api/v1/auth/oauth2/token_info', { token: access.access_token }).then((ti) => assert.isFalse(ti.data.active)),
      ]);
    });
  });
});
