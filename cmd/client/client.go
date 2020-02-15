package main

import (
	"database/sql"
	"fmt"
	"github.com/AbduvokhidovRustamzhon/managers-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

func main() {
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

	fmt.Fprintln(os.Stdout, "Добро пожаловать!")
	log.Print("start operations loop")
	operationsLoop(-1, db, unauthorizedOperations, unauthorizedOperationsLoop)
	log.Print("finish operations loop")
	log.Print("finish application")
}

func operationsLoop(userId int64, db *sql.DB, commands string, loop func(db *sql.DB, cmd string, userId int64) (exit bool)) {
	for {
		fmt.Println(commands)
		var cmd string
		_, err := fmt.Scan(&cmd)
		if err != nil {
			log.Fatalf("Can't read input: %v", err)
		}
		if exit := loop(db, strings.TrimSpace(cmd), userId); exit {
			return
		}
	}
}

func unauthorizedOperationsLoop(db *sql.DB, cmd string, userId int64) (exit bool) {
	switch cmd {
	case "1":
		id, ok, err := handleLoginForClient(db)
		if err != nil {
			fmt.Println("Неправильно введён логин или пароль")
			log.Printf("can't handle login: %v", err)
			return false
		}
		if !ok {
			fmt.Println("Неправильно введён логин или пароль. Попробуйте ещё раз.")
			//unauthorizedOperationsLoop(db, "1")
			//Graceful shutdown
			return false
		}
		userId = id
		operationsLoop(userId, db, authorizedOperations, authorizedOperationsLoop)
	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}

func authorizedOperationsLoop(db *sql.DB, cmd string, userId int64) (exit bool) {
	switch cmd {
	// TODO: may be log fatal
	case "1":
		listBalance, err := core.GetBalanceList(db, userId)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printClientBalance(listBalance)

	case "2":
		operationsLoop(userId, db, transactionOperations, transactionOperationsLoop)
	case "3":
		atms, err := core.GetAllAtms(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printAtm(atms)

	case "4":
		operationsLoop(userId,db,serviceOperations,serviceOperationsLoop)
		//err := servicePaying(db)
		//if err != nil {
		//	log.Printf("can't pay fo service: %v", err)
		//	return true
		//}

	case "q":
		return true
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}
	return false
}


func printAtm(atms []core.ATM) {
	for _, atm := range atms {
		fmt.Printf(
			"id: %d, name: %s, street:%s\n",
			atm.Id,
			atm.Name,
			atm.Address,
		)
	}
}

func printClientBalance(listBalance []core.Client)  {
	for _, clientAccounts := range listBalance {
		fmt.Printf(
			"id: %d, name: %s, balanceNumber: %d, balance:%d\n",
			clientAccounts.Id,
			clientAccounts.Name,
			clientAccounts.BalanceNumber,
			clientAccounts.Balance,
		)
	}
}

func printServiceList(serviceList []core.Services)  {
	for _, listService := range serviceList {
		fmt.Printf(
			"id: %d, name: %s, price:%d\n",
			listService.Id,
			listService.Name,
			listService.Price,

		)
	}
}

func handleLoginForClient(db *sql.DB) (id int64,ok bool, err error) {
	fmt.Println("Введите ваш логин и пароль")
	var login string
	fmt.Print("Логин: ")
	_, err = fmt.Scan(&login)
	if err != nil {
		return -1,false, err
	}

	var password string
	fmt.Print("Пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil {
		return -1,false, err
	}

	id, ok, err = core.LoginUser(login, password, db)
	if err != nil {
		return -1,false, err
	}

	return id,ok, err
}

func transaction(db *sql.DB)(err error)  {

	var myPhoneNumber int64
	fmt.Print("Введите свой номер телефон: ")
	_, err = fmt.Scan(&myPhoneNumber)
	if err != nil {
		return err
	}
	var phoneNumber int64
	fmt.Print("Введите номер телефон клиента: ")
	_, err = fmt.Scan(&phoneNumber)
	if err != nil {
		return err
	}
	err = core.CheckByPhoneNumber(phoneNumber, db)
	if err != nil {
		fmt.Println("fddsdf")
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}


	err = core.TransactionMinus(core.Client{
		Id:            0,
		Balance:       balance,
		PhoneNumber:   myPhoneNumber,
	}, db)

	if err != nil {

		fmt.Println("Извините у вас мало денег")
		return err
	}else {
		if  myPhoneNumber == phoneNumber {
			fmt.Println("Схожие номера ")
			authorizedOperationsLoop(db,"2",1)
		}}

	err = core.TransactionPlus(phoneNumber,balance, db)
	if err != nil {
		return err
	}
	fmt.Println("Счет клиента успешно добавлен!")
	return nil
}

func transactionByBalanceNumber(db *sql.DB)(err error)  {

	var mybalanceNumber uint64
	fmt.Print("Введите номер своего баланса: ")
	_, err = fmt.Scan(&mybalanceNumber)
	if err != nil {
		return err
	}
	var balanceNumber uint64
	fmt.Print("Введите номер баланс клиента: ")
	_, err = fmt.Scan(&balanceNumber)
	if err != nil {
		return err
	}
	err = core.CheckByBalanceNumber(balanceNumber, db)
	if err !=nil{
		fmt.Println("щшфырыфр")
		return err
	}
	var balance uint64
	fmt.Print("Введите пополняемую сумму: ")
	_, err = fmt.Scan(&balance)
	if err != nil {
		return err
	}


	err = core.TransactionBalanceNumberMinus(core.Client{
		Id: 0,
		Balance: balance,
		BalanceNumber: mybalanceNumber,
	}, db)

	if err != nil {
		fmt.Println("Извините у вас мало денег")
		return err
	}else {
		if mybalanceNumber == balanceNumber {
			fmt.Println("Так не надо")
			authorizedOperationsLoop(db,"3",1)
		}}

	err = core.TransactionBalanceNumberPlus(balanceNumber,balance, db)
	if err != nil {
		return err
	}
	fmt.Println("Счет клиента успешно добавлен!")
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



func serviceOperationsLoop(db *sql.DB, cmd string, userId int64) bool {
	switch cmd {
	case "1":
		serviceList, err := core.GetAllServices(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			return true // TODO: may be log fatal
		}
		printServiceList(serviceList)

	case "2":
		err := servicePaying(db)
		if err != nil {
			log.Printf("can't pay fo service: %v", err)
			return true
		}

	case "q":
		operationsLoop(userId, db, authorizedOperations, authorizedOperationsLoop)
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}



func transactionOperationsLoop(db *sql.DB, cmd string, userId int64) bool {
	switch cmd {
	case "1":
		err:=transaction(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			operationsLoop(userId, db, transactionOperations, transactionOperationsLoop)
			return true // TODO: may be log fatal
		}
	case "2":
		err:=transactionByBalanceNumber(db)
		if err != nil {
			log.Printf("can't get all products: %v", err)
			operationsLoop(userId, db, transactionOperations, transactionOperationsLoop)
			return true // TODO: may be log fatal
		}

	case "q":
		operationsLoop(userId, db, authorizedOperations, authorizedOperationsLoop)
	default:
		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
	}

	return false
}
