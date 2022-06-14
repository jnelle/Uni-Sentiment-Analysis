package api

import (
	"errors"
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"semantic.analysis.fom/models"
)

var client = fasthttp.Client{
	Dial: fasthttpproxy.FasthttpHTTPDialer(os.Getenv("PROXY")),
}

const imdbBaseURL string = "imdb-api.com"

func GetTopMovies() (*models.Top250Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("/en/API/Top250Movies/" + os.Getenv("APIKEY"))
	req.URI().SetScheme("https")
	req.SetHost(imdbBaseURL)

	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch chartrating\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		return nil, errors.New(errorMessage)
	}

	var topMovies *models.Top250Response
	err = sonic.Unmarshal(response, &topMovies)
	if err != nil {
		return nil, err
	}

	return topMovies, nil
}

func GetCommentsIMDB(titleID string) (*models.CommentsResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("/en/API/Reviews/" + os.Getenv("APIKEY") + "/" + titleID + "?limit=2500")
	req.URI().SetScheme("https")
	req.SetHost(imdbBaseURL)

	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch chartrating\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		panic(errorMessage)
	}

	var comments *models.CommentsResponse
	err = sonic.Unmarshal(response, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetMovieMetaData(titleID string) (*models.MetadataResponse, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("/en/API/Title/" + os.Getenv("APIKEY") + "/" + titleID + "/Wikipedia,")
	req.URI().SetScheme("https")
	req.SetHost(imdbBaseURL)

	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch chartrating\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		panic(errorMessage)
	}

	var metadata *models.MetadataResponse
	err = sonic.Unmarshal(response, &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func doRequest(req *fasthttp.Request) ([]byte, fasthttp.ResponseHeader, error) {

	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, fasthttp.ResponseHeader{}, err

	}
	return resp.Body(), resp.Header, nil
}
