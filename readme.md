just an experiment forward some data from other source (eg: topic, queue, stream, etc) use case using SSE. for this example i use grpc stream.
was inspired from EventSource repo in here https://github.com/EventSource/eventsource/tree/master then recreate in go. so run the server-grpc
first then grpc-client after that (very sorry for inconsistent naming). please notice SSE is 1 direction. you can find more information about SSE at https://developer.mozilla.org/en-US/docs/Web/API/EventSource. off course you can optimize in a lot of aspect on the forwarder (grpc-client) so it will has more consistency.

## command

simply run

    cd server-grpc
    go run .
    cd ../grpc-client
    go run .
    open http://localhost:8000      # Browser client 

![Screenshot](/browserclient.png)

some note to my self for my skill issue (i hate protoc) :
protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative *.proto



