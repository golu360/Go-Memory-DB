package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Database struct {
	data map[string]string
}

func (d Database) print() {
	fmt.Println("Database:", d.data)
}

func (d Database) delete(key string) {
	delete(d.data, key)
}

var data Database

func init() {
	data = Database{
		map[string]string{},
	}
}

func main() {
	li, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	log.Output(1, "Server Started on PORT "+os.Args[1])
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Panic(err)
			continue
		}
		go handle(conn, &data)
	}
}

func handle(con net.Conn, d *Database) {
	defer con.Close()
	err := con.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatalln("CONN TIMEOUT")
	}
	commands := map[string]string{
		"SET":    "SET",
		"GET":    "GET",
		"DELETE": "DELETE",
	}
	scanner := bufio.NewScanner(con)
	if scanner.Scan() {
		line := scanner.Text()
		command := strings.Fields(line)

		if _, present := commands[command[0]]; !present {
			fmt.Fprint(con, "Invalid Command Passed\n")
		} else {
			switch command[0] {
			case "SET":
				d.data[command[1]] = command[2]
				fmt.Fprint(con, "Key: "+command[1]+" Value: "+command[2]+" is SET\n")
				fmt.Printf("Data Set\n")
				d.print()
			case "GET":
				fmt.Fprintf(con, d.data[command[1]])
				fmt.Printf("Data Get\n")
				d.print()
			case "DELETE":
				d.delete(command[1])
				fmt.Fprintf(con, "Deleted Key Successfully\n")
				fmt.Printf("Data Set\n")
				d.print()
			}
		}
	}
}
