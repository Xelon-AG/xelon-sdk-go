package xelon

import (
	"context"
	"fmt"
	"iter"
	"net/http"
)

// paginatedResponse is a generic "container" used to unmarshal any paginated list response
// from Xelon REST API. It assumes the response body has a "data" field, containing a slice
// of items and a "meta" field for pagination.
type paginatedResponse[T any] struct {
	Data []T   `json:"data,omitempty"`
	Meta *Meta `json:"meta"`
}

func newPaginator[T any](ctx context.Context, client *Client, pathURL string, opts *ListOptions) (iter.Seq2[T, *Response], func() error) {
	var iterErr error
	seq := func(yield func(item T, resp *Response) bool) {
		if opts == nil {
			opts = &ListOptions{}
		}
		if opts.Page == 0 {
			opts.Page = 1
		}
		if opts.PerPage == 0 {
			opts.PerPage = 10
		}

		for {
			select {
			// if the context has been canceled, the context's error is more useful
			case <-ctx.Done():
				iterErr = ctx.Err()
				return
			default:
			}

			path, err := addOptions(pathURL, opts)
			if err != nil {
				iterErr = fmt.Errorf("failed to construct URL with options: %w", err)
				return
			}
			req, err := client.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				iterErr = fmt.Errorf("failed to prepare paginated request: %w", err)
				return
			}

			page := new(paginatedResponse[T])

			resp, err := client.Do(ctx, req, page)
			if err != nil {
				iterErr = err
				return
			}
			if m := page.Meta; m != nil {
				resp.Meta = m
			}

			for _, item := range page.Data {
				if !yield(item, resp) {
					// stop iteration if the consumer stops
					return
				}
			}

			if resp.Meta == nil || opts.Page >= resp.Meta.LastPage {
				// no more next pages, exit from pagination
				break
			}

			opts.Page++
		}
	}

	return seq, func() error { return iterErr }
}
