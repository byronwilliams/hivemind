package main

import (
    "fmt"
    "net"
    "os"
    "time"
    "strconv"
    "log"

    // "io"
    "net/http"
    "golang.org/x/net/websocket"
)

/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}

type Temp struct {
    id string
    temp float32
}
type Temps map[string]float32

type appContext struct {
    temps Temps
    mychan chan *Temp
}

func main() {
    temps := Temps{}

    context := &appContext{temps: temps} // Simplified for this example

    go listenForTemps(context)
    pushTemps(context)
}

func listenForTemps(cxt *appContext) {
    /* Lets prepare a address at any address at port 10001*/
    ServerAddr,err := net.ResolveUDPAddr("udp",":32210")
    CheckError(err)

    /* Now listen at selected port */
    ServerConn, err := net.ListenUDP("udp", ServerAddr)
    CheckError(err)
    defer ServerConn.Close()

    cxt.mychan = make(chan *Temp)

    buf := make([]byte, 1024)

    for {
        // n,addr,err := ServerConn.ReadFromUDP(buf)
        _,_,err := ServerConn.ReadFromUDP(buf)
        // fmt.Println("Received ",string(buf[0:n]), " from ",addr)

        if err != nil {
            log.Println("Error: ",err)
        }

        theId := string(buf[14:30])
        theTemp, err := strconv.ParseFloat(string(buf[33:38]), 32)
        cxt.temps[theId] = float32(theTemp)

        log.Println(theId + ":" + strconv.FormatFloat(theTemp, 'f', 2, 32))

        cxt.mychan <- &Temp{theId, float32(theTemp)}

        if err != nil {
            log.Println("Error: ",err)
        }
    }
}

func Echo(cxt *appContext) websocket.Handler {
    return func(ws *websocket.Conn) {
        var err error

        for {
            t := <- cxt.mychan
            sout := t.id + ":" + strconv.FormatFloat(float64(t.temp), 'f', 2, 32)

            if err = websocket.Message.Send(ws, sout); err != nil {
                fmt.Println("Can't send")
            }
            time.Sleep(1 * time.Second)
        }
    }
}

func pushTemps(cxt *appContext) {
    http.Handle("/", Echo(cxt))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
