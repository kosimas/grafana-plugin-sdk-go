package backend

import (
	"context"
	"net/http"

	"github.com/kosimas/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/kosimas/grafana-plugin-sdk-go/genproto/pluginv2"
)

// dataSDKAdapter adapter between low level plugin protocol and SDK interfaces.
type dataSDKAdapter struct {
	queryDataHandler QueryDataHandler
}

func newDataSDKAdapter(handler QueryDataHandler) *dataSDKAdapter {
	return &dataSDKAdapter{
		queryDataHandler: handler,
	}
}

func withHeaderMiddleware(ctx context.Context, headers http.Header) context.Context {
	if len(headers) > 0 {
		ctx = httpclient.WithContextualMiddleware(ctx,
			httpclient.MiddlewareFunc(func(opts httpclient.Options, next http.RoundTripper) http.RoundTripper {
				if !opts.ForwardHTTPHeaders {
					return next
				}

				return httpclient.RoundTripperFunc(func(qreq *http.Request) (*http.Response, error) {
					// Only set a header if it is not already set.
					for k, v := range headers {
						if qreq.Header.Get(k) == "" {
							for _, vv := range v {
								qreq.Header.Add(k, vv)
							}
						}
					}
					return next.RoundTrip(qreq)
				})
			}))
	}
	return ctx
}

func (a *dataSDKAdapter) QueryData(ctx context.Context, req *pluginv2.QueryDataRequest) (*pluginv2.QueryDataResponse, error) {
	ctx = propagateTenantIDIfPresent(ctx)
	parsedRequest := FromProto().QueryDataRequest(req)
	ctx = withHeaderMiddleware(ctx, parsedRequest.GetHTTPHeaders())
	resp, err := a.queryDataHandler.QueryData(ctx, parsedRequest)
	if err != nil {
		return nil, err
	}

	return ToProto().QueryDataResponse(resp)
}
