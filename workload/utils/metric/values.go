package metric

import (
	"encoding/json"
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"math"
	"time"
)

const (
	TypeFloat ValueType = iota
	TypeSize
	TypeDuration
)

type Value struct {
	Type  ValueType
	Value interface{}
}

func (v *Value) UnmarshalJSON(bytes []byte) error {
	var val struct {
		Type  ValueType
		Value interface{}
	}
	if err := json.Unmarshal(bytes, &val); err != nil {
		return err
	}
	v.Type = val.Type
	v.Value = val.Value
	switch v.Type {
	// Any better way? If leave the unmarshal work to encoding/json,
	// int64 would be formatted to float numbers.
	case TypeSize:
		v.Value = int64(v.Value.(float64))
	case TypeDuration:
		v.Value = time.Duration(v.Value.(float64))
	default:
	}
	return nil
}

func (v *Value) String() string {
	switch v.Type {
	case TypeFloat:
		return fmt.Sprintf("%.2f", v.Value)
	case TypeSize:
		return ToIec(v.Value.(int64))
	case TypeDuration:
		return v.Value.(time.Duration).String()
	default:
		return fmt.Sprintf("<error type> %v", v.Value)
	}
}

func (v *Value) Float() float64 {
	switch v.Type {
	case TypeFloat:
		return v.Value.(float64)
	case TypeSize:
		return float64(v.Value.(int64))
	case TypeDuration:
		return float64(v.Value.(time.Duration))
	default:
		log.Warn("unsupported type casting to float", zap.Reflect("type", v.Type))
		return math.NaN()
	}
}

func Float(f float64) Value {
	return Value{
		Type:  TypeFloat,
		Value: f,
	}
}

func Size(size int64) Value {
	return Value{
		Type:  TypeSize,
		Value: size,
	}
}

func Duration(d time.Duration) Value {
	return Value{
		Type:  TypeDuration,
		Value: d,
	}
}

func (v Value) Named(name string) Metric {
	return Metric{
		Name:  name,
		Value: v,
	}
}

// ComparePercentDiff compares the percent diff of other according to their Float value.
// returns: (v / other) - 1
func (v Value) ComparePercentDiff(other Value) float64 {
	return v.Float()/other.Float() - 1
}
