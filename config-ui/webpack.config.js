/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
// DEVELOPMENT ONLY WEBPACK CONFIG
const path = require('path')
const webpack = require('webpack')
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin
const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyPlugin = require('copy-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const ESLintPlugin = require('eslint-webpack-plugin')
const Dotenv = require('dotenv-webpack')

module.exports = (env = {}) => {
  const optionalPlugins = []

  // Only appends bundle analyzer when forced via env
  if (env.ANALYZE_BUILD === 1) {
    optionalPlugins.push(new BundleAnalyzerPlugin())
  }

  return {
    entry: path.resolve(__dirname, './src/index.js'),
    mode: 'development',
    module: {
      rules: [
        {
          test: /\.jsx?$/,
          use: ['babel-loader'],
          exclude: [/node_modules/, /packages/, /cypress/, /^config$/],
        },
        {
          test: /\.css$/,
          use: [
            'style-loader',
            'css-loader'
          ]
        },
        {
          test: /\.scss$/,
          use: [
            'style-loader',
            'css-loader',
            'resolve-url-loader',
            // 'sass-loader'
            {
              loader: 'sass-loader',
              options: {
                implementation: require('node-sass'),
                sourceMap: true,
                additionalData: '@import "@/styles/theme.scss";',
              }
            }
          ]
        },
        {
          test: /\.html$/,
          use: ['html-loader']
        },
        {
          test: /\.(ttf|eot)$/,
          use: {
            loader: 'file-loader',
            options: {
              name: 'fonts/[hash].[ext]',
            },
          },
        },
        {
          test: /\.(?:png|jpe?g|gif|ttf|woff|woff2)$/,
          loader: 'url-loader',
          options: {
            limit: 10 * 1024,
          },
        },
        {
          test: /\.svg$/,
          use: ['@svgr/webpack', 'url-loader'],
        },
        {
          test: /\.po$/,
          use: ['@lingui/loader']
        }
      ]
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src/'),
        '@config': path.resolve(__dirname, './config/'),
      },
      extensions: ['*', '.js', '.jsx']
    },
    output: {
      path: path.resolve(__dirname, './dist'),
      filename: 'bundle.js',
      publicPath: '/'
    },
    plugins: [
      new webpack.DefinePlugin({
        'process.env': {
          LOCAL: true,
        },
      }),
      new Dotenv(),
      new CleanWebpackPlugin(),
      new webpack.HotModuleReplacementPlugin(),
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, './src/index.html'),
        filename: 'index.html',
        favicon: path.resolve(__dirname, './src/images/favicon.ico'),
      }),
      new webpack.EnvironmentPlugin({ ...process.env }),
      new CopyPlugin({
        patterns: [
          { from: 'src/images/logo.svg', to: 'logo.svg' },
        ],
      }),
      new ESLintPlugin({
        fix: true,
        context: path.resolve(__dirname, './'),
        exclude: ['dist', 'packages', 'cypress', 'config', 'node_modules']
      }),
      ...optionalPlugins
    ],
    devServer: {
      hot: true,
      port: 4000,
      host: '0.0.0.0',
      historyApiFallback: true,
      proxy: {
        '/api': { target: 'http://[::1]:8080', pathRewrite: { '^/api': '' }, changeOrigin: true },
      }
    },
    devtool: 'source-map'
  }
}
