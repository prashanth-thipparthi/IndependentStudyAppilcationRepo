package main;

import (
        "fmt"
        "time"
        "net" 
        "strings"
)

const (
//        S_CONN_HOST = "192.168.43.48"
        FACE_DETECT_HOST="0.0.0.0"
        FACE_DETECT_PORT="8180"
        S_CONN_PORT = "8181"
        S_CONN_TYPE = "tcp"
)

func Check(e error, s string) {
        if e != nil {
                fmt.Println(s)
                panic(e)
        }
}

func SendText(conn net.Conn, typeOfData string, data string) {
  //      defer conn.Close()
//        conn.Write([]byte(typeOfData))
        conn.Write([]byte(data+":::"))
}

func RecvText(conn net.Conn, typeOfData string) string {
        //defer conn.Close()
        fmt.Println("About to receive ", typeOfData)
        buffer := make([]byte, 64)
        conn.Read(buffer)
        data := strings.Split(string(buffer), ":::")[0]
        //data := string(buffer)
        fmt.Println("data:"+data)
        return data
}

func main() {

    data := "/tmp/data/1587012662_fd_image.jpg" 
    fmt.Println("connecting to:"+FACE_DETECT_HOST)
    time.Sleep(100000)
    conn, err := net.Dial(S_CONN_TYPE, FACE_DETECT_HOST+":"+FACE_DETECT_PORT)
    defer conn.Close()
    Check(err, "Unable to connect to server")
    SendText(conn, "string", data)
    data = RecvText(conn, "string")
    fmt.Println("Received data: ",data)
//    return data
/*
    filename := "/tmp/data/1587012662_fd_image.jpg"
    _, err := os.Open(filename)
    if err != nil {
        fmt.Println("cannot open file")
        panic(err)
    }*/
}
