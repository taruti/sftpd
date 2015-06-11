package main

func main() {
	go RunServer("127.0.0.1:2023", synthetic{})
	go RunServer("127.0.0.1:2024", rfs{})
	<-make(chan int)
}
