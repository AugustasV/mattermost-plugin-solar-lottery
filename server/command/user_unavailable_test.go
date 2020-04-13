// Copyright (c) 2019-present Mattermost, Inc. All Rights Reserved.
// See License for license information.
package command

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/sl"
	"github.com/mattermost/mattermost-plugin-solar-lottery/server/utils/types"
	"github.com/stretchr/testify/require"
)

func TestCommandUserUnavailable(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		SL, store := getTestSL(t, ctrl)

		// test-user is in PST
		outmd, err := runCommand(t, SL, `
		/lotto user unavailable -s 2025-01-01T11:00 -f 2025-01-02T09:30
		`)
		require.NoError(t, err)
		require.Equal(t, "added unavailable event personal: 2025-01-01T11:00 to 2025-01-02T09:30 to @test-user-username", outmd.String())

		user := sl.NewUser("")
		err = store.Entity(sl.KeyUser).Load("test-user", user)
		require.NoError(t, err)
		require.Len(t, user.Calendar, 1)
		require.EqualValues(t,
			sl.Unavailable{
				Interval: types.Interval{
					Start:  types.NewTime(time.Date(2025, 1, 1, 19, 0, 0, 0, time.UTC)),
					Finish: types.NewTime(time.Date(2025, 1, 2, 17, 30, 0, 0, time.UTC)),
				},
				Reason: sl.ReasonPersonal,
			},
			*user.Calendar[0])

		err = runCommands(t, SL, `
				/lotto user unavailable -s 2025-02-01 -f 2025-02-03
				/lotto user unavailable -s 2025-02-07 -f 2025-02-10
				/lotto user unavailable -s 2025-06-28 -f 2025-07-05
			`)
		require.NoError(t, err)

		out := sl.OutCalendar{
			Users: sl.NewUsers(),
		}
		_, err = runJSONCommand(t, SL, `
				/lotto user unavailable --clear -s 2025-01-30T10:00 -f 2025-02-08T11:00
				`, &out)
		users := out.Users
		require.NoError(t, err)
		require.EqualValues(t, []string{"test-user"}, users.TestIDs())
		require.Equal(t, 2, len(users.Get("test-user").Calendar))
		require.EqualValues(t,
			&sl.Unavailable{
				Interval: types.Interval{
					Start:  types.NewTime(time.Date(2025, 1, 1, 19, 0, 0, 0, time.UTC)),
					Finish: types.NewTime(time.Date(2025, 1, 2, 17, 30, 0, 0, time.UTC)),
				},
				Reason: sl.ReasonPersonal,
			},
			users.Get("test-user").Calendar[0])
		require.EqualValues(t,
			&sl.Unavailable{
				Interval: types.Interval{
					Start:  types.NewTime(time.Date(2025, 6, 28, 7, 0, 0, 0, time.UTC)),
					Finish: types.NewTime(time.Date(2025, 7, 5, 7, 0, 0, 0, time.UTC)),
				},
				Reason: sl.ReasonPersonal,
			},
			users.Get("test-user").Calendar[1])
	})
}