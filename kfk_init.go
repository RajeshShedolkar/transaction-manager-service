package main

import (
	"fmt"
	"transaction-manager/kfkmig"
)

func main(){
	fmt.Println("")
	kfkmig.BootstrapKafkaTopics()
}