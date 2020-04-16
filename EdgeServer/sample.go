package main;

import (
        "fmt"
        "os"
)
func main() {
    filename := "/tmp/data/1587012662_fd_image.jpg"
    _, err := os.Open(filename)
    if err != nil {
        fmt.Println("cannot open file")
        panic(err)
    }
}
