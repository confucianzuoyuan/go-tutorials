package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Make a new ToDo type that is a typed collection of fields
// (Title and Status), both of which are of type string
type ToDo struct {
	Title, Status string
}

type EditToDo struct {
	Title, NewTitle, NewStatus string
}

type Task int

// Declare variable 'todoSlice' that is a slice made up of
// type ToDo items
var todoSlice []ToDo

// GetToDo takes a string type and returns a ToDo
func (t *Task) GetToDo(title string, reply *ToDo) error {
	var found ToDo
	// Range statement that iterates over todoArray
	// 'v' is the value of the current iterateee
	for _, v := range todoSlice {
		if v.Title == title {
			found = v
		}
	}
	// found will either be the found ToDo or a zerod ToDo
	*reply = found
	return nil
}

func (t *Task) GetSlice(title string, reply *[]ToDo) error {
	*reply = todoSlice
	return nil
}

// MakeToDo takes a ToDo type and appends to the todoArray
func (t *Task) MakeToDo(todo ToDo, reply *ToDo) error {
	todoSlice = append(todoSlice, todo)
	*reply = todo
	return nil
}

// EditToDo takes a string type and a ToDo type and edits an item in the todoArray
func (t *Task) EditToDo(todo EditToDo, reply *ToDo) error {
	var edited ToDo
	// 'i' is the index in the array and 'v' the value
	for i, v := range todoSlice {
		if v.Title == todo.Title {
			todoSlice[i] = ToDo{todo.NewTitle, todo.NewStatus}
			edited = ToDo{todo.NewTitle, todo.NewStatus}
		}
	}
	// edited will be the edited ToDo or a zeroed ToDo
	*reply = edited
	return nil
}

// DeleteToDo takes a ToDo type and deletes it from todoArray
func (t *Task) DeleteToDo(todo ToDo, reply *ToDo) error {
	var deleted ToDo
	for i, v := range todoSlice {
		if v.Title == todo.Title && v.Status == todo.Status {
			// Delete ToDo by appending the items before it and those
			// after to the todoArray variable
			todoSlice = append(todoSlice[:i], todoSlice[i+1:]...)
			deleted = todo
			break
		}
	}
	*reply = deleted
	return nil
}

func main() {
	task := new(Task)
	// Publish the receivers methods
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 1234)
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}
