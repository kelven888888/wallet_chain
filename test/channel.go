package test

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Person struct {
	Name    string
	Age     int
	Address string
}

func WriteChan(PersonChan chan Person) {
	for i := 0; i < 10; i++ {
		PersonChan <- Person{
			Name:    "Person" + strconv.Itoa(rand.Intn(100)),
			Age:     rand.Intn(100),
			Address: "成都市成华大道" + strconv.Itoa(rand.Intn(100)) + "号",
		}
		fmt.Println("write data=", i)
	}
	close(PersonChan)
}
func ReadChan(PersonChan chan Person, exitChan chan bool) {
	i := 0
	for {
		i++
		v, ok := <-PersonChan
		if !ok {
			break
		}
		fmt.Printf("第%v个perosn数据的Name=%v,Age=%v,Address=%v\n", i, v.Name, v.Age, v.Address)
	}
	exitChan <- true
	close(exitChan)

}

func main() {

	PersonChan := make(chan Person, 100)
	exitChan := make(chan bool, 1)
	go WriteChan(PersonChan)
	go ReadChan(PersonChan, exitChan)
	for {
		_, ok := <-exitChan
		if !ok {
			break
		}
	}
}
