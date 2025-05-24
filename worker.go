package main

import (
	"fmt"
	"sync"
)

func Worker(ps *PaymentSystem, ch <-chan Transaction, wg *sync.WaitGroup) {
    defer wg.Done()
    for t := range ch {
        ProcessingTransactions(ps, t)
    }
}

func transactions() {

	ps := &PaymentSystem{
		Users: []*User{
			{ID: 1, Name: "Spongebob Squarepants", Balance: 1000.0},
			{ID: 2, Name: "Patrick Star", Balance: 500.0},
		},
		Transactions: []Transaction{
			{FromID: 1, ToID: 2, Amount: 200.0},
			{FromID: 2, ToID: 1, Amount: 50.0},
		},
	}

	ch := make(chan Transaction, len(ps.Transactions))

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {

		wg.Add(1)
		go Worker(ps, ch, &wg)
		
	}

	for _, transaction := range ps.Transactions {
		ch <- transaction
	}

	close(ch)
	wg.Wait()
	

	// Подсказка
	fmt.Println("Создаю UserID: 1 с балансом 1000")
	fmt.Println("Создаю UserID: 2 с балансом 500")

	// Подсказка
	fmt.Println("Перевожу с UserID: 1 на UserID: 2 сумму в размере 200")
	fmt.Println("Перевожу с UserID: 2 на UserID: 1 сумму в размере 50")

	// DONE: Добавляем транзакции
	ps.Transactions = []Transaction{
		{FromID: 1, ToID: 2, Amount: 200.0},
		{FromID: 2, ToID: 1, Amount: 50.0},
	}

	fmt.Println("Итого")
	fmt.Printf("У первого пользователя должно получиться 850, а получилось %f\n", ps.Users[0].Balance)
	fmt.Printf("У второго пользователя должно получиться 650, а получилось %f", ps.Users[1].Balance)
}
