package api

import (
	"errors"
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
			renew: func(_ string, _ int) (*Secret, error) {
				return renewedSecret, nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          time.Second,
			name:                 "short_increment_duration",
			leaseDurationSeconds: 60,
			incrementSeconds:     10,
			renew: func(_ string, _ int) (*Secret, error) {
				return renewedSecret, nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          5 * time.Second,
			name:                 "one_error",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Secret, error) {
				if caseOneErrorCount == 0 {
					caseOneErrorCount++
					return nil, errors.New("renew failure")
				}
				return renewedSecret, nil
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "many_errors",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Secret, error) {
				if caseManyErrorsCount == 3 {
					return renewedSecret, nil
				}
				caseManyErrorsCount++
				return nil, errors.New("renew failure")
			},
			expectError:   nil,
			expectRenewal: true,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "only_errors",
			leaseDurationSeconds: 15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Secret, error) {
				return nil, errors.New("renew failure")
			},
			expectError:   nil,
			expectRenewal: false,
		},
		{
			maxTestTime:          15 * time.Second,
			name:                 "negative_lease_duration",
			leaseDurationSeconds: -15,
			incrementSeconds:     15,
			renew: func(_ string, _ int) (*Secret, error) {
				return renewedSecret, nil
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
					if r.Secret != renewedSecret {
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
