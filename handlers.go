package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/skip2/go-qrcode"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.New("login").Parse(LoginTemplate)
	data := struct{ Error string }{Error: r.URL.Query().Get("error")}
	t.Execute(w, data)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.New("register").Parse(RegisterTemplate)
		t.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	_, exists := s.userStore.GetUser(email)
	if exists {
		http.Redirect(w, r, "/?error=既に登録済みのメールアドレスです", http.StatusSeeOther)
		return
	}

	user, err := s.userStore.CreateUser(email)
	if err != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	setupURL := fmt.Sprintf("/setup?email=%s&secret=%s",
		url.QueryEscape(email),
		url.QueryEscape(user.Secret))
	http.Redirect(w, r, setupURL, http.StatusSeeOther)
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
		http.Redirect(w, r, "/?error=ユーザーが見つかりません。先にアカウントを作成してください", http.StatusSeeOther)
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
	secret := r.URL.Query().Get("secret")
	email := r.URL.Query().Get("email")

	if email == "" || secret == "" {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	t, _ := template.New("setup").Parse(SetupTemplate)
	data := struct {
		Email  string
		Secret string
	}{
		Email:  email,
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

func (s *Server) qrHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email parameter required", http.StatusBadRequest)
		return
	}

	user, exists := s.userStore.GetUser(email)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	qrURL := user.getQRCodeURL("TOTPApp", email)
	
	qrCode, err := qrcode.New(qrURL, qrcode.Medium)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	qrCode.Write(256, w)
}