const HOST = process.env.SATEST_HOST || 'localhost:9002';

module.exports = {
  apiSharedKey: 'super-secret',
  simpleAuthKey: 'your-super-secret-token',
  baseURL: `http://${HOST}`,
};
