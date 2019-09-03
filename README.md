JSON vs Protocal Buffer

    Parsing JSON is CPU intensive and JSON takes more memory
    Parsig PRotocal Buffers(Binary Format) is les CPU intensive because it's closer to machine represents data
    Protocol Buffers allows faster and more efficient communication

    gRPC leverages HTTP/2 as a backbone of communications/uses only HTTP/2
    HTTP/2 is comparatively very fast than HTTP/1.1
    Eg: https://imagekit.io/demo/http2-vs-http1

    gRPC servers are asynchronous by default
    Types:
    - Unary
    - Server Streaming
    - Client Streaming
    - Bi Directional Streaming

HTTP/2 vs HTTP/1.1

    HTTP/2 is the newer standard for internet communications that addresses common pitfall of HTTP/1.1
    on modern webpages

    HTTP/1.1:
    - Released in 1997
    - Opens a new TCP connection to server at each request
    - Does not compress the Headers(which are plain text - heavy size)
    - Only works with Request/Response mechanism

    HTTP/2:
    - Released in 2015
    - Developed by google under the name SPDY
    - Supports multiplexing that means client and server can push messages in parallel over same TCP connection
    - Reduces latency
    - Server can push miltiple messages with single request from client / saves round trips
    - Supports compressed Headers
    - HTTP/2 is binary -> an extremely great match for Protocol Buffers
    - Secure -> SSL is not required but recommended by default)

Installations Required:

    Proctoc Setup:
    Mac: brew install protobuf
    Download the windows archive: https://github.com/google/protobuf/releases

    go gRPC: go get -u google.golang.org/grpc
    go Protocol Buffer: go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
    Mac: Add them to bash_profile
        - vim ~/.bash_profile
        - Add: 
            export GO_PATH=~/go
            export PATH=$PATH:/$GO_PATH/bin
        - source ~/.bash_profile

    if there are issues with Protocol Buffers:
    - Download the protobuf-2.4.1 from https://protobuf.googlecode.com/files/protobuf-2.4.1.tar.gz
    - Extract the tar.gz file.
    - $cd ~/Downloads/protobuf-2.4.1
    - $./configure
    - $make
    - $make check
    - $sudo make install
    - $which protoc
    - $protoc --version

Running the samples:

    Greeting:
    - To generate the greet.pb.go: "protoc greet/greetpb/greet.proto --go_out=plugins=grpc:."
    - Running the server: go run greet/greet_server/server.go
    - Running the client: go run greet/greet_client/client.go

    Calculator:
    - To genrate the calculator.pb.go: "protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:."
    - Running the server: go run calculator/calculator_server/server.go
    - Running the client: go run calculator/calculator_client/client.go
