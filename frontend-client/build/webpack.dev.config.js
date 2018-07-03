const merge = require('webpack-merge');
const baseWebpackConfig = require('./webpack.base.config');

module.exports = merge(baseWebpackConfig, {
    devtool: 'source-map',
    stats: {
        colors: true,
    },
});
