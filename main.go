package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/bbalet/stopwords"
	"github.com/cdipaolo/sentiment"
	"github.com/euskadi31/go-tokenizer"
	_ "github.com/joho/godotenv/autoload"
	"github.com/jszwec/csvutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"semantic.analysis.fom/api"
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
var t = tokenizer.New()

func main() {
	// 	Restore restores a pre-trained models from a binary asset this is the preferable method of generating a model
	// This basically wraps RestoreModels.
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	// The make built-in function allocates and initializes an object of type slice, map, or chan
	multipleGenreMap := make(map[string]PairList)

	// The make built-in function allocates and initializes an object of type slice, map, or chan
	singleGenreMap := make(map[string]PairList)

	// map which contains counted words for each genre
	// ex:
	// {
	//	"Drama": {
	//		"Hello": 2,
	// 		}
	// }
	singleGenreWordMap := make(map[string]map[string]int)

	// map which contains counted words for each genre
	// ex:
	// {
	//	"Drama": {
	//		"Hello": 2,
	// 		}
	// }
	multipleGenreWordMap := make(map[string]map[string]int)

	// get all movies
	allMovies, err := getAllMovies()
	if err != nil {
		log.Fatalln(err)
	}

	// hypothesis 1,2,3 count words for each genre
	countGenreWords(allMovies, singleGenreWordMap, multipleGenreWordMap)

	// sort counted words ascending
	sortGenreMaps(singleGenreWordMap, singleGenreMap)
	sortGenreMaps(multipleGenreWordMap, multipleGenreMap)

	// start sentiment analysis
	doSentimentAnalysis(singleGenreMap, sentimentModel)
	doSentimentAnalysis(multipleGenreMap, sentimentModel)

	writeToCSV(singleGenreMap, "output_hypothese12.csv")
	writeToCSV(multipleGenreMap, "output_hypothese3.csv")

}

func doSentimentAnalysis(genreMap map[string]PairList, sentimentModel sentiment.Models) {
	for key := range genreMap {
		for x := range genreMap[key] {
			results := sentimentModel.SentimentAnalysis(genreMap[key][x].Key, sentiment.English)
			genreMap[key][x].Sentiment = scoreToString(results.Score)
		}
	}
}

func sortGenreMaps(genreWordMap map[string]map[string]int, genreMap map[string]PairList) {
	for key, val := range genreWordMap {
		for k, v := range val {
			genreMap[key] = append(genreMap[key], Pair{Key: k, Value: v})
		}
		sort.Slice(genreMap[key], func(i, j int) bool {
			return genreMap[key][i].Value > genreMap[key][j].Value
		})
	}
}

func countGenreWords(allMovies []*models.IMDBMovie, singleGenreWordMap map[string]map[string]int, multipleGenreWordMap map[string]map[string]int) {
	for x := range allMovies {
		tmpMovieData := &models.TitleMapData{
			Rank:   allMovies[x].Rank,
			Title:  allMovies[x].Title,
			Rating: allMovies[x].Rating,
		}
		for y := range allMovies[x].Genres {
			tmpMovieData.Genres = append(tmpMovieData.Genres, allMovies[x].Genres[y].Value)
		}

		if len(allMovies[x].Genres) == 1 {
			movieComments, _ := getSpecificMovieComments(allMovies[x].ID)
			for i := range movieComments {
				countedWords := countWords(tokenizeAndRemoveStopWordsFromComment(movieComments[i].Content))
				for key, value := range countedWords {
					if _, ok := singleGenreWordMap[allMovies[x].Genres[0].Value]; !ok {
						singleGenreWordMap[allMovies[x].Genres[0].Value] = make(map[string]int)
					}
					singleGenreWordMap[allMovies[x].Genres[0].Value][key] = singleGenreWordMap[allMovies[x].Genres[0].Value][key] + value
				}
			}
		}

		if len(allMovies[x].Genres) > 1 {
			movieComments, _ := getSpecificMovieComments(allMovies[x].ID)
			for i := range movieComments {
				countedWords := countWords(tokenizeAndRemoveStopWordsFromComment(movieComments[i].Content))
				for key, value := range countedWords {
					if _, ok := multipleGenreWordMap[allMovies[x].Genres[0].Value]; !ok {
						multipleGenreWordMap[allMovies[x].Genres[0].Value] = make(map[string]int)
					}
					multipleGenreWordMap[allMovies[x].Genres[0].Value][key] = multipleGenreWordMap[allMovies[x].Genres[0].Value][key] + value
				}
			}

		}
	}
}

func writeDataToDB() {
	result, err := api.GetTopMovies()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, y := range result.Items {
		wg.Add(1)
		go func(x models.Items) {
			defer wg.Done()

			metadata, err := api.GetMovieMetaData(x.ID)
			if err != nil {
				log.Fatalln(err)
			}
			movie := &models.IMDBMovie{
				ID:                       x.ID,
				Rating:                   x.ImDbRating,
				RatingCount:              x.ImDbRatingCount,
				Rank:                     x.Rank,
				Year:                     x.Year,
				Crew:                     x.Crew,
				Title:                    x.Title,
				Budget:                   metadata.BoxOffice.Budget,
				CumulativeWorldwideGross: metadata.BoxOffice.CumulativeWorldwideGross,
				Genres:                   metadata.GenreList,
				ReleaseDate:              metadata.ReleaseDate,
				RunTime:                  metadata.RuntimeMins,
			}
			IMDBMovieList = append(IMDBMovieList, movie)

			_, err = lib.MongoDBGetIMDBCollection().InsertOne(context.TODO(), movie, nil)
			if err != nil {
				panic(err)
			}
		}(y)
	}
	wg.Wait()

	for x := range IMDBMovieList {
		comments, err := api.GetCommentsIMDB(IMDBMovieList[x].ID)
		if err != nil {
			log.Fatalln(err)
		}

		for i := range comments.Items {
			comment := &models.IMDBComments{
				ID:      comments.ImDbID,
				Rate:    comments.Items[i].Rate,
				Helpful: comments.Items[i].Helpful,
				Title:   comments.Items[i].Title,
				Content: comments.Items[i].Content,
				Date:    comments.Items[i].Date,
			}
			_, err := lib.MongoDBGetCommentsCollection().InsertOne(context.TODO(), comment, nil)
			if err != nil {
				log.Fatalln(err)
			}
		}

	}
}

func getAllMovies() ([]*models.IMDBMovie, error) {
	cursor, err := lib.MongoDBGetIMDBCollection().Find(context.TODO(), bson.D{{}}, options.Find())
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

func getAllComments() ([]*models.IMDBComments, error) {
	cursor, err := lib.MongoDBGetCommentsCollection().Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	var commentResult []*models.IMDBComments
	err = cursor.All(context.TODO(), &commentResult)
	if err != nil {
		return nil, err
	}

	return commentResult, nil
}

func getSpecificMovieComments(imdbID string) ([]*models.IMDBComments, error) {
	cursor, err := lib.MongoDBGetCommentsCollection().Find(context.TODO(), bson.M{"id": imdbID}, options.Find())
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

func tokenizeAndRemoveStopWordsFromComment(comment string) []string {
	r, _ := regexp.Compile(`[\p{P}\p{S}]+`)
	stopwords.LoadStopWordsFromFile("stopwords_new.txt", "en", "\n")
	cleanedSpecialCharsComment := r.ReplaceAllString(comment, " ")
	cleanedComment := stopwords.CleanString(cleanedSpecialCharsComment, "en", false)
	tokens := t.Tokenize(cleanedComment)

	// doc, err := prose.NewDocument(cleanedComment)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return tokens
}

func countWords(tokens []string) map[string]int {
	tmpWordCountMap := make(map[string]int, 1)
	for _, token := range tokens {
		// skip single characters
		// skip strings with special characters
		if len(token) == 1 || strings.Contains("[$&+,:;=?@#|'<>.-^*()%!]\"", token) {
			continue
		}

		tmpWordCountMap[token] = tmpWordCountMap[token] + 1
	}
	return tmpWordCountMap
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// write sentiment results to csv
func writeToCSV(pairMap map[string]PairList, filename string) {
	var csvData []*models.CSVModel

	for key := range pairMap {

		for i := range pairMap[key] {
			pair := &models.CSVModel{
				Genre: key,
			}
			pair.Amount = pairMap[key][i].Value
			pair.Word = pairMap[key][i].Key
			pair.Sentiment = pairMap[key][i].Sentiment
			csvData = append(csvData, pair)
			pair = nil
		}
	}

	b, err := csvutil.Marshal(csvData)
	if err != nil {
		log.Fatalln(err)
	}

	destination, err := os.Create(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer destination.Close()

	destination.Write(b)

}

func scoreToString(sentimentScore uint8) string {
	var result string

	switch sentimentScore {
	case 1:
		result = "POSITIV"
	case 0:
		result = "NEGATIV"
	}
	return result
}
