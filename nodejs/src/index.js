const SendMsg = require('./pipe/pipe.js');
SendMsg('haha', function(response){
    console.log(response)
})

client.on('data', function(data) {
    console.log('Client: on data:', data.toString());
    client.end('Thanks!');
});

client.on('end', function() {
    console.log('Client: on end');
})

//client.write('Hello, server! Love, Client.\n');
//client.write('2222\n');
client.write('1\n');
