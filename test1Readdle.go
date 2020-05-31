package main

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"net/http"
	"io/ioutil"
    "encoding/json"
)

func main() {
	
	urlString := "https://date.nager.at/Api/v1/Get/UA/2020"
	data, err := loadHolidayData(urlString)
	errorAlarm(err)

    holidays, err := parseHolidayData(data)
	errorAlarm(err)

	//определяем какой сегодня день
	nowDay := time.Now()
	//nowDay := time.Date(2020, 12, 28, 20, 34, 58, 651387237, time.UTC)
	fmt.Println("Today is", nowDay.Weekday(),",", nowDay.Month(), ",", nowDay.Day())

	message := makeHolidayMassage(nowDay, holidays)
	fmt.Println(message)
}

// define data structure 
type holiday struct {
	Date        string      `json:"date"`
	LocalName   string      `json:"localName"`
	Name        string      `json:"name"`
	CountryCode string      `json:"countryCode"`
	Fixed       bool        `json:"fixed"`
	Global      bool        `json:"global"`
	Counties    string      `json:"counties"`
	LaunchYear  int         `json:"launchYear"`
	Type        string      `json:"type"`
}

// получаем данные с сайта
func loadHolidayData(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	errorAlarm(err)

	dataFromUrl, err := ioutil.ReadAll(resp.Body)
	errorAlarm(err)
	defer resp.Body.Close()
	isJsonDataValidAlarm (dataFromUrl)

	return dataFromUrl, nil
}

// превращаем json-данные в массив определённого типа
func parseHolidayData(jsonData []byte) ([]holiday, error) {
	//создаем массив json-обьектов нужного типа
    var allHolidaysArray []holiday

    //заполняем массив json-данными
	err := json.Unmarshal(jsonData, &allHolidaysArray)
	errorAlarm(err)

	return allHolidaysArray, nil
}

//определяет, ждут ли нас длинные выходные,
//если да - возвращает строку с датами начала и конца длинных выходных 
//если выходные короткие - пустую строку
func isLongWeekendMassage(tmpDateOfHoliday time.Time) (massage string) {
    switch tmpDateOfHoliday.Weekday() {
	    case time.Friday, time.Saturday:
	    	lastDayOff := tmpDateOfHoliday.AddDate(0, 0, 2)

	    	return ("\nThe weekend will last 3 days: "+ 
	    		tmpDateOfHoliday.Month().String() + ", "+ strconv.Itoa(tmpDateOfHoliday.Day()) +
	    		"  -  " +
	    		lastDayOff.Month().String() + ", "+ strconv.Itoa(lastDayOff.Day()))
	        
	    case time.Sunday:
	    	firstDayOff := tmpDateOfHoliday.AddDate(0, 0, -1)
	    	lastDayOff := tmpDateOfHoliday.AddDate(0, 0, 1)

	    	return ("\nThe weekend will last 3 days: "+ 
	    		firstDayOff.Month().String()  + ", "+ strconv.Itoa(firstDayOff.Day()) +
	    		"  -  " +
	    		lastDayOff.Month().String()  + ", "+ strconv.Itoa(lastDayOff.Day()))
	        
	    case time.Monday:
	        firstDayOff := tmpDateOfHoliday.AddDate(0, 0, -2)
	    	return ("\nThe weekend will last 3 days: "+ 
	    		firstDayOff.Month().String()  + ", "+ strconv.Itoa(firstDayOff.Day()) +
	    		"  -  " +
	    		tmpDateOfHoliday.Month().String() + ", "+ strconv.Itoa(tmpDateOfHoliday.Day()))
	        
	    default:
	        return " "
    }
}

// возвращает строку - сообщение с праздником
func makeCongratulationMessage(holidayDay holiday, tmpDateOfHoliday time.Time, today bool) (massage string){
	massageString := ""
	if(today == true){
		massageString += "Today is a holiday! It is " + holidayDay.Name
	} else {
		massageString += ("The next holiday is " + holidayDay.Name +
    		" (" +
    		tmpDateOfHoliday.Weekday().String() + ", " +
    		tmpDateOfHoliday.Month().String() + ", " +
    		strconv.Itoa(tmpDateOfHoliday.Day()) +
    		")")
	}
	massageString += isLongWeekendMassage(tmpDateOfHoliday)
    return massageString
}

// возвращает строку - финальный результат
func makeHolidayMassage(now time.Time,  holidays []holiday) (massage string) {

	for _, kHoliday := range holidays {
    	
    	DateOfHoliday := parseDateFromString(now, kHoliday)
    	//
    	/*в Json-файле праздники отсортированы в правильном порядке, 
    		поэтому когда мы встречаем первую дату, которая равна сегодняшнему дню
    		(или является будущей, относительно нашей даты) 
    		выводим сообщение о празднике и выходим из цикла
    	*/ 
    	if now.Before(DateOfHoliday) {
    		return makeCongratulationMessage(kHoliday, DateOfHoliday, false)
    		
    	} else if now.Equal(DateOfHoliday) {
    		return makeCongratulationMessage(kHoliday, DateOfHoliday, true)    		
    	}
    }
    return "There are no holidays until the end of the year."
}

//превращаем Json-дату в формате строки в дату типа Date()
func parseDateFromString(now time.Time, holidayDay holiday) (time.Time) {
	/* данные в Json-файле это строка типа "yyyy-mm-dd"
    		разбиваем её на массив из трех ячеек [yyyy, mm, dd] 
    		и создаем обьект типа Дата (вместо параметров которых мы не знаем - используем параметры текущей даты)
    */
    tmpDateArr := strings.Split(holidayDay.Date, "-")
    tmpMonth, _ := strconv.Atoi(tmpDateArr[1])
    tmpDay, _ := strconv.Atoi(tmpDateArr[2])
    tmpDateOfHoliday := time.Date(
    				now.Year(), time.Month(tmpMonth), tmpDay,
    				now.Hour(), now.Minute(), now.Second(),
    				now.Nanosecond(), now.Location())
    return tmpDateOfHoliday
}

// проверка не являются ли введенные данные пустыми
func errorAlarm(b interface{}) {
    if b != nil {
      fmt.Print(b)
    }
}

// проверяем, являются ли введенные json-данные валидными
func isJsonDataValidAlarm (data []byte) {
    if(json.Valid([]byte(data)) == false){
    	panic("invalid json data")
    }
}
