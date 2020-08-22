module.exports = {
  title: "Simple Auth",
  description: "Simple White-Labeled Authentication Provider",
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Guide', link: '/guide/' },
      { text: 'Github', link: 'https://google.com' }
    ],
    sidebar: [
      '/',
      {
        title: 'Login Providers',
        path: '/login',
        children: [
          '/login/simple',
          '/login/oidc',
        ],
      },
      '/config',
    ]
  },
}