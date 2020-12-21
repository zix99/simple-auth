module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Github', link: 'https://google.com' }
    ],
    sidebar: [
      {
        title: 'Simple Auth',
        path: '/',
        children: [
          '/quickstart',
          '/config',
          '/customization',
          '/cli',
        ],
      },
      {
        title: 'Login Providers',
        path: '/login',
        children: [
          '/login/local',
          '/login/oidc',
        ],
      },
      {
        title: 'Authenticators',
        path: '/authenticators',
        children: [
          '/authenticators/simple',
          '/authenticators/vouch',
        ],
      },
      {
        title: 'Cookbooks',
        path: '/cookbooks',
        children: [
          '/cookbooks/gateway',
          '/cookbooks/nginx-auth-request',
        ],
      },
    ]
  },
}