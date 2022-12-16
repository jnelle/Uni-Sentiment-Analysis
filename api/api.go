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
	// AcquireRequest returns an empty Request instance from request pool.
	// The returned Request instance may be passed to ReleaseRequest when it is no longer needed. This allows Request recycling, reduces GC pressure and usually improves performance.
	req := fasthttp.AcquireRequest()

	// SetRequestURI sets RequestURI.
	req.SetRequestURI("/en/API/Top250Movies/" + os.Getenv("APIKEY"))

	// SetScheme sets URI scheme, i.e. http, https, ftp, etc.
	req.URI().SetScheme("https")

	// SetHost sets host for the request.
	req.SetHost(imdbBaseURL)

	// SetUserAgent sets User-Agent header value-
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	// SetMethodBytes sets HTTP request method.
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch top movies\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		return nil, errors.New(errorMessage)
	}

	// deserialize json to golang object
	var topMovies *models.Top250Response
	err = sonic.Unmarshal(response, &topMovies)
	if err != nil {
		return nil, err
	}

	return topMovies, nil
}

func GetCommentsIMDB(titleID string) (*models.CommentsResponse, error) {
	// AcquireRequest returns an empty Request instance from request pool.
	// The returned Request instance may be passed to ReleaseRequest when it is no longer needed. This allows Request recycling, reduces GC pressure and usually improves performance.
	req := fasthttp.AcquireRequest()

	// SetRequestURI sets RequestURI.
	req.SetRequestURI("/en/API/Reviews/" + os.Getenv("APIKEY") + "/" + titleID + "?limit=2500")

	// SetScheme sets URI scheme, i.e. http, https, ftp, etc.
	req.URI().SetScheme("https")

	// SetHost sets host for the request.
	req.SetHost(imdbBaseURL)

	// SetUserAgent sets User-Agent header value-
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	// SetMethodBytes sets HTTP request method.
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	// start http request
	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch chartrating\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		panic(errorMessage)
	}

	// deserialize json to golang object
	var comments *models.CommentsResponse
	err = sonic.Unmarshal(response, &comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func GetMovieMetaData(titleID string) (*models.MetadataResponse, error) {
	// AcquireRequest returns an empty Request instance from request pool.
	// The returned Request instance may be passed to ReleaseRequest when it is no longer needed. This allows Request recycling, reduces GC pressure and usually improves performance.
	req := fasthttp.AcquireRequest()

	// SetRequestURI sets RequestURI.
	req.SetRequestURI("/en/API/Title/" + os.Getenv("APIKEY") + "/" + titleID + "/Wikipedia,")

	// SetScheme sets URI scheme, i.e. http, https, ftp, etc.
	req.URI().SetScheme("https")

	// SetHost sets host for the request.
	req.SetHost(imdbBaseURL)

	// SetUserAgent sets User-Agent header value-
	req.Header.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	// SetMethodBytes sets HTTP request method.
	req.Header.SetMethodBytes([]byte(fasthttp.MethodGet))

	response, header, err := doRequest(req)
	if err != nil || header.StatusCode() != 200 {
		errorMessage := fmt.Sprintf("Couldn't fetch movie metadata\tstatus code: %d\t error: %v", header.StatusCode(), string(response))
		panic(errorMessage)
	}

	// deserialize json to golang object
	var metadata *models.MetadataResponse
	err = sonic.Unmarshal(response, &metadata)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func doRequest(req *fasthttp.Request) ([]byte, fasthttp.ResponseHeader, error) {

	// 	ReleaseRequest returns req acquired via AcquireRequest to request pool.
	// It is forbidden accessing req and/or its' members after returning it to request pool.
	defer fasthttp.ReleaseRequest(req)

	// 	AcquireResponse returns an empty Response instance from response pool.
	// The returned Response instance may be passed to ReleaseResponse when it is no longer needed. This allows Response recycling, reduces GC pressure and usually improves performance.
	resp := fasthttp.AcquireResponse()

	// 	ReleaseResponse return resp acquired via AcquireResponse to response pool.
	// It is forbidden accessing resp and/or its' members after returning it to response pool.
	defer fasthttp.ReleaseResponse(resp)

	// 	Do performs the given http request and fills the given http response.
	// Request must contain at least non-zero RequestURI with full url (including scheme and host) or non-zero Host header + RequestURI.
	// Client determines the server to be requested in the following order:
	// - from RequestURI if it contains full url with scheme and host;
	// - from Host header otherwise.
	err := client.Do(req, resp)
	if err != nil {
		return nil, fasthttp.ResponseHeader{}, err

	}
	// Body returns response body.
	return resp.Body(), resp.Header, nil
}
