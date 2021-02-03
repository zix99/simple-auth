const { assert } = require('chai');
const speakeasy = require('speakeasy');
const http = require('../../http');
const config = require('../../config');

const tfaSecret = 'ORDRZHDCYXU435ETZCIQ====';

const routes = [
  ['GET', '/api/v1/account', null],
  ['GET', '/api/v1/account/audit', null],
  ['GET', '/api/v1/local', null],
  ['POST', '/api/v1/local/password', { newpassword: 'bla' }],
  ['GET', '/api/v1/local/2fa', null],
  ['GET', '/api/v1/local/2fa/qrcode', { secret: tfaSecret }],
  ['POST', '/api/v1/local/2fa', { secret: tfaSecret, code: speakeasy.totp({ secret: tfaSecret, encoding: 'base32' }) }],
  ['DELETE', '/api/v1/local/2fa', { code: speakeasy.totp({ secret: tfaSecret, encoding: 'base32' }) }],
  ['GET', '/api/v1/auth/oauth2', null],
  ['POST', '/api/v1/auth/oauth2/grant', { client_id: 'testid', response_type: 'code', redirect_uri: 'http://example.com/redirect' }],
  ['DELETE', '/api/v1/auth/oauth2/token', { client_id: 'testid' }],
];

describe('route-surface#private', () => {
  let testUser = null;
  before(() => {
    return http.post('/api/v1/account', {
      username: 'privtest',
      password: 'test-pass',
      email: 'tps@example.com',
    }, {
      headers: {
        Authorization: `SharedKey ${config.apiSharedKey}`,
      },
    }).then((resp) => {
      testUser = resp.data;
    });
  });

  routes.forEach((route) => {
    const [method, url, payload] = route;

    it(`${method} ${url}: Should return 401 if no auth provided`, () => {
      return http({
        validateStatus: () => true,
        method,
        url,
      }).then((resp) => {
        assert.equal(resp.status, 401);
      });
    });

    it(`${method} ${url}: Should return 401 if wrong auth provided`, () => {
      return http({
        validateStatus: () => true,
        method,
        url,
        headers: {
          Authorization: 'SharedKey made-up',
        },
      }).then((resp) => {
        assert.equal(resp.status, 401);
      });
    });

    if (payload !== null) {
      it(`${method} ${url}: Should return 400 if auth was provided with no body`, () => {
        return http({
          validateStatus: () => true,
          method,
          url,
          headers: {
            Authorization: `SharedKey ${config.apiSharedKey}`,
            'X-Account-UUID': testUser.id,
          },
        }).then((resp) => {
          assert.equal(resp.status, 400);
        });
      });
    }

    it(`${method} ${url}: Should return 200 when a correct payload was provided`, () => {
      const req = {
        method,
        url,
        headers: {
          Authorization: `SharedKey ${config.apiSharedKey}`,
          'X-Account-UUID': testUser.id,
        },
      };

      if (method === 'GET' || method === 'DELETE') {
        req.params = payload;
      } else {
        req.data = payload;
      }

      return http(req).then((resp) => {
        assert.equal(resp.status, 200);
      });
    });
  });
});
