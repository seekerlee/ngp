const net = require('net');
const os = require('os');
const _ = require('lodash');

var PIPE_NAME = "secrettunnel";
var PIPE_PATH
const platform = os.platform()
if(platform == 'win32') {
    PIPE_PATH = "\\\\.\\pipe\\" + PIPE_NAME
} else {
    PIPE_PATH = _.trimEnd(os.tmpdir(), '/') + '/' + PIPE_NAME
}

var client = net.connect(PIPE_PATH, function() {
    console.log('Client: on connection');
})

var msg = ''
client.on('data', function(data) {
    msg += data
});

client.on('end', function() {
    console.log('Client: on end');
    console.log(msg);
})

function MsgToBuf(msg) {
    var len = Buffer.byteLength(msg)
    var bb = Buffer.alloc(4 + len);
    bb.writeUInt32BE(len)
    bb.write(msg,4)
    return bb
}

client.write(MsgToBuf('hello!!你好'));
