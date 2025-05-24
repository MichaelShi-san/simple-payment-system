package main

import (
	"fmt"
	"sync"
)

type User struct {
	ID      int
	Name    string
	Balance float64
	mu      sync.Mutex
}

type Transaction struct {
	FromID int
	ToID   int
	Amount float64
}

type PaymentSystem struct {
	Users        []*User
	Transactions []Transaction
}

func AddUser(ps *PaymentSystem, user *User) {
	ps.Users = append(ps.Users, user)
}

func AddTransaction(ps *PaymentSystem, transaction Transaction) {
	ps.Transactions = append(ps.Transactions, transaction)
}

func Deposit(user *User, amount float64) {
	user.mu.Lock()
	defer user.mu.Unlock()
	
	if amount > 0 {
		user.Balance += amount
	} else {
		fmt.Println("Депозит не может быть отрицательным")
	}
}

func Withdraw(user *User, amount float64) bool {
	user.mu.Lock()
	defer user.mu.Unlock()
	
	if amount > 0 && user.Balance >= amount {
		user.Balance -= amount
		return true
	}
	fmt.Println("Недостаточно средств")
	return false
}

func Transfer(fromUser *User, toUser *User, amount float64) bool {
    if fromUser == nil || toUser == nil {
        fmt.Println("Пользователи не найдены")
        return false
    }

    first, second := fromUser, toUser
    if fromUser.ID > toUser.ID {
        first, second = toUser, fromUser
    }

    first.mu.Lock()
    second.mu.Lock()
    defer first.mu.Unlock()
    defer second.mu.Unlock()

    if amount <= 0 {
        fmt.Println("Сумма перевода должна быть положительной")
        return false
    }

    if fromUser.Balance < amount {
        fmt.Println("Недостаточно средств")
        return false
    }

    fromUser.Balance -= amount
    toUser.Balance += amount
    return true
}

func ProcessingTransactions(ps *PaymentSystem, transaction Transaction) {
    var fromUser, toUser *User

    for _, user := range ps.Users {
        if user.ID == transaction.FromID {
            fromUser = user
        }
        if user.ID == transaction.ToID {
            toUser = user
        }
    }

    if fromUser == nil || toUser == nil {
        fmt.Println("Пользователи не найдены")
        return
    }

    Transfer(fromUser, toUser, transaction.Amount)
}

func main() {
	user1 := &User{ID: 1, Name: "Spongebob Squarepants", Balance: 1000.0}
	user2 := &User{ID: 2, Name: "Patrick Star", Balance: 500.0}

	paymentSystem := &PaymentSystem{
		Users: []*User{user1, user2},
	}

	Deposit(user1, 500.0)
	Withdraw(user2, 100.0)

	fmt.Printf("User 1: ID=%d, Name=%s, Balance=%.2f\n", user1.ID, user1.Name, user1.Balance)
	fmt.Printf("User 2: ID=%d, Name=%s, Balance=%.2f\n", user2.ID, user2.Name, user2.Balance)

	paymentSystem.Transactions = []Transaction{
		{FromID: 1, ToID: 2, Amount: 100.0},
		{FromID: 2, ToID: 1, Amount: 50.0},
		{FromID: 1, ToID: 2, Amount: 2000.0},
	}

	for _, transaction := range paymentSystem.Transactions {
		ProcessingTransactions(paymentSystem, transaction)
	}

	fmt.Println("После выполнения транзакций:")
	fmt.Printf("User 1: Balance = %.2f\n", user1.Balance)
	fmt.Printf("User 2: Balance = %.2f\n", user2.Balance)
}