package main;

import (
        "fmt"
        "io"
        "net"
        "os"
        "strconv"
        "strings"
//        "bufio"
)

import (
        //      "bytes"
//        "errors"
//        "fmt"
        "time"
        "image/color"
        "gocv.io/x/gocv"

//        "github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"

 //       "github.com/edgexfoundry/go-mod-core-contracts/models"

 //       "github.com/edgexfoundry/app-functions-sdk-go/appcontext"
 //       "github.com/edgexfoundry/app-functions-sdk-go/appsdk"
 //       "github.com/almutawm/netConn"
)

const (
        CONN_HOST = "0.0.0.0"
        CONN_PORT = "8181"
        CONN_TYPE = "tcp"
        BUFFERSIZE = 1024
)
func Check(e error, s string) {
        if e != nil {
                fmt.Println(s)
                panic(e)
        }
}

func faceDetection(img gocv.Mat) {

        defer img.Close()

        xmlFile := "haarcascade_frontalface_alt.xml"

        // color for the rect when faces detected
        //blue := color.RGBA{0, 0, 255, 0}
        red := color.RGBA{255, 0, 0, 0}

        // load classifier to recognize faces
        classifier := gocv.NewCascadeClassifier()
        defer classifier.Close()
        fmt.Printf("one")
        if !classifier.Load(xmlFile) {
                fmt.Printf("Error reading cascade file: %v\n", xmlFile)
                return
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
                return
        }
        fmt.Println("Just saved " + fileName)
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
func main() {
        var data string 
        conn, err := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
        Check(err, "Unable to create file")
        data = RecvFile(conn,"/home/pi/")
        fmt.Println("Received file: ",data)
        var fileName string
        fileName = data
        fmt.Println("About to call ConnectToSend() to send file ", fileName)
        img := gocv.IMRead(fileName, gocv.IMReadColor )
        if img.Empty() {
        fmt.Println("Unable to read Image file")
        //   return nil
        } else {
           fmt.Println("About to detect face")
           //go faceDetection(img)
           faceDetection(img) 
        }

}
