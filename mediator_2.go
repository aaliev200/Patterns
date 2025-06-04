package main

import (
	"fmt"
	"time"
)

// Интерфейсы компонентов системы
type OrderComponent interface {
	Notify(event string, data interface{})
	GetID() string
}

type OrderMediator interface {
	Register(component OrderComponent)
	Notify(sender OrderComponent, event string, data interface{})
}

// Конкретный посредник
type OrderProcessingMediator struct {
	components map[string]OrderComponent
}

func NewOrderProcessingMediator() *OrderProcessingMediator {
	return &OrderProcessingMediator{
		components: make(map[string]OrderComponent),
	}
}

func (m *OrderProcessingMediator) Register(component OrderComponent) {
	m.components[component.GetID()] = component
}

func (m *OrderProcessingMediator) Notify(sender OrderComponent, event string, data interface{}) {
	for id, component := range m.components {
		if id != sender.GetID() {
			component.Notify(event, data)
		}
	}
}

// Сущность Заказ
type Order struct {
	id       string
	status   string
	mediator OrderMediator
}

func NewOrder(id string, mediator OrderMediator) *Order {
	order := &Order{
		id:       id,
		status:   "pending",
		mediator: mediator,
	}
	mediator.Register(order)
	return order
}

func (o *Order) GetID() string {
	return o.id
}

func (o *Order) Notify(event string, data interface{}) {
	switch event {
	case "payment_confirmed":
		o.status = "processing"
		fmt.Printf("Order %s: Status changed to processing\n", o.id)
		o.mediator.Notify(o, "order_processing", o.id)
	case "inventory_reserved":
		o.status = "ready_to_ship"
		fmt.Printf("Order %s: Status changed to ready_to_ship\n", o.id)
	}
}

func (o *Order) Pay() {
	fmt.Printf("Order %s: Initiating payment\n", o.id)
	o.mediator.Notify(o, "payment_initiated", o.id)
}

func (o *Order) GetStatus() string {
	return o.status
}

// Сущность Платежная система
type PaymentSystem struct {
	id       string
	mediator OrderMediator
}

func NewPaymentSystem(id string, mediator OrderMediator) *PaymentSystem {
	payment := &PaymentSystem{
		id:       id,
		mediator: mediator,
	}
	mediator.Register(payment)
	return payment
}

func (p *PaymentSystem) GetID() string {
	return p.id
}

func (p *PaymentSystem) Notify(event string, data interface{}) {
	if event == "payment_initiated" {
		fmt.Printf("PaymentSystem %s: Processing payment for order %s\n", p.id, data.(string))
		time.Sleep(time.Second) // Симуляция обработки платежа
		p.mediator.Notify(p, "payment_confirmed", data)
	}
}

func (p *PaymentSystem) ProcessPayment(orderID string) {
	fmt.Printf("PaymentSystem %s: Manually processing payment for %s\n", p.id, orderID)
	p.mediator.Notify(p, "payment_initiated", orderID)
}

// Сущность Склад
type Inventory struct {
	id       string
	mediator OrderMediator
}

func NewInventory(id string, mediator OrderMediator) *Inventory {
	inventory := &Inventory{
		id:       id,
		mediator: mediator,
	}
	mediator.Register(inventory)
	return inventory
}

func (i *Inventory) GetID() string {
	return i.id
}

func (i *Inventory) Notify(event string, data interface{}) {
	if event == "order_processing" {
		fmt.Printf("Inventory %s: Reserving items for order %s\n", i.id, data.(string))
		time.Sleep(500 * time.Millisecond) // Симуляция резервирования
		i.mediator.Notify(i, "inventory_reserved", data)
	}
}

func (i *Inventory) CheckAvailability(orderID string) bool {
	fmt.Printf("Inventory %s: Checking availability for order %s\n", i.id, orderID)
	return true // Просто для примера
}

func main() {
	// Создаем посредника
	mediator := NewOrderProcessingMediator()

	// Инициализируем компоненты
	order := NewOrder("ORD001", mediator)
	paymentSystem := NewPaymentSystem("PAY001", mediator)
	inventory := NewInventory("INV001", mediator)

	// Используем компоненты явно
	fmt.Println("Initial order status:", order.GetStatus())

	if inventory.CheckAvailability("ORD001") {
		order.Pay()
	}

	// Даем время на обработку всех событий
	time.Sleep(2 * time.Second)

	fmt.Println("Final order status:", order.GetStatus())
	paymentSystem.ProcessPayment("ORD001") // Дополнительная демонстрация использования

	time.Sleep(2 * time.Second)
}
