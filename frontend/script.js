
var socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = function (e) {
    console.log("[open] Connection established");
    socket.send("Hello, server!");
};

socket.onmessage = function (event) {
    console.log(`[message] Data received from server: ${event.data}`);
};

socket.onclose = function (event) {
    if (event.wasClean) {
        console.log(`[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`);
    } else {
        console.error('[close] Connection died');
    }
};

socket.onerror = function (error) {
    console.error(`[error] ${error.message}`);
};


document.getElementById('send-button').addEventListener('click', function () {
    var message = document.getElementById('message-input').value;
    if (message) {
        // Here you would send the message to the server
        console.log('Message sent:', message);
        document.getElementById('message-input').value = ''; // Clear input
        socket.send(message);
    }
});

// Additional functionality to retrieve and display messages will be added here
