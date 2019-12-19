package octosql

import (
	"math"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v Value
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "int test with a negative number",
			args: args{
				v: MakeInt(-18249),
			},
		},
		{
			name: "int test with a positive number",
			args: args{
				v: MakeInt(1587129),
			},
		},
		{
			name: "int test with a zero",
			args: args{
				v: MakeInt(0),
			},
		},
		{
			name: "string test",
			args: args{
				v: MakeString("ala ma kota i psa"),
			},
		},
		{
			name: "string test with weird values",
			args: args{
				v: MakeString(string([]byte{28, 192, 0, 123, 11, 99, 243, 172, 111, 3, 4, 5})),
			},
		},
		{
			name: "tuple test 1",
			args: args{
				v: MakeTuple([]Value{MakeString("ala"), MakeInt(17), MakeInt(99)}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Value

			bytes := (&tt.args.v).SortedMarshal()

			err := b.SortedUnmarshal(bytes)
			if err != nil {
				t.Errorf("SortedMarshal() error = %v, wantErr false", err)
				return
			}

			if !reflect.DeepEqual(tt.args.v, b) {
				t.Errorf("Values are not equal!")
			}
		})
	}
}

func TestMarshalContinuity(t *testing.T) {
	type args struct {
		values []Value
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "strings test",
			args: args{
				values: []Value{
					MakeString("ala ma kota"),
					MakeString("ala ma psa"),
					MakeString("bartek ma papugi"),
				},
			},
		},
		{
			name: "ints test",
			args: args{
				values: []Value{
					MakeInt(math.MinInt64),
					MakeInt(-456),
					MakeInt(-189),
					MakeInt(-10),
					MakeInt(0),
					MakeInt(1),
					MakeInt(3),
					MakeInt(17),
					MakeInt(24287),
					MakeInt(math.MaxInt64),
				},
			},
		},
		{
			name: "floats test",
			args: args{
				values: []Value{
					// TODO - nie dziala dla ujemnych :<
					MakeFloat(-124.0001),
					MakeFloat(-123.9998),
					MakeFloat(-1.01),
					MakeFloat(0.0000001),
					MakeFloat(1.01),
					MakeFloat(3),
					MakeFloat(2345.5432),
					MakeFloat(24287.111111),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argCount := len(tt.args.values)
			marshaled := make([][]byte, argCount)

			for i := 0; i < argCount; i++ {
				bytes := sortedMarshal(&tt.args.values[i])
				marshaled[i] = bytes
			}

			if !isIncreasing(marshaled) {
				t.Errorf("The marshaled values aren't increasing!")
			}
		})
	}
}

/* Auxiliary functions */
func isIncreasing(b [][]byte) bool {
	for i := 0; i < len(b)-1; i++ {
		if compareByteSlices(b[i], b[i+1]) != -1 {
			return false
		}
	}

	return true
}

func compareByteSlices(b1, b2 []byte) int {
	b1Length := len(b1)
	b2Length := len(b2)
	length := intMin(b1Length, b2Length)

	for i := 0; i < length; i++ {
		if b1[i] < b2[i] {
			return -1
		} else if b1[i] > b2[i] {
			return 1
		}
	}

	if b1Length > b2Length {
		return 1
	} else if b2Length == b1Length {
		return 0
	} else {
		return -1
	}
}

func intMin(x, y int) int {
	if x >= y {
		return y
	}

	return x
}
