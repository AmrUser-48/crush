package gemini

import (
	"context"
	"net/http"

	"charm.land/fantasy"
)

type callUAKey struct{}
type callHeadersKey struct{}

func withCallUA(ctx context.Context, call fantasy.Call) context.Context {
	if call.UserAgent != "" {
		ctx = context.WithValue(ctx, callUAKey{}, call.UserAgent)
	}
	if len(call.Headers) > 0 {
		ctx = context.WithValue(ctx, callHeadersKey{}, call.Headers)
	}
	return ctx
}

func withObjectCallUA(ctx context.Context, call fantasy.ObjectCall) context.Context {
	if call.UserAgent != "" {
		ctx = context.WithValue(ctx, callUAKey{}, call.UserAgent)
	}
	if len(call.Headers) > 0 {
		ctx = context.WithValue(ctx, callHeadersKey{}, call.Headers)
	}
	return ctx
}

func wrapHTTPClient(client *http.Client) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &http.Client{
		Transport:     &uaTransport{base: transport},
		CheckRedirect: client.CheckRedirect,
		Jar:           client.Jar,
		Timeout:       client.Timeout,
	}
}

type uaTransport struct {
	base http.RoundTripper
}

func (t *uaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ua, _ := req.Context().Value(callUAKey{}).(string)
	headers, _ := req.Context().Value(callHeadersKey{}).(map[string]string)
	if ua != "" || len(headers) > 0 {
		req = req.Clone(req.Context())
		if ua != "" {
			req.Header.Set("User-Agent", ua)
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	return t.base.RoundTrip(req)
}
