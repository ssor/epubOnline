var tasks = [
    { src: "", dest: "release", type: "file_or_dir_ignore" }, // clean
    { src: "conf/config_default.json", dest: "release/conf/config_default.json", type: "file_or_dir" },
    { src: "static", dest: "release/static", type: "file_or_dir" },
    { src: "views", dest: "release/views", type: "file_or_dir" },
    { src: "", dest: "release/views/index.html", type: "file_or_dir_ignore" },
    { src: "GOOS=linux GOARCH=amd64 go build -o ./release/epub_online_linux64 main.go", dest: "", type: "cmd" },
]

const fs = require('fs-extra');
const R = require('ramda');
const path = require('path');
const spawn = require('child_process');

var mkdir = task => {
    if (task.dest.length <= 0 || task.src.length <= 0) {
        return task
    }
    var dir_name = path.dirname(task.dest)
    if (!fs.existsSync(dir_name)) {
        fs.mkdirSync(dir_name)
    }
    return task
};

var doCopyTask = task => {
    if (task.type == "file_or_dir") {
        fs.copySync(task.src, task.dest)
    }
    return task
}

var doIgnoreTask = task => {
    if (task.type == "file_or_dir_ignore") {
        fs.removeSync(task.dest)
    }
    return task
}

var doCmd = task => {
    if (task.type == "cmd") {
        var output = spawn.execSync(task.src)
        if (output.toString().length > 0) {
            console.log("cmd-> ", output.toString())
        }
    }
    return task
}

var doTask = R.compose(mkdir, doCopyTask, doIgnoreTask, doCmd)
R.forEach(doTask, tasks)