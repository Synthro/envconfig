// Copyright (c) 2013 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.

package envconfig

import (
	"os"
	"testing"
	"time"
)

type Specification struct {
	Embedded
	EmbeddedButIgnored           `envconfig:"-"`
	Debug                        bool
	Port                         int
	Rate                         float32
	User                         string
	TTL                          uint32
	Timeout                      time.Duration
	AdminUsers                   []string
	MagicNumbers                 []int
	MultiWordVar                 string
	SomePointer                  *string
	MultiWordVarWithAlt          string `envconfig:"MULTI_WORD_VAR_WITH_ALT"`
	MultiWordVarWithLowerCaseAlt string `envconfig:"multi_word_var_with_lower_case_alt"`
	NoPrefixWithAlt              string `envconfig:"SERVICE_HOST"`
	Ignored                      string `envconfig:"-"`
	Nested                       Nested
	NestedPointer                *Nested
}

type Embedded struct {
	Enabled             bool
	EmbeddedPort        int
	MultiWordVar        string
	MultiWordVarWithAlt string `envconfig:"MULTI_WITH_DIFFERENT_ALT"`
	EmbeddedAlt         string `envconfig:"EMBEDDED_WITH_ALT"`
	EmbeddedIgnored     string `ignored:"true"`
}

type EmbeddedButIgnored struct {
	FirstEmbeddedButIgnored  string
	SecondEmbeddedButIgnored string
}

type Nested struct {
	Child struct {
		Body    string
		Enabled bool
		Ignored string `envconfig:"-"`
	}
	ChildAlt struct {
		BodyAlt    string `envconfig:"body"`
		EnabledAlt bool   `envconfig:"enabled"`
		Ignored    string `envconfig:"-"`
	} `envconfig:"child_alt"`
	Body    string
	Enabled bool
	Ignored string `envconfig:"-"`
}

func TestProcess(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("ENV_CONFIG_TIMEOUT", "2m")
	os.Setenv("ENV_CONFIG_ADMINUSERS", "John,Adam,Will")
	os.Setenv("ENV_CONFIG_MAGICNUMBERS", "5,10,20")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	os.Setenv("ENV_CONFIG_TTL", "30")
	os.Setenv("ENV_CONFIG_IGNORED", "was-not-ignored")
	err := Process("env_config", &s)
	if err != nil {
		t.Error(err.Error())
	}
	if s.NoPrefixWithAlt != "127.0.0.1" {
		t.Errorf("expected %v, got %v", "127.0.0.1", s.NoPrefixWithAlt)
	}
	if !s.Debug {
		t.Errorf("expected %v, got %v", true, s.Debug)
	}
	if s.Port != 8080 {
		t.Errorf("expected %d, got %v", 8080, s.Port)
	}
	if s.Rate != 0.5 {
		t.Errorf("expected %f, got %v", 0.5, s.Rate)
	}
	if s.TTL != 30 {
		t.Errorf("expected %d, got %v", 30, s.TTL)
	}
	if s.User != "Kelsey" {
		t.Errorf("expected %s, got %s", "Kelsey", s.User)
	}
	if s.Timeout != 2*time.Minute {
		t.Errorf("expected %s, got %s", 2*time.Minute, s.Timeout)
	}
	if len(s.AdminUsers) != 3 ||
		s.AdminUsers[0] != "John" ||
		s.AdminUsers[1] != "Adam" ||
		s.AdminUsers[2] != "Will" {
		t.Errorf("expected %#v, got %#v", []string{"John", "Adam", "Will"}, s.AdminUsers)
	}
	if len(s.MagicNumbers) != 3 ||
		s.MagicNumbers[0] != 5 ||
		s.MagicNumbers[1] != 10 ||
		s.MagicNumbers[2] != 20 {
		t.Errorf("expected %#v, got %#v", []int{5, 10, 20}, s.MagicNumbers)
	}
	if s.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestParseErrorBool(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Debug" {
		t.Errorf("expected %s, got %v", "Debug", v.FieldName)
	}
	if s.Debug != false {
		t.Errorf("expected %v, got %v", false, s.Debug)
	}
}

func TestParseErrorFloat32(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_RATE", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Rate" {
		t.Errorf("expected %s, got %v", "Rate", v.FieldName)
	}
	if s.Rate != 0 {
		t.Errorf("expected %v, got %v", 0, s.Rate)
	}
}

func TestParseErrorInt(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_PORT", "string")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "Port" {
		t.Errorf("expected %s, got %v", "Port", v.FieldName)
	}
	if s.Port != 0 {
		t.Errorf("expected %v, got %v", 0, s.Port)
	}
}

func TestParseErrorUint(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_TTL", "-30")
	err := Process("env_config", &s)
	v, ok := err.(*ParseError)
	if !ok {
		t.Errorf("expected ParseError, got %v", v)
	}
	if v.FieldName != "TTL" {
		t.Errorf("expected %s, got %v", "TTL", v.FieldName)
	}
	if s.TTL != 0 {
		t.Errorf("expected %v, got %v", 0, s.TTL)
	}
}

func TestErrInvalidSpecification(t *testing.T) {
	m := make(map[string]string)
	err := Process("env_config", &m)
	if err != ErrInvalidSpecification {
		t.Errorf("expected %v, got %v", ErrInvalidSpecification, err)
	}
}

func TestUnsetVars(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("USER", "foo")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// If the var is not defined the non-prefixed version should not be used
	// unless the struct tag says so
	if s.User != "" {
		t.Errorf("expected %q, got %q", "", s.User)
	}
}

func TestAlternateVarNames(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_LOWER_CASE_ALT", "baz")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	// Setting the alt version of the var in the environment has no effect if
	// the struct tag is not supplied
	if s.MultiWordVar != "" {
		t.Errorf("expected %q, got %q", "", s.MultiWordVar)
	}

	// Setting the alt version of the var in the environment correctly sets
	// the value if the struct tag IS supplied
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %q, got %q", "bar", s.MultiWordVarWithAlt)
	}

	// Alt value is not case sensitive and is treated as all uppercase
	if s.MultiWordVarWithLowerCaseAlt != "baz" {
		t.Errorf("expected %q, got %q", "baz", s.MultiWordVarWithLowerCaseAlt)
	}
}

func TestPointerFieldBlank(t *testing.T) {
	var s Specification
	os.Clearenv()
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.SomePointer != nil {
		t.Errorf("expected <nil>, got %2", *s.SomePointer)
	}
}

func TestMustProcess(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_DEBUG", "true")
	os.Setenv("ENV_CONFIG_PORT", "8080")
	os.Setenv("ENV_CONFIG_RATE", "0.5")
	os.Setenv("ENV_CONFIG_USER", "Kelsey")
	os.Setenv("SERVICE_HOST", "127.0.0.1")
	MustProcess("env_config", &s)

	defer func() {
		if err := recover(); err != nil {
			return
		}

		t.Error("expected panic")
	}()
	m := make(map[string]string)
	MustProcess("env_config", &m)
}

func TestEmbeddedStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_ENABLED", "true")
	os.Setenv("ENV_CONFIG_EMBEDDEDPORT", "1234")
	os.Setenv("ENV_CONFIG_MULTIWORDVAR", "foo")
	os.Setenv("ENV_CONFIG_MULTI_WORD_VAR_WITH_ALT", "bar")
	os.Setenv("ENV_CONFIG_MULTI_WITH_DIFFERENT_ALT", "baz")
	os.Setenv("ENV_CONFIG_EMBEDDED_WITH_ALT", "foobar")
	os.Setenv("ENV_CONFIG_SOMEPOINTER", "foobaz")
	os.Setenv("ENV_CONFIG_EMBEDDED_IGNORED", "was-not-ignored")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if !s.Enabled {
		t.Errorf("expected %v, got %v", true, s.Enabled)
	}
	if s.EmbeddedPort != 1234 {
		t.Errorf("expected %d, got %v", 1234, s.EmbeddedPort)
	}
	if s.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.MultiWordVar)
	}
	if s.Embedded.MultiWordVar != "foo" {
		t.Errorf("expected %s, got %s", "foo", s.Embedded.MultiWordVar)
	}
	if s.MultiWordVarWithAlt != "bar" {
		t.Errorf("expected %s, got %s", "bar", s.MultiWordVarWithAlt)
	}
	if s.Embedded.MultiWordVarWithAlt != "baz" {
		t.Errorf("expected %s, got %s", "baz", s.Embedded.MultiWordVarWithAlt)
	}
	if s.EmbeddedAlt != "foobar" {
		t.Errorf("expected %s, got %s", "foobar", s.EmbeddedAlt)
	}
	if *s.SomePointer != "foobaz" {
		t.Errorf("expected %s, got %s", "foobaz", *s.SomePointer)
	}
	if s.EmbeddedIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestEmbeddedButIgnoredStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_FIRSTEMBEDDEDBUTIGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_SECONDEMBEDDEDBUTIGNORED", "was-not-ignored")
	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if s.FirstEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
	if s.SecondEmbeddedButIgnored != "" {
		t.Errorf("expected empty string, got %#v", s.Ignored)
	}
}

func TestEmbeddedDurationStruct(t *testing.T) {
	var s struct {
		Timeout struct {
			time.Duration
		}
	}
	os.Clearenv()
	os.Setenv("ENV_CONFIG_TIMEOUT", "5s")

	if err := Process("env_config", &s); err != nil {
		t.Fatal(err)
	}

	if s.Timeout.String() != "5s" {
		t.Errorf("expected %s, got %s", "5s", s.Timeout)
	}
}

func TestNestedStruct(t *testing.T) {
	var s Specification
	os.Clearenv()
	os.Setenv("ENV_CONFIG_NESTED_BODY", "body1")
	os.Setenv("ENV_CONFIG_NESTED_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTED_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_BODY", "body2")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_ALT_BODY", "body3")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_ALT_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTED_CHILD_ALT_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_BODY", "body1")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_BODY", "body2")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_IGNORED", "was-not-ignored")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_ALT_BODY", "body3")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_ALT_ENABLED", "true")
	os.Setenv("ENV_CONFIG_NESTEDPOINTER_CHILD_ALT_IGNORED", "was-not-ignored")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}
	if !s.Nested.Enabled {
		t.Errorf("expected %v, got %v", true, s.Nested.Enabled)
	}
	if s.Nested.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Nested.Ignored)
	}
	if s.Nested.Body != "body1" {
		t.Errorf("expected %s, got %s", "body1", s.Nested.Body)
	}
	if !s.Nested.Child.Enabled {
		t.Errorf("expected %v, got %v", true, s.Nested.Child.Enabled)
	}
	if s.Nested.Child.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Nested.Child.Ignored)
	}
	if s.Nested.Child.Body != "body2" {
		t.Errorf("expected %s, got %s", "body2", s.Nested.Child.Body)
	}
	if !s.Nested.ChildAlt.EnabledAlt {
		t.Errorf("expected %v, got %v", true, s.Nested.ChildAlt.EnabledAlt)
	}
	if s.Nested.ChildAlt.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.Nested.ChildAlt.Ignored)
	}
	if s.Nested.ChildAlt.BodyAlt != "body3" {
		t.Errorf("expected %s, got %s", "body3", s.Nested.ChildAlt.BodyAlt)
	}
	if s.NestedPointer == nil {
		t.Fatalf("expected non-nil pointer")
	}
	if !s.NestedPointer.Enabled {
		t.Errorf("expected %v, got %v", true, s.NestedPointer.Enabled)
	}
	if s.NestedPointer.Ignored != "" {
		t.Errorf("expected empty string, got %#v", s.NestedPointer.Ignored)
	}
	if s.NestedPointer.Body != "body1" {
		t.Errorf("expected %s, got %s", "body1", s.NestedPointer.Body)
	}
	if !s.NestedPointer.Child.Enabled {
		t.Errorf("expected %v, got %v", true, s.NestedPointer.Child.Enabled)
	}
	if s.NestedPointer.Child.Body != "body2" {
		t.Errorf("expected %s, got %s", "body2", s.NestedPointer.Child.Body)
	}
	if !s.NestedPointer.ChildAlt.EnabledAlt {
		t.Errorf("expected %v, got %v", true, s.NestedPointer.ChildAlt.EnabledAlt)
	}
	if s.NestedPointer.ChildAlt.BodyAlt != "body3" {
		t.Errorf("expected %s, got %s", "body3", s.NestedPointer.ChildAlt.BodyAlt)
	}
}

func TestNonPointerFailsProperly(t *testing.T) {
	var s Specification
	os.Clearenv()

	err := Process("env_config", s)
	if err != ErrInvalidSpecification {
		t.Errorf("non-pointer should fail with ErrInvalidSpecification, was instead %s", err)
	}
}

func TestCustomDecoder(t *testing.T) {
	s := struct {
		Foo string
		Bar bracketed
	}{}

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.Foo != "foo" {
		t.Errorf("foo: expected 'foo', got %q", s.Foo)
	}

	if string(s.Bar) != "[bar]" {
		t.Errorf("bar: expected '[bar]', got %q", string(s.Bar))
	}
}

func TestCustomDecoderWithPointer(t *testing.T) {
	s := struct {
		Foo string
		Bar *bracketed
	}{}

	// Decode would panic when b is nil, so make sure it
	// has an initial value to replace.
	var b bracketed = "initial_value"
	s.Bar = &b

	os.Clearenv()
	os.Setenv("ENV_CONFIG_FOO", "foo")
	os.Setenv("ENV_CONFIG_BAR", "bar")

	if err := Process("env_config", &s); err != nil {
		t.Error(err.Error())
	}

	if s.Foo != "foo" {
		t.Errorf("foo: expected 'foo', got %q", s.Foo)
	}

	if string(*s.Bar) != "[bar]" {
		t.Errorf("bar: expected '[bar]', got %q", string(*s.Bar))
	}
}

type bracketed string

func (b *bracketed) UnmarshalText(text []byte) error {
	*b = bracketed("[" + string(text) + "]")
	return nil
}
