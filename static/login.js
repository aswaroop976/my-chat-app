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
        .then(data => console.log(data));
});

//document.getElementById('signUpForm').addEventListener('click', function (e) {
//    e.preventDefault();
//    var uname = document.getElementById('newUsername').value;
//    var passwd = document.getElementById('newPassword').value;
//    var em = document.getElementById('email').value;
//    const url = 'http://localhost:8080/signup'; // Replace '/api' with your specific endpoint
//    const data = {
//        username: uname,
//        password: passwd,
//        email: em
//    };
//    console.log(data)
//    fetch(url, {
//        method: 'POST', // Specify the method
//        headers: {
//            'Content-Type': 'application/json', // Specify the content type
//        },
//        body: JSON.stringify(data), // Convert the JavaScript object to a JSON string
//    })
//
//
//});
