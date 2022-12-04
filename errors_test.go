package scpi

import (
	"errors"
	"testing"
)

func TestInvalidProtocolError(t *testing.T) {
	err := InvalidProtocolError("foo")
	if got, want := err.Error(), "invalid protocol foo"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestInvalidFormatError(t *testing.T) {
	err := InvalidFormatError("foo")
	if got, want := err.Error(), "invalid format: foo"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestCommandError(t *testing.T) {
	err := &CommandError{
		cmd:  "foo",
		code: -101,
		msg:  "invalid character",
	}

	if got, want := err.Code(), -101; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
	if got, want := err.Error(), "'foo' returned -101: invalid character"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestConfirmError(t *testing.T) {
	tests := map[string]struct {
		in   map[string]string
		want error
	}{
		"NoError": {
			in: map[string]string{
				"cmd":    "*CLS",
				"errRes": "+0,\"No error\"",
			},
			want: nil,
		},
		"InvalidFormat": {
			in: map[string]string{
				"cmd":    "foo",
				"errRes": "foo, bar, baz",
			},
			want: InvalidFormatError("foo, bar, baz"),
		},
		"CommandError": {
			in: map[string]string{
				"cmd":    "foo",
				"errRes": "-101,\"Invalid character\"",
			},
			want: &CommandError{
				cmd:  "foo",
				code: -101,
				msg:  "invalid character",
			},
		},
		"CommandErrorWithoutQuotes": {
			in: map[string]string{
				"cmd":    "foo",
				"errRes": "-101, Invalid character",
			},
			want: &CommandError{
				cmd:  "foo",
				code: -101,
				msg:  "invalid character",
			},
		},
	}

	for n, tt := range tests {
		t.Run(n, func(t *testing.T) {
			err := confirmError(tt.in["cmd"], tt.in["errRes"])
			if got, want := err, tt.want; !(errors.Is(got, want)) {
				cast, _ := want.(*CommandError)
				if !errors.As(got, &cast) {
					t.Fatalf("got %+v, want %+v", got, want)
				}
			}
		})
	}
}
