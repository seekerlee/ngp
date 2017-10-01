var net = require('net');

var PIPE_NAME = "secrettunnel";
var PIPE_PATH = "\\\\.\\pipe\\" + PIPE_NAME;


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
