Simple chat app, where users can login/signup, and then connect to a global chatroom where they can send messages. Each user connects to the server using websockets, and the backend server which handles
the websocket connections is written in Go. 
Inspired from Yik Yak, need to work on anonymity, and also fixing some CORS, and user authentication issues.
Plan on integrating a sql database soon, mainly for user authentication and possibly storing messages
Very much a work in progress!
