package main

import (
	"log"
	"net"
	"sync"
	"time"
)

// Server struct holds the TCP listener and a channel for handling connections.
type Server struct {
	listener    net.Listener   // TCP listener to accept incoming client connections.
	connChannel chan net.Conn  // Channel for sending accepted connections to worker goroutines.
	wg          sync.WaitGroup // WaitGroup to wait for all worker goroutines to finish.
}

// NewServer initializes a new Server.
func NewServer(port string) (*Server, error) {
	// Create a new TCP listener.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, err
	}

	// Initialize and return a new Server.
	return &Server{
		listener:    listener,
		connChannel: make(chan net.Conn), // Create a new channel for handling connections.
	}, nil
}

// Start runs the server.
func (s *Server) Start() {
	// Start 10 worker goroutines.
	for i := 0; i < 10; i++ {
		s.wg.Add(1)   // Increment the WaitGroup counter.
		go s.worker() // Start a new worker.
	}

	// This goroutine accepts incoming client connections.
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				// If an error occurs, stop accepting new connections.
				log.Println("Stopping connection loop.")
				close(s.connChannel) // Close the connection channel.
				return
			}
			s.connChannel <- conn // Send the new connection to the worker pool.
		}
	}()

}

// Stop shuts down the server.
func (s *Server) Stop() {
	s.listener.Close() // Close the TCP listener.
	s.wg.Wait()        // Wait for all worker goroutines to finish.
}

// Worker goroutine to handle connections.
func (s *Server) worker() {
	defer s.wg.Done() // Decrement the WaitGroup counter when the goroutine completes.
	for conn := range s.connChannel {
		s.do(conn) // Process each incoming connection.
	}
}

// Do performs the actual data processing for each connection.
func (s *Server) do(conn net.Conn) {
	buf := make([]byte, 1024) // Buffer to read incoming data.

	// Set timeouts for read and write operations.
	conn.SetReadDeadline(time.Now().Add(120 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	// Read incoming data into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		log.Println(err) // Log any errors.
		return
	}

	// Simulate data processing delay.
	time.Sleep(50 * time.Microsecond)

	// Send a response back to the client.
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello, World!\r\n"))
	conn.Close() // Close the connection.
}

func main() {
	server, err := NewServer(":1729")
	if err != nil {
		log.Fatal(err) // If initialization fails, terminate the program.
	}

	// Start the server.
	server.Start()
}
