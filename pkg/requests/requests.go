package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
)

const (
	ContentTypeJSON = "application/json; charset=utf-8"                  // content type json
	ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8" // context type form
	ContentTypeText = "text/plain; charset=utf-8"                        // content type text

	HeaderXRequestID = "X-Request-ID" // header field for request id
)

var (
	// defaultClient default http client with some optimize connection configuration
	defaultClient = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     time.Second * 10,
		},
	}
)

// request hold request context / params / data
type requests struct {
	client  *http.Client
	ctx     context.Context
	method  string
	headers map[string]string
	uri     string
	query   url.Values
	form    url.Values
	body    io.Reader
	retry   RetryStrategy
}

// Response hold response context & data with debug information
type Response struct {
	err      error
	escape   time.Duration
	response *http.Response
	cnt      int
}

// New create new requests instance
func New() *requests {
	return &requests{
		client:  defaultClient,
		ctx:     context.Background(),
		method:  http.MethodGet,
		headers: make(map[string]string, 8),
		uri:     "",
	}
}

// WithClient replace default http client
func (t *requests) WithClient(client *http.Client) *requests {
	t.client = client
	return t
}

// WithContext set context
func (t *requests) WithContext(ctx context.Context) *requests {
	t.ctx = ctx
	return t
}

// Retry set retry strategy
func (t *requests) Retry(retry RetryStrategy) *requests {
	t.retry = retry
	return t
}

// Method set http request method
func (t *requests) Method(method string) *requests {
	t.method = method
	return t
}

// Uri set uri
func (t *requests) Uri(uri string) *requests {
	t.uri = uri
	return t
}

// Get set get method with uri
func (t *requests) Get(uri string) *requests {
	return t.Method(http.MethodGet).Uri(uri)
}

// Post set post method with uri
func (t *requests) Post(uri string) *requests {
	return t.Method(http.MethodPost).Uri(uri)
}

// Delete set delete method with uri
func (t *requests) Delete(uri string) *requests {
	return t.Method(http.MethodDelete).Uri(uri)
}

// Put set put method with uri
func (t *requests) Put(uri string) *requests {
	return t.Method(http.MethodPut).Uri(uri)
}

// Patch set patch method with uri
func (t *requests) Patch(uri string) *requests {
	return t.Method(http.MethodPatch).Uri(uri)
}

// Query set query params
func (t *requests) Query(query url.Values) *requests {
	t.query = query
	return t
}

// Form set form params
func (t *requests) Form(form url.Values) *requests {
	t.form = form
	return t
}

// JSONBody set json data as body and set request content type is JSON
func (t *requests) JSONBody(data interface{}) *requests {
	b, _ := json.Marshal(data)
	return t.ContentType(ContentTypeJSON).Data(b)
}

// Body set body as io reader from stream request
func (t *requests) Body(body io.Reader) *requests {
	t.body = body
	return t
}

// Data set body is raw bytes data
func (t *requests) Data(data []byte) *requests {
	return t.Body(bytes.NewReader(data))
}

// ContentType set content type
func (t *requests) ContentType(contentType string) *requests {
	return t.AddHeader("Content-Type", contentType)
}

// UserAgent set ua
func (t *requests) UserAgent(userAgent string) *requests {
	return t.AddHeader("User-Agent", userAgent)
}

// RequestId set request id pass to target (endpoint)
func (t *requests) RequestId(requestId string) *requests {
	return t.AddHeader(HeaderXRequestID, requestId)
}

// AddHeader add request header
func (t *requests) AddHeader(key, value string) *requests {
	t.headers[key] = value
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

	req, err := http.NewRequest(t.method, uri, body)
	if err != nil {
		return nil, err
	}

	if span := opentracing.SpanFromContext(t.ctx); span != nil {
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		_ = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	}

	if requestId := t.ctx.Value(HeaderXRequestID); requestId != nil {
		req.Header.Set("x-request-id", fmt.Sprint(requestId))
	}

	for k, v := range t.headers {
		req.Header.Set(k, v)
	}

	return req, nil
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
		r.cnt += 1
		rsp, err := t.client.Do(req)
		if err != nil {
			if t.retry != nil {
				backoff := t.retry.NextBackoff()
				if backoff > 0 {
					// release, prepare for next call
					if rsp != nil {
						_, _ = ioutil.ReadAll(rsp.Body)
						_ = rsp.Body.Close()
					}
					time.Sleep(backoff)
					continue
				}
			}
		}

		r.escape = time.Since(start)
		r.response = rsp
		if rsp != nil {
			r.response.Request.Body = io.NopCloser(t.body)
		}
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

	if t.response != nil && t.response.Request != nil {
		if b, err := httputil.DumpRequest(t.response.Request, body); err == nil {
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
		_, _ = io.Copy(ioutil.Discard, t.response.Body)
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
