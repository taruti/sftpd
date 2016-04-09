package main

func main() {
	go RunServerHighLevel("127.0.0.1:2023", synthetic{})
	go RunServerLowLevel("127.0.0.1:2024", readOnlyDirFs{})
	<-make(chan int)
}
