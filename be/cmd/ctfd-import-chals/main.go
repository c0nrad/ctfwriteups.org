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
	Name  string
	Title string

	Description string
	Solves      int
	Category    string
	Categories  []string
}

func LoadChallenges() []Challenge {
	chalFile, err := os.Open("challenges.json")
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
	env := os.Getenv("ENV")
	if env == "" {
		env = "prod"
	}

	config.InitLogger()
	config.InitEnv(env)
	datastore.InitDatabase()

	challenges := LoadChallenges()
	for _, c := range challenges {
		fmt.Println(c.Name, c.Category, c.Solves)
	}

	// os.Exit(0)

	me, err := datastore.GetUserByEmail(context.Background(), "c0nrad@c0nrad.io")
	if err != nil {
		panic(err)
	}

	ctfName := "PotluckCTF 2023"
	ctf, err := datastore.GetCTFByName(context.TODO(), ctfName)
	if err != nil {
		panic(err)
	}

	for i := range challenges {
		challenges[i].Category = strings.ToLower(challenges[i].Category)
	}

	isMissingCategory := false
	for i := range challenges {
		challenges[i].Category = strings.ToLower(challenges[i].Category)
		if challenges[i].Category == "" && len(challenges[i].Categories) > 0 {
			challenges[i].Category = challenges[i].Categories[0]
		}
		if challenges[i].Name == "" && challenges[i].Title != "" {
			challenges[i].Name = challenges[i].Title
		}
		if !slices.Contains(ctf.Categories, challenges[i].Category) {
			fmt.Println("missing category", challenges[i].Category)
			ctf.Categories = append(ctf.Categories, challenges[i].Category)
			isMissingCategory = true
		}
	}

	fmt.Println("Saving CTF", ctf)
	fmt.Println("ENv", env)
	fmt.Println("Challenge Count", len(challenges))
	time.Sleep(10 * time.Second)

	if isMissingCategory {
		err = datastore.UpdateCTF(context.Background(), *ctf)
		if err != nil {
			panic(err)
		}
	}

	for _, c := range challenges {
		chal := models.Challenge{
			ID:          primitive.NewObjectID(),
			TS:          time.Now(),
			Name:        c.Name,
			Solves:      c.Solves,
			Category:    strings.ToLower(c.Category),
			CTFID:       ctf.ID,
			SubmitterID: me.ID,
		}

		fmt.Printf("%+v\n", chal)

		existingChal, err := datastore.GetChallengeByName(context.Background(), c.Name, c.Category, ctf.ID.Hex())
		if err != nil {
			fmt.Println("New Chal!", chal)

			err = datastore.SaveChallenge(context.Background(), chal)
			if err != nil {
				panic(err)
			}
			continue
		} else {
			fmt.Println("chal already exist")
			if existingChal.Solves != chal.Solves {
				fmt.Println("updating solves")
				existingChal.Solves = chal.Solves
				err = datastore.UpdateChallenge(context.Background(), *existingChal)
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
