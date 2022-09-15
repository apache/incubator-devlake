module.exports = {
  extends: ['mints'],
  plugins: ['header'],
  rules: {
    'header/header': ['warn', '.file-headerrc'],
  },
};
