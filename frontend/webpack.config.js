const { VueLoaderPlugin } = require("vue-loader");
const TerserPlugin = require('terser-webpack-plugin');

module.exports = {
  entry: [ './main.ts' ],
  output: {
    path: `${__dirname}/../public/`,
    filename: 'gion.js'
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: 'vue-loader'
      },
      {
        test: /\.tsx?$/,
        loader: 'ts-loader',
        options: {
          appendTsSuffixTo: [/\.vue$/],
        },
        exclude: /node_modules/,
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/,
        options: {
          presets: [
            ["@babel/preset-env",
              {
                "useBuiltIns": "usage",
                corejs: 3
              }
            ],
          ],
        },
      },
    ]
  },
  plugins: [
    new VueLoaderPlugin()
  ],
  resolve: {
    extensions: ['.ts', '.js'],
    alias: {
      vue: "vue/dist/vue.esm-bundler.js"
    }
  },
  optimization: {
    minimizer: [
      new TerserPlugin({
        extractComments: false
      })
    ]
  },
  target: ['web', 'es5'],
}

