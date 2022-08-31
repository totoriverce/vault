package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
)

func TestRenewer_NewRenewer(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		i    *RenewerInput
		e    *Renewer
		err  bool
	}{
		{
			name: "nil",
			i:    nil,
			e:    nil,
			err:  true,
		},
		{
			name: "missing_secret",
			i: &RenewerInput{
				Secret: nil,
			},
			e:   nil,
			err: true,
		},
		{
			name: "default_grace",
			i: &RenewerInput{
				Secret: &Secret{},
			},
			e: &Renewer{
				secret: &Secret{},
			},
			err: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := client.NewRenewer(tc.i)
			if (err != nil) != tc.err {
				t.Fatal(err)
			}

			if v == nil {
				return
			}

			// Zero-out channels because reflect
			v.client = nil
			v.random = nil
			v.doneCh = nil
			v.renewCh = nil
			v.stopCh = nil

			if diff := deep.Equal(tc.e, v); diff != nil {
				t.Error(diff)
			}
		})
	}
}

// renewedResponse returns an empty http.Response
// for testing TestLifetimeWatcher
func renewedResponse(age string) *Response {
	headers := http.Header{}

	if age != "" {
		headers.Set("Age", age)
	}

	return &Response{
		Response: &http.Response{
			Header: headers,
			Body:   io.NopCloser(strings.NewReader("{}")),
		},
	}
}

func TestLifetimeWatcher(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	// Note that doRenewWithOptions starts its loop with an initial renewal.
	// This has a big impact on the particulars of the following cases.

	renewedSecret := &Secret{}
	var caseOneErrorCount int
	var caseManyErrorsCount int
	cases := []struct {
		maxTestTime          time.Duration
		name                 string
		leaseDurationSeconds int
		incrementSeconds     int
		renew                renewFunc
		expectError          error
		expectRenewal        bool
	}{
		{
			maxTestTime:          time.Second,
			name:                 "no_error",
			leaseDurationSeconds: 60,
			incrementSeconds:     60,
			renew: func(_ string, _ int) (*Response, error) {
				return renewedResponse(""), nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          time.Second,
			name:                 "short_increment_duration",
			leaseDurationSeconds: 60,
			incrementSeconds:     10,
			renew: func(_ string, _ int) (*Response, error) {
				return renewedResponse(""), nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          5 * time.Second,
			name:                 "one_error",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Response, error) {
				if caseOneErrorCount == 0 {
					caseOneErrorCount++
					return nil, fmt.Errorf("renew failure")
				}
				return renewedResponse(""), nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "many_errors",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Response, error) {
				if caseManyErrorsCount == 3 {
					return renewedResponse(""), nil
				}
				caseManyErrorsCount++
				return nil, fmt.Errorf("renew failure")
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "only_errors",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Response, error) {
				return nil, fmt.Errorf("renew failure")
			},
			expectError:   nil,
			expectRenewal: false,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "negative_lease_duration",
			leaseDurationSeconds: -15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Response, error) {
				return renewedResponse(""), nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
				Secret: &Secret{
					LeaseDuration: tc.leaseDurationSeconds,
				},
				Increment: tc.incrementSeconds,
			})
			if err != nil {
				t.Fatal(err)
			}

			doneCh := make(chan error, 1)
			go func() {
				doneCh <- v.doRenewWithOptions(false, false,
					tc.leaseDurationSeconds, "myleaseID", tc.renew, time.Second)
			}()
			defer v.Stop()

			receivedRenewal := false
			receivedDone := false
		ChannelLoop:
			for {
				select {
				case <-time.After(tc.maxTestTime):
					t.Fatalf("renewal didn't happen")
				case r := <-v.RenewCh():
					if !tc.expectRenewal {
						t.Fatal("expected no renewals")
					}
					if !reflect.DeepEqual(r.Secret, renewedSecret) {
						t.Fatalf("expected secret %v, got %v", renewedSecret, r.Secret)
					}
					receivedRenewal = true
					if !receivedDone {
						continue ChannelLoop
					}
					break ChannelLoop
				case err := <-doneCh:
					receivedDone = true
					if tc.expectError != nil && !errors.Is(err, tc.expectError) {
						t.Fatalf("expected error %q, got: %v", tc.expectError, err)
					}
					if tc.expectError == nil && err != nil {
						t.Fatalf("expected no error, got: %v", err)
					}
					if tc.expectRenewal && !receivedRenewal {
						// We might have received the stop before the renew call on the channel.
						continue ChannelLoop
					}
					break ChannelLoop
				}
			}

			if tc.expectRenewal && !receivedRenewal {
				t.Fatalf("expected at least one renewal, got none.")
			}
		})
	}
}

func agingRenew() func(string, int) (*Response, error) {
	origin := time.Now()
	headers := http.Header{}

	return func(leaseID string, increment int) (*Response, error) {
		headers.Set("Age", strconv.Itoa(int(time.Since(origin))))
		return &Response{
			Response: &http.Response{
				Header: headers,
				Body:   io.NopCloser(strings.NewReader(fmt.Sprintf(`{"lease_id": "%s", "increment": %d}`, leaseID, increment))),
			},
		}, nil
	}
}

// TestLifeTimeWatcherCached specifically tests how LifeTimeWatchers react to
// cached renew responses.
func TestLifeTimeWatcherCached(t *testing.T) {
	t.Parallel()

	client, err := NewClient(DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		// the full duration a test should run for
		// the maxTestTime should relate to the number
		// of renewals expected.
		maxTestTime time.Duration
		// corresponds to the lease duration of a secret
		// the lease duration should indicate how many
		// "maximum" renewals we should expect during a maxTestTime.
		leaseDurationSeconds time.Duration
		// indicates how long renewals should request
		// also determines how many renewals we should expect during the maxTestTime.
		incrementSeconds time.Duration
		// the initial age in seconds for a renewal response.
		age         int
		expectError error
	}{
		{
			name:                 "origin",
			maxTestTime:          20 * time.Second,
			leaseDurationSeconds: 2 * time.Second,
			incrementSeconds:     1 * time.Second,
			age:                  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// from the initial lease how many renewals will we execute
			// +1 is because the renewal loop always runs the first time a lifetimewatcher starts
			expectedRenewals := int(((tc.maxTestTime - tc.leaseDurationSeconds) / tc.incrementSeconds) + 1)

			lw, err := client.NewLifetimeWatcher(&LifetimeWatcherInput{
				Secret: &Secret{
					LeaseDuration: int(tc.leaseDurationSeconds.Seconds()),
				},
				Increment: int(tc.incrementSeconds.Seconds()),
			})

			if err != nil {
				t.Fatal(err)
			}

			doneCh := make(chan error, 1)

			go func() {
				doneCh <- lw.doRenewWithOptions(false, false, int(tc.leaseDurationSeconds.Seconds()), tc.name, agingRenew(), tc.incrementSeconds)
			}()
			defer lw.Stop()

			//count the renewals that occur during a test.
			actualRenewals := 0
		ChannelLoop:
			for {
				select {
				case <-time.After(tc.maxTestTime):
					break ChannelLoop
				case <-lw.RenewCh():
					actualRenewals += 1
					continue ChannelLoop
				case loopErr := <-doneCh:
					if loopErr != nil {
						t.Fatal(loopErr)
						return
					}
				}
			}

			if expectedRenewals != actualRenewals {
				t.Fatalf("expected %d renewals, %d occurred", expectedRenewals, actualRenewals)
			}

		})
	}
}
