package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/nlopes/slack"
	"math/rand"
	"strings"
	"time"
)

const ReminderText string = "You are up this week to wrangle the PR's/Issues for Sig Docs! :tada:"

var GithubUsernameToSlackUsernameMap = map[string]string{
	"zparnold": "U7EFYPNUR",
}

func main() {

	c := colly.NewCollector()

	//grab the first table on the page (ASSUMPTION! The first table is the most recent table to be processed)
	c.OnHTML("tbody:first-of-type", func(e *colly.HTMLElement) {
		//iterate over rows
		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			htmlString, err := row.DOM.Html()
			if err != nil {
				panic(err)
			} else {
				//preprocess input from webpage
				entities := getEntitiesFromHtml(htmlString)

				//be a courteous user and set a random delay
				rand.Seed(time.Now().UnixNano())
				rDelay := rand.Intn(3000)
				time.Sleep(time.Millisecond * time.Duration(rDelay))

				//This assumes the table format is: DATE, USER
				setSlackReminder(entities[0], entities[1])
			}
		})
	})

	err := c.Visit("https://github.com/kubernetes/website/wiki/PR-Wranglers")
	if err != nil{
		panic(err)
	}

}

//This is a cheater function to speed up the performance of the loops above
func getEntitiesFromHtml(s string) (r []string) {
	//remove nil at end of string
	s = strings.Replace(s, "<nil>", "", -1)

	//remove @ from usernames
	s = strings.Replace(s, "@", "", -1)

	//remove <td> tags
	s = strings.Replace(s, "<td>", "", -1)
	s = strings.Replace(s, "</td>", "", -1)

	r = strings.Split(s, "\n")

	fmt.Println(r)
	return r
}

func setSlackReminder(date, toWhom string) {
	api := slack.New("YOUR_TOKEN_HERE", slack.OptionDebug(true))
	slackUser := GithubUsernameToSlackUsernameMap[toWhom]
	d := fmt.Sprintf("%s at 9am", date)
	reminder, err := api.AddUserReminder(slackUser, ReminderText, d)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("Reminder set: ", reminder.ID)
	}
}