const path = require('path');
// eslint-disable-next-line
const webpack = require('webpack');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');

const goLoader = function(source) {
    var callback = this.async();
    var options = loaderUtils.getOptions(this);
    var command = exec(options.script, function(err, result) {
        if (err) return callback(err);
        callback(null, result);
    });
    command.stdin.write(source);
    command.stdin.end();
};

module.exports = {
    entry: path.resolve('src', 'app.jsx'),
    output: {
        path: path.resolve('../dist'),
        filename: 'app.bundle.js',
    },
    module: {
        rules: [
            {
                test: /\.go?$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: '[name].wasm',
                        },
                    },
                    'go-wasm',
                ],
            },
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
    resolveLoader: {
        modules: ['node_modules', path.resolve(__dirname, '../loaders')],
    },
};
