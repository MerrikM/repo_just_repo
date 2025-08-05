package main

import "fmt"

func main() {
	action := Action{
		Human: Human{Name: "Алексей", Age: 22},
		Skill: "Встраивание структур",
	}

	action.Greetings()

	action.SomeAction()
}

type Human struct {
	Name string
	Age  int
}

type Action struct {
	Human
	Skill string
}

func (h *Human) Greetings() {
	fmt.Printf("Привет, меня зовут %s и мне %d года!\n", h.Name, h.Age)
}

func (a *Action) SomeAction() {
	fmt.Printf("%s делает L1.1: %s\n", a.Name, a.Skill)
}
