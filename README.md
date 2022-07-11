### go-delayqueue client  

A client of the [go-delayqueue](https://github.com/raymondmars/go-delayqueue) for golang project.   

### Features to do 
- [x] Push message to delayqueue.   
- [x] Delete a message from delayqueue.   
- [x] Send Ping test message. 
- [ ] Upate a message in the delayqueue. 

### How to use  
```go 
import (
	delayclient "godelayqueueclient"
)
testHost := "127.0.0.1"    
fromNowDelaySeconds := 600   

// notify by HTTP  
client := delayclient.NewClientWithHostAndPort(testHost, "3450")
resp, err := client.Push(fromNowDelaySeconds, 
                          delayclient.HTTP, 
                          fmt.Sprintf("https://google.com", testHost), 
                          "hello,world")
if err != nil {
	t.Log(err)
} else {
	t.Log(resp)
  t.Log(resp.TaskId)
}

// notify by third queue (such as rabbitmq)  
var PubSubQueueName = "pubsub-queue"
resp, err := client.Push(fromNowDelaySeconds, delayclient.SubPub, PubSubQueueName, "hello,world")
if err != nil {
	t.Log(err)
} else {
	t.Log(resp)
	t.Log(resp.TaskId)
}
```   
### Contributing  
Anyone is welcome to submit pull requests and suggestions, issues.   

### License  
See [LICENSE](./LICENSE)

