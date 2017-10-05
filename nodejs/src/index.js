const SendMsg = require('./pipe/pipe.js');
SendMsg('haha', function(response){
    console.log(response)
})
