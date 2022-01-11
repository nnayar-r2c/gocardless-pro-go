package gocardless

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

var _ = query.Values
var _ = bytes.NewBuffer
var _ = json.NewDecoder
var _ = errors.New

// CurrencyExchangeRateService manages currency_exchange_rates
type CurrencyExchangeRateService struct {
	endpoint string
	token    string
	client   *http.Client
}

// CurrencyExchangeRate model
type CurrencyExchangeRate struct {
	Rate   string `url:"rate,omitempty" json:"rate,omitempty"`
	Source string `url:"source,omitempty" json:"source,omitempty"`
	Target string `url:"target,omitempty" json:"target,omitempty"`
	Time   string `url:"time,omitempty" json:"time,omitempty"`
}

// CurrencyExchangeRateListParams parameters
type CurrencyExchangeRateListParams struct {
	After     string `url:"after,omitempty" json:"after,omitempty"`
	Before    string `url:"before,omitempty" json:"before,omitempty"`
	CreatedAt struct {
		Gt  string `url:"gt,omitempty" json:"gt,omitempty"`
		Gte string `url:"gte,omitempty" json:"gte,omitempty"`
		Lt  string `url:"lt,omitempty" json:"lt,omitempty"`
		Lte string `url:"lte,omitempty" json:"lte,omitempty"`
	} `url:"created_at,omitempty" json:"created_at,omitempty"`
	Limit  int    `url:"limit,omitempty" json:"limit,omitempty"`
	Source string `url:"source,omitempty" json:"source,omitempty"`
	Target string `url:"target,omitempty" json:"target,omitempty"`
}

// CurrencyExchangeRateListResult response including pagination metadata
type CurrencyExchangeRateListResult struct {
	CurrencyExchangeRates []CurrencyExchangeRate `json:"currency_exchange_rates"`
	Meta                  struct {
		Cursors struct {
			After  string `url:"after,omitempty" json:"after,omitempty"`
			Before string `url:"before,omitempty" json:"before,omitempty"`
		} `url:"cursors,omitempty" json:"cursors,omitempty"`
		Limit int `url:"limit,omitempty" json:"limit,omitempty"`
	} `json:"meta"`
}

// List
// Returns a [cursor-paginated](#api-usage-cursor-pagination) list of all
// exchange rates.
func (s *CurrencyExchangeRateService) List(ctx context.Context, p CurrencyExchangeRateListParams, opts ...RequestOption) (*CurrencyExchangeRateListResult, error) {
	uri, err := url.Parse(fmt.Sprintf(s.endpoint + "/currency_exchange_rates"))
	if err != nil {
		return nil, err
	}

	o := &requestOptions{
		retries: 3,
	}
	for _, opt := range opts {
		err := opt(o)
		if err != nil {
			return nil, err
		}
	}

	var body io.Reader

	v, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", uri.String(), body)
	if err != nil {
		return nil, err
	}
	req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("GoCardless-Version", "2015-07-06")
	req.Header.Set("GoCardless-Client-Library", "gocardless-pro-go")
	req.Header.Set("GoCardless-Client-Version", "1.0.0")
	req.Header.Set("User-Agent", userAgent)

	for key, value := range o.headers {
		req.Header.Set(key, value)
	}

	client := s.client
	if client == nil {
		client = http.DefaultClient
	}

	var result struct {
		Err *APIError `json:"error"`
		*CurrencyExchangeRateListResult
	}

	err = try(o.retries, func() error {
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		err = responseErr(res)
		if err != nil {
			return err
		}

		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			return err
		}

		if result.Err != nil {
			return result.Err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if result.CurrencyExchangeRateListResult == nil {
		return nil, errors.New("missing result")
	}

	return result.CurrencyExchangeRateListResult, nil
}

type CurrencyExchangeRateListPagingIterator struct {
	cursor         string
	response       *CurrencyExchangeRateListResult
	params         CurrencyExchangeRateListParams
	service        *CurrencyExchangeRateService
	requestOptions []RequestOption
}

func (c *CurrencyExchangeRateListPagingIterator) Next() bool {
	if c.cursor == "" && c.response != nil {
		return false
	}

	return true
}

func (c *CurrencyExchangeRateListPagingIterator) Value(ctx context.Context) (*CurrencyExchangeRateListResult, error) {
	if !c.Next() {
		return c.response, nil
	}

	s := c.service
	p := c.params
	p.After = c.cursor

	uri, err := url.Parse(fmt.Sprintf(s.endpoint + "/currency_exchange_rates"))

	if err != nil {
		return nil, err
	}

	o := &requestOptions{
		retries: 3,
	}
	for _, opt := range c.requestOptions {
		err := opt(o)
		if err != nil {
			return nil, err
		}
	}

	var body io.Reader

	v, err := query.Values(p)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", uri.String(), body)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("GoCardless-Version", "2015-07-06")
	req.Header.Set("GoCardless-Client-Library", "gocardless-pro-go")
	req.Header.Set("GoCardless-Client-Version", "1.0.0")
	req.Header.Set("User-Agent", userAgent)

	for key, value := range o.headers {
		req.Header.Set(key, value)
	}
	client := s.client
	if client == nil {
		client = http.DefaultClient
	}

	var result struct {
		Err *APIError `json:"error"`
		*CurrencyExchangeRateListResult
	}

	err = try(o.retries, func() error {
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		err = responseErr(res)

		if err != nil {
			return err
		}

		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			return err
		}

		if result.Err != nil {
			return result.Err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if result.CurrencyExchangeRateListResult == nil {
		return nil, errors.New("missing result")
	}

	c.response = result.CurrencyExchangeRateListResult
	c.cursor = c.response.Meta.Cursors.After
	return c.response, nil
}

func (s *CurrencyExchangeRateService) All(ctx context.Context,
	p CurrencyExchangeRateListParams,
	opts ...RequestOption) *CurrencyExchangeRateListPagingIterator {
	return &CurrencyExchangeRateListPagingIterator{
		params:         p,
		service:        s,
		requestOptions: opts,
	}
}