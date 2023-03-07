package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"                  // content type json
	ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8" // context type form
	ContentTypeText = "text/plain; charset=utf-8"                        // content type text

	HeaderXRequestID = "X-Request-ID" // header field for request id
)

// Requests request interface
type Requests interface {
	// WithClient replace default http client
	WithClient(client *http.Client) Requests

	// WithContext set context
	WithContext(ctx context.Context) Requests

	// Retry set retry strategy
	Retry(retry RetryStrategy) Requests

	// Method set http request method
	Method(method string) Requests

	// Uri set uri
	Uri(uri string) Requests

	// Get set get method with uri
	Get(uri string) Requests

	// Post set post method with uri
	Post(uri string) Requests

	// Delete set delete method with uri
	Delete(uri string) Requests

	// Put set put method with uri
	Put(uri string) Requests

	// Patch set patch method with uri
	Patch(uri string) Requests

	// Query set query params
	Query(query url.Values) Requests

	// Form set form params
	Form(form url.Values) Requests

	// JSONBody set json data as body and set request content type is JSON
	JSONBody(data interface{}) Requests

	// Body set body as io reader from stream request
	Body(body io.Reader) Requests

	// Data set body is raw bytes data
	Data(data []byte) Requests

	// ContentType set content type
	ContentType(contentType string) Requests

	// UserAgent set ua
	UserAgent(userAgent string) Requests

	// RequestId set request id pass to target (endpoint)
	RequestId(requestId string) Requests

	// AddHeader add request header
	AddHeader(key, value string) Requests

	// Do requests
	Do() *Response
}

var (
	// defaultClient default http client with some optimize connection configuration
	defaultClient = &http.Client{
		Timeout: time.Second * 5,
		Transport: otelhttp.NewTransport(
			&http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     time.Second * 10,
			},
		),
	}
)

// request hold request context / params / data
type requests struct {
	client *http.Client
	ctx    context.Context
	method string
	header http.Header
	uri    string
	query  url.Values
	form   url.Values
	body   io.Reader
	retry  RetryStrategy
}

// Response hold response context & data with debug information
type Response struct {
	request  *http.Request
	err      error
	escape   time.Duration
	response *http.Response
	cnt      int
}

// New create new requests instance
func New() Requests {
	return &requests{
		client: defaultClient,
		ctx:    context.Background(),
		method: http.MethodGet,
		header: make(http.Header),
		uri:    "",
	}
}

// WithClient replace default http client
func (t *requests) WithClient(client *http.Client) Requests {
	c := *client
	c.Transport = otelhttp.NewTransport(c.Transport)
	t.client = &c
	return t
}

// WithContext set context
func (t *requests) WithContext(ctx context.Context) Requests {
	t.ctx = ctx
	return t
}

// Retry set retry strategy
func (t *requests) Retry(retry RetryStrategy) Requests {
	t.retry = retry
	return t
}

// Method set http request method
func (t *requests) Method(method string) Requests {
	t.method = method
	return t
}

// Uri set uri
func (t *requests) Uri(uri string) Requests {
	t.uri = uri
	return t
}

// Get set get method with uri
func (t *requests) Get(uri string) Requests {
	return t.Method(http.MethodGet).Uri(uri)
}

// Post set post method with uri
func (t *requests) Post(uri string) Requests {
	return t.Method(http.MethodPost).Uri(uri)
}

// Delete set delete method with uri
func (t *requests) Delete(uri string) Requests {
	return t.Method(http.MethodDelete).Uri(uri)
}

// Put set put method with uri
func (t *requests) Put(uri string) Requests {
	return t.Method(http.MethodPut).Uri(uri)
}

// Patch set patch method with uri
func (t *requests) Patch(uri string) Requests {
	return t.Method(http.MethodPatch).Uri(uri)
}

// Query set query params
func (t *requests) Query(query url.Values) Requests {
	t.query = query
	return t
}

// Form set form params
func (t *requests) Form(form url.Values) Requests {
	t.form = form
	return t
}

// JSONBody set json data as body and set request content type is JSON
func (t *requests) JSONBody(data interface{}) Requests {
	b, _ := json.Marshal(data)
	return t.ContentType(ContentTypeJSON).Data(b)
}

// Body set body as io reader from stream request
func (t *requests) Body(body io.Reader) Requests {
	t.body = body
	return t
}

// Data set body is raw bytes data
func (t *requests) Data(data []byte) Requests {
	return t.Body(bytes.NewBuffer(data))
}

// ContentType set content type
func (t *requests) ContentType(contentType string) Requests {
	return t.AddHeader("Content-Type", contentType)
}

// UserAgent set ua
func (t *requests) UserAgent(userAgent string) Requests {
	return t.AddHeader("User-Agent", userAgent)
}

// RequestId set request id pass to target (endpoint)
func (t *requests) RequestId(requestId string) Requests {
	return t.AddHeader(HeaderXRequestID, requestId)
}

// AddHeader add request header
func (t *requests) AddHeader(key, value string) Requests {
	t.header.Add(key, value)
	return t
}

func (t *requests) buildRequest() (*http.Request, error) {
	u, err := url.Parse(t.uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range t.query {
		for _, it := range v {
			q.Add(k, it)
		}
	}

	u.RawQuery = q.Encode()
	uri := u.String()

	var body io.Reader
	if len(t.form) > 0 {
		body = strings.NewReader(t.form.Encode())
	} else if t.body != nil {
		var buf bytes.Buffer
		body = io.TeeReader(t.body, &buf)
		defer func() { t.body = &buf }()
	}

	req, err := http.NewRequestWithContext(t.ctx, t.method, uri, body)
	if err != nil {
		return nil, err
	}

	if len(t.header) > 0 {
		req.Header = t.header
	}

	if requestId := t.ctx.Value(HeaderXRequestID); requestId != nil {
		req.Header.Set("x-request-id", fmt.Sprint(requestId))
	}

	return req, nil
}

func (t *requests) drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

// Do requests
func (t *requests) Do() *Response {
	start := time.Now()
	r := &Response{}
	for {
		req, err := t.buildRequest()
		r.err = err
		if err != nil {
			return r
		}
		var body io.ReadCloser
		body, req.Body, err = t.drainBody(req.Body)
		if err != nil {
			r.err = err
		}
		r.cnt += 1
		rsp, err := t.client.Do(req)
		if err != nil {
			if t.retry != nil {
				backoff := t.retry.NextBackoff()
				if backoff > 0 {
					// release, prepare for next call
					if rsp != nil {
						_, _ = io.ReadAll(rsp.Body)
						_ = rsp.Body.Close()
					}
					time.Sleep(backoff)
					continue
				}
			}
		}

		// copy back
		req.Body = body
		r.request = req

		r.escape = time.Since(start)
		r.response = rsp
		r.err = err
		if err != nil {
			return r
		}
		return r
	}
}

// Err get error
func (t *Response) Err() error {
	return t.err
}

// Dump dump request & response
func (t *Response) Dump(body bool) map[string]interface{} {
	dump := make(map[string]interface{})
	dump["escape"] = t.escape.Milliseconds()
	dump["cnt"] = t.cnt

	if t.err != nil {
		dump["error"] = t.err
	}

	if t.request != nil {
		if b, err := httputil.DumpRequestOut(t.request, body); err == nil {
			dump["request"] = string(b)
		}
	}

	if t.response != nil {
		if b, err := httputil.DumpResponse(t.response, body); err == nil {
			dump["response"] = string(b)
		}
	}

	return dump
}

// StatusCode get response status code
func (t *Response) StatusCode() int {
	if t.response != nil {
		return t.response.StatusCode
	}
	return 0
}

// Close release response
func (t *Response) Close() {
	if t.response != nil {
		_, _ = io.Copy(io.Discard, t.response.Body)
		_ = t.response.Body.Close()
	}
}

// RawResponse get raw response (http.Response)
func (t *Response) RawResponse() *http.Response {
	return t.response
}

// Header get header
func (t *Response) Header() http.Header {
	if t.response != nil {
		return t.response.Header
	}
	return nil
}

// JSON Unmarshal response as JSON
func (t *Response) JSON(obj interface{}) error {
	b, err := t.Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

// Text get response as plain text
func (t *Response) Text() (string, error) {
	b, err := t.Bytes()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Bytes get response as bytes data
func (t *Response) Bytes() ([]byte, error) {
	return io.ReadAll(t.response.Body)
}
