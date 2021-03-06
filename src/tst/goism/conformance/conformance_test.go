package conformance

import (
	"testing"
	"tst/goism"
)

func init() {
	goism.LoadPackage("conformance")
}

func Test1Ops(t *testing.T) {
	testCalls(t, goism.CallTests{
		// Int ops.
		"add1Int 1":     "2",
		"addInt 1 2 3":  "6",
		"sub1Int 1":     "0",
		"subInt 3 2 1":  "0",
		"mulInt 2 2 20": "80",
		"quoInt 20 2 2": "5",
		"gtInt 2 1":     "t",
		"ltInt 2 1":     "nil",
		"incInt 2":      "3",
		"decInt 2":      "1",
		// Float ops.
		"add1Float 1.0":         "2.0",
		"addFloat 1.1 2.2 3.3":  "6.6",
		"sub1Float 1.5":         "0.5",
		"subFloat 3.5 2.5 1.5":  "-0.5",
		"mulFloat 2.0 2.0 20.0": "80.0",
		"quoFloat 20.0 2.0 2.0": "5.0",
		"gtFloat 2.0 0":         "t",
		"ltFloat 2.0 0":         "nil",
		"incFloat 2.0":          "3.0",
		"decFloat 2.0":          "1.0",
		// String ops.
		`concatStr "a" "b" "c"`: `"abc"`,
		`ltStr "abc" "foo"`:     "t",
		`ltStr "foo" "abc"`:     "nil",
		`eqStr "foo" "foo"`:     "t",
		`eqStr "foo" "bar"`:     "nil",
	})
}

func Test2GlobalVars(t *testing.T) {
	table := []struct {
		name          string
		valueExpected string
	}{
		{"var1", "1"},
		{"var2", "2"},
		{"var3", "3"},
		{"var4", "4"},
		{"var5", "5"},
		{"var6", "6"},
	}

	for _, row := range table {
		res := evalVar(row.name)
		if res != row.valueExpected {
			t.Errorf("%s=>%s (want %s)", row.name, res, row.valueExpected)
		}
	}
}

func Test3MultiResult(t *testing.T) {
	table := []struct {
		call           string
		names          []string
		valuesExpected []string
	}{
		{"return2", []string{"r2_1", "r2_2"}, q("a", "b")},
		{"return3", []string{"r3_2", "r3_3"}, q("b", "c")},
		{"return4", []string{"r4_1", "r4_3", "r4_4"}, q("a", "c", "d")},
	}

	for _, row := range table {
		evalCall(row.call) // For side effects
		for i, name := range row.names {
			res := evalVar(name)
			if res != row.valuesExpected[i] {
				t.Errorf("%s=>%s (want %s)", name, res, row.valuesExpected[i])
			}
		}
	}
}

func Test4Goto(t *testing.T) {
	testCalls(t, goism.CallTests{
		"testGoto 10":              "10",
		"testGotoOutBlock 10":      "10",
		"testGotoTwice 10":         "10",
		"testGotoChain 10":         "10",
		"testGotoBack 10":          "10",
		"testGotoScopes1 10":       "10",
		"testGotoScopes2 10":       "10",
		"testGotoScopes3 10":       "10",
		"testGotoBackAndScopes 10": "10",
	})
}

func Test5If(t *testing.T) {
	testCalls(t, goism.CallTests{
		"testIfTrue 10":       "10",
		"testIfFalse 10":      "10",
		"testIfZero 0":        "t",
		"testIfZero 1":        "nil",
		"testIfElse1 0":       `"0"`,
		"testIfElse1 1":       `"1"`,
		"testIfElse1 2":       `"2"`,
		"testIfElse1 3":       `"x"`,
		"testIfElse2 0":       `"0"`,
		"testIfElse2 1":       `"1"`,
		"testIfElse2 2":       `"2"`,
		"testIfElse2 3":       `"x"`,
		"testNestedIfZero 0":  "t",
		"testNestedIfZero 1":  "nil",
		"testIfInitDef 10":    "10",
		"testIfInitAssign 10": "10",
		"testAnd t t t":       "t",
		"testAnd t t nil":     "nil",
		"testAnd nil nil nil": "nil",
		"testOr t t t":        "t",
		"testOr t t nil":      "t",
		"testOr nil nil nil":  "nil",
	})
}

func Test6Arrays(t *testing.T) {
	testCalls(t, goism.CallTests{
		"testArrayLit 10": "10",
		// {"testKeyedArrayLit 10", "10"}, #REFS: 73
		"testArrayZeroVal 10":      "10",
		"testArrayUpdate 10":       "10",
		"testArrayCopyOnAssign 10": "10",
	})
}

func Test7Switch(t *testing.T) {
	testCalls(t, goism.CallTests{
		"stringifyInt3 0": `"0"`,
		"stringifyInt3 1": `"1"`,
		"stringifyInt3 2": `"2"`,
		"stringifyInt3 3": `"x"`,
		"stringifyInt4 0": `"0"`,
		"stringifyInt4 1": `"1"`,
		"stringifyInt4 2": `"2"`,
		"stringifyInt4 3": `"x"`,
	})
}

func Test8Slices(t *testing.T) {
	testCalls(t, goism.CallTests{
		"sliceLen $sliceOf3":   "3",
		"sliceLen $sliceOf4_5": "4",
		"sliceCap $sliceOf3":   "3",
		"sliceCap $sliceOf4_5": "5",
	})
}

func Test9For(t *testing.T) {
	testCalls(t, goism.CallTests{
		`testFor 10`:               "10",
		`testFor 1`:                "1",
		`testWhile 10`:             "10",
		`testWhile 1`:              "1",
		`testForBreak 10`:          "10",
		`testForBreak 1`:           "1",
		`testForContinue 10`:       "10",
		`testForContinue 1`:        "1",
		`testNestedFor 10`:         "10",
		`testNestedFor 1`:          "1",
		`testNestedForBreak 10`:    "10",
		`testNestedForBreak 1`:     "1",
		`testNestedForContinue 10`: "10",
		`testNestedForContinue 1`:  "1",
		`testForScopes1 10`:        "10",
		`testForScopes2 10`:        "10",
		`testNestedForScopes1 10`:  "10",
		`testNestedForScopes2 10`:  "10",
	})
}

func Test10Range(t *testing.T) {
	testCalls(t, goism.CallTests{
		"sumArray1": "6",
	})
}

func Test11Maps(t *testing.T) {
	testCalls(t, goism.CallTests{
		"testMapMake 10":      "10",
		"testMapNilLookup 10": "10",
		"testMapUpdate 10":    "10",
		"testMapDelete 10":    "10",
		"testMapLen 10":       "10",
	})
}

func Test12Strings(t *testing.T) {
	testCalls(t, goism.CallTests{
		// Single-byte strings.
		`stringGet "abc" 0`: chr('a'),
		`stringGet "abc" 1`: chr('b'),
		`stringGet "abc" 2`: chr('c'),
		`stringLen "abc"`:   "3",
		`stringLen ""`:      "0",
		// Multi-byte strings.
		`stringLen "丂a♞a"`:   "8",
		`stringGet "丂a♞a" 0`: chr(rune("丂"[0])),
		`stringGet "丂a♞a" 1`: chr(rune("丂"[1])),
		`stringGet "丂a♞a" 2`: chr(rune("丂"[2])),
		`stringGet "丂a♞a" 3`: chr('a'),
		`stringGet "丂a♞a" 4`: chr(rune("♞"[0])),
		`stringGet "丂a♞a" 5`: chr(rune("♞"[1])),
		`stringGet "丂a♞a" 6`: chr(rune("♞"[2])),
		`stringGet "丂a♞a" 7`: chr('a'),
		`stringLen "𠜱a𠺢"`:    "9",
		`stringGet "𠜱a𠺢" 0`:  chr(rune("𠜱"[0])),
		`stringGet "𠜱a𠺢" 1`:  chr(rune("𠜱"[1])),
		`stringGet "𠜱a𠺢" 2`:  chr(rune("𠜱"[2])),
		`stringGet "𠜱a𠺢" 3`:  chr(rune("𠜱"[3])),
		`stringGet "𠜱a𠺢" 4`:  chr('a'),
		`stringGet "𠜱a𠺢" 5`:  chr(rune("𠺢"[0])),
		`stringGet "𠜱a𠺢" 6`:  chr(rune("𠺢"[1])),
		`stringGet "𠜱a𠺢" 7`:  chr(rune("𠺢"[2])),
		`stringGet "𠜱a𠺢" 8`:  chr(rune("𠺢"[3])),
	})
}

func TestCombined(t *testing.T) {
	testCalls(t, goism.CallTests{
		"factorial 0": "1",
		"factorial 2": "2",
		"factorial 3": "6",

		"isAlpha ?a": "t",
		"isAlpha ?g": "t",
		"isAlpha ?z": "t",
		"isAlpha ?A": "t",
		"isAlpha ?G": "t",
		"isAlpha ?Z": "t",
		"isAlpha ?@": "nil",
		"isAlpha ?1": "nil",
		"isAlpha 0":  "nil",

		"max4 0 0 0 0":               "0",
		"max4 1.1 1.2 1.3 1.4":       "1.4",
		"max4 -1 1 2 -2":             "2",
		"max4 0.1 0.01 0.001 0.0001": "0.1",

		`replace "hello" ?l ?d`: `"heddo"`,
		`replace "hello" ?a ?z`: `"hello"`,

		`substring "hello" -1 -1`: `"hello"`,
		`substring "hello" 0 0`:   `""`,
		`substring "hello" 1 -1`:  `"ello"`,
		`substring "hello" -1 4`:  `"hell"`,
		`substring "hello" 1 4`:   `"ell"`,
	})
}
