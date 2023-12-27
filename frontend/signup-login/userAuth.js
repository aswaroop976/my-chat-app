document.getElementById('loginForm').addEventListener('submit', function (e) {
    e.preventDefault();

    var username = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var message = document.getElementById('message');
// need to change all of this below:    
    if (username === "admin" && password === "password") {
        message.innerHTML = "Login successful!";
        message.style.color = "green";
    } else {
        message.innerHTML = "Login failed. Invalid username or password.";
        message.style.color = "red";
    }
});

document.getElementById('signUpButton').addEventListener('click', function () {
    window.location.href = 'signup.html';
    var username = document.getElementById('newUsername').value;
    var password = document.getElementById('newPassword').value;
    var email = document.getElementById('email').value;
});
