package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"strings"
	"time"
)

//Обьявление аргументов и команд.
var (
	app            = kingpin.New("notification", "A notification of tasks application.")
	register       = app.Command("register", "Register a new task.")
	registerDate   = register.Arg("date", "Date of the task. Format: 'yyyy-MM-dd hh:mm:ss'").String()
	registerTask   = register.Arg("task", "Description of task.").String()
	announcement   = app.Command("announcement", "Notification mode")
	announcementIn = announcement.Arg("announcementIn", "In how long make announcement").String()
)

//Структура того, что читаем с базы данных.
type taskstr struct {
	id   int
	date string
	task string
}

//Считывает таски в структуру.
func ReadTasks(db *sql.DB) (tasksw []taskstr, err error) {
	rows, err := db.Query("select * from golang.tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []taskstr{}
	for rows.Next() {
		t := taskstr{}
		err := rows.Scan(&t.id, &t.date, &t.task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

//Соединение с базой данных.
func Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golang")
	return db, err
}

// Описание формата времени.
const timeLayout = "2006-01-02T15:04:05"

// Вывод уведомления о таске в указанное время.
func callAt(callTime string, task string, durationIn time.Duration) (time.Duration, error) {
	//Переменная для времменого промежутка до таска, которую задал юзер.
	var durationUser time.Duration = 0
	//Разбираем время таска, не придумал лучше способа как это сделать с тем форматом что приходит с базы данных.
	callTime = strings.ReplaceAll(callTime, " ", "T")
	loc, _ := time.LoadLocation("Europe/Kiev")
	ctime, err := time.ParseInLocation(timeLayout, callTime, loc)
	if err != nil {
		return 0, err
	}
	//Вычисляем временной промежуток до таска.
	duration := ctime.Sub(time.Now().Local())
	//Вычисляем временной промежуток до таска который задал юзер.
	durationUser = ctime.Sub(time.Now().Local().Add((durationIn)))
	//Выводим уведомление через столько таск относительно момента запуска программы.
	fmt.Println(task, "is in ", duration.String())
	//Создам горутин.
	go func() {
		//Если времени до таска больше чем задал юзер(то есть можно вывести уведомление, что таск за "время которое задал юзер") то слипаем горутин на это время.
		if durationUser > 0 {
			time.Sleep(durationUser)
			//Выводим уведомление что таск через "время которое задал юзер".
			fmt.Println(task, "is in "+durationIn.String())
			//Слипаем на время до самого таска.
			time.Sleep(duration - durationUser)
		} else {
			//Если времени меньше, то тоже слипаем но на время до самого таска.
			time.Sleep(duration)
		}
		//Выводим уведомление что таск прямо сейчас.
		fmt.Println(task, "is right now")
	}()
	return duration, nil
}
func main() {
	db, err := Connect()
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var durMax time.Duration = 0
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Register user
	case register.FullCommand():
		_, err := db.Exec("INSERT INTO tasks (date, task) VALUES(?,?)", *registerDate, *registerTask)
		if err != nil {
			fmt.Println(err)
		}
	case announcement.FullCommand():
		tasks, err := ReadTasks(db)
		if err != nil {
			fmt.Println(err)
		}
		for _, t := range tasks {
			hr, _ := time.ParseDuration(*announcementIn)
			dur, err := callAt(t.date, t.task, hr)
			if err != nil {
				fmt.Println(err)
			}
			//Вычесляем самый поздний таск(таск до которого самый большой промежуток времени).
			if dur > durMax {
				durMax = dur
			}
		}
		//Продолжаем работу программы до того момента пока не наступит самый поздний таск.
		time.Sleep(durMax)
	}
}
