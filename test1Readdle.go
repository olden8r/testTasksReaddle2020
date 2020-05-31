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
	// get json-data
	urlString := "https://date.nager.at/Api/v1/Get/UA/2020"
	data, err := loadHolidayData(urlString)
	errorAlarm(err)

	// parse json-data
    holidays, err := parseHolidayData(data)
	errorAlarm(err)

	//define what day is today
	nowDay := time.Now()
	fmt.Println("Today is", nowDay.Weekday(),",", nowDay.Month(), ",", nowDay.Day())

	//print message with result
	message := makeHolidayMessage(nowDay, holidays)
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

// get data from API
func loadHolidayData(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	errorAlarm(err)

	dataFromUrl, err := ioutil.ReadAll(resp.Body)
	errorAlarm(err)
	defer resp.Body.Close()
	isJsonDataValidAlarm (dataFromUrl)

	return dataFromUrl, nil
}

// convert json-data into array of type holiday
func parseHolidayData(jsonData []byte) ([]holiday, error) {
	//create an array of type holiday
    var allHolidaysArray []holiday

    //fill an array with json-data
	err := json.Unmarshal(jsonData, &allHolidaysArray)
	errorAlarm(err)

	return allHolidaysArray, nil
}

//determines if a long weekend awaits us,
//if yes - returns a string with the start and end dates of the long weekend
//if the weekend is short - an empty string
func isLongWeekendMessage(tmpDateOfHoliday time.Time) (message string) {
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

// returns a string - a congratulation message about holiday
func makeCongratulationMessage(holidayDay holiday, tmpDateOfHoliday time.Time, today bool) (message string){
	messageString := ""
	if(today == true){
		messageString += "Today is a holiday! It is " + holidayDay.Name
	} else {
		messageString += ("The next holiday is " + holidayDay.Name +
    		" (" +
    		tmpDateOfHoliday.Weekday().String() + ", " +
    		tmpDateOfHoliday.Month().String() + ", " +
    		strconv.Itoa(tmpDateOfHoliday.Day()) +
    		")")
	}
	messageString += isLongWeekendMessage(tmpDateOfHoliday)
    return messageString
}

// returns a string - the final result
func makeHolidayMessage(now time.Time,  holidays []holiday) (message string) {
	for _, kHoliday := range holidays {	
		
    	DateOfHoliday := parseDateFromString(now, kHoliday)

    	// in the json file, the holidays are sorted in the correct order,
    	// so when we meet the first date, which is equal to today date
    	// or is bigger than our date
		// display a message about the holiday and exit the cycle

    	if now.Before(DateOfHoliday) {
    		return makeCongratulationMessage(kHoliday, DateOfHoliday, false)
    		
    	} else if now.Equal(DateOfHoliday) {
    		return makeCongratulationMessage(kHoliday, DateOfHoliday, true)    		
    	}
    }
    return "There are no holidays until the end of the year."
}

//convert Json-data from string format into format time.Date()
func parseDateFromString(now time.Time, holidayDay holiday) (time.Time) {

	// the data in the json file is a string like "yyyy-mm-dd",
	// we break it into an array of three cells [yyyy, mm, dd]
	// and create an object of type Date
	//instead of the unknown parameters- we use the parameters of the current date

    tmpDateArr := strings.Split(holidayDay.Date, "-")
    tmpMonth, _ := strconv.Atoi(tmpDateArr[1])
    tmpDay, _ := strconv.Atoi(tmpDateArr[2])
    tmpDateOfHoliday := time.Date(
    				now.Year(), time.Month(tmpMonth), tmpDay,
    				now.Hour(), now.Minute(), now.Second(),
    				now.Nanosecond(), now.Location())
    return tmpDateOfHoliday
}

// check if the entered data is empty
func errorAlarm(b interface{}) {
    if b != nil {
      fmt.Print(b)
    }
}

// check if the entered json-data is valid
func isJsonDataValidAlarm (data []byte) {
    if(json.Valid([]byte(data)) == false){
    	panic("invalid json data")
    }
}
