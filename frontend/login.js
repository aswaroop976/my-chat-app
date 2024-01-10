document.getElementById('loginForm').addEventListener('submit', function (e) {
    e.preventDefault();

    var uname = document.getElementById('username').value;
    var passwd = document.getElementById('password').value;
    // Define the URL of the local server
    const url = 'http://localhost:8080/login'; // Replace '/api' with your specific endpoint

    // Define the data you want to send
    const data = {
        username: uname,
        password: passwd
    };
    console.log(data)
    // Use the fetch API to send the POST request
    fetch(url, {
        method: 'POST', // Specify the method
        headers: {
            'Content-Type': 'application/json', // Specify the content type
        },
        body: JSON.stringify(data), // Convert the JavaScript object to a JSON string
    }).then(response => response.json())
        .then(data => {
            console.log(data.content)
            if (data.content == "Allow login") {
                window.location.href = 'http://localhost:8080/chat'
            }
        })
});

