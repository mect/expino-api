package main

func sendUpdate() {
	io.BroadcastTo("update", "update")
}
