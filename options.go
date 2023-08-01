package rekuest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// Option for HTTP request options
type Option func(*option) error

// RedirectInterceptorFunc adalah signature function untuk
type RedirectInterceptorFunc func(*http.Request, []*http.Request) error

// ResponseHeaderCapturer for header capture when using header capture option
type ResponseHeaderCapturer struct {
	Header http.Header
}

type customErrorCapturer struct {
	capturer    any
	successCode int
}

type option struct {
	header                    http.Header
	query                     url.Values
	customPayload             any
	requestDump               io.Writer
	responseDump              io.Writer
	redirectInterceptor       RedirectInterceptorFunc
	httpResponseHeaderCapture *ResponseHeaderCapturer
	customErrResponse         *customErrorCapturer
	customClient              *http.Client
	context                   context.Context
}

// WithHeader for optional header when making request
func WithHeader(key, val string) Option {
	return func(o *option) error {
		if o.header == nil {
			o.header = http.Header{}
		}

		o.header.Add(key, val)
		return nil
	}
}

// WithQuery for optional query string when making request
func WithQuery(key, val string) Option {
	return func(o *option) error {
		if o.query == nil {
			o.query = url.Values{}
		}

		o.query.Add(key, val)
		return nil
	}
}

// WithRequestDump for dumping request
func WithRequestDump(buf io.Writer) Option {
	return func(o *option) error {
		o.requestDump = buf
		return nil
	}
}

// WithResponseDump for dumping response
func WithResponseDump(buf io.Writer) Option {
	return func(o *option) error {
		o.responseDump = buf
		return nil
	}
}

// WithJSON for adding JSON payload to the request
func WithJSON(payload any) Option {
	return func(o *option) error {
		o.customPayload = payload
		return nil
	}
}

// WithHTTPRedirectIntercept for adding interceptor when the request is being
// redirected. If you want to cancel the redirect request, make the function
// return http.ErrUseLastResponse. This will instruct the net/http package to
// not execute the redirect request and instead returns the last HTTP response
// with the body un-closed.
func WithHTTPRedirectIntercept(c RedirectInterceptorFunc) Option {
	return func(o *option) error {
		o.redirectInterceptor = c
		return nil
	}
}

// WithHTTPResponseHeaderCapture for capturing the response headers
func WithHTTPResponseHeaderCapture(h *ResponseHeaderCapturer) Option {
	return func(o *option) error {
		o.httpResponseHeaderCapture = h
		return nil
	}
}

// WithCustomErrorResponse for adding a custom response capturer. For example
// in normal scenario (HTTP 200), an API returns an object of type A, but on
// failure, it returns B. You would put the a pointer of type B as 'capture',
// and http.StatusOK as 'successStatus'. This will capture the custom error
// response in that B object you provided
func WithCustomErrorResponse(capture any, successStatus int) Option {
	return func(o *option) error {
		if capture == nil {
			return errors.Join(ErrInvalidOption, fmt.Errorf("WithCustomErrorResponse needs a pointer type to capture the output"))
		}

		if reflect.TypeOf(capture).Kind() != reflect.Pointer {
			return errors.Join(ErrInvalidOption, fmt.Errorf("WithCustomErrorResponse needs a pointer type to capture the output"))
		}

		o.customErrResponse = &customErrorCapturer{
			capturer:    capture,
			successCode: successStatus,
		}
		return nil
	}
}

// WithCustomHTTPClient instructs the function to execute the request using
// this client. If not provided, use a bare minimum client instead
func WithCustomHTTPClient(c *http.Client) Option {
	return func(o *option) error {
		o.customClient = c
		return nil
	}
}

// WithContext instructs the function to use this context when executing the
// request. Otherwise, context.Background() is used
func WithContext(c context.Context) Option {
	return func(o *option) error {
		o.context = c
		return nil
	}
}
