document.getElementById('loginForm').addEventListener('submit', function (e) {
    e.preventDefault();

    var uname = document.getElementById('username').value;
    var passwd = document.getElementById('password').value;
    var message = document.getElementById('message');
    // Define the URL of the local server
    const url = 'http://localhost:8080/form'; // Replace '/api' with your specific endpoint

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
    })
});

document.getElementById('signUpButton').addEventListener('click', function () {
    window.location.href = 'signup.html';
    var username = document.getElementById('newUsername').value;
    var password = document.getElementById('newPassword').value;
    var email = document.getElementById('email').value;
});
