package rekuest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
)

// HTTPRequest for making HTTP request to an endpoint. This function accepts
// multiple options for customizing its behaviour. If the HTTP request does
// reach the upstream server, the httpStatus return value will be filled,
// otherwise it won't.
func HTTPRequest(method string, urlEndpoint string, resp any, options ...Option) (httpStatus int, e error) {
	if resp != nil && reflect.TypeOf(resp).Kind() != reflect.Pointer {
		return 0, errors.Join(ErrInvalidOption, fmt.Errorf("resp must be a pointer type. If you meant to ignore the response, use nil instead"))
	}

	// parse options
	appliedOptions := &option{}
	for _, apply := range options {
		err := apply(appliedOptions)
		if err != nil {
			return 0, err
		}
	}

	// parse url
	target, err := url.Parse(urlEndpoint)
	if err != nil {
		return 0, err
	}

	if appliedOptions.query != nil {
		target.RawQuery = appliedOptions.query.Encode()
	}

	// parse payload
	payload := bytes.NewBuffer([]byte{})
	if appliedOptions.customPayload != nil {
		if err != nil {
			return 0, err
		}

		if err := json.NewEncoder(payload).Encode(appliedOptions.customPayload); err != nil {
			return 0, err
		}
	}

	// generate context
	var ctx context.Context
	if appliedOptions.context != nil {
		ctx = appliedOptions.context
	} else {
		ctx = context.Background()
	}

	// create request object
	httpReq, err := http.NewRequestWithContext(ctx, method, target.String(), payload)
	if err != nil {
		return 0, err
	}

	// set headers
	if appliedOptions.header != nil {
		httpReq.Header = appliedOptions.header
	}

	// check dump request
	if appliedOptions.requestDump != nil {
		b, err := httputil.DumpRequest(httpReq, true)
		if err != nil {
			fmt.Println("[HTTP] Cannot dump request:", err)
		} else {
			fmt.Fprintf(appliedOptions.requestDump, "REQUEST ======\n%s\n\n", string(b))
		}
	}

	// create custom client
	var client *http.Client
	if appliedOptions.customClient != nil {
		client = appliedOptions.customClient
	} else {
		client = new(http.Client)
	}

	// apply custom redirect interceptor option
	if appliedOptions.redirectInterceptor != nil {
		client.CheckRedirect = appliedOptions.redirectInterceptor
	}

	// do the request
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return 0, err
	}

	// dump response
	if appliedOptions.responseDump != nil {
		b, err := httputil.DumpResponse(httpResp, true)
		if err != nil {
			fmt.Println("[HTTP] Cannot dump response:", err)
		} else {
			fmt.Fprintf(appliedOptions.responseDump, "RESPONSE ======\n%s\n\n", string(b))
		}
	}

	// optional jika capture response header
	if appliedOptions.httpResponseHeaderCapture != nil {
		appliedOptions.httpResponseHeaderCapture.Header = httpResp.Header.Clone()
	}

	// early return, if the user only cares the status code
	if resp == nil {
		return httpResp.StatusCode, nil
	}

	// parse body, also includes the copy if an error happens while parsing to JSON
	var (
		respByte        = bytes.NewBuffer([]byte{})
		respByteCopy    = bytes.NewBuffer([]byte{})
		respByteErrCopy = bytes.NewBuffer([]byte{})
	)
	b, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return httpResp.StatusCode, err
	}

	if _, err := io.MultiWriter(respByte, respByteCopy, respByteErrCopy).Write(b); err != nil {
		return httpResp.StatusCode, err
	}

	// if user applies the custom error response option, use that to capture
	// the failure response
	if appliedOptions.customErrResponse != nil && httpResp.StatusCode != appliedOptions.customErrResponse.successCode {
		body := respByteErrCopy.Bytes()
		err := json.Unmarshal(body, appliedOptions.customErrResponse.capturer)
		if err != nil {
			resp := respByteCopy.Bytes()
			return httpResp.StatusCode, errors.Join(err, fmt.Errorf("response: %s", string(resp)))
		}

		// exits, providing the http status & unmarshaled error response body
		return httpResp.StatusCode, nil
	}

	if err := json.Unmarshal(respByte.Bytes(), resp); err != nil {
		resp := respByteCopy.Bytes()
		return httpResp.StatusCode, errors.Join(err, fmt.Errorf("response: %s", string(resp)))
	}

	return httpResp.StatusCode, nil
}
