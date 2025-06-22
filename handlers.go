package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("login").Parse(LoginTemplate)
	data := struct{ Error string }{Error: r.URL.Query().Get("error")}
	t.Execute(w, data)
}

func (s *Server) verifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	code := r.FormValue("code")

	if email == "" || code == "" {
		http.Redirect(w, r, "/?error=メールアドレスとTOTPコードを入力してください", http.StatusSeeOther)
		return
	}

	user, exists := s.userStore.GetUser(email)
	if !exists {
		var err error
		user, err = s.userStore.CreateUser(email)
		if err != nil {
			http.Redirect(w, r, "/?error=ユーザー作成エラー", http.StatusSeeOther)
			return
		}

		qrURL := user.getQRCodeURL("TOTPApp", email)

		setupURL := fmt.Sprintf("/setup?email=%s&qrUrl=%s&secret=%s",
			url.QueryEscape(email),
			url.QueryEscape(qrURL),
			url.QueryEscape(user.Secret))
		http.Redirect(w, r, setupURL, http.StatusSeeOther)
		return
	}

	ctx := context.Background()
	if user.Verify(ctx, code) {
		http.Redirect(w, r, "/success?email="+email, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/?error=TOTPコードが正しくありません", http.StatusSeeOther)
	}
}

func (s *Server) setupHandler(w http.ResponseWriter, r *http.Request) {
	qrUrl := r.URL.Query().Get("qrUrl")
	secret := r.URL.Query().Get("secret")
	email := r.URL.Query().Get("email")

	t, _ := template.New("setup").Parse(SetupTemplate)
	data := struct {
		Email  string
		QRUrl  string
		Secret string
	}{
		Email:  email,
		QRUrl:  qrUrl,
		Secret: secret,
	}
	t.Execute(w, data)
}

func (s *Server) successHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")

	t, _ := template.New("success").Parse(SuccessTemplate)
	data := struct{ Email string }{Email: email}
	t.Execute(w, data)
}