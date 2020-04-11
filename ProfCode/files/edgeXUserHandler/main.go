package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	//	"io"
	"strconv"
	"strings"
	"github.com/almutawm/netConn"
)

var pathToApps string = "/home/mm/go/src/github.com/edgexfoundry-holding/app-service-examples/app-services/"
var clientIP string = ""
var fileName string = ""
var REMOTEPORT string = "8877"
func errCheck(e error, s string) bool {
	if e != nil {
		fmt.Println(s)
		return true
	} else {
		return false
	}
}
func main() {
	netConn.RunTCPServer(handleRequestCmd)
}

// Handles incoming requests.
func handleRequestCmd(conn net.Conn, i int) {
	clientIP = strings.Split(conn.RemoteAddr().String(), ":")[0]
	fmt.Println("client ip is ", clientIP)
	// Make a buffer to hold incoming data.
	cmdbuf := make([]byte, 5)
	_ , err := conn.Read(cmdbuf)

	if errCheck(err, "Problem getting command") {return}

	cmdStr := string(cmdbuf)

	fmt.Println("Received command ", cmdStr)
	go runCommand(cmdStr, conn)


}

func runCommand(command string, conn net.Conn) {
	switch command {
	case "CMD00":
		cmd := exec.Command("./app-service", clientIP, REMOTEPORT)
		cmd.Dir = pathToApps + "image-forwarder"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if errCheck(err, "problem running" + command) {return}
		pid := cmd.Process.Pid
		fmt.Println("Pid of the process is is %d", pid)
		pidStr := strconv.FormatInt(int64(pid),10)
		fmt.Println("Pid of the process is (string) ", pidStr)
		netConn.SendText(conn, "PID00", pidStr)
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD01":
		cmd := exec.Command("./app-service", clientIP, REMOTEPORT)
		cmd.Dir = pathToApps + "face-detection"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if errCheck(err, "problem running" + command) {return}
		pid := cmd.Process.Pid
		fmt.Println("Pid of the process is is ", pid)
		netConn.SendText(conn, "PID00", strconv.FormatInt(int64(pid),10))
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD02":
		cmd := exec.Command("./app-service", clientIP, REMOTEPORT)
		cmd.Dir = pathToApps + "face-blur"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if errCheck(err, "problem running" + command) {return}
		pid := cmd.Process.Pid
		fmt.Println("Pid of the process is is ", pid)
		netConn.SendText(conn, "PID00", strconv.FormatInt(int64(pid),10))
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD03":
		cmd := exec.Command("./app-service", clientIP, REMOTEPORT)
		cmd.Dir = pathToApps + "face-blur"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if errCheck(err, "problem running" + command) {return}
		pid := cmd.Process.Pid
		fmt.Println("Pid of the process is is ", pid)
		netConn.SendText(conn, "PID00", strconv.FormatInt(int64(pid),10))
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD04":
		cmd := exec.Command("./app-service", clientIP, REMOTEPORT)
		cmd.Dir = pathToApps + "face-blur"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if errCheck(err, "problem running" + command) {return}
		pid := cmd.Process.Pid
		fmt.Println("Pid of the process is is ", pid)
		netConn.SendText(conn, "PID00", strconv.FormatInt(int64(pid),10))
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD09":
		cmd := exec.Command("ls", "-lah")
		cmd.Dir = "/home/mm/go/src/github.com"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if errCheck(err, "problem running" + command) {return}
		conn.Close()
		fmt.Println("Finished from command ")

	case "CMD11": //send to edgeX coreCommand to fetch new image
		cmd := exec.Command("curl", "http://localhost:48082/api/v1/device/e1c51b0a-2cee-4336-a910-b4bf0c83a4c6/command/a0379de6-8134-4dd8-b8a0-40f70e56d90b")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if errCheck(err, "problem running" + command) {return}
		netConn.SendText(conn, "DONE0", "fetch image")
		conn.Close()
		fmt.Println("Finished from command ", command)

	case "CMD22":
		appPidStr := netConn.RecvText(conn, "PID")
		appPid, _ := strconv.ParseInt(strings.Trim(appPidStr, ":"), 10, 0)
		p, err := os.FindProcess(int(appPid))
		if errCheck(err, "problem finding the process, inside " + command) {return}
		err = p.Kill()
		if errCheck(err, "problem killing the process, inside" + command) {return}
		netConn.SendText(conn, "DONE1", "App Service stopped")
		conn.Close()
		fmt.Println("Finished from process %d", appPid)
	}
}
