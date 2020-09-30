package login

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/GSA/grace-log-parser/handler/modules"
	"github.com/GSA/grace-log-parser/handler/modules/email"
)

// call email.Register to register the email.login sub-module
func init() {
	l := &Login{}
	email.Register(l)
	fmt.Printf("registered %s sub-module", l.Name())
}

// Login handles processing ConsoleLogin events
type Login struct {
	tmpl *template.Template
}

// Name returns the name of this email sub-module
func (l *Login) Name() string {
	return "email.login"
}

// Begin is not implemented for Login
func (l *Login) Begin() error {
	data := `<table>
{{ if .UserIdentity }}
	<tr><th colspan="2">UserIdentity</th></tr>
	{{ if gt (len .UserIdentity.AccessKeyID) 0 }}<tr><td>AccessKeyID</td><td>{{.UserIdentity.AccessKeyID}}</td></tr>{{end}}
	{{ if gt (len .UserIdentity.AccountID) 0 }}<tr><td>AccountID</td><td>{{.UserIdentity.AccountID}}</td></tr>{{end}}
	{{ if gt (len .UserIdentity.Arn) 0 }}<tr><td>Arn</td><td>{{.UserIdentity.Arn}}</td></tr>{{end}}
{{- if .UserIdentity.SessionContext }}{{ if .UserIdentity.SessionContext.Attributes }}
	{{ if .UserIdentity.SessionContext.Attributes }}<tr><td>MFAAuthenticated</td><td>{{.UserIdentity.SessionContext.Attributes.MFAAuthenticated}}</td></tr>{{end}}
{{ end }}{{ end -}}
	{{ if gt (len .UserIdentity.PrincipalID) 0 }}<tr><td>PrincipalID</td><td>{{.UserIdentity.PrincipalID}}</td></tr>{{end}}
	{{ if gt (len .UserIdentity.Type) 0 }}<tr><td>Type</td><td>{{.UserIdentity.Type}}</td></tr>{{end}}
	{{ if gt (len .UserIdentity.UserName) 0 }}<tr><td>UserName</td><td>{{.UserIdentity.UserName}}</td></tr>{{end}}
{{ end }}
</table>`
	t := template.New("login")
	tmpl, err := t.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse login template: %v", err)
	}
	l.tmpl = tmpl
	return nil
}

// Process returns the custom payload for a particular matched event
// or returns a modules.NotApplicableErr
func (l *Login) Process(evt *modules.Event) (string, error) {
	if evt.Type != "ConsoleLogin" {
		return "", modules.NotApplicableErr{}
	}

	buf := &bytes.Buffer{}
	err := l.tmpl.Execute(buf, evt)
	if err != nil {
		return "", fmt.Errorf("error creating login detail: %v", err)
	}

	return buf.String(), nil
}

// End is not implemented for Login
func (l *Login) End() error {
	return nil
}
