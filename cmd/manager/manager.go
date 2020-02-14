package main

import (
	"database/sql"
	"fmt"
	"github.com/AbduvokhidovRustamzhon/managers-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strings"
)

const errorLogin = "Неправильно введён логин или пароль. Попробуйте ещё раз."
const incorrectCommand = "Вы выбрали неверную команду: %s\n"

// TODO: для тех, кто хочет попробовать, можете использовать структуры и методы:
type manager struct {
	db  *sql.DB
	out io.Writer
	in  io.Reader
}

func newManagerCLI(db *sql.DB, out io.Writer, in io.Reader) *manager {
	return &manager{db: db, out: out, in: in}
}

// Writer, Reader

func main() {
	// os.Stdin, os.Stout, os.Stderr, File
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Print("start application")
	log.Print("open db")
	db, err := sql.Open("sqlite3", "db.sqlite")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer func() {
		log.Print("close db")
		if err := db.Close(); err != nil {
			log.Fatalf("can't close db: %v", err)
		}
	}()
	err = core.Init(db)
	if err != nil {
		log.Fatalf("can't init db: %v", err)
	}

	fmt.Fprintln(os.Stdout, "Добро пожаловать в наше приложение")
	log.Print("start operations loop")
	operationsLoop(db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}

func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err) // %v - natural ...
		}
		if exit := loop(db, strings.TrimSpace(cmd)); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		ok, err := handleLogin(db)
		if err != nil {
			log.Printf("can't handle login: %v", err)
			fmt.Println(errorLogin)
			return false
		}
		if !ok {
			fmt.Println(errorLogin)
			//unauthorizedOperationsLoop(db, "1")
			//Graceful shutdown
			return false
		}
		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
	case "2": // TODO:  список банкоматов
		products, err := core.GetAllProducts(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printProducts(products)
	case "q":
		return true
	default:
		fmt.Printf(incorrectCommand, cmd)
	}

	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
	switch cmd {
	case "1":
		products, err := core.GetAllProducts(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printProducts(products) // Ctrl + Alt + M
	case "2":
		err := handleSale(db)
		if err != nil {
			log.Printf("can't add sale: %v", err)
			return true
		}
	case "3":
		err := handleClient(db)
		if err != nil {
			log.Printf("can't add client: %v", err)
		}
	case "4":
		err := handleATM(db)
		if err != nil {
			log.Printf("can't add client: %v", err)
		}
	case "5":
		err := handleService(db)
		if err != nil {
			log.Printf("can't add service: %v", err)
		}
	case "6":
		err := updateBalance(db)
		if err != nil {
			log.Printf("can't update balance: %v", err)
		}
	case "7":
		fmt.Println("w") //TODO
	case "8":
		err := ExportOperationsLoop
		if err != nil {
			log.Printf("can't open function: %v", err)
		}
	case "9":
		err := servicePaying(db)
		if err != nil {
			log.Printf("can't pay for service")
		}

	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}

func printProducts(products []core.Product) {
	for _, product := range products {
		fmt.Printf(
			"id: %d, name: %s, price: %d, qty: %d\n",
			product.Id,
			product.Name,
			product.Price,
			product.Qty,
		)
	}
}

func handleLogin(db *sql.DB) (ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return false, err
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return false, err
	}

	ok, err = core.Login(login, password, db)
	if err != nil {
		return false, err
	}

	return ok, err
}

func handleSale(db *sql.DB) (err error) {
	fmt.Println("Введите данные")
	var id int64
	fmt.Print("Id продукта: ")
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}
	var qty int64
	fmt.Print("Количество: ")
	_, err = fmt.Scan(&qty)
	if err != nil {
		return err
	}

	err = core.Sale(id, qty, db)
	if err != nil {
		return err
	}

	return nil
}

func handleClient(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные клиента")
	var name string
	fmt.Print("Введите имя: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	var login string
	fmt.Print("Введите логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return err
	}

	var password string
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return err
	}

	var passportSeries string
	fmt.Print("Введите вашу серию пасспорта: ")
	_, err = fmt.Scan(&passportSeries)
	if err != nil {
		return err
	}

	var numberPhone int
	fmt.Print("Введите телефон-номер клиента: ")
	_, err = fmt.Scan(&numberPhone)
	if err != nil {
		return err
	}

	var balance uint64
	fmt.Print("Введите начальный баланс: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	var balanceNumber int64
	fmt.Print("Введите номер счета: ")
	_, err = fmt.Scan(&balanceNumber)
	if err != nil {
		return err
	}

	err = core.AddUser(name, login, password, passportSeries, numberPhone, balance, balanceNumber, db)
	if err != nil {
		return err
	}
	fmt.Println("Клиент успешно добавлен!")
	return nil
}

func handleATM(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные банкомата")
	var name string
	fmt.Print("Введите имя: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	var address string
	fmt.Print("Введите адресс: ")
	_, err = fmt.Scan(&address)
	if err != nil {
		return err
	}

	err = core.AddAtm(name, address, db)
	if err != nil {
		return err
	}
	fmt.Println("Банкомат успешно добавлен!")
	return nil
}

func handleService(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные услуги")
	var name string
	fmt.Print("Введите название услуги: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	var price int64
	fmt.Print("Введите цену услуги: ")
	_, err = fmt.Scan(&price)
	if err != nil {
		return err
	}

	err = core.AddService(name, price, db)
	if err != nil {
		return err
	}
	fmt.Println("Услуга успешно добавлена!")
	return nil
}

func updateBalance(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные клиента")
	var id int64
	fmt.Print("Введите счет клиента: ")
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}
	var balance int64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	err = core.UpdateBalanceClient(id, balance, db)
	if err != nil {
		return err
	}
	fmt.Println("Счет клиента успешно добавлен!")
	return nil
}

func handleCard(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные счета: ")
	var name string
	fmt.Print("Введите имя: ")
	_, err = fmt.Scan(&name)
	if err != nil {
		return err
	}
	var cardBalance string
	fmt.Print("Введите адрес: ")
	_, err = fmt.Scan(&cardBalance)
	if err != nil {
		return err
	}

	var cardUserId string
	fmt.Print("Введите адрес: ")
	_, err = fmt.Scan(&cardUserId)
	if err != nil {
		return err
	}

	//	err = core.AddCard(name, cardBalance, cardUserId, db)
	if err != nil {
		return err
	}
	fmt.Println("Банкомат успешно добавлен!")
	return nil
}

func servicePaying(db *sql.DB) (err error) { // dobavka klienta
	fmt.Println("Введите данные услуги: ")
	var id int64
	fmt.Print("Введите ID услуги: ")
	_, err = fmt.Scan(&id)
	if err != nil {
		return err
	}
	var balance int64
	fmt.Print("Введите сумму которую вы собираетесь заплатить: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}

	fmt.Println("Подтвердите что это вы!")
	if login, yes := areYouSure(db); yes {
		err := core.UpdateBalanceClientForService(login, balance, db)
		if err != nil {
			return err
		}

		err = core.PayForService(id, balance, db)
		if err != nil {
			return err
		}


		fmt.Println("Услуга успешно оплачена!")
		return nil
	} else {
		fmt.Println("Операция отменена")
		return nil
	}

}

func areYouSure(db *sql.DB) (string, bool) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err := fmt.Scan(&login)
	if err != nil {
		return login, false
	}
	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return login, false
	}

	_, err = core.Login(login, password, db)
	if err != nil {
		return login, false
	}

	return login, true
}
