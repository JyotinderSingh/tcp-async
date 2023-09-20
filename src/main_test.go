package main

import (
	"bufio"
	"net"
	"os"
	"sync"
	"testing"
)

var server *Server // Assuming YourServerType is the type of your server.

func TestMain(m *testing.M) {
	// Setup: Initialize the server once for all tests.
	var err error
	server, err = NewServer(":1729")
	if err != nil {
		os.Exit(1)
	}
	go server.Start()

	// Run all tests.
	code := m.Run()

	// Teardown: Stop the server.
	server.Stop()

	// Exit.
	os.Exit(code)
}

// Test for single client
func TestServer_SingleClient(t *testing.T) {
	testClient(t, nil)
}

// Test for multiple clients
func TestServer_MultipleClients(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go testClient(t, &wg)
	}

	wg.Wait()
}

func testClient(t *testing.T, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	// Establish a TCP connection to the server.
	conn, err := net.Dial("tcp", "localhost:1729")
	if err != nil {
		t.Errorf("Could not connect to server: %v", err)
		return
	}

	// Send a GET request to the server.
	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n\r\n"))
	if err != nil {
		t.Errorf("Could not write to connection: %v", err)
		return
	}

	// Read the server response using a buffered reader.
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		t.Errorf("Could not read from connection: %v", err)
		return
	}
	if line != "HTTP/1.1 200 OK\r\n" {
		t.Errorf("Unexpected response: %s", line)
	}
	conn.Close()
}
