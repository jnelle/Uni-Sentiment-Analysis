package main

import (
	"context"
	"log"

	"github.com/cdipaolo/sentiment"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"semantic.analysis.fom/lib"
	"semantic.analysis.fom/models"
)

type Pair struct {
	Key       string
	Value     int
	Sentiment string
}

type PairList []Pair

var IMDBMovieList []*models.IMDBMovie

func main() {
	// 	Restore restores a pre-trained models from a binary asset this is the preferable method of generating a model
	// This basically wraps RestoreModels.
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	// get all movies
	allMovies, err := getAllMovies()
	if err != nil {
		log.Fatalln(err)
	}
	errGrp, ctx := errgroup.WithContext(context.Background())

	for i := range allMovies {
		currentMovie := allMovies[i]
		errGrp.Go(func() error {
			return doSentimentAnalysis(currentMovie, sentimentModel, ctx)
		})

	}
	err = errGrp.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// iterate through genreMap and calculate score
func doSentimentAnalysis(currentMovie *models.IMDBMovie, sentimentModel sentiment.Models, ctx context.Context) error {
	mongo := lib.MongoDBProcessedReviews()

	currentMovieComments, err := getSpecificMovieComments(currentMovie.ID)
	if err != nil {
		return err
	}
	for n := range currentMovieComments {
		currentComment := currentMovieComments[n]
		results := sentimentModel.SentimentAnalysis(currentComment.ReviewContent, sentiment.English)
		currentComment.SentimentRating = results.Score
		_, err = mongo.InsertOne(ctx, currentComment)
		if err != nil {
			return err
		}
	}
	return nil

}

// get all movies from mongodb
func getAllMovies() ([]*models.IMDBMovie, error) {
	cursor, err := lib.MongoDBIMDBCollection().Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	var allMovies []*models.IMDBMovie
	err = cursor.All(context.TODO(), &allMovies)
	if err != nil {
		return nil, err
	}
	return allMovies, nil
}

// get comments for specific movie from mongodb
func getSpecificMovieComments(imdbID string) ([]*models.IMDBComments, error) {
	cursor, err := lib.MongoDBCommentsCollection().Find(context.TODO(), bson.M{"id": imdbID}, options.Find())
	if err != nil {
		return nil, err
	}

	var movieComments []*models.IMDBComments
	err = cursor.All(context.TODO(), &movieComments)
	if err != nil {
		return nil, err
	}

	return movieComments, nil
}

// tokenize and remove stopwords from comment
// func tokenizeAndRemoveStopWordsFromComment(comment string) []string {
// 	// extract single comment without any special character
// 	r, _ := regexp.Compile(`[\p{P}\p{S}]+`)
// 	stopwords.LoadStopWordsFromFile("stopwords_new.txt", "en", "\n")
// 	// delete whitespaces
// 	cleanedSpecialCharsComment := r.ReplaceAllString(comment, " ")
// 	cleanedComment := stopwords.CleanString(cleanedSpecialCharsComment, "en", false)
// 	tokens := t.Tokenize(cleanedComment)

// 	return tokens
// }
