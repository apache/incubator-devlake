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

/* eslint-disable import/no-extraneous-dependencies */
const path = require('path')
const webpack = require('webpack')
const BundleAnalyzerPlugin = require('webpack-bundle-analyzer').BundleAnalyzerPlugin
const HtmlWebpackPlugin = require('html-webpack-plugin')
const CopyPlugin = require('copy-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')
const ESLintPlugin = require('eslint-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const CssMinimizerPlugin = require('css-minimizer-webpack-plugin')
const TerserPlugin = require('terser-webpack-plugin')

module.exports = (env = {}) => {
  const optionalPlugins = []

  if (env.ANALYZE_BUILD === 1) {
    optionalPlugins.push(new BundleAnalyzerPlugin())
  }

  return {
    entry: path.resolve(__dirname, 'src/index.js'),
    mode: 'production',
    devtool: 'source-map',
    module: {
      rules: [
        {
          test: /\.jsx?$/,
          exclude: [/node_modules/, /packages/, /cypress/, /^config$/],
          use: {
            loader: 'babel-loader',
            options: {
              babelrc: false,
              cacheDirectory: false,
              presets: ['@babel/preset-env', '@babel/preset-react', '@babel/preset-flow'],
              plugins: ['@babel/plugin-transform-runtime']
            }
          },
        },
        {
          test: /\.css$/,
          use: [
            // !WARNING! We only use style-loader for DEV MODE!
            // 'style-loader',
            MiniCssExtractPlugin.loader,
            {
              loader: 'css-loader',
              options: {
                modules: true,
                importLoaders: 1,
              }
            }
          ]
        },
        {
          test: /\.scss$/,
          use: [
            {
              loader: MiniCssExtractPlugin.loader,
            },
            {
              loader: require.resolve('css-loader'),
              options: {
                importLoaders: 1,
              },
            },
            {
              loader: require.resolve('postcss-loader'),
              options: {
                postcssOptions: {
                  plugins: [require('autoprefixer'), require('cssnano')({ preset: 'default' })],
                },
              },
            },
            {
              loader: 'resolve-url-loader'
            },
            {
              loader: 'sass-loader',
              options: {
                implementation: require('node-sass'),
                sourceMap: true,
                additionalData: '@import "@/styles/theme.scss";',
              }
            }
          ],
          // sideEffects: true
        },
        {
          test: /\.html$/,
          use: ['html-loader']
        },
        {
          test: /\.(png|gif|jpe?g)$/,
          loader: require.resolve('file-loader'),
          options: {
            name: '[name].[ext]?[hash]',
            outputPath: 'assets/',
          },
        },
        {
          test: /\.(eot|ttf|woff|woff2|svg)$/,
          loader: require.resolve('file-loader'),
          include: [
            path.resolve(__dirname, './src/fonts/'),
            path.resolve(__dirname, './node_modules/')
          ],
          options: {
            name: '[name].[ext]',
            outputPath: 'fonts/',
            esModule: false,
          },
        },
        {
          test: /\.svg$/,
          exclude: path.resolve(__dirname, './src/fonts'),
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
      // modules: ['node_modules'],
      extensions: ['*', '.js', '.jsx', '.scss']
    },
    output: {
      path: path.resolve(__dirname, './dist'),
      filename: '[name].[hash].js',
      publicPath: '/'
    },
    optimization: {
      minimize: true,
      minimizer: [
        new CssMinimizerPlugin(),
        new TerserPlugin({
          terserOptions: {
            // cache: true,
            parallel: true,
            sourceMap: false,
            compress: {
              drop_console: true,
            }
          }
        })
      ],
    },
    plugins: [
      new CleanWebpackPlugin(),
      new webpack.HotModuleReplacementPlugin(),
      new MiniCssExtractPlugin({ filename: '[name].[hash].css' }),
      new HtmlWebpackPlugin({
        template: path.resolve(__dirname, './src/index-production.html'),
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
        context: path.resolve(__dirname, './'),
        exclude: ['dist', 'packages', 'cypress', 'config', 'node_modules']
      }),
      ...optionalPlugins
    ]
  }
}
