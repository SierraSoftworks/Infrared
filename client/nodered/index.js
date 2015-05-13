var http = require('http');
var url = require('url');

exports.register = function(server, nodeType, nodeDetails, callback) {
	
	if(typeof nodeDetails !== 'object') return setTimeout(function() {
		callback(new Error("Expected nodeDetails argument to be a hashmap"));
	}, 0);
	
	if(!nodeDetails.hostname) return setTimeout(function() {
		callback(new Error("No hostname provided in nodeDetails argument"));
	}, 0);
	
	if(!nodeDetails.port) return setTimeout(function() {
		callback(new Error("No port provided in nodeDetails argument"));
	}, 0);
	
	var hostUrl = url.parse(server);
	var req = http.request({
		hostname: hostUrl.hostname,
		port: hostUrl.port || 80,
		method: 'POST',
		path: '/api/v1/' + nodeType,
		headers: {
			"Content-Type": "application/json"
		}
	}, function(res) {
		res.setEncoding("utf8");
		
		var data = "";
		
		res.on("data", function(chunk) {
			data += chunk;
		});
		
		res.on("end", function() {
			if(!callback) return;
			
			var body = JSON.parse(data);
			
			if(res.statusCode !== 200) {
				var err = new Error(body.message);
				err.error = err.error;
				err.code = err.code;
				
				callback(err);
			} else {
				callback(null, body);
			}
		});
	});
	
	req.on('error', function(err) {
		callback(err);
		callback = null;
	});
	
	req.write(JSON.stringify(nodeDetails));
	req.end();
};