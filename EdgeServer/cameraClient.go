package main;

import (
        "fmt"
        "io"
        "net"
        "os"
        "strconv"
        "strings"
)

import (
        //      "bytes"
//        "errors"
        "time"
        "image/color"
        "gocv.io/x/gocv"
)

const (
//        S_CONN_HOST = "192.168.43.48"
        S_CONN_PORT = "8181"
        S_CONN_TYPE = "tcp"
        CONN_HOST = "0.0.0.0"
        CONN_PORT = "8180"
        CONN_TYPE = "tcp"
        BUFFERSIZE = 1024
)

//var ip = []string{"192.168.43.48",""}

func Check(e error, s string) {
        if e != nil {
                fmt.Println(s)
                panic(e)
        }
}

func faceDetection(img gocv.Mat)(string) {

        defer img.Close()

        xmlFile := "/IndependentStudyAppilcationRepo/ProfCode/files/haarcascade_frontalface_alt.xml"

        // color for the rect when faces detected
        //blue := color.RGBA{0, 0, 255, 0}
        red := color.RGBA{255, 0, 0, 0}

        // load classifier to recognize faces
        classifier := gocv.NewCascadeClassifier()
        defer classifier.Close()
        fmt.Printf("one")
        if !classifier.Load(xmlFile) {
                fmt.Printf("Error reading cascade file: %v\n", xmlFile)
                return ""
        }
        fmt.Printf("two")
        // detect faces
        rects := classifier.DetectMultiScale(img)
        fmt.Printf("found %d faces\n", len(rects))
        fmt.Printf("three")

        // draw a rectangle around each face on the original image,
        // along with text identifying as "Human"
        for _, r := range rects {
                gocv.Rectangle(&img, r, red, 2)
        }
        fmt.Printf("four")

        fileName := "/tmp/" + strconv.FormatInt(time.Now().Unix(),10)+"_fd_image.jpg"
        b := gocv.IMWrite(fileName, img)
        if (!b) {
                fmt.Println("Writing Mat to file failed")
                return ""
        }
        fmt.Println("Just saved " + fileName)
        return fileName
 //       netConn.ConnectToSend(conn_host, conn_port, "FILE0", fileName)
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
func SendFile(conn net.Conn, filename string) {

	defer conn.Close()
	file, err := os.Open(filename)

	Check(err, "Unable to open file, exiting")

	fileInfo, err := file.Stat()
	Check(err, "Unable to get file Stat, exiting")

	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 32)

	//conn.Write([]byte("FILE0"))
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
	fmt.Println("File: ", fileName ," has been sent, file size: ", fileSize, ", number of bytes sent: ", nBytes)
}

func RecvText(conn net.Conn,  typeOfData string) string {
	//defer conn.Close()
	fmt.Println("About to receive ", typeOfData)
	buffer := make([]byte, 64)
	conn.Read(buffer)
	//data := strings.Split(string(buffer), ":")[0]
	data := string(buffer)
        fmt.Println("data:"+data)
        return data
}

func handleRequest(clientCon net.Conn) {
    var msg string
    var processedFileName string
    processedFileName = "" 
    msg = RecvText(clientCon,  "string")
    options := strings.Split(msg, ",")
    //rasperryIpIndex, err := strconv.Atoi(options[0]) // if we want to  assign rpi's indices and use their index we can uncomment this line
    //Check(err, "Unable to convert string to integer")
    //fileName := getImageFromRaspberrypi(ip[rasperryIpIndex])
    fileName := getImageFromRaspberrypi(strings.TrimSpace(options[0]))
    img := gocv.IMRead(fileName, gocv.IMReadColor )
    if img.Empty() {
       fmt.Println("Unable to read Image file")
    } else {
       fmt.Println("processing the image option:"+options[1])
       //go faceDetection(img)
      options[1] = "facedetection" 
       switch options[1]{
          case "facedetection":
              processedFileName = faceDetection(img)
          case "faceblur":
              processedFileName = faceDetection(img)        
          default:
              fmt.Println("Invalid option")       
      }
    } 
    SendFile(clientCon, processedFileName)
}

func getImageFromRaspberrypi(S_CONN_HOST string) string {
        var data string
        fmt.Println("connecting to:"+S_CONN_HOST) 
        conn, err := net.Dial(S_CONN_TYPE, S_CONN_HOST+":"+S_CONN_PORT)
        Check(err, "Unable to connect to server")
        home, err := os.UserHomeDir()
        Check(err, "Unable to get home directory")
        data = RecvFile(conn,home+"/")
        fmt.Println("Received file: ",data)
        return data 
}

func main() {
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
                defer conn.Close()
                Check(err, "Error when trying to accept")
                fmt.Println("Accepted connection ", i)
                // Handle connections
                handleRequest(conn)
        }

}
