package sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

type TestCase struct {
	test           string
	expectedRaw    string
	expectedValues []string
}

func TestSqlAdaptor(t *testing.T) {
	g := NewGomegaWithT(t)
	t.Run("test sql adaptor success", func(t *testing.T) {
		type ExampleDBStruct struct {
			Identity uint    `filter:"%"`
			Name     string  `filter:"*"`
			Email    *string `filter:"#;="`
			Age      uint8   `filter:"=;>;"`
		}
		testCases := []TestCase{
			// Test * rule
			{
				test:           `(name=duckhue01 AND name%duckhue) OR name#"(duckhue01,duckhue02)"`,
				expectedRaw:    "((name=? AND name LIKE ?) OR name IN (?, ?))",
				expectedValues: []string{"duckhue01", "%duckhue%", "duckhue01", "duckhue02"},
			},
			// Test for an empty quoted string.
			{
				test:           "(name=\"\" AND email=bob-dylan@aol.com) OR age > 1",
				expectedRaw:    "((name=? AND email=?) OR age>?)",
				expectedValues: []string{"", "bob-dylan@aol.com", "1"},
			},
			// Test % rule
			{
				test:           `identity%1`,
				expectedRaw:    "identity LIKE ?",
				expectedValues: []string{"%1%"},
			},
			// Test # rule
			{
				test:           "email#duckhue",
				expectedRaw:    "email IN (?)",
				expectedValues: []string{"duckhue"},
			},
		}
		for _, testCase := range testCases {
			sa := NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
			response, err := sa.Parse(testCase.test)
			g.Expect(err).To(BeNil(), fmt.Sprintf("failed case: %s", testCase.test))
			g.Expect(response.Raw).To(Equal(testCase.expectedRaw), fmt.Sprintf("failed case raw: %s", testCase.test))
			g.Expect(response.Values).To(Equal(testCase.expectedValues), fmt.Sprintf("failed case values: %s", testCase.test))
		}
	})
	t.Run("test sql adaptor failure", func(t *testing.T) {
		type ExampleDBStruct struct {
			ID           uint
			Name         string  `filter:"=;>;%;#"`
			Email        *string `filter:"=;>"`
			Age          uint8   `filter:"=;>;#"`
			Birthday     *time.Time
			MemberNumber sql.NullString
			ActivatedAt  sql.NullTime
			CreatedAt    time.Time
			UpdatedAt    time.Time
		}
		testCases := []TestCase{
			{
				test: "(name=max AND invalidField=wow) OR age > 1",
			},
			{
				test: "id = wow",
			},
			{
				test: "age = wow",
			},
			{
				test: "name = default AND",
			},
			{
				test: "name",
			},
			{
				test: "name = default AND age",
			},
		}
		for _, testCase := range testCases {
			sa := NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
			_, err := sa.Parse(testCase.test)
			g.Expect(err).ToNot(BeNil(), fmt.Sprintf("failed case: %s", testCase.test))
		}
	})
	t.Run("test FieldParseValidatorFromStruct", func(t *testing.T) {
		type ExampleDBStruct struct {
			ID    uint
			Name  string  `filter:"=;>;%;#"`
			Email *string `filter:"=;>"`
			Age   uint8   `filter:"=;>;#"`
		}
		defaultFields := FieldParseValidatorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
		_, ok := defaultFields["name"]
		g.Expect(ok).To(Equal(true))
	})
}

func FuzzSQLAdaptor(f *testing.F) {
	type ExampleDBStruct struct {
		ID    uint
		Name  string  `filter:"=;>;%;#"`
		Email *string `filter:"=;>"`
		Age   uint8   `filter:"=;>;#"`
	}
	testcases := []string{"(name=max AND invalidField=wow) OR age > 1", "(name=max AND email=bob-dylan@aol.com) OR age > 1", "id = wow", "(name%max AND email=\"bob-dylan@aol.com\") OR age > 1"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, s string) {
		sa := NewDefaultAdaptorFromStruct(reflect.ValueOf(&ExampleDBStruct{}))
		response, err := sa.Parse(s)
		if err != nil && response != nil {
			t.Errorf("expected nil response when err is not nil")
		}
	})
}
