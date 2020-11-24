module.exports = {
  env: {
    mocha: true,
  },
  extends: [
    'airbnb-base',
  ],
  parserOptions: {
    ecmaVersion: 11,
    sourceType: 'module',
  },
  rules: {
    'max-len': ['warn', { code: 120 }],
    'arrow-body-style': 'off',
  },
};
