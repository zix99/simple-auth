const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

const routes = [
  ['POST', '/api/v1/account', { username: 'skeediddle', password: 'test-pass', email: 'test-create@example.com' }],
  ['POST', '/api/v1/account/check', { username: 'skeedoo' }],
  ['POST', '/api/v1/auth/session', { username: 'skeediddle', password: 'test-pass' }],
  ['DELETE', '/api/v1/auth/session'],
  ['POST', '/api/v1/auth/onetime', { email: 'test-create@example.com' }],
];

describe('route-surface#public', () => {
  routes.forEach((route) => {
    const [method, url, payload] = route;

    it(`${method} ${url}: Should return 403 (CSRF fail) if no auth provided`, () => {
      return http({
        validateStatus: () => true,
        method,
        url,
      }).then((resp) => {
        assert.equal(resp.status, 403);
      });
    });

    if (payload) {
      it(`${method} ${url}: Should return 400 if auth was provided with no body`, () => {
        return http({
          validateStatus: () => true,
          method,
          url,
          headers: {
            Authorization: `SharedKey ${config.apiSharedKey}`,
          },
        }).then((resp) => {
          assert.equal(resp.status, 400);
        });
      });
    }

    it(`${method} ${url}: Should return 200 when a correct payload was provided`, () => {
      return http({
        method,
        url,
        headers: {
          Authorization: `SharedKey ${config.apiSharedKey}`,
        },
        data: payload,
      }).then((resp) => {
        assert.equal(Math.trunc(resp.status / 100) * 100, 200);
      });
    });
  });
});
