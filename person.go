package main

import "fmt"

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p Person) String() string {
	return fmt.Sprintf("id: %d, name: %s, email: %s", p.ID, p.Name, p.Email)
}
