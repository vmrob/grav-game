const fs = require('fs');
const os = require('os');
const path = require('path');
const util = require('util');
const exec = util.promisify(require('child_process').exec);

module.exports = function(source) {
    const callback = this.async();

    util.promisify(fs.mkdtemp)(path.join(os.tmpdir(), 'go-wasm-'))
        .then((dir) => {
            const inputPath = path.join(dir, 'wasm.go');
            const outputPath = path.join(dir, 'wasm.wasm');
            util.promisify(fs.writeFile)(inputPath, source)
                .then(() => {
                    return exec(`GOOS=js GOARCH=wasm go build -o ${outputPath} ${path.relative('', dir)}`); 
                })
                .then(() => {
                    return util.promisify(fs.readFile)(outputPath);
                })
                .then((buf) => {
                    util.promisify(fs.unlink)(inputPath)
                        .then(() => {
                            return util.promisify(fs.unlink)(outputPath);
                        })
                        .then(() => {
                            return util.promisify(fs.rmdir)(dir);
                        })
                        .then(() => {
                            callback(null, buf);
                        })
                        .catch(callback);
                })
                .catch(callback);
        })
        .catch(callback);
};
