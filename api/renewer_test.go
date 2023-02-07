package api

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
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
					return nil, fmt.Errorf("renew failure")
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
			renew: func(_ string, _ int) (*Secret, error) {
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

// randDuration calculates a random duration for use in property based testing.
func randDuration(r *rand.Rand, max int64) time.Duration {
	return time.Duration(r.Int63n(max))
}

// sleepLessThanRemainingLease tests that "calculateSleepDuration" will always return a value less than
// the remaining lease duration given a random leaseDuration, priorDuration, remainingLeaseDuration, and increment.
// Inputs are generated so that:
// leaseDuration > priorDuration > remainingLeaseDuration
// and leaseDuration > increment
func sleepLessThanRemainingLease(r *rand.Rand, leaseDuration, priorDuration, remainingLeaseDuration time.Duration, increment int) bool {
	lw := LifetimeWatcher{
		grace:     0,
		increment: increment,
		random:    r,
	}

	lw.calculateGrace(leaseDuration, time.Duration(increment))

	// ensure that we sleep for less than the remaining lease.
	return lw.calculateSleepDuration(remainingLeaseDuration, priorDuration) < remainingLeaseDuration
}

// TestCalcSleepPeriod uses property based testing to evaluate the calculateSleepDuration
// function of LifeTimeWatchers, but also incidentally tests "calculateGrace".
// This is on account of "calculateSleepDuration" performing the "calculateGrace"
// function in particular instances.
// Both of these functions support the vital functionality of the LifeTimeWatcher
// and therefore should be tested rigorously.
func TestCalcSleepPeriod(t *testing.T) {
	c := quick.Config{
		MaxCount: 1000,
		Values: func(values []reflect.Value, r *rand.Rand) {
			// total lease duration
			leaseDuration := randDuration(r, math.MaxInt64)
			priorDuration := randDuration(r, int64(leaseDuration))
			remainingLeaseDuration := randDuration(r, int64(priorDuration))

			values[0] = reflect.ValueOf(r)
			values[1] = reflect.ValueOf(leaseDuration)

			// prior lease duration
			values[2] = reflect.ValueOf(priorDuration)
			// remaining lease duration
			values[3] = reflect.ValueOf(remainingLeaseDuration)

			// increment
			// integer truncation... could be interesting.
			values[4] = reflect.ValueOf(r.Intn(int(leaseDuration)))
		},
	}

	if err := quick.Check(sleepLessThanRemainingLease, &c); err != nil {
		t.Error(err)
	}
}
