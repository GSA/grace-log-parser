package login

import (
	"strings"
	"testing"

	"github.com/GSA/grace-log-parser/handler/modules"
)

func TestLogin(t *testing.T) {
	l := &Login{}
	if err := l.Begin(); err != nil {
		t.Fatal(err)
		return
	}

	tt := map[string]struct {
		evt    *modules.Event
		detail string
		err    error
	}{
		"test1": {
			evt:    test1evt,
			detail: test1detail,
			err:    nil,
		},
		"test2": {
			evt:    test2evt,
			detail: test2detail,
			err:    nil,
		},
		"test3": {
			evt:    test3evt,
			detail: test3detail,
			err:    nil,
		},
		"test4": {
			evt:    test4evt,
			detail: test4detail,
			err:    nil,
		},
		"test5": {
			evt:    test5evt,
			detail: "",
			err:    &modules.NotApplicableErr{},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			d, err := l.Process(tc.evt)
			if err != nil && !modules.IsNotApplicable(err) {
				t.Fatalf("error value not expected: %v", err)
				return
			}
			if err != nil && modules.IsNotApplicable(err) && modules.IsNotApplicable(tc.err) {
				return
			}
			a := stripWhitespace(d)
			b := stripWhitespace(tc.detail)
			if a != b {
				t.Fatalf("detail is invalid, expected %q got %q", b, a)
				return
			}
		})
	}
}

func stripWhitespace(s0 string) string {
	s1 := strings.Replace(s0, "\n", "", -1)
	s2 := strings.Replace(s1, "\t", "", -1)
	return s2
}

var (
	test1evt    = &modules.Event{Type: "ConsoleLogin"}
	test1detail = `<table></table>`

	test2evt = &modules.Event{
		Type: "ConsoleLogin",
		UserIdentity: &modules.UserIdentity{
			AccessKeyID: "accessKeyID",
		},
	}
	test2detail = `<table><tr><th colspan="2">UserIdentity</th></tr><tr>
					<td>AccessKeyID</td><td>accessKeyID</td></tr></table>`

	test3evt = &modules.Event{
		Type: "ConsoleLogin",
		UserIdentity: &modules.UserIdentity{
			SessionContext: &modules.SessionContext{
				Attributes: &modules.SessionAttributes{
					MFAAuthenticated: "true",
				},
			},
		},
	}
	test3detail = `<table><tr><th colspan="2">UserIdentity</th></tr>
					<tr><td>MFAAuthenticated</td><td>true</td></tr></table>`

	test4evt = &modules.Event{
		Type: "ConsoleLogin",
		UserIdentity: &modules.UserIdentity{
			AccessKeyID: "accessKeyID",
			AccountID:   "accountID",
			PrincipalID: "principalID",
			Arn:         "arn",
			Type:        "type",
			UserName:    "username",
			SessionContext: &modules.SessionContext{
				Attributes: &modules.SessionAttributes{
					MFAAuthenticated: "true",
				},
			},
		},
	}
	test4detail = `<table><tr><th colspan="2">UserIdentity</th></tr><tr><td>AccessKeyID</td><td>accessKeyID</td></tr>
					<tr><td>AccountID</td><td>accountID</td></tr><tr><td>Arn</td><td>arn</td></tr><tr><td>MFAAuthenticated</td><td>true</td></tr>
					<tr><td>PrincipalID</td><td>principalID</td></tr><tr><td>Type</td><td>type</td></tr>
					<tr><td>UserName</td><td>username</td></tr></table>`

	test5evt = &modules.Event{
		Type: "NotConsoleLogin",
	}
)
