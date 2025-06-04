package main

import "fmt"

// Mediator интерфейс определяет методы для взаимодействия между пользователями
type Mediator interface {
    SendMessage(message string, sender *User)
}

// ChatMediator — конкретный медиатор для чата
type ChatMediator struct {
    users []*User
}

// NewChatMediator создает новый экземпляр медиатора
func NewChatMediator() *ChatMediator {
    return &ChatMediator{users: make([]*User, 0)}
}

// AddUser добавляет пользователя в чат
func (m *ChatMediator) AddUser(user *User) {
    m.users = append(m.users, user)
}

// SendMessage рассылает сообщение всем пользователям, кроме отправителя
func (m *ChatMediator) SendMessage(message string, sender *User) {
    for _, user := range m.users {
        if user != sender {
            user.Receive(message)
        }
    }
}

// User представляет пользователя чата
type User struct {
    name    string
    mediator Mediator
}

// NewUser создает нового пользователя с именем и ссылкой на медиатор
func NewUser(name string, mediator Mediator) *User {
    return &User{name: name, mediator: mediator}
}

// Send отправляет сообщение через медиатор
func (u *User) Send(message string) {
    fmt.Printf("%s отправил: %s\n", u.name, message)
    u.mediator.SendMessage(message, u)
}

// Receive получает сообщение
func (u *User) Receive(message string) {
    fmt.Printf("%s получил: %s\n", u.name, message)
}

func main() {
    // Создаем медиатор
    mediator := NewChatMediator()

    // Создаем пользователей и добавляем их в чат
    alice := NewUser("Alice", mediator)
    bob := NewUser("Bob", mediator)
    charlie := NewUser("Charlie", mediator)

    mediator.AddUser(alice)
    mediator.AddUser(bob)
    mediator.AddUser(charlie)

    // Пользователи отправляют сообщения через медиатор
    alice.Send("Привет, всем!")
    bob.Send("Привет, Alice!")
}
