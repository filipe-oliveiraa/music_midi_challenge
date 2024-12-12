package typ

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Duration time.Duration

var emptyDuration Duration

func (d Duration) GetValue(v string) (any, error) {
	dv, err := time.ParseDuration(v)
	if err != nil {
		return emptyDuration, nil
	}

	return Duration(dv), nil
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	fmt.Println(time.Duration(d).String())
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmpDuration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		*d = Duration(time.Duration(tmpDuration))
		return nil
	default:
		return errors.New("invalid duration")
	}
}
