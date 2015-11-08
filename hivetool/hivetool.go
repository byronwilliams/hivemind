package main

import (
        "log"
        "bufio"
        "fmt"
        "time"
        "strconv"
        "strings"
        "net"

        "github.com/tarm/serial"
        "github.com/syndtr/goleveldb/leveldb"
)

func main() {
        db, err := leveldb.OpenFile("../hivetool.db", nil)
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

            if k == "cfg" && v == "time_unset" {
                t := "T" + strconv.FormatInt(time.Now().Unix(), 10) + "\n"
                log.Println(t)
                _, err = s.Write([]byte(t))
                if err != nil {
                        log.Fatal(err)
                }
            } else if k == "log" {
                log.Println(v)
            } else {
                if(k != "elapsed") {
                    go writeToDb(db, k, v)
                    go writeToMind(k, v)
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

func writeToMind(k string, v string) {
    ServerAddr,err := net.ResolveUDPAddr("udp","c27f1925-8697-4d56-b83f-d8afd6a7192a.pub.cloud.scaleway.com:32210")
    CheckError(err)

    LocalAddr, err := net.ResolveUDPAddr("udp", "0.0.0.0:0")
    CheckError(err)

    Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
    CheckError(err)

    buf := []byte(k + " - " + v)
    _,err = Conn.Write(buf)

    Conn.Close()
}

func parse(s string) (string, string) {
    parts := strings.Split(strings.TrimSpace(s), ":")
    timestamp := parts[0]
    csv := parts[1]
    csvparts := strings.Split(csv, ",")

    if timestamp == "cfg" || timestamp == "elapsed" || timestamp == "log" {
        return timestamp, csvparts[0]
    }


    return timestamp + csvparts[0] + csvparts[1], csvparts[2]
}

func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
