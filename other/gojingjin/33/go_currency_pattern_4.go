package main

// func main() {
// 	done := spawnGroup(5, worker, 30)
// 	println("spawn a group of workers")

// 	timer := time.NewTimer(time.Second * 5)
// 	defer timer.Stop()

// 	select {
// 	case <-timer.C:
// 		println("wait group workers exit timeout!")
// 	case <-done:
// 		println("group workers done.")
// 	}
// }

// func foo() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	defer ticker.Stop()

// 	select {
// 	case <-ticker.C:
// 		println("tiker up.")
// 	}
// }
