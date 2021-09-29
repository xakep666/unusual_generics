package unusual_generics_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xakep666/unusual_generics"
)

var testZone = time.FixedZone("test", int(time.Hour.Seconds()))

type noLocationProvider struct{}

func (noLocationProvider) TimeLayout() string { return "2006 Jan _2 15:04:05 -0700" }

func (noLocationProvider) TimeLocation() *time.Location { return nil }

type withLocationProvider struct{}

func (withLocationProvider) TimeLayout() string { return "2006 Jan _2 15:04:05" }

func (withLocationProvider) TimeLocation() *time.Location { return testZone }

func TestTimeFormatted(t *testing.T) {
	type timeValues struct {
		TimeA unusual_generics.TimeFormatted[noLocationProvider]   `json:"time_a"`
		TimeB unusual_generics.TimeFormatted[withLocationProvider] `json:"time_b"`
	}

	t.Run("MarshalJSON", func(t *testing.T) {
		b, err := json.Marshal(timeValues{
			TimeA: unusual_generics.FromTime[noLocationProvider](
				time.Date(2021, time.September, 26, 10, 0, 0, 0, testZone),
			),
			TimeB: unusual_generics.FromTime[withLocationProvider](
				time.Date(2021, time.September, 26, 10, 0, 0, 0,
					time.FixedZone("test2", 2*int(time.Hour.Seconds())),
				),
			),
		})
		require.NoError(t, err)

		assert.JSONEq(t, `{
	"time_a": "2021 Sep 26 10:00:00 +0100",
	"time_b": "2021 Sep 26 09:00:00"
}`, string(b))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		var values timeValues

		require.NoError(t, json.Unmarshal([]byte(`{
	"time_a": "2021 Sep 26 10:00:00 +0100",
	"time_b": "2021 Sep 26 09:00:00"
}`), &values))

		assert.Equal(t,
			time.Date(2021, time.September, 26, 9, 0, 0, 0, time.UTC),
			values.TimeA.Time().UTC(),
		)

		_, off := values.TimeA.Time().Zone()
		assert.Equal(t, int(time.Hour.Seconds()), off)

		assert.Equal(t,
			time.Date(2021, time.September, 26, 9, 0, 0, 0, testZone),
			values.TimeB.Time(),
		)
	})
}
