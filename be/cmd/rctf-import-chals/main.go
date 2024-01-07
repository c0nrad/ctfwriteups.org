package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/c0nrad/ctfwriteups/config"
	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Challenge struct {
	Name        string
	Description string
	Solves      int
	Category    string
}

func LoadChallenges() []Challenge {
	chalFile, err := os.Open("chals.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(chalFile)
	var challenges []Challenge
	err = decoder.Decode(&challenges)
	if err != nil {
		panic(err)
	}
	return challenges
}

func main() {
	env := "prod"

	config.InitLogger()
	config.InitEnv(env)
	datastore.InitDatabase()

	challenges := LoadChallenges()
	for _, c := range challenges {
		fmt.Println(c.Name, c.Category, c.Solves)
	}

	me, err := datastore.GetUserByEmail(context.Background(), "c0nrad@c0nrad.io")
	if err != nil {
		panic(err)
	}

	ctfName := "BackdoorCTF 2023"
	ctf, err := datastore.GetCTFByName(context.TODO(), ctfName)
	if err != nil {
		panic(err)
	}

	for i := range challenges {
		challenges[i].Category = strings.ToLower(challenges[i].Category)
	}

	isMissingCategory := false
	for _, c := range challenges {
		if !slices.Contains(ctf.Categories, c.Category) {
			fmt.Println("missing category", c.Category)
			ctf.Categories = append(ctf.Categories, c.Category)
			isMissingCategory = true
		}
	}

	if isMissingCategory {
		err = datastore.UpdateCTF(context.Background(), *ctf)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Saving CTF", ctf)
	fmt.Println("ENv", env)
	fmt.Println("Challenge Count", len(challenges))
	time.Sleep(10 * time.Second)

	for _, c := range challenges {
		chal := models.Challenge{
			ID:          primitive.NewObjectID(),
			TS:          time.Now(),
			Name:        c.Name,
			Solves:      c.Solves,
			Category:    c.Category,
			CTFID:       ctf.ID,
			SubmitterID: me.ID,
		}

		existingChal, err := datastore.GetChallengeByName(context.Background(), c.Name, c.Category, ctf.ID.Hex())
		if err != nil {
			fmt.Println("New Chal!", chal)

			err = datastore.SaveChallenge(context.Background(), chal)
			if err != nil {
				panic(err)
			}

			err = datastore.IncrementChallengeCount(context.Background(), ctf.ID.Hex())
			if err != nil {
				panic(err)
			}

			continue
		} else {
			if existingChal.Solves != chal.Solves {
				existingChal.Solves = chal.Solves
				err = datastore.UpdateChallenge(context.Background(), *existingChal)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
