module.exports = {
    entry: '../test/script.src.js',
    output: {
        path: __dirname + '/../test',
        filename: 'script.js',
    },
    module: {
        rules: [{
            test: /\.js$/,
            exclude: /(node_modules|bower_components)/,
            use: {
                loader: 'babel-loader'
            }
        }]
    }
}