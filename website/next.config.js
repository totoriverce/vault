const withHashicorp = require('@hashicorp/nextjs-scripts')

module.exports = withHashicorp({
  transpileModules: ['is-absolute-url', '@hashicorp/react-.*'],
})({
  svgo: { plugins: [{ removeViewBox: false }] },
  rewrites: () => [
    {
      source: '/api/:path*',
      destination: '/api-docs/:path*',
    },
  ],
  redirects: () => [
    {
      source: '/intro',
      destination: '/intro/getting-started',
      permanent: false,
    },
  ],
  env: {
    HASHI_ENV: process.env.HASHI_ENV || 'development',
    SEGMENT_WRITE_KEY: 'OdSFDq9PfujQpmkZf03dFpcUlywme4sC',
    BUGSNAG_CLIENT_KEY: '07ff2d76ce27aded8833bf4804b73350',
    BUGSNAG_SERVER_KEY: 'fb2dc40bb48b17140628754eac6c1b11',
  },
})
