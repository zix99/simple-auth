const repo = 'zix99/simple-auth';
const docsUrl = "https://simple-auth.surge.sh";
const repoUrl = `https://github.com/${repo}`;
const fileUrl = `${repoUrl}/tree/master`

module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  head: [
    ['link', { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" }],
    ['script', { src: 'https://stats.zdyn.net/umami.js', 'data-website-id': 'ec52632a-21d3-400c-ad85-cab25b75dcf8', async: true, defer: true }],
  ],
  themeConfig: {
    repo,
    docsDir: 'docs',
    editLinks: true,
    editLinkText: 'Help us improve this page!',
    repoUrl,
    fileUrl,
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
          '/download',
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
          '/authenticators/oauth2',
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
          '/cookbooks/scaling',
          '/cookbooks/decodejwt',
          '/cookbooks/signingkey-pair',
          '/cookbooks/login-redirect',
          '/cookbooks/restrictcreateuser',
          '/cookbooks/tls',
          '/cookbooks/prometheus',
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