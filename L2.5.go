package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

/*

Прогограмма выведет строку "error", т.к. структура customError реализует интерфейс error
Сама переменная err будет равна nil, но функция test() возвращает структуру customError,
которая, в свою очередь, реализует метод Error(), т.е. тут произошла следующая ситуация,
что указатель на data = nil, а указатель на itab != nil (потому что известен тип, это *customError), соответственно, проверить переменную err на равенство nil не получится
err != nil проверяет сам интерфейс, а не только указатель на data, если itab != nil, то интерфейс считается ненулевым

*/
