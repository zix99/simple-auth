module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quickstart', link: 'quickstart' },
      { text: 'API Docs', link: 'http://simple-auth.surge.sh/api' },
    ],
    sidebar: [
      {
        title: 'Simple Auth',
        path: '/',
        collapsable: false,
        children: [
          '/quickstart',
          '/config',
          '/customization',
          '/email',
          '/cli',
        ],
      },
      {
        title: 'Login Providers',
        path: '/login',
        collapsable: false,
        children: [
          '/login/local',
          '/login/oidc',
        ],
      },
      {
        title: 'Authenticators',
        path: '/authenticators',
        collapsable: false,
        children: [
          '/authenticators/simple',
          '/authenticators/vouch',
        ],
      },
      {
        title: 'Access Layer',
        path: '/access',
        collapsable: false,
        children: [
          '/access/cookie',
          '/access/gateway',
        ],
      },
      {
        title: 'Cookbooks',
        path: '/cookbooks',
        children: [
          '/cookbooks/gateway',
          '/cookbooks/nginx-auth-request',
          '/cookbooks/signingkey-pair',
        ],
      },
      {
        title: 'API Docs',
        path: 'http://simple-auth.surge.sh/api',
      },
    ]
  },
}