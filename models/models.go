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
	ObjectID primitive.ObjectID `bson:"_id,omitempty"`
	ID       string             `json:"id" bson:"id"`
	Rate     string             `json:"rate" bson:"rate"`
	Helpful  string             `json:"helpful" bson:"helpful"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	Date     string             `json:"date" bson:"date"`
}

type ChartratingResponse []struct {
	ID          string  `json:"id"`
	ChartRating float64 `json:"chartRating"`
}

type Top250Response struct {
	Items        []Items `json:"items"`
	ErrorMessage string  `json:"errorMessage"`
}
type Items struct {
	ID              string `json:"id"`
	Rank            string `json:"rank"`
	Title           string `json:"title"`
	FullTitle       string `json:"fullTitle"`
	Year            string `json:"year"`
	Image           string `json:"image"`
	Crew            string `json:"crew"`
	ImDbRating      string `json:"imDbRating"`
	ImDbRatingCount string `json:"imDbRatingCount"`
}

type CommentsResponse struct {
	ImDbID       string          `json:"imDbId"`
	Title        string          `json:"title"`
	FullTitle    string          `json:"fullTitle"`
	Type         string          `json:"type"`
	Year         string          `json:"year"`
	Items        []CommentsItems `json:"items"`
	ErrorMessage string          `json:"errorMessage"`
}

type CommentsItems struct {
	Username        string `json:"username"`
	UserURL         string `json:"userUrl"`
	ReviewLink      string `json:"reviewLink"`
	WarningSpoilers bool   `json:"warningSpoilers"`
	Date            string `json:"date"`
	Rate            string `json:"rate"`
	Helpful         string `json:"helpful"`
	Title           string `json:"title"`
	Content         string `json:"content"`
}

type MetadataResponse struct {
	ID               string         `json:"id"`
	Title            string         `json:"title"`
	OriginalTitle    string         `json:"originalTitle"`
	FullTitle        string         `json:"fullTitle"`
	Type             string         `json:"type"`
	Year             string         `json:"year"`
	Image            string         `json:"image"`
	ReleaseDate      string         `json:"releaseDate"`
	RuntimeMins      string         `json:"runtimeMins"`
	RuntimeStr       string         `json:"runtimeStr"`
	Plot             string         `json:"plot"`
	PlotLocal        string         `json:"plotLocal"`
	PlotLocalIsRtl   bool           `json:"plotLocalIsRtl"`
	Awards           string         `json:"awards"`
	Directors        string         `json:"directors"`
	DirectorList     []DirectorList `json:"directorList"`
	Writers          string         `json:"writers"`
	WriterList       []WriterList   `json:"writerList"`
	Stars            string         `json:"stars"`
	StarList         []StarList     `json:"starList"`
	ActorList        []ActorList    `json:"actorList"`
	FullCast         interface{}    `json:"fullCast"`
	Genres           string         `json:"genres"`
	GenreList        []GenreList    `json:"genreList"`
	Companies        string         `json:"companies"`
	CompanyList      []CompanyList  `json:"companyList"`
	Countries        string         `json:"countries"`
	CountryList      []CountryList  `json:"countryList"`
	Languages        string         `json:"languages"`
	LanguageList     []LanguageList `json:"languageList"`
	ContentRating    string         `json:"contentRating"`
	ImDbRating       string         `json:"imDbRating"`
	ImDbRatingVotes  string         `json:"imDbRatingVotes"`
	MetacriticRating string         `json:"metacriticRating"`
	Ratings          interface{}    `json:"ratings"`
	Wikipedia        Wikipedia      `json:"wikipedia"`
	Posters          interface{}    `json:"posters"`
	Images           interface{}    `json:"images"`
	Trailer          interface{}    `json:"trailer"`
	BoxOffice        BoxOffice      `json:"boxOffice"`
	Tagline          string         `json:"tagline"`
	Keywords         string         `json:"keywords"`
	KeywordList      []string       `json:"keywordList"`
	Similars         []Similars     `json:"similars"`
	TvSeriesInfo     interface{}    `json:"tvSeriesInfo"`
	TvEpisodeInfo    interface{}    `json:"tvEpisodeInfo"`
	ErrorMessage     interface{}    `json:"errorMessage"`
}

type DirectorList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WriterList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StarList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ActorList struct {
	ID          string `json:"id"`
	Image       string `json:"image"`
	Name        string `json:"name"`
	AsCharacter string `json:"asCharacter"`
}

type GenreList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CompanyList struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CountryList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type LanguageList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PlotShort struct {
	PlainText string `json:"plainText"`
	HTML      string `json:"html"`
}

type PlotFull struct {
	PlainText string `json:"plainText"`
	HTML      string `json:"html"`
}

type Wikipedia struct {
	ImDbID          string    `json:"imDbId"`
	Title           string    `json:"title"`
	FullTitle       string    `json:"fullTitle"`
	Type            string    `json:"type"`
	Year            string    `json:"year"`
	Language        string    `json:"language"`
	TitleInLanguage string    `json:"titleInLanguage"`
	URL             string    `json:"url"`
	PlotShort       PlotShort `json:"plotShort"`
	PlotFull        PlotFull  `json:"plotFull"`
	ErrorMessage    string    `json:"errorMessage"`
}

type BoxOffice struct {
	Budget                   string `json:"budget"`
	OpeningWeekendUSA        string `json:"openingWeekendUSA"`
	GrossUSA                 string `json:"grossUSA"`
	CumulativeWorldwideGross string `json:"cumulativeWorldwideGross"`
}

type Similars struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Image      string `json:"image"`
	ImDbRating string `json:"imDbRating"`
}

type TitleMapData struct {
	Genres  []string
	Title   string
	Rank    string
	Rating  string
	Comment string
}

type WordCounts struct {
	Word    string
	Amount  int
	MovieID string
}

type CSVModel struct {
	Word      string `csv:"word"`
	Amount    int    `csv:"amount"`
	Genre     string `csv:"genre"`
	Sentiment string `csv:"sentiment"`
}
