package http

import (
	"net/http"
	"github.com/open-falcon/mail-provider/config"
	"github.com/toolkits/web/param"
	"gopkg.in/gomail.v2"
	"crypto/tls"
	"strings"
	)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		m := gomail.NewMessage()
		m.SetHeader("From", cfg.Smtp.From)
		m.SetHeader("To", strings.Split(tos, ",")...)//send email to multipul persons
		m.SetHeader("Subject", param.MustString(r, "subject"))
		m.SetBody("text/html", param.MustString(r, "content"))
		d := gomail.Dialer{Host: cfg.Smtp.Addr, Port: cfg.Smtp.Port, Username: cfg.Smtp.Username, Password: cfg.Smtp.Password, SSL:false}
		d.TLSConfig = &tls.Config{InsecureSkipVerify: cfg.Smtp.InsecureSkipVerify}
		if err := d.DialAndSend(m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
