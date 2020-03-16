package netConn

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"bufio"
)


const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8181"
	CONN_TYPE = "tcp"
	BUFFERSIZE = 1024
)

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}

func Check(e error, s string) {
	if e != nil {
		fmt.Println(s)
		panic(e)
	}
}

func RunTCPServer(TCPRequesthandler func(c net.Conn, n int)) {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	Check(err, "Unable to listen on port ...")
	i:=0
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		i+=1
		// Listen for an incoming connection.
		conn, err := l.Accept()
		Check(err, "Error when trying to accept")
		fmt.Println("Accepted connection ", i)
		// Handle connections
		go TCPRequesthandler(conn, i)
	}
}

func RequestHSendFiles(conn net.Conn, i int) {
	defer conn.Close()

	if (i%2 == 0) {
		SendFile(conn, "2.jpg")
	} else {
		SendFile(conn, "1.jpg")
	}
	return
}

func RequestHText(conn net.Conn, i int) {
	fmt.Println("Inside handling Request: ", i)
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_ , err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("received: ", string(buf))
	// Send a response back to person contacting us.
	conn.Write([]byte("hi your number is" + strconv.Itoa(i)))
	// Close the connection when you're done with it.
	conn.Close()
}

func ConnectToSend(conn_host string, conn_port string, typeOfData string, data string) {

	conn, err := net.Dial(CONN_TYPE, conn_host+":"+conn_port)

	Check(err, "Unable to to establish a connection")

	defer conn.Close()

	fmt.Println("Connected to servere ")

	if typeOfData == "FILE0" {
		SendFile(conn, data)
	} else {
		if len(typeOfData) == 5 {
			SendText(conn, typeOfData, data)
		} else {
			fmt.Println("Nothing was sent, typeOfData must be 5 characters long")
		}
	}

	fmt.Println("Closing connection!")
}

func SendText(conn net.Conn, typeOfData string, data string) {
	defer conn.Close()
	conn.Write([]byte(typeOfData))
	conn.Write([]byte(data+";;;;"))
}

func SendFile(conn net.Conn, filename string) {

	defer conn.Close()
	file, err := os.Open(filename)

	Check(err, "Unable to open file, exiting")

	fileInfo, err := file.Stat()
	Check(err, "Unable to get file Stat, exiting")

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 32)

	conn.Write([]byte("FILE0"))
	conn.Write([]byte(fileSize))
	conn.Write([]byte(fileName))

	sendBuffer := make([]byte, BUFFERSIZE)
	fmt.Println("Start sending file!")
	nBytes := 0
	for {
		_, err = file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		n, _ := conn.Write(sendBuffer)
		nBytes += n
	}
	fmt.Println("File: ", fileName ," has been sent, file size: ", fileSize , ", number of bytes sent: ", nBytes)
}

func RecvText(conn net.Conn,  typeOfData string) string {
	defer conn.Close()
	fmt.Println("About to receive ", typeOfData)
	buffer := make([]byte, 64)
	conn.Read(buffer)
	data := strings.Split(string(buffer), ":")[0]
	return data
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
	fileName := path + strings.Trim(string(bufferFileName), ":")

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

func ConnectToRecv(conn_host string, conn_port string, typeOfData string, data string) string {

	conn, err := net.Dial(CONN_TYPE, conn_host+":"+conn_port)

	Check(err, "Unable to to establish a connection")

	returnData := ""
	if (typeOfData == "FILE") {
		returnData = RecvFile(conn, data) // data is path where to save the received file
	} else {
		returnData = RecvText(conn, typeOfData)
	}
	return returnData
}

func ConnectSendUserInput(conn_host string, conn_port string) {

	conn, err := net.Dial(CONN_TYPE, conn_host+":"+conn_port)

	Check(err, "Unable to to establish a connection")

	defer conn.Close()

	fmt.Println("Connected to servere")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	text, _ := reader.ReadString('\n')
	// send to socket
	fmt.Fprintf(conn, text + "\n")
	// listen for reply
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println("Message from server: ", message )
	conn.Close()
}
