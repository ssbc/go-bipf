package bipf_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/boreq/go-bipf"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	var testCases = []struct {
		Name   string
		Thing  any
		Binary string
	}{
		{
			Name:   "byte",
			Thing:  byte(100),
			Binary: "2264000000",
		},
		{
			Name:   "byte_pointer",
			Thing:  p(byte(100)),
			Binary: "2264000000",
		},
		{
			Name:   "int8",
			Thing:  int8(100),
			Binary: "2264000000",
		},
		{
			Name:   "int8_pointer",
			Thing:  p(int8(100)),
			Binary: "2264000000",
		},
		{
			Name:   "uint8",
			Thing:  uint8(100),
			Binary: "2264000000",
		},
		{
			Name:   "uint8_pointer",
			Thing:  p(uint8(100)),
			Binary: "2264000000",
		},
		{
			Name:   "int64",
			Thing:  int64(100),
			Binary: "2264000000",
		},
		{
			Name:   "int64_pointer",
			Thing:  p(int64(100)),
			Binary: "2264000000",
		},
		{
			Name:   "uint64",
			Thing:  uint64(100),
			Binary: "2264000000",
		},
		{
			Name:   "uint64_pointer",
			Thing:  p(uint64(100)),
			Binary: "2264000000",
		},
		{
			Name:   "int32",
			Thing:  int32(100),
			Binary: "2264000000",
		},
		{
			Name:   "int32_pointer",
			Thing:  p(int32(100)),
			Binary: "2264000000",
		},
		{
			Name:   "uint32",
			Thing:  uint32(100),
			Binary: "2264000000",
		},
		{
			Name:   "uint32_pointer",
			Thing:  p(uint32(100)),
			Binary: "2264000000",
		},
		{
			Name:   "int64",
			Thing:  int64(100),
			Binary: "2264000000",
		},
		{
			Name:   "int64_pointer",
			Thing:  p(int64(100)),
			Binary: "2264000000",
		},
		{
			Name:   "uint64",
			Thing:  uint64(100),
			Binary: "2264000000",
		},
		{
			Name:   "uint64_pointer",
			Thing:  p(uint64(100)),
			Binary: "2264000000",
		},
		{
			Name:   "int",
			Thing:  int(100),
			Binary: "2264000000",
		},
		{
			Name:   "int_pointer",
			Thing:  p(int(100)),
			Binary: "2264000000",
		},
		{
			Name:   "uint",
			Thing:  uint(100),
			Binary: "2264000000",
		},
		{
			Name:   "uint_pointer",
			Thing:  p(uint(100)),
			Binary: "2264000000",
		},

		{
			Name:   "float32",
			Thing:  float32(1.5),
			Binary: "43000000000000f83f",
		},
		{
			Name:   "float32_pointer",
			Thing:  p(float32(1.5)),
			Binary: "43000000000000f83f",
		},
		{
			Name:   "float64",
			Thing:  float64(1.5),
			Binary: "43000000000000f83f",
		},
		{
			Name:   "float64_pointer",
			Thing:  p(float64(1.5)),
			Binary: "43000000000000f83f",
		},

		{
			Name:   "true",
			Thing:  true,
			Binary: "0e01",
		},
		{
			Name:   "true_pointer",
			Thing:  p(true),
			Binary: "0e01",
		},
		{
			Name:   "false",
			Thing:  false,
			Binary: "0e00",
		},
		{
			Name:   "false_pointer",
			Thing:  p(false),
			Binary: "0e00",
		},

		{
			Name:   "nil_map",
			Thing:  map[int]struct{}(nil),
			Binary: "06",
		},
		{
			Name:   "nil_map_pointer",
			Thing:  p(map[int]struct{}(nil)),
			Binary: "06",
		},
		{
			Name:   "nil_slice",
			Thing:  []struct{}(nil),
			Binary: "06",
		},
		{
			Name:   "nil_slice_pointer",
			Thing:  p([]struct{}(nil)),
			Binary: "06",
		},
		{
			Name:   "nil",
			Thing:  nil,
			Binary: "06",
		},

		{
			Name:   "bytes",
			Thing:  []byte{0xDE, 0xAD, 0xBE, 0xEF},
			Binary: "21deadbeef",
		},

		{
			Name:   "empty_string",
			Thing:  "",
			Binary: "00",
		},
		{
			Name:   "empty_string_pointer",
			Thing:  p(""),
			Binary: "00",
		},
		{
			Name:   "string_hello",
			Thing:  "hello",
			Binary: "2868656c6c6f",
		},
		{
			Name:   "string_möterhead",
			Thing:  "möterhead",
			Binary: "506dc3b674657268656164",
		},

		{
			Name:   "empty array []",
			Thing:  [0]struct{}{},
			Binary: "04",
		},
		{
			Name:   "empty slice []",
			Thing:  []struct{}{},
			Binary: "04",
		},

		{
			Name:   "array of numbers 1 to 9",
			Thing:  [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			Binary: "ec02220100000022020000002203000000220400000022050000002206000000220700000022080000002209000000",
		},
		{
			Name:   "slice of numbers 1 to 9",
			Thing:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			Binary: "ec02220100000022020000002203000000220400000022050000002206000000220700000022080000002209000000",
		},

		{
			Name:   "empty object struct",
			Thing:  struct{}{},
			Binary: "05",
		},
		{
			Name:   "empty object map",
			Thing:  map[string]int{},
			Binary: "05",
		},
		{
			Name: "object with one key-value pair struct",
			Thing: struct {
				Foo bool `bipf:"foo"`
			}{Foo: true},
			Binary: "3518666f6f0e01",
		},
		{
			Name:   "object with one key-value pair map",
			Thing:  map[string]bool{"foo": true},
			Binary: "3518666f6f0e01",
		},
		{
			Name:   "{\"1\":true}",
			Thing:  map[string]bool{"1": true},
			Binary: "2508310e01",
		},

		{
			Name: "array with various types",
			Thing: []any{
				-1,
				map[string]bool{
					"foo": true,
				},
				map[string]any{
					"data": []int{222, 173, 190, 239},
				},
				map[string]any{
					"type": "Buffer",
				},
			},
			Binary: "b40322ffffffff3518666f6f0e01dd012064617461a40122de00000022ad00000022be00000022ef00000065207479706530427566666572",
		},

		// todo order of map keys is not guaranteed so this test is hard to write
		//{
		//	Name: "package.json",
		//	Json: "7b226e664d65223a2262697066222c226465736372697074696f6e223a2262696e61727920696e2d706c66436520666f726d6174222c2276657273696f6e223a22312e352e31222c22686f6d6570664765223a2268747470733a2f2f6769746875622e636f6d2f737362632f62697066222c227265706f7369746f7279223a7b2274797065223a22676974222c2275726c223a226769743a2f2f6769746875622e636f6d2f737362632f626970662e676974227d2c22646570656e64656e63696573223a7b22766172696e74223a225e352e302e30227d2c22646576446570656e64656e63696573223a7b2266664b6572223a225e352e352e31222c2274617065223a225e342e392e30227d2c2273637269707473223a7b2274657374223a226e6f646520746573742f696e6465782e6a73202626206e6f646520746573742f636f6d706172652e6a73202626206e6f646520746573742f66697874757265732e6a73227d2c22617574686f72223a22446f6d696e69632054617272203c646f6d696e69632e7461727240676d66496c2e636f6d3e2028687474703a2f2f646f6d696e6963746172722e636f6d29222c226c6963656e7365223a224d4954227d",
		//	Thing: map[string]any{
		//		"name":        "bipf",
		//		"description": "binary in-place format",
		//		"version":     "1.5.1",
		//		"homepage":    "https://github.com/ssbc/bipf",
		//		"repository": map[string]string{
		//			"type": "git",
		//			"url":  "git://github.com/ssbc/bipf.git",
		//		},
		//		"dependencies": map[string]string{
		//			"varint": "^5.0.0",
		//		},
		//		"devDependencies": map[string]string{
		//			"faker": "^5.5.1",
		//			"tape":  "^4.9.0",
		//		},
		//		"scripts": map[string]string{
		//			"test": "node test/index.js && node test/compare.js && node test/fixtures.js",
		//		},
		//		"author":  "Dominic Tarr <dominic.tarr@gmail.com> (http://dominictarr.com)",
		//		"license": "MIT",
		//	},
		//	Binary: "dd18206e664d652062697066586465736372697074696f6eb00642696e61727920696e2d706c66436520666f726d61743876657273696f6e28312e352e3140686f6d6570664765e00648747470733a2f2f6769746875622e636f6d2f737362632f62697066507265706f7369746f7279ed022074797065186769741875726cf0064769743a2f2f6769746875622e636f6d2f737362632f626970662e67697460646570656e64656e636965737530766172696e74305e352e302e3078646576446570656e64656e63696573cd012866664b6572305e352e352e312074617065305e342e392e303873637269707473d504207465737498046e6f646520746573742f696e6465782e6a73202626206e6f646520746573742f636f6d706172652e6a73202626206e6f646520746573742f66697874757265732e6a7330617574686f72f003446f6d696e69632054617272203c646f6d696e69632e7461727240676d66496c2e636f6d3e2028687474703a2f2f646f6d696e6963746172722e636f6d29386c6963656e7365184d4954",
		//},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, err := bipf.Marshal(testCase.Thing)
			require.NoError(t, err)

			gotHex := hex.EncodeToString(got)

			require.Equal(t, testCase.Binary, gotHex)
		})
	}
}

func TestUnmarshalNative(t *testing.T) {
	t.Run("byte", func(t *testing.T) {
		var v byte
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, byte(100), v)
	})

	t.Run("byte_pointer", func(t *testing.T) {
		var v *byte
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, p(byte(100)), v)
	})

	t.Run("int8", func(t *testing.T) {
		var v int8
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, int8(100), v)
	})

	t.Run("uint8", func(t *testing.T) {
		var v uint8
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, uint8(100), v)
	})

	t.Run("int64", func(t *testing.T) {
		var v int64
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, int64(100), v)
	})

	t.Run("uint64", func(t *testing.T) {
		var v uint64
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, uint64(100), v)
	})

	t.Run("int32", func(t *testing.T) {
		var v int32
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, int32(100), v)
	})

	t.Run("uint32", func(t *testing.T) {
		var v uint32
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, uint32(100), v)
	})

	t.Run("int64", func(t *testing.T) {
		var v int64
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, int64(100), v)
	})

	t.Run("uint64", func(t *testing.T) {
		var v uint64
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, uint64(100), v)
	})

	t.Run("int", func(t *testing.T) {
		var v int
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, int(100), v)
	})

	t.Run("uint", func(t *testing.T) {
		var v uint
		err := bipf.Unmarshal(h("2264000000"), &v)
		require.NoError(t, err)
		require.Equal(t, uint(100), v)
	})

	t.Run("float32", func(t *testing.T) {
		var v float32
		err := bipf.Unmarshal(h("43000000000000f83f"), &v)
		require.NoError(t, err)
		require.Equal(t, float32(1.5), v)
	})

	t.Run("float64", func(t *testing.T) {
		var v float64
		err := bipf.Unmarshal(h("43000000000000f83f"), &v)
		require.NoError(t, err)
		require.Equal(t, float64(1.5), v)
	})

	t.Run("true", func(t *testing.T) {
		var v bool
		err := bipf.Unmarshal(h("0e01"), &v)
		require.NoError(t, err)
		require.Equal(t, true, v)
	})

	t.Run("false", func(t *testing.T) {
		var v bool
		err := bipf.Unmarshal(h("0e00"), &v)
		require.NoError(t, err)
		require.Equal(t, false, v)
	})

	t.Run("empty string", func(t *testing.T) {
		var v string
		err := bipf.Unmarshal(h("00"), &v)
		require.NoError(t, err)
		require.Equal(t, "", v)
	})

	t.Run("string: hello", func(t *testing.T) {
		var v string
		err := bipf.Unmarshal(h("2868656c6c6f"), &v)
		require.NoError(t, err)
		require.Equal(t, "hello", v)
	})

	t.Run("string: möterhead", func(t *testing.T) {
		var v string
		err := bipf.Unmarshal(h("506dc3b674657268656164"), &v)
		require.NoError(t, err)
		require.Equal(t, "möterhead", v)
	})

	t.Run("nil map", func(t *testing.T) {
		var v map[int]struct{}
		err := bipf.Unmarshal(h("06"), &v)
		require.NoError(t, err)
		require.Equal(t, (map[int]struct{})(nil), v)
	})

	t.Run("nil", func(t *testing.T) {
		var v *string
		err := bipf.Unmarshal(h("06"), &v)
		require.NoError(t, err)
		require.Equal(t, (*string)(nil), v)
	})

	t.Run("bytes", func(t *testing.T) {
		var v []byte
		err := bipf.Unmarshal(h("21deadbeef"), &v)
		require.NoError(t, err)
		require.Equal(t, []byte{0xDE, 0xAD, 0xBE, 0xEF}, v)
	})
}

func TestUnmarshalSlicesAndArrays(t *testing.T) {
	t.Run("nil slice", func(t *testing.T) {
		var v []struct{}
		err := bipf.Unmarshal(h("06"), &v)
		require.NoError(t, err)
		require.Equal(t, ([]struct{})(nil), v)
	})

	t.Run("empty slice", func(t *testing.T) {
		var v []struct{}
		err := bipf.Unmarshal(h("04"), &v)
		require.NoError(t, err)
		require.Equal(t, ([]struct{})(nil), v)
	})

	t.Run("empty slice", func(t *testing.T) {
		var v []struct{}
		err := bipf.Unmarshal(h("04"), &v)
		require.NoError(t, err)
		require.Equal(t, ([]struct{})(nil), v)
	})

	t.Run("slice of numbers 1 to 9", func(t *testing.T) {
		var v []int
		err := bipf.Unmarshal(h("ec02220100000022020000002203000000220400000022050000002206000000220700000022080000002209000000"), &v)
		require.NoError(t, err)
		require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, v)
	})

	t.Run("array of numbers 1 to 9", func(t *testing.T) {
		var v [9]int
		err := bipf.Unmarshal(h("ec02220100000022020000002203000000220400000022050000002206000000220700000022080000002209000000"), &v)
		require.NoError(t, err)
		require.Equal(t, [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, v)
	})
}

func TestUnmarshalStructs(t *testing.T) {
	t.Run("empty object struct", func(t *testing.T) {
		v := struct{}{}

		err := bipf.Unmarshal(h("05"), &v)
		require.NoError(t, err)
	})

	t.Run("empty object map", func(t *testing.T) {
		v := struct{}{}

		err := bipf.Unmarshal(h("05"), &v)
		require.NoError(t, err)
	})

	t.Run("object with one key-value pair struct", func(t *testing.T) {
		v := struct {
			Foo bool `bipf:"foo"`
		}{}

		err := bipf.Unmarshal(h("3518666f6f0e01"), &v)
		require.NoError(t, err)
	})

	t.Run("object with one key-value pair map", func(t *testing.T) {
		v := map[string]bool{}

		err := bipf.Unmarshal(h("3518666f6f0e01"), &v)
		require.NoError(t, err)
	})
}

func TestComplexStruct(t *testing.T) {
	v := newComplexStruct()

	bipfBytes, err := bipf.Marshal(v)
	require.NoError(t, err)

	expected := newComplexStruct()
	// since BIPF only encodes int32 and we don't have type information in the
	// struct when using "any" this field will unmarshal as int32
	expected.Map["4"] = int32(expected.Map["4"].(int64))
	expected.Struct.Map["4"] = int32(expected.Struct.Map["4"].(int64))
	(*expected.MapPtr)["4"] = int32((*expected.MapPtr)["4"].(int64))
	(*expected.Struct.MapPtr)["4"] = int32((*expected.Struct.MapPtr)["4"].(int64))
	// since BIPF only encodes float64 and we don't have type information in the
	// struct when using "any" this field will unmarshal as float64
	expected.Map["5"] = float64(expected.Map["5"].(float32))
	(*expected.MapPtr)["5"] = float64((*expected.MapPtr)["5"].(float32))
	expected.Struct.Map["5"] = float64(expected.Struct.Map["5"].(float32))
	(*expected.Struct.MapPtr)["5"] = float64((*expected.Struct.MapPtr)["5"].(float32))

	var unmarshaled complexStruct
	err = bipf.Unmarshal(bipfBytes, &unmarshaled)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(expected, unmarshaled))
}

func TestWithoutTagsNamesAreCaseInsensitive(t *testing.T) {
	a := struct {
		Name string `bipf:"name"`
	}{
		Name: "a",
	}

	aBytes, err := bipf.Marshal(a)
	require.NoError(t, err)

	b := struct {
		Name string `bipf:"Name"`
	}{
		Name: "b",
	}

	bBytes, err := bipf.Marshal(b)
	require.NoError(t, err)

	require.NotEqual(t, aBytes, bBytes)

	t.Run("a", func(t *testing.T) {
		target := struct {
			Name string
		}{}

		err = bipf.Unmarshal(aBytes, &target)
		require.EqualValues(t, "a", target.Name)

	})

	t.Run("b", func(t *testing.T) {
		target := struct {
			Name string
		}{}

		err = bipf.Unmarshal(bBytes, &target)
		require.EqualValues(t, "b", target.Name)
	})
}

func TestComplexStructIntoAny(t *testing.T) {
	v := newComplexStruct()

	bipfBytes, err := bipf.Marshal(v)
	require.NoError(t, err)

	expected := map[any]any{
		"SimpleString":    "string",
		"HardString":      "łó\"śð",
		"SimpleStringPtr": "string",
		"HardStringPtr":   "łó\"śð",

		"Int32":    int32(123),
		"Int64":    int32(123),
		"Int32Ptr": int32(123),
		"Int64Ptr": int32(123),

		"Float32":    float64(float32(1.123)),
		"Float64":    float64(1.123),
		"Float32Ptr": float64(float32(1.123)),
		"Float64Ptr": float64(1.123),
		"Map": map[any]any{
			"1": "string",
			"2": "łó\"śð",
			"3": int32(123),
			"4": int32(123),
			"5": float64(float32(1.123)),
			"6": float64(1.123),
		},
		"MapPtr": map[any]any{
			"1": "string",
			"2": "łó\"śð",
			"3": int32(123),
			"4": int32(123),
			"5": float64(float32(1.123)),
			"6": float64(1.123),
		},

		"Slice":    []any{"string", "łó\"śð", int32(123), float64(1.123)},
		"SlicePtr": []any{"string", "łó\"śð", int32(123), float64(1.123)},

		"Bytes":    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		"BytesPtr": []byte{0xDE, 0xAD, 0xBE, 0xEF},

		"Struct": map[any]any{
			"SimpleString": "string",
			"HardString":   "łó\"śð",

			"Int32": int32(123),
			"Int64": int32(123),

			"Float32": float64(float32(1.123)),
			"Float64": float64(1.123),

			"Map": map[any]any{
				"1": "string",
				"2": "łó\"śð",
				"3": int32(123),
				"4": int32(123),
				"5": float64(float32(1.123)),
				"6": float64(1.123),
			},
			"MapPtr": map[any]any{
				"1": "string",
				"2": "łó\"śð",
				"3": int32(123),
				"4": int32(123),
				"5": float64(float32(1.123)),
				"6": float64(1.123),
			},

			"Slice":    []any{"string", "łó\"śð", int32(123), float64(1.123)},
			"SlicePtr": []any{"string", "łó\"śð", int32(123), float64(1.123)},

			"Bytes":    []byte{0xDE, 0xAD, 0xBE, 0xEF},
			"BytesPtr": []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
	}

	var unmarshaled any
	err = bipf.Unmarshal(bipfBytes, &unmarshaled)
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(expected, unmarshaled))
}

func BenchmarkSimpleStruct(b *testing.B) {
	v := newSimpleStruct()

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}

	bipfBytes, err := bipf.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("json_marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(v)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("json_unmarshal", func(b *testing.B) {
		var target simpleStruct
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(jsonBytes, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("bipf_marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := bipf.Marshal(v)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("bipf_unmarshal", func(b *testing.B) {
		var target simpleStruct
		for i := 0; i < b.N; i++ {
			err := bipf.Unmarshal(bipfBytes, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkComplexStruct(b *testing.B) {
	v := newComplexStruct()

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}

	bipfBytes, err := bipf.Marshal(v)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("json_marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(v)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("json_unmarshal", func(b *testing.B) {
		var target complexStruct
		for i := 0; i < b.N; i++ {
			err := json.Unmarshal(jsonBytes, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("bipf_marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := bipf.Marshal(v)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("bipf_unmarshal", func(b *testing.B) {
		var target complexStruct
		for i := 0; i < b.N; i++ {
			err := bipf.Unmarshal(bipfBytes, &target)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

type simpleStruct struct {
	String  string
	Int64   int64
	Float64 float64
	Slice   []string
	Bytes   []byte
}

func newSimpleStruct() simpleStruct {
	str := "string"
	int_64 := int64(123)
	float_64 := float64(1.123)

	return simpleStruct{
		String:  str,
		Int64:   int_64,
		Float64: float_64,
		Slice:   []string{str},
		Bytes:   []byte{0xDE, 0xAD, 0xBE, 0xEF},
	}
}

type complexStruct struct {
	SimpleString    string
	HardString      string
	SimpleStringPtr *string
	HardStringPtr   *string

	Int32    int32
	Int64    int64
	Int32Ptr *int32
	Int64Ptr *int64

	Float32    float32
	Float64    float64
	Float32Ptr *float32
	Float64Ptr *float64

	Map    map[string]any
	MapPtr *map[string]any

	Slice    []any
	SlicePtr *[]any

	Bytes    []byte
	BytesPtr *[]byte

	Struct struct {
		SimpleString string
		HardString   string

		Int32 int32
		Int64 int64

		Float32 float32
		Float64 float64

		Map    map[string]any
		MapPtr *map[string]any

		Slice    []any
		SlicePtr *[]any

		Bytes    []byte
		BytesPtr *[]byte
	}
}

func newComplexStruct() complexStruct {
	simpleString := "string"
	hardString := "łó\"śð"
	int_32 := int32(123)
	int_64 := int64(123)
	float_32 := float32(1.123)
	float_64 := float64(1.123)

	m := func() map[string]any {
		return map[string]any{
			"1": simpleString,
			"2": hardString,
			"3": int_32,
			"4": int_64,
			"5": float_32,
			"6": float_64,
		}
	}

	return complexStruct{
		SimpleString:    simpleString,
		HardString:      hardString,
		SimpleStringPtr: p(simpleString),
		HardStringPtr:   p(hardString),

		Int32:    int_32,
		Int64:    int_64,
		Int32Ptr: p(int_32),
		Int64Ptr: p(int_64),

		Float32:    float_32,
		Float64:    float_64,
		Float32Ptr: p(float_32),
		Float64Ptr: p(float_64),

		Map:      m(),
		MapPtr:   p(m()),
		Slice:    []any{simpleString, hardString, int_32, float_64},
		SlicePtr: p([]any{simpleString, hardString, int_32, float_64}),

		Bytes:    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		BytesPtr: p([]byte{0xDE, 0xAD, 0xBE, 0xEF}),

		Struct: struct {
			SimpleString string
			HardString   string
			Int32        int32
			Int64        int64
			Float32      float32
			Float64      float64
			Map          map[string]any
			MapPtr       *map[string]any
			Slice        []any
			SlicePtr     *[]any
			Bytes        []byte
			BytesPtr     *[]byte
		}{
			SimpleString: simpleString,
			HardString:   hardString,
			Int32:        int_32,
			Int64:        int_64,
			Float32:      float_32,
			Float64:      float_64,
			Map:          m(),
			MapPtr:       p(m()),
			Slice:        []any{simpleString, hardString, int_32, float_64},
			SlicePtr:     p([]any{simpleString, hardString, int_32, float_64}),
			Bytes:        []byte{0xDE, 0xAD, 0xBE, 0xEF},
			BytesPtr:     p([]byte{0xDE, 0xAD, 0xBE, 0xEF}),
		},
	}
}

func ExampleMarshal() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}

	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	b, err := bipf.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Println(hex.EncodeToString(b))
	// Output:
	// 9d031049442201000000204e616d65205265647330436f6c6f7273c401384372696d736f6e185265642052756279304d61726f6f6e
}

func ExampleUnmarshal() {
	bipfBlob, err := hex.DecodeString("9d031049442201000000204e616d65205265647330436f6c6f7273c401384372696d736f6e185265642052756279304d61726f6f6e")
	if err != nil {
		fmt.Println("error:", err)
	}

	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}

	var group ColorGroup

	err = bipf.Unmarshal(bipfBlob, &group)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", group)
	// Output:
	// {ID:1 Name:Reds Colors:[Crimson Red Ruby Maroon]}
}

func h(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func p[T any](v T) *T {
	return &v
}
