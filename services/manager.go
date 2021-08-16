package services

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Connection Manager
type Connection struct {
	Members         int
	containerId     string
	cancelDelete    chan bool
	inDeletionQueue bool
	mu              sync.RWMutex
}

type ConnectionManager struct {
	connections        *sync.Map
	OnDeleteConnection func(string) error
}

func (m *ConnectionManager) GetConnection(id string) (*Connection, error) {
	if val, ok := m.connections.Load(id); ok {
		connection, ok := val.(*Connection)
		if !ok {
			error_msg := "Could not cast connection manager value to connection type"
			log.Fatalln(error_msg)
			return nil, errors.New(error_msg)
		}
		return connection, nil
	}
	return nil, nil
}

func (m *ConnectionManager) StartedContainer(id string) error {
	value, loaded := m.connections.LoadOrStore(id, &Connection{
		Members:         0,
		containerId:     id,
		cancelDelete:    make(chan bool),
		inDeletionQueue: true,
	})
	if loaded {
		log.Printf("Connection %s already exists!\n", id)
		return nil
	}
	connection := value.(*Connection)
	go m.deleteAfterInterval(connection)
	return nil
}

// Container ID should be full length
func (m *ConnectionManager) NewConnection(id string) error {
	log.Printf("ConnectionManager: Received New Connection request of %s\n", id)

	value, ok := m.connections.Load(id)
	if !ok {
		return fmt.Errorf("no container with ID %s found", id)
	}

	connection, ok := value.(*Connection)
	if !ok {
		error_msg := "could not cast connection manager value to connection type"
		log.Fatalln(error_msg)
		return errors.New(error_msg)
	}

	connection.mu.Lock()
	defer connection.mu.Unlock()

	if connection.inDeletionQueue {
		connection.cancelDelete <- true
		connection.inDeletionQueue = false
	}
	connection.Members++

	return nil
}

func (m *ConnectionManager) RemoveConnection(id string) error {
	log.Printf("ConnectionManager: Received Remove Connection request of %s\n", id)
	connection, err := m.GetConnection(id)
	if err != nil {
		return err
	}
	if connection == nil {
		log.Printf("Connection %s does not exist\n", id)
		return nil
	}
	connection.mu.Lock()
	defer connection.mu.Unlock()

	connection.Members--

	if connection.Members <= 0 {
		connection.Members = 0
		if !connection.inDeletionQueue {
			connection.inDeletionQueue = true
			go m.deleteAfterInterval(connection)
		}
	}
	return nil
}

func (m *ConnectionManager) deleteAfterInterval(connection *Connection) {
	log.Printf("started deletion period of: %s\n", connection.containerId)

	select {
	case <-time.After(1 * time.Minute):
		connection.mu.Lock()
		if !connection.inDeletionQueue {
			return
		}
		connection.inDeletionQueue = false
		connection.mu.Unlock()
		m.deleteImmediately(connection)
	case <-connection.cancelDelete:
		log.Printf("cancelled deletion of %s\n", connection.containerId)
	}
}

func (m *ConnectionManager) deleteImmediately(connection *Connection) {
	errStop := m.OnDeleteConnection(connection.containerId)
	if errStop != nil {
		log.Fatalf("error deleting connection %s: %v \n", connection.containerId, errStop)
	}
	m.connections.Delete(connection.containerId)
}

func NewConnectionManager(onDeleteConnection func(string) error) *ConnectionManager {
	return &ConnectionManager{
		connections:        &sync.Map{},
		OnDeleteConnection: onDeleteConnection,
	}
}

var ContainerManager *ConnectionManager

func init() {
	ContainerManager = NewConnectionManager(StopContainer)
}
