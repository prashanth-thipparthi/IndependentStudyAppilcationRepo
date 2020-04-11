//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	//	"bytes"
	"errors"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"
	"strconv"
	"gocv.io/x/gocv"

	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"

	"github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/almutawm/netConn"
)

const (
	serviceKey = "image-forwarder"
)

var counter int = 0
var conn_host string = ""
var conn_port string = ""
func main() {

	conn_host = os.Args[1]
	conn_port = os.Args[2]
	fmt.Println("connection  ", conn_host, ":", conn_port)

	// 1) First thing to do is to create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		edgexSdk.LoggingClient.Error(fmt.Sprintf("SDK initialization failed: %v\n", err))
		os.Exit(-1)
	}

	// 2) shows how to access the application's specific configuration settings.
	appSettings := edgexSdk.ApplicationSettings()
	if appSettings == nil {
		edgexSdk.LoggingClient.Error("No application settings found")
		os.Exit(-1)
	}

	valueDescriptorList, ok := appSettings["ValueDescriptors"]
	if !ok {
		edgexSdk.LoggingClient.Error("ValueDescriptors application setting not found")
		os.Exit(-1)
	}

	// 3) Since our FilterByValueDescriptor Function requires the list of ValueDescriptor's we would
	// like to search for, we'll go ahead create that list from the corresponding configuration setting.
	valueDescriptorList = strings.Replace(valueDescriptorList, " ", "", -1)
	valueDescriptors := strings.Split(valueDescriptorList, ",")
	edgexSdk.LoggingClient.Info(fmt.Sprintf("Filtering for %v value descriptors...", valueDescriptors))

	// 4) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(valueDescriptors).FilterByValueDescriptor,
		processImages,
	)

	// 5) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err := edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

func processImages(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	event, ok := params[0].(models.Event)
	if !ok {
		return false, errors.New("processImages didn't receive expect models.Event type")

	}

	for _, reading := range event.Readings {

		fileName := reading.Value
		fmt.Println("About to call ConnectToSend() to send file ", fileName)
		img := gocv.IMRead(fileName, gocv.IMReadColor )
		if img.Empty() {
			fmt.Println("Unable to read Image file")
			return false, nil
		} else {
			go faceDetection(img)
		}

	}

	return false, nil
}


func faceDetection(img gocv.Mat) {

	defer img.Close()

	xmlFile := "../../cascade/haarcascade_frontalface_alt.xml"

	// color for the rect when faces detected
	//blue := color.RGBA{0, 0, 255, 0}
	red := color.RGBA{255, 0, 0, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	// detect faces
	rects := classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	// draw a rectangle around each face on the original image,
	// along with text identifying as "Human"
	for _, r := range rects {
		gocv.Rectangle(&img, r, red, 2)
	}

	fileName := "/tmp/" + strconv.FormatInt(time.Now().Unix(),10)+"_fd_image.jpg"
	b := gocv.IMWrite(fileName, img)
	if (!b) {
		fmt.Println("Writing Mat to file failed")
		return
	}
	fmt.Println("Just saved " + fileName)
	netConn.ConnectToSend(conn_host, conn_port, "FILE0", fileName)
}
