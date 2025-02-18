package c5

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"testing"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/mabels/c5-envelope/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type SimpleEnvelopeSuite struct {
	suite.Suite
	mockedSvalFn *SvalFnMock
}

func (s *SimpleEnvelopeSuite) SetupTest() {
	s.mockedSvalFn = &SvalFnMock{}
}

var mtimer = &mockedTimer{}

// ####################
// ## sortKeys tests ##
// ####################
func (s *SimpleEnvelopeSuite) TestSimpleTAsLiteral() {
	props := SimpleEnvelopeProps{
		T: 4711,
		Data: map[string]interface{}{
			"kind": "Kind",
			"data": map[string]interface{}{"Hallo": 1},
		},
	}
	n := NewSimpleEnvelope(&props)
	assert.Equal(s.T(), n.AsEnvelope().Data.Kind, "Kind")
	assert.Equal(s.T(), n.AsEnvelope().Data.Data, map[string]interface{}{"Hallo": 1})
	assert.Equal(s.T(), n.AsEnvelope().T, float64(4711))
}

func (s *SimpleEnvelopeSuite) TestSimpleTAsObj() {
	now := time.Now()
	pay := PayloadT{}
	FromDictPayloadT(map[string]interface{}{
		"kind": "Kind",
		"data": map[string]interface{}{"Hallo": 1},
	}, &pay)
	props := SimpleEnvelopeProps{
		T:    now,
		Data: pay,
	}
	n := NewSimpleEnvelope(&props)
	assert.Equal(s.T(), n.AsEnvelope().Data.Kind, "Kind")
	assert.Equal(s.T(), n.AsEnvelope().Data.Data, map[string]interface{}{"Hallo": 1})
	assert.Equal(s.T(), n.AsEnvelope().T, float64(now.UnixMilli()))
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithString() {
	strValue := "string"
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		return assert.Equal(s.T(), strValue, (sVal.val.(JsonValType)).Val)
	}))
	SortKeys(strValue, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithDate() {
	t := time.Now()
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		return assert.Equal(s.T(), t, (sVal.val.(JsonValType)).Val)
	}))
	SortKeys(t, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithNumber() {
	n := 78
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		return assert.Equal(s.T(), n, (sVal.val.(JsonValType)).Val)
	}))
	SortKeys(n, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithBoolean() {
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		return assert.Equal(s.T(), true, (sVal.val.(JsonValType)).Val)
	}))
	SortKeys(true, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithArrayOfEmpty() {
	var emptySlice []int
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_START, sVal.outState.String())
		} else if funcCallIdx == 2 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_END, sVal.outState.String())
		}
		funcCallIdx++
		return result
	}))
	SortKeys(emptySlice, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithArrayOf_1_2() {
	ar := []int{1, 2}
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_START, sVal.outState.String())
		} else if funcCallIdx == 2 {
			result = assert.Equal(s.T(), ar[0], (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 3 {
			result = assert.Equal(s.T(), ar[1], (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 4 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_END, sVal.outState.String())
		}
		funcCallIdx++
		return result
	}))
	SortKeys(ar, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithArrayOf_1_2_3_4() {
	ar := [][]int{{1, 2}, {3, 4}}
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 || funcCallIdx == 2 || funcCallIdx == 6 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_START, sVal.outState.String())
		} else if funcCallIdx == 3 {
			result = assert.Equal(s.T(), ar[0][0], (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 4 {
			result = assert.Equal(s.T(), ar[0][1], (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 5 || funcCallIdx == 9 || funcCallIdx == 10 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), ARRAY_END, sVal.outState.String())
		} else if funcCallIdx == 7 {
			result = assert.Equal(s.T(), ar[1][0], (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 8 {
			result = assert.Equal(s.T(), ar[1][1], (sVal.val.(JsonValType)).Val)
		}
		funcCallIdx++
		return result
	}))
	SortKeys(ar, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithObjOfEmptyObj() {
	var obj struct{}
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_START, sVal.outState.String())
		} else if funcCallIdx == 2 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_END, sVal.outState.String())
		}
		funcCallIdx++
		return result
	}))
	SortKeys(obj, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithObjOfObj_Y_1_X_2() {
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_START, sVal.outState.String())
		} else if funcCallIdx == 2 {
			result = assert.Equal(s.T(), "x", sVal.attribute)
		} else if funcCallIdx == 3 {
			result = assert.Equal(s.T(), 2, (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 4 {
			result = assert.Equal(s.T(), "y", sVal.attribute)
		} else if funcCallIdx == 5 {
			result = assert.Equal(s.T(), 1, (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 6 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_END, sVal.outState.String())
		}
		funcCallIdx++
		return result
	}))
	SortKeys(struct {
		Y int `json:"y"`
		X int `json:"x"`
	}{Y: 1, X: 2}, s.mockedSvalFn.Execute)
}

func (s *SimpleEnvelopeSuite) TestSortWithOutWithObjOfObj_Y_B_1_A_2() {
	funcCallIdx := 1
	s.mockedSvalFn.On("Execute", mock.MatchedBy(func(sVal SVal) bool {
		result := false
		if funcCallIdx == 1 || funcCallIdx == 3 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_START, sVal.outState.String())
		} else if funcCallIdx == 2 {
			result = assert.Equal(s.T(), "y", sVal.attribute)
		} else if funcCallIdx == 4 {
			result = assert.Equal(s.T(), "a", sVal.attribute)
		} else if funcCallIdx == 5 {
			result = assert.Equal(s.T(), 2, (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 6 {
			result = assert.Equal(s.T(), "b", sVal.attribute)
		} else if funcCallIdx == 7 {
			result = assert.Equal(s.T(), 1, (sVal.val.(JsonValType)).Val)
		} else if funcCallIdx == 8 || funcCallIdx == 9 {
			result = assert.Equal(s.T(), nil, sVal.val)
			result = result && assert.Equal(s.T(), OBJECT_END, sVal.outState.String())
		}
		funcCallIdx++
		return result
	}))

	type Obj struct {
		B int `json:"b"`
		A int `json:"a"`
	}
	SortKeys(struct {
		Y Obj `json:"y"`
	}{Y: Obj{
		B: 1,
		A: 2,
	}}, s.mockedSvalFn.Execute)
}

// #########################
// ## JSONCollector tests ##
// #########################
func (s *SimpleEnvelopeSuite) TestJSONCollectorEmptyObj() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	var obj struct{}
	SortKeys(obj, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "{}", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorEmptyArray() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	var emptySlice []int
	SortKeys(emptySlice, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollector_X_Y_1_Z_x_Y_Z() {
	type Obj struct {
		Y int    `json:"y"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}

	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	SortKeys(struct {
		X Obj      `json:"x"`
		Y struct{} `json:"y"`
		Z []int    `json:"z"`
	}{
		X: Obj{
			Y: 1,
			Z: "x",
		},
		Y: emptypObj,
		Z: emptySlice,
	}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "{\"x\":{\"y\":1,\"z\":\"x\"},\"y\":{},\"z\":[]}", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorArray_xx() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	SortKeys([]string{"xx"}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[\"xx\"]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorArray_1_2() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	SortKeys([]interface{}{1, "2"}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[1,\"2\"]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollector_1_2_A() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	SortKeys([]interface{}{1, []string{"2", "A"}, "E"}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[1,[\"2\",\"A\"],\"E\"]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorIndent2EmptyObj() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, NewJsonProps(2, ""))
	var obj struct{}
	SortKeys(obj, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "{}", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorIndent2ArrayEmpty() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, NewJsonProps(2, ""))
	var emptySlice []int
	SortKeys(emptySlice, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollectorIndent2_X_Y_1_Z_x() {
	type Obj struct {
		Y int    `json:"y"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}

	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, NewJsonProps(2, ""))
	SortKeys(struct {
		X Obj      `json:"x"`
		Y struct{} `json:"y"`
		Z []int    `json:"z"`
	}{
		X: Obj{
			Y: 1,
			Z: "x",
		},
		Y: emptypObj,
		Z: emptySlice,
	}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "{\n  \"x\": {\n    \"y\": 1,\n    \"z\": \"x\"\n  },\n  \"y\": {},\n  \"z\": []\n}", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollector_Indent2_xx() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, NewJsonProps(2, ""))
	SortKeys([]string{"xx"}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[\n  \"xx\"\n]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollector_Indent2_array_1_2() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, NewJsonProps(2, ""))
	SortKeys([]interface{}{1, "2"}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[\n  1,\n  \"2\"\n]", out)
}

func (s *SimpleEnvelopeSuite) TestJSONCollector_1_date444() {
	out := ""
	col := NewJsonCollector(func(str string) {
		out += str
	}, nil)
	SortKeys([]interface{}{1, time.UnixMilli(444).UTC()}, func(prob SVal) {
		col.Append(prob)
	})
	assert.Equal(s.T(), "[1,\"1970-01-01T00:00:00.444Z\"]", out)
}

// #########################
// ## HashCollector tests ##
// #########################
func (s *SimpleEnvelopeSuite) TestHashCollector_date() {
	h := NewHashCollector()
	SortKeys(time.UnixMilli(444).UTC(), func(prob SVal) {
		h.Append(prob)
	})
	assert.Equal(s.T(), "DzYqv3YaniBJWwqrNBn4534oTe4nL14TqcfVCguf9Yyv", h.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_X_1_Y2() {
	h := NewHashCollector()
	SortKeys(struct {
		X int
		Y int
	}{
		X: 1,
		Y: 2,
	}, func(prob SVal) {
		h.Append(prob)
	})
	assert.Equal(s.T(), "DkPs9C3fYabdDLFxqMh4ZoTNr3xD1xYGFvYnJioF7V6H", h.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_1() {
	h := NewHashCollector()
	type Obj struct {
		Y int    `json:"y"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}
	SortKeys(struct {
		X Obj       `json:"x"`
		Y struct{}  `json:"y"`
		Z []int     `json:"z"`
		D time.Time `json:"d"`
	}{
		X: Obj{
			Y: 1,
			Z: "x",
		},
		Y: emptypObj,
		Z: emptySlice,
		D: time.UnixMilli(444).UTC(),
	}, func(prob SVal) {
		h.Append(prob)
	})
	assert.Equal(s.T(), "5PvJAWGkaKAHax6tsaKGfPYm6JfXxZs15wRTDpSKaZ2G", h.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_2() {
	h := NewHashCollector()
	type Obj struct {
		Y int    `json:"y"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}
	SortKeys(struct {
		X    Obj       `json:"x"`
		Y    struct{}  `json:"y"`
		Z    []int     `json:"z"`
		Date time.Time `json:"date"`
	}{
		X: Obj{
			Y: 2,
			Z: "x",
		},
		Y:    emptypObj,
		Z:    emptySlice,
		Date: time.UnixMilli(444).UTC(),
	}, func(prob SVal) {
		h.Append(prob)
	})
	assert.Equal(s.T(), "ECVWfmcNaUGkgvPZe7CojrnRNULxNczKXU8PGns6UDvr", h.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_3() {
	h := NewHashCollector()
	type Obj struct {
		X int    `json:"x"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}
	SortKeys(struct {
		X    Obj       `json:"x"`
		Y    struct{}  `json:"y"`
		Z    []int     `json:"z"`
		Date time.Time `json:"date"`
	}{
		X: Obj{
			X: 1,
			Z: "x",
		},
		Y:    emptypObj,
		Z:    emptySlice,
		Date: time.UnixMilli(444).UTC(),
	}, func(prob SVal) {
		h.Append(prob)
	})
	assert.Equal(s.T(), "EoYNGMtap1k9iEAGeVtHmJwpMjQLKWJmR27SG6aC9fSg", h.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_4() {
	h1 := NewHashCollector()
	type Obj struct {
		X int    `json:"x"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}
	SortKeys(struct {
		X    Obj       `json:"x"`
		Y    struct{}  `json:"y"`
		Z    []int     `json:"z"`
		Date time.Time `json:"date"`
	}{
		X: Obj{
			X: 1,
			Z: "x",
		},
		Y:    emptypObj,
		Z:    emptySlice,
		Date: time.UnixMilli(444).UTC(),
	}, func(prob SVal) {
		h1.Append(prob)
	})

	h2 := NewHashCollector()
	SortKeys(struct {
		Date time.Time `json:"date"`
		X    Obj       `json:"x"`
		Y    struct{}  `json:"y"`
		Z    []int     `json:"z"`
	}{
		Date: time.UnixMilli(444).UTC(),
		X: Obj{
			X: 1,
			Z: "x",
		},
		Y: emptypObj,
		Z: emptySlice,
	}, func(prob SVal) {
		h2.Append(prob)
	})

	assert.Equal(s.T(), h1.Digest(), h2.Digest())
}

func (s *SimpleEnvelopeSuite) TestHashCollector_3_InternalUpdate() {
	hashCalculator := sha256.New()

	type Obj struct {
		R int    `json:"r"`
		Z string `json:"z"`
	}
	var emptySlice []int
	var emptypObj struct{}
	expectedArgs := []string{"date", "1970-01-01T00:00:00.444Z", "x", "r", "1", "z", "u", "y", "z"}
	mck := &mocks.Hash{}
	idx := 0
	mck.On("Write", mock.MatchedBy(func(p []byte) bool {
		hashCalculator.Write(p)
		idx++
		return expectedArgs[idx-1] == string(p)
	})).Return(1, nil)

	t := struct {
		X    Obj       `json:"x"`
		Y    struct{}  `json:"y"`
		Z    []int     `json:"z"`
		Date time.Time `json:"date"`
	}{
		X: Obj{
			R: 1,
			Z: "u",
		},
		Y:    emptypObj,
		Z:    emptySlice,
		Date: time.UnixMilli(444).UTC(),
	}
	collector := &HashCollector{mck}
	SortKeys(t, func(prob SVal) {
		collector.Append(prob)
	})

	var nilBytes []byte
	mck.On("Sum", nilBytes).Return([]byte{})
	collector.Digest()
	mck.AssertNumberOfCalls(s.T(), "Sum", 1)

	assert.Equal(s.T(), "CwEMjUHV6BpDS7AGBAYqjY6qMKE6xC8Z56H5T2ZuUuXe", base58.Encode(hashCalculator.Sum(nil)))
}

func (s *SimpleEnvelopeSuite) TestSimpleHash() {
	type Data struct {
		Name string `json:"name"`
		Date string `json:"date"`
	}

	type KindData struct {
		Kind string `json:"kind"`
		Data Data   `json:"data"`
	}

	hashC := NewHashCollector()
	SortKeys(KindData{
		Kind: "test",
		Data: Data{
			Name: "object",
			Date: "2021-05-20",
		},
	}, func(sval SVal) {
		hashC.Append(sval)
	})
	assert.Equal(s.T(), "5zWhdtvKuGob1FbW9vUGPQKobcLtYYr5wU8AxQRVraeB", hashC.Digest())
}

// ##########################
// ## SimpleEnvelope tests ##
// ##########################
func (s *SimpleEnvelopeSuite) TestSerialization() {
	typ := SampleNameDate{}
	FromDictSampleNameDate(map[string]interface{}{
		"name": "object",
		"date": "2021-05-20",
	}, &typ)
	// s.Assertions.Equal(typ.ToDict(), []string{})
	props := &SimpleEnvelopeProps{
		ID:  "1624140000000-4a2a6fb97b3afe6a7ca4c13457c441664c7f6a6c2ea7782e1f2dea384cf97cb8",
		Src: "test case",
		Data: PayloadT1{
			Data: typ.ToDict(),
			Kind: "test",
		},
		Dst: []string{},
		T:   time.UnixMilli(444),
		TTL: 10,
	}
	se := NewSimpleEnvelope(props)
	assert.JSONEq(s.T(), *se.AsJson(), `{"data":{"data":{"date":"2021-05-20","name":"object"},"kind":"test"},"dst":[],"id":"1624140000000-4a2a6fb97b3afe6a7ca4c13457c441664c7f6a6c2ea7782e1f2dea384cf97cb8","src":"test case","t":444,"ttl":10,"v":"A"}`)
}

type mockedTimer struct{}

func (*mockedTimer) Now() time.Time {
	return time.UnixMilli(1624140000000)
}

func (s *SimpleEnvelopeSuite) TestSerializationWithHash() {
	typ := SampleNameDate{}
	FromDictSampleNameDate(map[string]interface{}{
		"name": "object",
		"date": "2021-05-20",
	}, &typ)
	props := &SimpleEnvelopeProps{
		Src: "test case",
		Data: PayloadT1{
			Kind: "test",
			Data: typ.ToDict(),
		},
		Dst:           []string{},
		TTL:           10,
		TimeGenerator: mtimer,
	}
	se := NewSimpleEnvelope(props)
	assert.JSONEq(s.T(), *se.AsJson(), `{"data":{"data":{"date":"2021-05-20","name":"object"},"kind":"test"},"dst":[],"id":"1624140000000-BbYxQMurpUmj1W6E4EwYM79Rm3quSz1wwtNZDSsFt1bp","src":"test case","t":1624140000000,"ttl":10,"v":"A"}`)
}

func (s *SimpleEnvelopeSuite) TestSerializationWithIndent() {
	b := []byte(`{"data":{"data":{"date":"2021-05-20","name":"object"},"kind":"test"},"dst":[],"id":"1624140000000-BbYxQMurpUmj1W6E4EwYM79Rm3quSz1wwtNZDSsFt1bp","src":"test case","t":1624140000000,"ttl":10,"v":"A"}`)
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	assert.NoError(s.T(), err)

	typ := SampleNameDate{}
	FromDictSampleNameDate(map[string]interface{}{
		"name": "object",
		"date": "2021-05-20",
	}, &typ)
	props := &SimpleEnvelopeProps{
		Src: "test case",
		Data: PayloadT1{
			Kind: "test",
			Data: typ.ToDict(),
		},
		Dst:           []string{},
		TTL:           10,
		JsonProp:      NewJsonProps(2, ""),
		TimeGenerator: mtimer,
	}
	se := NewSimpleEnvelope(props)
	assert.Equal(s.T(), *se.AsJson(), out.String())
}

func (s *SimpleEnvelopeSuite) TestMissingDataInEnvelope() {
	typ := SampleY{Y: 4}
	message := &SimpleEnvelopeProps{
		Src: "test case",
		Data: PayloadT1{
			Kind: "kind",
			Data: typ.ToDict(),
		},
		TimeGenerator: mtimer,
	}
	se := NewSimpleEnvelope(message)

	var ref EnvelopeT
	assert.NoError(s.T(), json.Unmarshal([]byte(*se.AsJson()), &ref))

	env := NewSimpleEnvelope(&SimpleEnvelopeProps{
		ID:            ref.ID,
		Src:           ref.Src,
		Dst:           ref.Dst,
		T:             time.UnixMilli(int64(ref.T)),
		TTL:           int(ref.TTL),
		Data:          ref.Data,
		JsonProp:      nil,
		TimeGenerator: mtimer,
	})

	envData := env.AsEnvelope()
	assert.Equal(s.T(), message.Data.(PayloadT1).Kind, envData.Data.Kind)

	yEnv := EnvelopeT{}
	ok := FromDictEnvelopeT(env.AsEnvelope().ToDict(), &yEnv)
	//fmt.Fprintln(os.Stderr, ok)
	// fmt.Fprintln(os.Stderr, yEnv)
	assert.Nil(s.T(), ok)

	mapVal := env.AsEnvelope().ToDict()["data"].(map[string]interface{})["data"].(map[string]interface{})
	// assert.True(s.T(), ok)

	yVal := SampleY{}
	FromDictSampleY(yEnv.Data.Data, &yVal)
	assert.EqualValues(s.T(), yVal.Y, mapVal["y"])
}

func TestSimpleEnvelopeSuite(t *testing.T) {
	suite.Run(t, new(SimpleEnvelopeSuite))
}

type SvalFnMock struct {
	mock.Mock
}

// Execute provides a mock function with given fields: prob
func (_m *SvalFnMock) Execute(prob SVal) {
	_m.Called(prob)
}
