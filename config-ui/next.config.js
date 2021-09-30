module.exports = {
  pageExtensions: ['page.jsx', 'js'],

  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback.fs = false;
    }
    return config;
  },
}
