package rest

import (
	"fmt"
	"log"
	"time"
)

var productsFixtures = []product{
	{ID: "1", Title: "Bose Headphone", Description: "Bose Headphone", PriceInCents: 330 * 100},
	{ID: "2", Title: "Panasonic BT PRO", Description: "Panasonic BT PRO", PriceInCents: 299 * 100},
	{ID: "3", Title: "JBL Masters Earphone", Description: "JBL Masters Earphone", PriceInCents: 279 * 100},
}

var usersFixtures = []user{
	{ID: "1", FirstName: "User1", LastName: "LastName1", DateOfBirth: yearsAgo(33)}, //time.Now().Add(-22 * 12 * 720 * time.Hour)},
}

func yearsAgo(years int) time.Time {
	today := time.Now()
	t, err := time.Parse("02/01/2006", fmt.Sprintf("%d/%d/%d", today.Day(), today.Month(), today.Year()-years))
	if err != nil {
		log.Println("fixtures:yearsAgo err:", err.Error())
	}
	return t
}
