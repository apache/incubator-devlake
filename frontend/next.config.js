module.exports = {
  // reactStrictMode: true,
  basePath: '/frontend',

  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback.fs = false;
    }
    return config;
  },
}
