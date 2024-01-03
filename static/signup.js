document.getElementById('signUpForm').addEventListener('submit', function (e) {
    e.preventDefault();
    var uname = document.getElementById('newUsername').value;
    var passwd = document.getElementById('newPassword').value;
    var em = document.getElementById('email').value;
    const url = 'http://localhost:8080/signup'; // Replace '/api' with your specific endpoint
    const data = {
        username: uname,
        password: passwd,
        email: em
    };
    console.log(data)
    fetch(url, {
        method: 'POST', // Specify the method
        headers: {
            'Content-Type': 'application/json', // Specify the content type
        },
        body: JSON.stringify(data), // Convert the JavaScript object to a JSON string
    })


});
