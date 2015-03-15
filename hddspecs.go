package main

type HDDSpecs struct {
	Date            string          `json:"date"`
	SerialNumber    string          `json:"serial_number"`
	Model           string          `json:"model"`
	CapacityBytes   int64           `json:"capacity_bytes"`
	Failure         int             `json:"failure"`
	SmartNormalized SmartNormalized `json:"smart_normalized",omitempty`
	SmartRaw        SmartRaw        `json:"smart_raw",omitempty`
}

type SmartNormalized struct {
	Value1 int `json:"1"`
	Value2 int `json:"2"`
	Value3 int `json:"3"`
	Value4 int `json:"4"`
	Value5 int `json:"5"`
	Value7 int `json:"7"`
	Value8 int `json:"8"`
	Value9 int `json:"9"`

	Value10 int `json:"10"`
	Value11 int `json:"11"`
	Value12 int `json:"12"`
	Value13 int `json:"13"`
	Value15 int `json:"15"`

	Value183 int `json:"183"`
	Value184 int `json:"184"`
	Value187 int `json:"185"`
	Value188 int `json:"188"`
	Value189 int `json:"189"`

	Value190 int `json:"190"`
	Value191 int `json:"191"`
	Value192 int `json:"192"`
	Value193 int `json:"193"`
	Value194 int `json:"194"`
	Value195 int `json:"195"`
	Value196 int `json:"196"`
	Value197 int `json:"197"`
	Value198 int `json:"198"`
	Value199 int `json:"199"`

	Value200 int `json:"200"`
	Value201 int `json:"201"`

	Value223 int `json:"223"`
	Value225 int `json:"225"`

	Value240 int `json:"240"`
	Value241 int `json:"241"`
	Value242 int `json:"242"`

	Value250 int `json:"250"`
	Value251 int `json:"251"`
	Value252 int `json:"252"`
	Value254 int `json:"254"`
	Value255 int `json:"255"`
}

type SmartRaw struct {
	Value1 int `json:"1"`
	Value2 int `json:"2"`
	Value3 int `json:"3"`
	Value4 int `json:"4"`
	Value5 int `json:"5"`
	Value7 int `json:"7"`
	Value8 int `json:"8"`
	Value9 int `json:"9"`

	Value10 int `json:"10"`
	Value11 int `json:"11"`
	Value12 int `json:"12"`
	Value13 int `json:"13"`
	Value15 int `json:"15"`

	Value183 int `json:"183"`
	Value184 int `json:"184"`
	Value187 int `json:"185"`
	Value188 int `json:"188"`
	Value189 int `json:"189"`

	Value190 int `json:"190"`
	Value191 int `json:"191"`
	Value192 int `json:"192"`
	Value193 int `json:"193"`
	Value194 int `json:"194"`
	Value195 int `json:"195"`
	Value196 int `json:"196"`
	Value197 int `json:"197"`
	Value198 int `json:"198"`
	Value199 int `json:"199"`

	Value200 int `json:"200"`
	Value201 int `json:"201"`

	Value223 int `json:"223"`
	Value225 int `json:"225"`

	Value240 int `json:"240"`
	Value241 int `json:"241"`
	Value242 int `json:"242"`

	Value250 int `json:"250"`
	Value251 int `json:"251"`
	Value252 int `json:"252"`
	Value254 int `json:"254"`
	Value255 int `json:"255"`
}
