package main;

import (
        "fmt"
//        "io"
        "net"
//        "os"
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
        CONN_HOST = "0.0.0.0"
        CONN_PORT = "8179"
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

        xmlFile := "/home/tnr/IndependentStudyAppilcationRepo/EdgeServer/face_detection/haarcascade_frontalface_alt.xml"

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

func SendText(conn net.Conn, data string) {
 //       defer conn.Close()
//        conn.Write([]byte(typeOfData))
        fmt.Println("sent text"+data)
        conn.Write([]byte(data+":::"))
}

func RecvText(conn net.Conn,  typeOfData string) string {
	//defer conn.Close()
	fmt.Println("About to receive ", typeOfData)
	buffer := make([]byte, 64)
	conn.Read(buffer)
	data := strings.Split(string(buffer), ":::")[0]
	//data := string(buffer)
        fmt.Println("data:"+data)
        return data
}

func handleRequest(clientCon net.Conn) {
    var msg string
    var processedFileName string
    processedFileName = "" 
    msg = RecvText(clientCon,  "string")
    fileName := strings.TrimSpace(msg)
    img := gocv.IMRead(fileName, gocv.IMReadColor)
    if img.Empty() {
       fmt.Println("Unable to read Image file")
    }
    processedFileName = faceDetection(img)
    SendText(clientCon, processedFileName)
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
