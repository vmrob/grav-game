const path = require('path');
// eslint-disable-next-line
const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: path.resolve('src', 'app.jsx'),
    output: {
        path: path.resolve('dist'),
        filename: 'app.bundle.js',
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    presets: ['react'],
                },
            },
            {
                test: /\.js$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    presets: ['es2015', 'react'],
                },
            },
            {
                test: /\.css$/,
                use: ExtractTextPlugin.extract({
                    use: ['css-loader'],
                }),
            }
        ],
    },
    stats: {
        colors: true,
    },
    plugins: [
        // Build html file, inject output files
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: 'index.html',
            inject: true,
        }),
        // Extract css into separate file
        new ExtractTextPlugin({
            filename: 'styles.css',
            allChunks: true,
        })
    ],
    resolve: {
        extensions: ['.js', '.jsx'],
    },
};
