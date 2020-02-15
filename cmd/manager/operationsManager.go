package main

//
//const authorizedOperations  = `Список доступных операций:
//	1. Добавить пользователя
//	2. Добавить счёт пользователю (тогда сразу в пользователе)
//	3. Добавить услуги (название)
//	4. Добавить банкомат
//	q. Выйти из приложения
//
//Введите команду
//`
//
//const unauthorizedOperations  = `Список доступных операций:
//	1. Авторизоваться
//	q. Выйти из приложения
//Введите команду
//`
//



const unauthorizedOperations = `Список доступных операций:
	1. Авторизация
	2. Список банкоматов
	q. Выйти из приложения

Введите команду`

const authorizedOperations = `Список доступных операций:
	1. Просмотр списка продуктов
	2. Продажа товара
	3. Добавить клиента
	4. Добавить банкомат
	5. Добавить услугу
	6. пополнить баланс клиента
	7. добавить карту
	8. Экспорт/Импорт
	9. Заплатить за услугу
	q. Выйти (разлогиниться)

Введите команду`

const exportImportOperationsLoop  = `Список доступных операций:
	1. Экспортировать список банокматов в JSON
	2. Экспортировать список клиентов в JSON
	3. Экспортировать список банкоматов в XML
	4. Экспортировать список клиентов в XML
	5. Импортировать список банкоматов в JSON
	6. Импортировать список клиентов в JSON
	7. Импортировать список банкоматов в XML
	8. Импортировать список клиентов в XML

	q. Назад

Введите команду`

//
//func ExportOperationsLoop(db *sql.DB, cmd string) (exit bool) {
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
//		operationsLoop(db, exportOperationsLoop, authorizedOperationsLoop)
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
