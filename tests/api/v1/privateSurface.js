const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

const routes = [
  ['GET', '/api/v1/2fa', null],
  ['GET', '/api/v1/2fa/qrcode', { secret: 'ORDRZHDCYXU435ETZCIQ====' }],
  ['GET', '/api/v1/account', null],
  ['GET', '/api/v1/account/audit', null],
  // ['POST', '/api/v1/2fa', { secret: 'ORDRZHDCYXU435ETZCIQ====', code: '123' }],
  // ['DELETE', '/api/v1/2fa', { code: '123' }],
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
            'X-Account-UUID': 'abcdefg', // FIXME, should be a valid user (but need to port create-user api first)
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
