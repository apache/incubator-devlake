module.exports = {
  // reactStrictMode: true,
  // basePath: '/',

  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback.fs = false;
    }
    return config;
  },
}
