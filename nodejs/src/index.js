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

function Uint32ToUint8Arr(num) {

    arr = new ArrayBuffer(4); // an Int32 takes 4 bytes
    view = new DataView(arr);
    view.setUint32(0, num, false); // byteOffset = 0; litteEndian = false
    return arr;
}

function Uint8ArrToInt(Uint8Arr) {
    var length = Uint8Arr.length;

    let buffer = Buffer.from(Uint8Arr);
    var result = buffer.readUIntBE(0, length);

    return result;
}

var client = net.connect(PIPE_PATH, function() {
    console.log('Client: on connection');
})

client.on('data', function(data) {
    console.log('Client: on data:', data.toString());
    //client.end('Thanks!');
});

client.on('end', function() {
    console.log('Client: on end');
})

function MsgToBuf(msg) {
    var len = Buffer.byteLength(msg)
    var bb = Buffer.alloc(4 + len);
    bb.writeUInt32BE(len)
    bb.write(msg,4)
    return bb
}

client.write(MsgToBuf('hello!!你好'));
