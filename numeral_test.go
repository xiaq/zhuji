package zhuji

import "testing"

var numeralCases = []struct {
	num       int64
	s         string
	canonical bool
}{
	{-100, "负一百", true},
	{0, "零", true},
	{1, "一", true},
	{10, "十", true},
	{11, "十一", true},
	{42, "四十二", true},
	{1001, "一千零一", true},
	{1234, "一千二百三十四", true},
	{20401, "二万零四百零一", true},
	{123456789, "一亿二千三百四十五万六千七百八十九", true},
	{22222, "两万两千两百两十两", false},
	{31, "卅一", false},
	{29, "廿九", false},
}

func TestNumeral(t *testing.T) {
	for _, tc := range numeralCases {
		num, size := ParseNumeral(tc.s)
		if num != tc.num || size != len(tc.s) {
			t.Errorf("ParseNumeral(\"%s\") => (%d, %d), want (%d, %d)",
				tc.s, num, size, tc.num, len(tc.s))
		}
		if tc.canonical {
			s := ToNumeral(tc.num)
			if s != tc.s {
				t.Errorf("ToNumeral(%d) => %s, want %s", tc.num, s, tc.s)
			}
		}
	}
}
