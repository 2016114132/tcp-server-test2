# TCP Echo Server - Test #2

This is an improved version of the basic TCP Echo Server written in Go. It supports multiple clients concurrently and includes features such as logging, input validation, timeout handling, custom command protocol, and personality-based responses.

## Clone the Repository

Clone this repository to your local machine:

```bash
git clone https://github.com/2016114132/tcp-server-test2.git
```
Go into the directory:

```bash
cd tcp-server-test2
```

## How to Run

Make sure you have Go installed. Then run the server:

```bash
go run main.go --port=4000
```

You can change the port by passing a different value to the `--port` flag:

```bash
go run main.go --port=5000
```

In another terminal, connect using `nc`:

```bash
nc localhost 4000
```

## Demo Video

Watch the video demonstration here: 
https://www.youtube.com/watch?v=E0FZ1ld5on0


## Reflection

- **Most educationally enriching feature:**  
  The most educationally enriching feature for me was implementing the command-line flag for port configuration using Go's built-in flag package. Even though I had used flag in previous assignments, I see its value in building configurable applications that can be customized at runtime without needing to modify the code.

- **Feature that required the most research:**  
  The feature that required the most research was implementing the command protocol. At first, I overthought the problem and believed I would need to implement a complex parser to handle commands like /time or /echo. I eventually realized that the solution was much simpler — I just needed to check if the message started with a /, then parse and handle the command accordingly.

