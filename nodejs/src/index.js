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

client.on('data', function(data) {
    console.log('Client: on data:', data.toString());
    client.end('Thanks!');
});

client.on('end', function() {
    console.log('Client: on end');
})

client.write('Hello, server! Love, Client.\n');
client.write('2222\n');
client.write('\n');
