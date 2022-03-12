package protocol

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PapaCharlie/go-restli/protocol/restlicodec"
)

const (
	RestLiProtocolVersion = "2.0.0"

	RestLiHeader_ID              = "X-RestLi-Id"
	RestLiHeader_Method          = "X-RestLi-Method"
	RestLiHeader_ProtocolVersion = "X-RestLi-Protocol-Version"
	RestLiHeader_ErrorResponse   = "X-RestLi-Error-Response"
)

type RestLiMethod int

// Disabled until https://github.com/golang/go/issues/45218 is resolved: go:generate stringer -type=RestLiMethod -trimprefix Method_
const (
	Method_Unknown = RestLiMethod(iota)

	Method_get
	Method_create
	Method_delete
	Method_update
	Method_partial_update

	Method_batch_get
	Method_batch_create
	Method_batch_delete
	Method_batch_update
	Method_batch_partial_update

	Method_get_all

	Method_action
	Method_finder
)

var RestLiMethodNameMapping = func() map[string]RestLiMethod {
	mapping := make(map[string]RestLiMethod)
	for m := Method_get; m <= Method_finder; m++ {
		mapping[m.String()] = m
	}
	return mapping
}()

type RestLiClient struct {
	*http.Client
	HostnameResolver
	// Whether or not missing fields in a restli response should cause a MissingRequiredFields error to be returned.
	// Note that even if the error is returned, the response will still be fully deserialized.
	StrictResponseDeserialization bool
}

func (c *RestLiClient) FormatQueryUrl(rp ResourcePath, query QueryParamsEncoder) (*url.URL, error) {
	path, err := rp.ResourcePath()
	if err != nil {
		return nil, err
	}

	if query != nil {
		var params string
		params, err = query.EncodeQueryParams()
		if err != nil {
			return nil, err
		}
		path += "?" + params
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	root := rp.RootResource()
	hostUrl, err := c.ResolveHostnameAndContextForQuery(root, u)
	if err != nil {
		return nil, err
	}

	resolvedPath := "/" + strings.TrimSuffix(strings.TrimPrefix(hostUrl.EscapedPath(), "/"), "/")

	if resolvedPath == "/" {
		return hostUrl.ResolveReference(u), nil
	}

	if idx := strings.Index(resolvedPath, "/"+root); idx >= 0 &&
		(len(resolvedPath) == idx+len(root)+1 || resolvedPath[idx+len(root)+1] == '/') {
		resolvedPath = resolvedPath[:idx]
	}

	return hostUrl.Parse(resolvedPath + u.RequestURI())
}

func SetJsonAcceptHeader(req *http.Request) {
	req.Header.Set("Accept", "application/json")
}

func SetJsonContentTypeHeader(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func SetRestLiHeaders(req *http.Request, method RestLiMethod) {
	req.Header.Set(RestLiHeader_ProtocolVersion, RestLiProtocolVersion)
	req.Header.Set(RestLiHeader_Method, method.String())
}

func newRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	httpMethod string,
	method RestLiMethod,
	body io.Reader,
) (*http.Request, error) {
	u, err := c.FormatQueryUrl(rp, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, u.String(), body)
	if err != nil {
		return nil, err
	}

	SetRestLiHeaders(req, method)
	return req, nil
}

// NewGetRequest creates a GET http.Request and sets the expected rest.li headers
func NewGetRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method RestLiMethod,
) (*http.Request, error) {
	req, err := newRequest(c, ctx, rp, query, http.MethodGet, method, http.NoBody)
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)

	return req, nil
}

// NewDeleteRequest creates a DELETE http.Request and sets the expected rest.li headers
func NewDeleteRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method RestLiMethod,
) (*http.Request, error) {
	req, err := newRequest(c, ctx, rp, query, http.MethodDelete, method, http.NoBody)
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)

	return req, nil
}

func NewCreateRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	method RestLiMethod,
	create restlicodec.Marshaler,
	readOnlyFields restlicodec.PathSpec,
) (*http.Request, error) {
	return NewJsonRequest(c, ctx, rp, query, http.MethodPost, method, create, readOnlyFields)
}

// NewJsonRequest creates an http.Request with the given HTTP method and rest.li method, and populates the body of the
// request with the given restlicodec.Marshaler contents (see RawJsonRequest)
func NewJsonRequest(
	c *RestLiClient,
	ctx context.Context,
	rp ResourcePath,
	query QueryParamsEncoder,
	httpMethod string,
	restLiMethod RestLiMethod,
	contents restlicodec.Marshaler,
	excludedFields restlicodec.PathSpec,
) (*http.Request, error) {
	writer := restlicodec.NewCompactJsonWriterWithExcludedFields(excludedFields)
	err := contents.MarshalRestLi(writer)
	if err != nil {
		return nil, err
	}

	size := writer.Size()
	req, err := newRequest(c, ctx, rp, query, httpMethod, restLiMethod, writer.ReadCloser())
	if err != nil {
		return nil, err
	}

	SetJsonAcceptHeader(req)
	SetJsonContentTypeHeader(req)

	req.ContentLength = int64(size)
	return req, nil
}

// Do is a very thin shim between the standard http.Client.Do. All it does it parse the response into a RestLiError if
// the RestLi error header is set. A non-nil Response with a non-nil error will only occur if http.Client.Do returns
// such values (see the corresponding documentation). Otherwise, the response will only be non-nil if the error is nil.
// All (and only) network-related errors will be of type *url.Error. Other types of errors such as parse errors will use
// different error types.
func (c *RestLiClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return res, err
	}

	err = IsErrorResponse(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DoAndUnmarshal calls Do and attempts to unmarshal the response into the given value. The response body will always be
// read to EOF and closed, to ensure the connection can be reused.
func DoAndUnmarshal[V any](
	c *RestLiClient,
	req *http.Request,
	unmarshaler restlicodec.GenericUnmarshaler[V],
) (v V, res *http.Response, err error) {
	data, res, err := c.do(req)
	if err != nil {
		return v, res, err
	}

	r, err := restlicodec.NewJsonReader(data)
	if err != nil {
		return v, res, err
	}
	v, err = unmarshaler(r)
	if _, mfe := err.(*restlicodec.MissingRequiredFieldsError); mfe && !c.StrictResponseDeserialization {
		err = nil
	}
	return v, res, err
}

// DoAndIgnore calls Do and drops the response's body. The response body will always be read to EOF and closed, to
// ensure the connection can be reused.
func DoAndIgnore(c *RestLiClient, req *http.Request) (*http.Response, error) {
	_, res, err := c.do(req)
	return res, err
}

func (c *RestLiClient) do(req *http.Request) ([]byte, *http.Response, error) {
	res, err := c.Do(req)
	if err != nil {
		return nil, res, err
	}

	if v := res.Header.Get(RestLiHeader_ProtocolVersion); v != RestLiProtocolVersion {
		return nil, nil, &UnsupportedRestLiProtocolVersion{v}
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, &url.Error{
			Op:  "ReadResponse",
			URL: req.URL.String(),
			Err: err,
		}
	}

	err = res.Body.Close()
	if err != nil {
		return nil, nil, &url.Error{
			Op:  "CloseResponse",
			URL: req.URL.String(),
			Err: err,
		}
	}

	return data, res, nil
}
