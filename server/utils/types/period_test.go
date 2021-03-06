// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPeriod(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		p := &Period{
			Period: EveryWeek,
		}

		// Accounts for daylight savings time
		start := MustParseTime("2020-03-03").In(PST)
		next := p.ForNumber(start, 1)
		require.Equal(t, "2020-03-09T16:00", next.String())
		require.Equal(t, "167h0m0s", next.Sub(start.Time).String())

		next2 := p.ForNumber(next, 1)
		require.Equal(t, "2020-03-16T16:00", next2.String())
		require.Equal(t, "168h0m0s", next2.Sub(next.Time).String())

		next = p.ForNumber(start, 10)
		require.Equal(t, "2020-05-11T16:00", next.String())
		require.Equal(t, "1679h0m0s", next.Sub(start.Time).String())
	})
}

func TestPeriodForNumber(t *testing.T) {
	for _, tc := range []struct {
		name     string
		period   Period
		start    Time
		steps    int
		expected Time
	}{
		{
			name:     "happy weekly",
			period:   Period{Period: EveryWeek},
			start:    MustParseTime("2020-02-04T17:00").In(PST),
			steps:    1,
			expected: MustParseTime("2020-02-11T17:00").In(PST),
		},
		{
			name:     "happy bi-weekly",
			period:   Period{Period: EveryTwoWeeks},
			start:    MustParseTime("2020-02-04T17:00").In(PST),
			steps:    3,
			expected: MustParseTime("2020-03-17T16:00").In(PST),
		},
		{
			name:     "happy weekly daylight savings",
			period:   Period{Period: EveryWeek},
			start:    MustParseTime("2020-03-04T17:00").In(PST),
			steps:    1,
			expected: MustParseTime("2020-03-11T16:00").In(PST),
		},
		{
			name:     "happy monthly 1",
			period:   Period{Period: EveryMonth},
			start:    MustParseTime("2020-01-04T17:00").In(PST),
			steps:    1,
			expected: MustParseTime("2020-02-04T17:00").In(PST),
		},
		{
			name:     "happy monthly 2 with leap",
			period:   Period{Period: EveryMonth},
			start:    MustParseTime("2020-01-04T17:00").In(PST),
			steps:    2,
			expected: MustParseTime("2020-03-04T17:00").In(PST),
		},
		{
			name:     "happy monthly 4 with leap and daylight savings",
			period:   Period{Period: EveryMonth},
			start:    MustParseTime("2020-01-04T17:00").In(PST),
			steps:    4,
			expected: MustParseTime("2020-05-04T16:00").In(PST),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := tc.period.ForNumber(tc.start, tc.steps)
			require.Equal(t, tc.expected, tt)
		})
	}
}

func TestPeriodForTime(t *testing.T) {
	for _, tc := range []struct {
		name          string
		period        Period
		start         Time
		now           Time
		expectedNum   int
		expectedStart Time
	}{
		{
			name:          "bi-weekly before",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-02-03T17:00").In(PST),
			expectedNum:   -1,
			expectedStart: Time{},
		}, {
			name:          "bi-weekly on start",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-02-04T17:00").In(PST),
			expectedNum:   0,
			expectedStart: MustParseTime("2020-02-04T17:00").In(PST),
		}, {
			name:          "bi-weekly in first",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-02-10T17:00").In(PST),
			expectedNum:   0,
			expectedStart: MustParseTime("2020-02-04T17:00").In(PST),
		}, {
			name:          "bi-weekly on finish first",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-02-18T17:00").In(PST),
			expectedNum:   1,
			expectedStart: MustParseTime("2020-02-18T17:00").In(PST),
		}, {
			name:          "bi-weekly future DST",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-04-15T17:00").In(PST),
			expectedNum:   5,
			expectedStart: MustParseTime("2020-04-14T16:00").In(PST),
		}, {
			name:          "bi-weekly future non-DST",
			period:        Period{Period: EveryTwoWeeks},
			start:         MustParseTime("2020-02-04T17:00").In(PST),
			now:           MustParseTime("2020-12-01T17:00").In(PST),
			expectedNum:   21,
			expectedStart: MustParseTime("2020-11-24T17:00").In(PST),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			num, start := tc.period.ForTime(tc.start, tc.now)
			require.Equal(t, tc.expectedNum, num)
			require.Equal(t, tc.expectedStart, start)
		})
	}
}
