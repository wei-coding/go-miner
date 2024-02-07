# Go-miner

A **[duino-coin](https://duinocoin.com/)** miner made in golang.

[![Go](https://img.icons8.com/color/48/000000/golang.png)](https://golang.org/)
*Check out [Go](https://golang.org/)*.
****
### Arguments:
* **Username** -> User to mine for.
* **Goroutines** -> Amount of goroutines to run in the background (can be thought of as threads).
* **Difficulty** -> NET, LOW, MEDIUM or HIGH mining difficulty.
* **Miner key** -> Your wallet miner key if you set.

Learn more about [goroutines (threads)](https://gobyexample.com/goroutines).

**You can use the miner with a command line interface:**
```bash
./miner <username (string)> <goroutines (integer)> <difficulty (NET, LOW, MEDIUM, HIGH)> <miner key>
```

****
### Todo:
* Add cache for storing user's credentials and execute without asking for them.
