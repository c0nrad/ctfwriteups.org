package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/c0nrad/ctfwriteups/config"
	"github.com/c0nrad/ctfwriteups/datastore"
	"github.com/c0nrad/ctfwriteups/models"
	"github.com/mailgun/mailgun-go/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const EDITION = 5

func main() {
	env := "prod"

	config.InitLogger()
	config.InitEnv(env)
	datastore.InitDatabase()

	newsletter := GenerateNewsletter(context.Background())
	fmt.Printf("%+v\n", newsletter)

	users, err := datastore.GetNewsletterSubscriptions(context.Background())
	if err != nil {
		panic(err)
	}

	newsletter.SentCount = len(users)

	// fmt.Println("Subject: ", subject)
	// fmt.Println(template)
	// fmt.Println("Sending email to ", len(users), " users. 10 seconds.")
	// time.Sleep(10 * time.Second)

	// SendEmail(context.Background(), []string{"c0nrad@c0nrad.io"

	// err = SendEmail(context.Background(), []string{"c0nrad@c0nrad.io"}, "CTFWriteups@m.ctfwriteups.org", newsletter.Subject, newsletter.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	for _, user := range users {
		fmt.Println(user.Email)
		err = SendEmail(context.Background(), []string{user.Email}, "CTFWriteups@m.ctfwriteups.org", newsletter.Subject, newsletter.Body)
		if err != nil {
			fmt.Println(err)
		}
	}

	_, err = datastore.Client.Database(datastore.DB).Collection(datastore.NewsletterCollection).InsertOne(context.Background(), newsletter)
	if err != nil {
		panic(err)
	}

}

func SendEmail(ctx context.Context, to []string, from, subject, body string) error {

	if !strings.Contains(from, "@") {
		from += "@" + config.GetConfig("MAILGUN_DOMAIN")
	}

	mg := mailgun.NewMailgun("m.ctfwriteups.org", config.GetConfig("MAILGUN_API_KEY"))
	message := mg.NewMessage(from, subject, body, to...)
	// message.Set(body)
	message.SetReplyTo("c0nrad@c0nrad.io")

	var err error

	var mes, id string
	mes, id, err = mg.Send(ctx, message)

	fmt.Println(mes, id, err)
	return err
}

func GenerateNewsletter(ctx context.Context) models.Newsletter {

	writeups, err := datastore.GetWriteups(context.Background())
	if err != nil {
		panic(err)
	}

	recentWriteups := []models.Writeup{}
	ctfs := []models.CTF{}
	for _, writeup := range writeups {
		if writeup.TS.After(time.Now().Add(-7 * 24 * time.Hour)) {
			recentWriteups = append(recentWriteups, writeup)
		}
	}

	cursor, err := datastore.Client.Database(datastore.DB).Collection(datastore.CTFCollection).Find(ctx, bson.M{"enddate": bson.M{"$gte": time.Now().AddDate(0, 0, -7).Unix()}})
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, &ctfs)
	if err != nil {
		panic(err)
	}

	// sort by votes
	sort.Slice(recentWriteups, func(i, j int) bool {
		return recentWriteups[i].VoteCount > recentWriteups[j].VoteCount
	})

	winners := recentWriteups[:5]

	newsletter := models.Newsletter{
		ID:        primitive.NewObjectID(),
		TS:        time.Now(),
		IsSent:    false,
		SentCount: 0,

		Edition:       EDITION,
		Writeups:      winners,
		CTFs:          ctfs,
		TotalWriteups: len(recentWriteups),
	}

	newsletter.Subject = fmt.Sprintf("CTFWritups.org Newsletter #%d: ", newsletter.Edition)
	for _, ctf := range newsletter.CTFs {
		newsletter.Subject += ctf.Name + ", "
	}
	newsletter.Subject = newsletter.Subject[:len(newsletter.Subject)-2]

	newsletter.Body = `Welcome to CTFWriteups.org newsletter #{{.Edition}}!

	This week's winners are:
	
	{{range $i, $writeup := .Writeups}}
	{{ $i }}.) {{$writeup.CTFName}} {{$writeup.Category}}/{{$writeup.ChallengeName}}
	{{$writeup.URL}}
	{{end}}
	
	You can find all the writeups for this week at ctfwriteups.org. If you're no longer interested in receiving these emails, you can unsubscribe at https://ctfwriteups.org/unsubscribe.
	
	Happy Hacking!
	c0nrad - Sloppy Joe Pirates`

	t, err := template.New("newsletter.tmpl").ParseFiles("newsletter.tmpl")
	if err != nil {
		panic(err)
	}

	var b strings.Builder
	err = t.Execute(&b, newsletter)
	if err != nil {
		panic(err)
	}
	newsletter.Body = b.String()

	return newsletter
}
