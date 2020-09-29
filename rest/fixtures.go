package rest

import (
	"fmt"
	"log"
	"time"
)

var productsFixtures = []product{
	{ID: "1", Title: "Product1", Description: "Product1 Desc: User1 has discount today!", PriceInCents: 10 * 100},
	{ID: "2", Title: "Product2", Description: "Product2 Desc", PriceInCents: 15 * 100},
	{ID: "3", Title: "Product3", Description: "Product3 Desc", PriceInCents: 12 * 100},
}

var usersFixtures = []user{
	{ID: "1", FirstName: "User1", LastName: "LastName1", DateOfBirth: time.Now().Add(-22 * 12 * 720 * time.Hour)},
}

func yearsAgo(years int) time.Time {
	today := time.Now()
	t, err := time.Parse("02/01/2006", fmt.Sprintf("%d/%d/%d", today.Day(), today.Month(), today.Year()-years))
	if err != nil {
		log.Println("fixtures:yearsAgo err:", err.Error())
	}
	return t
}
