# Broadcast Server
This is a naive CLI tool implementation of an broadcast server using WebSockets communication protocol. Handles multiple client connections and graceful server shutdown. The server will listen for upcoming client connection requests and upgrade them to a WebSocket connection, once the connection is established it will handle upcoming messages sent by any client and echo them to every connected client.


## Usage
To start the server run the command:
```
broadcast-server start
```

You can specify the port for the server to listen to with the `-port` flag. If port is not provided it defaults to `:8080`
```
broadcast-server start -port 3000
```

Start a new client with the command:
```
broadcast-server connect
```

Example: Sending messages as a client
```
$ broadcast-server connect
dial to <server-url> successful
<enter your message here and hit Enter>
> message sent
> received message: <the message you sent>

```

Example: Receiving messages from other clients
```
$ broadcast-server connect
dial to <server-url> successful
...
> received message: <an awesome message>
```

For a detailed description of the commands use `broadcast-server help`

## Acknowledgements
This app was inspired by [roadmap.sh's](https://roadmap.sh/projects/broadcast-server) project idea.
