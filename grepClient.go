// grep_client
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

const (
	PORT        = "8008"
	MASTER_LIST = "masterlist.txt"
)

func main() {
	ipList := []string{}
	file, _ := os.Open(MASTER_LIST)
	scanner := bufio.NewScanner(file)

	//Compile list of ip address from masterlist.txt
	for scanner.Scan() {
		var ip_content = scanner.Text()
		ip_content = ip_content + ":" + PORT
		ipList = append(ipList, ip_content)
	}

	t0 := time.Now()

	if len(os.Args) < 1 {
		fmt.Println("Please provide the string or regular expression. Syntax: go run grepClient.go <optional parameters -c/-w> <string/regular expression> ")
		os.Exit(1)
	} else {
		c := make(chan string)
		
		var serverInput string = ""
		for i := 1; i <  len(os.Args); i++ {
			serverInput += os.Args[i]
			if i != (len(os.Args) -1) {
				serverInput += " "
			}
			
		}
		
		// Connect to every server in masterlist.txt
		for i := 0; i < len(ipList); i++ {
			go writeToServer(ipList[i], serverInput, c)
		}
		// Print results from server and write to a file
		_, err := os.Stat("logGrep")
		if	os.IsNotExist(err) {
		  _, err := os.Create("logGrep")
		  if err != nil {
	        panic(err)
		  }
		} 
		f, err := os.OpenFile("logGrep", os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
				panic(err)
				}
	    defer f.Close()
    
		for i := 0; i < len(ipList); i++ {
			serverResult := <-c
			fmt.Println(serverResult)
			fmt.Println("----------")
			_, err = f.WriteString(serverResult)
		}
		f.Sync()
		w := bufio.NewWriter(f)
		w.Flush()
	}

	t1 := time.Now()
	fmt.Print("Function took: ")
	fmt.Println(t1.Sub(t0))
}

/*
 * Writes the key value regex patterns as one regex for the grep
 * to use
 * @param key the grep pattern for the key
 * @param value the grep pattern for the value
 * @return the single grep pattern to query key-values in log files
 */
func rewriteKeyAndValue(key string, value string) string {
	var newKey string
	var newValue string

	newKey = key
	newValue = value

	// Check ^ on key
	if strings.HasPrefix(newKey, "^") {
		newKey = newKey[1:len(newKey)]
	} else {
		newKey = "[^:]*" + newKey
	}

	// Check $ on key
	if strings.HasSuffix(newKey, "$") {
		newKey = newKey[0 : len(newKey)-1]
	} else {
		newKey = newKey + "[^:]*"
	}

	// Check ^ on value
	if strings.HasPrefix(newValue, "^") {
		newValue = newValue[1:len(newValue)]
	} else {
		newValue = "[^:]*" + newValue
	}

	// Check $ on value
	if strings.HasSuffix(newValue, "$") {
		newValue = newValue[0 : len(newValue)-1]
	} else {
		newValue = newValue + "[^:]*"
	}

	serverInput := "^" + newKey + ":" + newValue + "$"

	return serverInput
}

/*
 * Sends a message to a server, and returns the file into a channel
 * @param ipAddr string representation of the server's IP Address
 * @param message the message to be sent back to the server
 * @param c the channel for returning server messages
 */
func writeToServer(ipAddr string, message string, c chan string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipAddr)
	if err != nil {
		c <- err.Error()
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		c <- err.Error()
		return
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		c <- err.Error()
		return
	}

	result, err := ioutil.ReadAll(conn)
	if err != nil {
		c <- err.Error()
		return
	}

	c <- string(result)
}
