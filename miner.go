// Program to mine Duino-Coin.
package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var username string = " " // User to mine to.
var diff string = " "     // Possible safe values: MEDIUM, LOW, NET.
var x int = 1             // Goroutines count.
var addr string = "103.253.43.216:2335" // Pool's IP:Pool's port for v4.0. asia-node-1
var key = " "

// Shares
var accepted int = 0
var rejected int = 0

func work() {
	conn, _ := net.Dial("tcp", addr)
	buffer := make([]byte, 5)
	_, err := conn.Read(buffer)
	log.Println("Server is on version: " + string(buffer))

	if err != nil {
		log.Println("Servers might be down or a routine may have restarted, quitting routine.")
		return
	}

	for {
		// get raspberry pi temperature
		cmd := exec.Command("vcgencmd", "measure_temp")
		tempBuf, err := cmd.CombinedOutput()
		var temp string
		if (err == nil) {
			temp = string(tempBuf)
			temp = strings.Replace(temp, "temp=", "", -1)
			temp = strings.Replace(temp, "'C\n", "", -1)
		} else {
			temp = ""
		}

		// Requesting a job.
		_, err = conn.Write([]byte("JOB," + username + "," + diff + "," + key + "," + temp))

		if err != nil {
			log.Fatal("Error requesting job.")
		}

		// Making a buffer for the job.
		buffer := make([]byte, 2048)
		_, err = conn.Read(buffer) // Getting the jobs.

		if (err != nil) {
			log.Println(err)
			log.Fatal("Error getting the job.")
		}
		buffer = bytes.Trim(buffer, "\x00")
		
		job := strings.Split(strings.TrimSpace(string(buffer)), ",") // Parsing the job.
		if (len(job)!= 3) {
			continue
		}
		hash := job[0]
		goal := job[1]

		// Removes null bytes from job then converts it to an int.
		diff, _ := strconv.Atoi(job[2])

		start_time := time.Now()

		for i := 0; i <= diff * 100; i++ {
			h := sha1.New()
			h.Write([]byte(hash + strconv.Itoa(i))) // Hash
			nh := hex.EncodeToString(h.Sum(nil))
			if nh == goal {
				end_time := time.Since(start_time)
				hash_rate := float64(i) / float64(end_time.Seconds())
				// Sends the result of hash algorithm to the pool.
				_, err = conn.Write([]byte(strconv.Itoa(i) + "," + strconv.FormatFloat(hash_rate, 'f', 2, 64) + ",Go Miner on RPi,Go Miner"))

				if err != nil {
					log.Println("Error writing hash result")
					break
				}

				feedback_buffer := make([]byte, 20)
				_, err = conn.Read(feedback_buffer) // Reads response.

				if err != nil{
					log.Println("Error receiving feedback")
					log.Fatal(err)
				}

				feedback_buffer = bytes.Trim(feedback_buffer, "\x00")
				feedback := (strings.TrimSpace(string(feedback_buffer)))

				if feedback == "GOOD" || feedback == "BLOCK" {
					accepted++
				} else if feedback == "BAD" {
					rejected++
				} else if feedback == "INVU" {
					log.Fatal("Invalid username received in feedback")
				}
				break
			}
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]

	log.Println("GO miner v4.1 started... ")

	if len(argsWithoutProg) == 0 {
		log.Println("Enter your username:")
		fmt.Scan(&username)
		log.Println("How many goroutines do you want to start?")
		fmt.Scan(&x)
		log.Println("Select a difficulty, the possible values are LOW, MEDIUM, NET or EXTREME:")
		fmt.Scan(&diff)
		log.Println("Miner key:")
		fmt.Scan(&key)
	} else if len(argsWithoutProg) > 0 {
		// Passing command line interface's arguments.
		username = os.Args[1]
		x, _ = strconv.Atoi(os.Args[2])
		diff = os.Args[3]
		key = os.Args[4]
	}

	string_count := strconv.Itoa(x)

	log.Println("Username: " + username)
	log.Println("Goroutines count: " + string_count)
	log.Println("Difficulty: " + diff)

	for i := 0; i < x; i++ {
		go work()
		time.Sleep(100 * time.Millisecond)
	}

	for {
		log.Printf("Accepted shares: %d Rejected shares: %d\n", accepted, rejected)
		time.Sleep(10 * time.Second)
	}
}