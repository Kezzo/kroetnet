# kroetnet
Game state and position synchronization library

## Server
To run the server use `go run *.go`

## Docker

Run your build docker image with `docker run -p 2448:2448/udp -it kroetnet`

# Profiling the Match Server

Run the following command to get the profiles:

- `DEBUG=true CC=clang CXX=clang++ go run !(*_test).go -cpuprofile "cpu.prof" -memprofile "mem.prof"`

Observe the results with `go tool pprof`

