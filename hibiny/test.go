package hibiny

import "fmt"

//Test func
func Test() {

	fmt.Println("НОВОСТИ")
	news := GetNews()
	for _, n := range news {
		fmt.Println("Заголовок: ", n.Header)
		fmt.Println("Иконка: ", n.ImageURL)
		fmt.Println("Дата: ", n.Data)
		fmt.Println("Контент: ", n.Content)
	}
}
