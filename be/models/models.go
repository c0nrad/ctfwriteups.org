package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Writeup struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
	// Body  string

	// AuthorID    primitive.ObjectID `json:"authorID,omitempty"`
	SubmitterID primitive.ObjectID `json:"submitterID,omitempty"`

	VoteCount    int `json:"voteCount"`
	CommentCount int `json:"commentCount"`

	//db.writeups.updateMany({"ctfid": ObjectId("655f923e6d63b42d886d1c20")}, {"$set": {"ctfenddate": 1700351991}})
	CTFID      primitive.ObjectID `json:"ctfID,omitempty"`
	CTFName    string             `json:"ctfName,omitempty"`
	CTFEndDate int                `json:"ctfEndDate,omitempty"`

	ChallengeID       primitive.ObjectID `json:"challengeID,omitempty"`
	ChallengeName     string             `json:"challengeName,omitempty"`
	ChallengeCategory string             `json:"challengeCategory,omitempty"`

	Tags []string `json:"tags,omitempty"`
}

type Vote struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	WriteupID primitive.ObjectID `json:"writeupID,omitempty"`
	UserID    primitive.ObjectID `json:"userID,omitempty"`

	IsUpvote bool `json:"isUpvote,omitempty"`
}

type Seen struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	WriteupID primitive.ObjectID `json:"writeupID,omitempty"`
	UserID    primitive.ObjectID `json:"userID,omitempty"`

	Count int `json:"count,omitempty"`
}

type User struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`

	LoginMethod string `json:"login_method,omitempty"`

	IsAdmin     bool `json:"isAdmin,omitempty"`
	IsModerator bool `json:"isModerator,omitempty"`
}

type Comment struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	UserID   primitive.ObjectID `json:"userID,omitempty"`
	Username string             `json:"username,omitempty"`

	WriteupID       primitive.ObjectID `json:"writeupID,omitempty"`
	ParentCommentID primitive.ObjectID `json:"parentCommentID,omitempty"`

	Body  string `json:"body,omitempty"`
	Votes int    `json:"votes,omitempty"`
}

type Tag struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time

	Category string
	Value    string
}

type UserSession struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time

	UserID primitive.ObjectID
	Email  string

	Token string
}

type CTF struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	StartDate int `json:"startDate,omitempty"`
	EndDate   int `json:"endDate,omitempty"`

	SubmitterID primitive.ObjectID `json:"submitterID,omitempty"`

	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`

	Categories []string `json:"categories"`

	ChallengeCount int `json:"challengeCount"`
	WriteupCount   int `json:"writeupCount"`
}

type Challenge struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	SubmitterID primitive.ObjectID `json:"submitterID,omitempty"`

	CTFID primitive.ObjectID `json:"ctfID,omitempty"`

	Name     string `json:"name"`
	Category string `json:"category"`
	Solves   int    `json:"solves"`

	ShortDescription string   `json:"shortDescription,omitempty"`
	Tags             []string `json:"tags"`
}

type NewsletterSubscription struct {
	ID primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS time.Time          `json:"ts,omitempty"`

	Email      string `json:"email,omitempty"`
	IsVerified bool   `json:"isVerified,omitempty"`
}

type Newsletter struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	TS      time.Time          `json:"ts,omitempty"`
	Edition int                `json:"edition,omitempty"`

	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`

	IsSent    bool `json:"isSent,omitempty"`
	SentCount int  `json:"sentCount,omitempty"`

	TotalWriteups int       `json:"totalWriteups,omitempty"`
	Writeups      []Writeup `json:"winners,omitempty"`
	CTFs          []CTF     `json:"ctfs,omitempty"`
}
