const { assert } = require('chai');
const http = require('../../http');
const config = require('../../config');

const routes = [
  ['POST', '/api/v1/account/check', { username: 'skeedoo' }],
];

describe('route-surface#public', () => {
  routes.forEach((route) => {
    it('Should return 403 (CSRF fail) if no auth provided', () => {
      return http({
        validateStatus: () => true,
        method: route[0],
        url: route[1],
      }).then((resp) => {
        assert.equal(resp.status, 403);
      });
    });

    it('Should return 400 if auth was provided', () => {
      return http({
        validateStatus: () => true,
        method: route[0],
        url: route[1],
        headers: {
          Authorization: `SharedKey ${config.apiSharedKey}`,
        },
      }).then((resp) => {
        assert.equal(resp.status, 400);
      });
    });

    it('Should return 200 when a correct payload was provided', () => {
      return http({
        method: route[0],
        url: route[1],
        headers: {
          Authorization: `SharedKey ${config.apiSharedKey}`,
        },
        data: route[2],
      }).then((resp) => {
        assert.equal(resp.status, 200);
      });
    });
  });
});
