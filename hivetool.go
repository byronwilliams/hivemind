package main

import (
        "log"
        "bufio"
        // "fmt"
        "time"
        "strconv"
        "strings"

        "github.com/tarm/serial"
        "github.com/syndtr/goleveldb/leveldb"
)

func main() {
        db, err := leveldb.OpenFile("hivetool.db", nil)
        defer db.Close()

        c := &serial.Config{Name: "/dev/ttyACM0", Baud: 9600}
        s, err := serial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }

        reader := bufio.NewReader(s)

        for {
            reply, err := reader.ReadBytes('\n')
            if err != nil {
                panic(err)
            }

            k, v := parse(string(reply))

            if k == "cfg" {
                t := "T" + strconv.FormatInt(time.Now().Unix(), 10) + "\n"
                log.Println(t)
                _, err = s.Write([]byte(t))
                if err != nil {
                        log.Fatal(err)
                }
            } else {
                if(k != "elapsed") {
                    go writeToDb(db, k, v)
                } else {
                    log.Println(k + ":" + v)
                }
            }
        }
}

func writeToDb(db *leveldb.DB, k string, v string) {
    log.Println(k + " - " + v)
    db.Put([]byte(k), []byte(v), nil)
}

func parse(s string) (string, string) {
    parts := strings.Split(s, ":")
    timestamp := parts[0]
    csv := parts[1]
    csvparts := strings.Split(csv, ",")

    if timestamp == "cfg" || timestamp == "elapsed" {
        return timestamp, csvparts[0]
    }

    return timestamp + csvparts[0] + csvparts[1], csvparts[2]
}
