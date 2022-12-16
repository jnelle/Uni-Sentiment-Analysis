package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type IMDBMovie struct {
	ID                       string      `json:"id" bson:"_id"`
	Rating                   string      `bson:"rating"`
	RatingCount              string      `bson:"rating_count"`
	Title                    string      `bson:"title"`
	Crew                     string      `bson:"crew"`
	Rank                     string      `bson:"rank"`
	Year                     string      `bson:"year"`
	Budget                   string      `bson:"budget"`
	CumulativeWorldwideGross string      `bson:"cumulative_worldwide_gross"`
	Genres                   []GenreList `bson:"genre_list"`
	ReleaseDate              string      `bson:"release_date"`
	RunTime                  string      `bson:"runtime"`
}

type IMDBComments struct {
	ObjectID        primitive.ObjectID `bson:"_id,omitempty"`
	ID              string             `json:"id" bson:"id"`
	Rate            string             `json:"rate" bson:"rate"`
	Helpfulness     string             `json:"helpfulness" bson:"helpfulness"`
	Title           string             `json:"title" bson:"review_title"`
	ReviewContent   string             `json:"review" bson:"review"`
	Date            string             `json:"date" bson:"date"`
	Rating          uint8              `json:"rating" bson:"rating"`
	Author          string             `json:"author" bson:"author"`
	Year            uint16             `json:"year" bson:"year"`
	MovieName       string             `json:"name" bson:"name"`
	ReviewURL       string             `json:"url" bson:"url"`
	SentimentRating uint8              `json:"sentiment_rating" bson:"sentiment_rating"`
}

type GenreList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
