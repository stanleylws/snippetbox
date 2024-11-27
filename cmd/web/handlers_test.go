package main

import (
	"net/http"
	"net/url"
	"testing"

	"snippetbox.stanley.net/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name string
		urlPath string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID",
			urlPath: "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name: "Non-existent ID",
			urlPath: "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName = "Bob"
		validPassword = "validPa$$word"
		validEmail = "bob@example.com"
		formTag = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests :=[]struct {
		name		string
		username	string
		email		string
		password	string
		csrfToken	string
		wantCode	int
		wantFormTag	string
	} {
		{
			name:		"Valid submission",
			username:	validName,
			email:		validEmail,
			password:	validPassword,
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusSeeOther,
		},
		{
			name:		"Invalid CSRF Token",
			username:	validName,
			email:		validEmail,
			password:	validPassword,
			csrfToken:	"wrongToken",
			wantCode:	http.StatusBadRequest,
		},
		{
			name:		"Empty name",
			username:	"",
			email:		validEmail,
			password:	validPassword,
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:		"Empty email",
			username:	validName,
			email:		"",
			password:	validPassword,
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:		"Empty password",
			username:	validName,
			email:		validEmail,
			password:	"",
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:		"Invalid email",
			username:	validName,
			email:		"bob@example.",
			password:	validPassword,	
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:		"Short password",
			username:	validName,
			email:		validEmail,
			password:	"pa$$",	
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
		{
			name:		"Duplicate email",
			username:	validName,
			email:		"dupe@example.com",
			password:	validPassword,	
			csrfToken:	validCSRFToken,
			wantCode:	http.StatusUnprocessableEntity,
			wantFormTag: formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}