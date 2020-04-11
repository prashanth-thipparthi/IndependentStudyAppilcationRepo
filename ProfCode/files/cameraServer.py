#!/usr/bin/env python

import os
import SocketServer
from picamera import PiCamera
from time import sleep
from time import time

class MyTCPHandler(SocketServer.BaseRequestHandler):
    """
    The request handler class for our server.
    It is instantiated once per connection to the server, and must
    override the handle() method to implement communication to the
    client.
    """

    def handle(self):
        # self.request is the TCP socket connected to the client
        print "inside handle ...\n"
        filename ='img' + str(time()) + '.jpg'
        fullfilename = '/run/user/1000/' + filename
        camera = PiCamera()
	try:
    		camera.capture(fullfilename)
	finally:
    		camera.close()

        filesize_list = list(str(os.path.getsize(fullfilename)))
        buf = "::::::::::"
        buf_list = list(buf)
        for i in range(len(filesize_list)):
            buf_list[i] = filesize_list[i]

        fs = ""
        fileSize = fs.join(buf_list)

	filename_list = list(filename)
        buf = "::::::::::::::::::::::::::::::::"
        buf_list = list(buf)
        for i in range(len(filename_list)):
            buf_list[i] = filename_list[i]

        fn = ""
        fileName = fn.join(buf_list)

        print "Done taking Picture filename is " + fileName + ", filesize is" + fileSize
        self.request.send(fileSize)
	self.request.send(fileName)
        f = open(fullfilename,'rb')
        l = f.read(1024)
        while (l):
            self.request.send(l)
            l = f.read(1024)

        f.close()
	print  "Picture sent"

if __name__ == "__main__":
    HOST, PORT = "0.0.0.0", 8181

    # Create the server, binding to localhost on port 9999
    server = SocketServer.TCPServer((HOST, PORT), MyTCPHandler)

    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    print "Listening on " + HOST + ":" + str(PORT)
    server.serve_forever()
