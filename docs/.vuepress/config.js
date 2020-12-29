const docsUrl = "https://simple-auth.surge.sh";
const repoUrl = "https://github.com/zix99/simple-auth/tree/master";

module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  themeConfig: {
    repoUrl,
    docsUrl,
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quickstart', link: '/quickstart' },
      { text: 'API Docs', link: `${docsUrl}/api` },
      { text: 'Source', link: repoUrl },
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
        title: 'REST API',
        path: '/api',
        children: [
          { title: 'API Docs', path: `${docsUrl}/api`},
        ],
      },
    ]
  },
}