const docsUrl = "https://simple-auth.surge.sh";

module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  themeConfig: {
    repoUrl: "https://github.com/zix99/simple-auth/tree/master",
    docsUrl,
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quickstart', link: '/quickstart' },
      { text: 'API Docs', link: `${docsUrl}/api` },
    ],
    sidebar: [
      {
        title: 'Simple Auth',
        path: '/',
        collapsable: false,
        sidebarDepth: 2,
        children: [
          '/quickstart',
          '/config',
          '/customization',
          '/email',
          '/database',
          '/cli',
        ],
      },
      {
        title: 'Login Providers',
        path: '/login',
        sidebarDepth: 2,
        children: [
          '/login/local',
          '/login/oidc',
        ],
      },
      {
        title: 'Authenticators',
        path: '/authenticators',
        sidebarDepth: 2,
        children: [
          '/authenticators/simple',
          '/authenticators/vouch',
        ],
      },
      {
        title: 'Access Layer',
        path: '/access',
        sidebarDepth: 2,
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
          '/cookbooks/traefik',
          '/cookbooks/decodejwt',
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