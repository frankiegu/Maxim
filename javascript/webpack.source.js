const path = require('path')

module.exports = {
    entry: './src/maxim.js',
    output: {
        path: __dirname + '/dist',
        filename: 'maxim.js',
    },
    resolveLoader: {
        alias: {
            'coffee-loader': path.join(__dirname, '/lib/coffee2-loader.js'),
        }
    },
    module: {
        rules: [{
            test: /\.coffee$/,
            exclude: /(node_modules|bower_components)/,
            use: {
                loader: 'coffee-loader?babel-loader'
            }
        }]
    }
}