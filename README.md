# b0red

`b0red` is a simple `TCP` based chat server written in `Go`. It's designed to focus on simplicity and ease of use.

> [!INFO]
> This project is used to learn `Go` and networking protocols in general. Therefore, it's not intended for production use.

## Features

- Accept multiple `TCP` connections
- Broadcast messages to all connected `TCP` clients

## Usage

1. Run the server:
   ```bash
   go run main.go
   ```
2. Connect to the server using a `TCP` client (e.g., `netcat`):
   ```bash
   nc localhost 4000
   ```
3. The server will prompt you to enter your name - if you don't enter a name, the connection will be closed.

## Roadmap

- [ ] Add `unit` tests
- [ ] Add `integration` tests
- [ ] Add `Docker` development environment
- [ ] Add message persistence
- [ ] Add private messaging
