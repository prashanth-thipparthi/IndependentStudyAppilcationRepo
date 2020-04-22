package main;

import (
        "fmt"
        "io"
        "net"
        "os"
        "strconv"
        "strings"
//        "time"
//        "bufio"
)

const (
        CONN_HOST = "192.168.43.81"
        CONN_PORT = "8179"
        CONN_TYPE = "tcp"
        BUFFERSIZE = 1024
)
func Check(e error, s string) {
        if e != nil {
                fmt.Println(s)
                panic(e)
        }
}
func RecvFile(conn net.Conn, path string) string {

	defer conn.Close()

	fmt.Println("Connected to server, start receiving the file name and file size")
	bufferFileName := make([]byte, 32)
	bufferFileSize := make([]byte, 10)

	conn.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Println("the file size is ", fileSize)
	conn.Read(bufferFileName)
    fmt.Println("path received:"+string(bufferFileName))
	fileName := path +strings.Trim(string(bufferFileName), ":")
    fmt.Println("fileName:"+string(bufferFileName))
	//newFile, err := os.Create("img_" + strconv.Itoa(i) + ".jpg")
	newFile, err := os.Create(fileName)

	Check(err, "Unable to create file")
	defer newFile.Close()
	var receivedBytes int64 = 0
	var wr int64 = 0

	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			wr, _ = io.CopyN(newFile, conn, (fileSize - receivedBytes))
			receivedBytes += wr
			conn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		wr, _ = io.CopyN(newFile, conn, BUFFERSIZE)
		receivedBytes += wr
	}
	fmt.Println("Received file: ", fileName,", bytes: ", receivedBytes)
	return fileName 
}

func SendText(conn net.Conn, typeOfData string, data string) {
        
//        conn.Write([]byte(typeOfData))
        conn.Write([]byte(data))
}

func main(){
    var data string
    conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    Check(err, "Unable to create file")
    defer conn.Close()
    SendText(conn,"text","192.168.43.48,facedetection")
//    home, err := os.UserHomeDir()
//    Check(err, "Unable to get home directory")
    data = RecvFile(conn,"/home/tnr/")
    fmt.Println("Received file: ",data)
}
