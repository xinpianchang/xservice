package requests

import (
	"bytes"
	"context"
	"crypto/tls"
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
	ContentTypeJSON = "application/json; charset=utf-8"
	ContentTypeForm = "application/x-www-form-urlencoded; charset=utf-8"
	ContentTypeText = "text/plain; charset=utf-8"

	HeaderXRequestID = "X-Request-ID"
)

var (
	defaultClient = &http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     time.Second * 10,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
	}
)

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

type Response struct {
	err      error
	escape   time.Duration
	response *http.Response
	cnt      int
}

func New() *requests {
	return &requests{
		client:  defaultClient,
		ctx:     context.Background(),
		method:  http.MethodGet,
		headers: make(map[string]string, 8),
		uri:     "",
	}
}

func (t *requests) WithClient(client *http.Client) *requests {
	t.client = client
	return t
}

func (t *requests) WithContext(ctx context.Context) *requests {
	t.ctx = ctx
	return t
}

func (t *requests) Retry(retry RetryStrategy) *requests {
	t.retry = retry
	return t
}

func (t *requests) Method(method string) *requests {
	t.method = method
	return t
}

func (t *requests) Uri(uri string) *requests {
	t.uri = uri
	return t
}

func (t *requests) Get(uri string) *requests {
	return t.Method(http.MethodGet).Uri(uri)
}

func (t *requests) Post(uri string) *requests {
	return t.Method(http.MethodPost).Uri(uri)
}

func (t *requests) Delete(uri string) *requests {
	return t.Method(http.MethodDelete).Uri(uri)
}

func (t *requests) Put(uri string) *requests {
	return t.Method(http.MethodPut).Uri(uri)
}

func (t *requests) Patch(uri string) *requests {
	return t.Method(http.MethodPatch).Uri(uri)
}

func (t *requests) Query(query url.Values) *requests {
	t.query = query
	return t
}

func (t *requests) Form(form url.Values) *requests {
	t.form = form
	return t
}

func (t *requests) JSONBody(data interface{}) *requests {
	b, _ := json.Marshal(data)
	return t.ContentType(ContentTypeJSON).Data(b)
}

func (t *requests) Body(body io.Reader) *requests {
	t.body = body
	return t
}

func (t *requests) Data(data []byte) *requests {
	return t.Body(bytes.NewReader(data))
}

func (t *requests) ContentType(contentType string) *requests {
	return t.AddHeader("Content-Type", contentType)
}

func (t *requests) UserAgent(userAgent string) *requests {
	return t.AddHeader("User-Agent", userAgent)
}

func (t *requests) RequestId(requestId string) *requests {
	return t.AddHeader(HeaderXRequestID, requestId)
}

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

func (t *Response) Err() error {
	return t.err
}

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

func (t *Response) StatusCode() int {
	return t.response.StatusCode
}

func (t *Response) Close() {
	if t.response != nil {
		_, _ = ioutil.ReadAll(t.response.Body)
		_ = t.response.Body.Close()
	}
}

func (t *Response) RawResponse() *http.Response {
	return t.response
}

func (t *Response) JSON(obj interface{}) error {
	b, err := io.ReadAll(t.response.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}
