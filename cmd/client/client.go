package main
//
//import (
//"database/sql"
//"fmt"
//"github.com/AbduvokhidovRustamzhon/managers-core/pkg/core"
//_ "github.com/mattn/go-sqlite3"
//"io"
//"log"
//"os"
//"strings"
//)
//
//const errorLogin = "Неправильно введён логин или пароль. Попробуйте ещё раз."
//const incorrectCommand  = "Вы выбрали неверную команду: %s\n"
//// TODO: для тех, кто хочет попробовать, можете использовать структуры и методы:
//type manager struct {
//	db  *sql.DB
//	out io.Writer
//	in  io.Reader
//}
//
//func newManagerCLI(db *sql.DB, out io.Writer, in io.Reader) *manager {
//	return &manager{db: db, out: out, in: in}
//}
//
//// Writer, Reader
//
//func main() {
//	// os.Stdin, os.Stout, os.Stderr, File
//	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.SetOutput(file)
//	log.Print("start application")
//	log.Print("open db")
//	db, err := sql.Open("sqlite3", "db.sqlite")
//	if err != nil {
//		log.Fatalf("can't open db: %v", err)
//	}
//	defer func() {
//		log.Print("close db")
//		if err := db.Close(); err != nil {
//			log.Fatalf("can't close db: %v", err)
//		}
//	}()
//	err = core.Init(db)
//	if err != nil {
//		log.Fatalf("can't init db: %v", err)
//	}
//
//	fmt.Fprintln(os.Stdout, "Добро пожаловать в наше приложение")
//	log.Print("start operations loop")
//	operationsLoop(db, unauthorizedOperations, unauthorizedOperationsLoop)
//	log.Print("finish operations loop")
//	log.Print("finish application")
//}
//
//func operationsLoop(db *sql.DB, commands string, loop func(db *sql.DB, cmd string) bool) {
//	for {
//		fmt.Println(commands)
//		var cmd string
//		_, err := fmt.Scan(&cmd)
//		if err != nil {
//			log.Fatalf("Can't read input: %v", err) // %v - natural ...
//		}
//		if exit := loop(db, strings.TrimSpace(cmd)); exit {
//			return
//		}
//	}
//}
//
//func unauthorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
//	switch cmd {
//	case "1":
//		ok, err := handleLogin(db)
//		if err != nil {
//			log.Printf("can't handle login: %v", err)
//			fmt.Println(errorLogin)
//			return false
//		}
//		if !ok {
//			fmt.Println(errorLogin)
//			//unauthorizedOperationsLoop(db, "1")
//			//Graceful shutdown
//			return false
//		}
//		operationsLoop(db, authorizedOperations, authorizedOperationsLoop)
//	case "2":  // TODO:  список банкоматов
//		products, err := core.GetAllProducts(db)
//		if err != nil {
//			log.Printf("can't get all products: %v", err)
//			return true // TODO: may be log fatal
//		}
//		printProducts(products)
//	case "q":
//		return true
//	default:
//		fmt.Printf(incorrectCommand, cmd)
//	}
//
//	return false
//}
//
//func authorizedOperationsLoop(db *sql.DB, cmd string) (exit bool) {
//	switch cmd {
//	case "1":
//		products, err := core.GetAllProducts(db)
//		if err != nil {
//			log.Printf("can't get all products: %v", err)
//			return true // TODO: may be log fatal
//		}
//		printProducts(products) // Ctrl + Alt + M
//	case "2":
//		err := handleSale(db)
//		if err != nil {
//			log.Printf("can't add sale: %v", err)
//			return true
//		}
//	case "3":
//		fmt.Println("hello")
//
//	case "q":
//		return true
//	default:
//		fmt.Printf("Вы выбрали неверную команду: %s\n", cmd)
//	}
//	return false
//}
//
//func printProducts(products []core.Product) {
//	for _, product := range products {
//		fmt.Printf(
//			"id: %d, name: %s, price: %d, qty: %d\n",
//			product.Id,
//			product.Name,
//			product.Price,
//			product.Qty,
//		)
//	}
//}
//
//func handleLogin(db *sql.DB) (ok bool, err error) {
//	fmt.Println("Введите ваш логин и пароль")
//	var login string
//	fmt.Print("Логин: ")
//	_, err = fmt.Scan(&login)
//	if err != nil {
//		return false, err
//	}
//	var password string
//	fmt.Print("Пароль: ")
//	_, err = fmt.Scan(&password)
//	if err != nil {
//		return false, err
//	}
//
//	ok, err = core.Login(login, password, db)
//	if err != nil {
//		return false, err
//	}
//
//	return ok, err
//}
//
//func handleSale(db *sql.DB) (err error) {
//	fmt.Println("Введите данные продажи")
//	var id int64
//	fmt.Print("Id продукта: ")
//	_, err = fmt.Scan(&id)
//	if err != nil {
//		return err
//	}
//	var qty int64
//	fmt.Print("Количество: ")
//	_, err = fmt.Scan(&qty)
//	if err != nil {
//		return err
//	}
//
//	err = core.Sale(id, qty, db)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
