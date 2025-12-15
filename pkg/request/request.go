package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	scheme string
	host   string
	client *http.Client
	tr     *http.Transport
	debug  bool
}

func NewClient(scheme string, host string, timeout time.Duration, opts ...ReqOpt) *Client {
	req := &Client{
		scheme: scheme,
		host:   host,
		client: &http.Client{
			Timeout: timeout,
		},
		debug: false,
	}

	for _, opt := range opts {
		opt(req)
	}

	if req.tr != nil {
		req.client.Transport = req.tr
	}

	return req
}

func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

func (c *Client) SetTransport(tr *http.Transport) {
	c.client.Transport = tr
}

func sendRequest[T any](c *Client, method, path string, opts ...Opt) (*T, error) {
	ctx := &Ctx{}
	rid := uuid.NewString()

	for _, opt := range opts {
		opt(ctx)
	}

	urlStr := buildURL(c, path, ctx.query)
	if c.debug {
		log.Printf("[REQ:%s] url: %s", rid, urlStr)
	}

	body, contentType, err := buildBodyAndType(ctx)
	if err != nil {
		return nil, err
	}

	if c.debug && ctx.body != nil {
		bs, _ := json.Marshal(ctx.body)
		buf := &bytes.Buffer{}
		if err := json.Indent(buf, bs, "", "  "); err != nil {
			log.Printf("[REQ:%s] body: %s", rid, string(bs))
		} else {
			log.Printf("[REQ:%s] body: %s", rid, buf.String())
		}
	}

	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	for k, v := range ctx.header {
		req.Header.Add(k, v)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.debug {
		log.Printf("[REQ:%s] headers: %+v", rid, req.Header)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if c.debug {
		buf := &bytes.Buffer{}
		if err := json.Indent(buf, b, "", "  "); err != nil {
			log.Printf("[REQ:%s] resp: %s", rid, string(b))
		} else {
			log.Printf("[REQ:%s] resp: %s", rid, buf.String())
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var rr T
	if err := json.Unmarshal(b, &rr); err != nil {
		return nil, err
	}
	return &rr, nil
}

func buildURL(c *Client, path string, q Query) string {
	u := url.URL{Scheme: c.scheme, Host: c.host, Path: path}
	if len(q) > 0 {
		values := u.Query()
		for k, v := range q {
			values.Add(k, v)
		}
		u.RawQuery = values.Encode()
	}
	return u.String()
}

func buildBodyAndType(ctx *Ctx) (io.Reader, string, error) {
	if ctx.body == nil {
		return nil, "", nil
	}
	bs, err := json.Marshal(ctx.body)
	if err != nil {
		return nil, "", err
	}
	switch ctx.contentType {
	case "multipart/form-data":
		m := make(map[string]string)
		if err := json.Unmarshal(bs, &m); err != nil {
			return nil, "", err
		}
		buf := &bytes.Buffer{}
		w := multipart.NewWriter(buf)
		for k, v := range m {
			if err := w.WriteField(k, v); err != nil {
				return nil, "", err
			}
		}
		if err := w.Close(); err != nil {
			return nil, "", err
		}
		return buf, w.FormDataContentType(), nil
	case "application/x-www-form-urlencoded":
		m := make(map[string]string)
		if err := json.Unmarshal(bs, &m); err != nil {
			return nil, "", err
		}
		data := url.Values{}
		for k, v := range m {
			data.Add(k, v)
		}
		return strings.NewReader(data.Encode()), "application/x-www-form-urlencoded", nil
	default:
		return bytes.NewBuffer(bs), "application/json", nil
	}
}

func Get[T any](c *Client, path string, opts ...Opt) (*T, error) {
	return sendRequest[T](c, http.MethodGet, path, opts...)
}

func Post[T any](c *Client, path string, body any, opts ...Opt) (*T, error) {
	opts = append(opts, WithBody(body))
	return sendRequest[T](c, http.MethodPost, path, opts...)
}

func Put[T any](c *Client, path string, body any, opts ...Opt) (*T, error) {
	opts = append(opts, WithBody(body))
	return sendRequest[T](c, http.MethodPut, path, opts...)
}

func Delete[T any](c *Client, path string, opts ...Opt) (*T, error) {
	return sendRequest[T](c, http.MethodDelete, path, opts...)
}

func GetHeaderMap(header string) map[string]string {
	headerMap := make(map[string]string)
	for _, h := range strings.Split(header, "\n") {
		if key, value, ok := strings.Cut(h, "="); ok {
			headerMap[key] = value
		}
	}
	return headerMap
}
