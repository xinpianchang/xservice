package echox

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Timestamp_scan(t *testing.T) {
	now := time.Now()
	nowStr1 := now.Format("2006-01-02 15:04:05")
	nowStr2 := now.Format("2006-01-02T15:04:05")
	nowStr3 := now.Format("2006-01-02T15:04:05-0700")
	nowStr4 := fmt.Sprint(now.UnixNano() / int64(time.Millisecond))

	var target1 Timestamp

	tests := []struct {
		name   string
		target Timestamp
		src    interface{}
		equal  interface{}
	}{
		{"time", target1, &now, Timestamp(now)},
		{"time str1", target1, &nowStr1, Timestamp(now)},
		{"time str2", target1, &nowStr2, Timestamp(now)},
		{"time str3", target1, &nowStr3, Timestamp(now)},
		{"time str4", target1, &nowStr4, Timestamp(now)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.target.Scan(tt.src)
			assert.NoErrorf(t, err, "%v -> %v", tt.src, tt.target)
			assert.Equal(t, fmt.Sprint(tt.equal), fmt.Sprint(tt.target), "%v -> %v", tt.src, tt.target)
		})
	}
}

func Test_Timestamp_copy(t *testing.T) {
	now := time.Now()

	{
		type ts struct {
			Tp Timestamp
		}
		var a = &ts{Tp: Timestamp(now)}
		b, err := json.Marshal(a)
		require.NoError(t, err)
		require.True(t, !strings.HasPrefix(string(b), `{"Tp":"`))

		type ts2 struct {
			Tp Timestamp
		}
		tp := &ts2{Tp: Timestamp(now)}
		err = copier.Copy(a, tp)
		require.NoError(t, err)
		require.Equal(t, now, time.Time(a.Tp))
	}

	{
		type ts struct {
			Tp *Timestamp
		}

		type ts2 struct {
			Tp *Timestamp
		}

		tp := Timestamp(now)
		a := &ts{}
		b := &ts2{&tp}
		err := copier.Copy(a, b)
		require.NoError(t, err)
		// t.Logf("a:%v, b:%v", a, b)
		require.Equal(t, now, time.Time(*a.Tp))
	}

	{
		type ts struct {
			Tp Timestamp
		}

		type ts2 struct {
			Tp time.Time
		}

		a := &ts{}
		b := &ts2{now}
		err := copier.Copy(a, b)
		require.NoError(t, err)
		// t.Logf("a:%v, b:%v", a, b)
		require.Equal(t, now, time.Time(a.Tp))
	}

	{
		type ts struct {
			Tp *Timestamp
		}

		type ts2 struct {
			Tp time.Time
		}

		a := &ts{}
		b := &ts2{now}
		err := copier.Copy(a, b)
		require.NoError(t, err)
		// t.Logf("a:%v, b:%v", a, b)
		require.Equal(t, now, time.Time(*a.Tp))
	}

	{
		type ts struct {
			Tp *Timestamp
		}

		type ts2 struct {
			Tp *time.Time
		}

		a := &ts{}
		b := &ts2{&now}
		err := copier.Copy(a, b)
		require.NoError(t, err)
		// t.Logf("a:%v, b:%v", a, b)
		require.Equal(t, now, time.Time(*a.Tp))
	}

}
