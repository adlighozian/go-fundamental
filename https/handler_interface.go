package https

import (
	"net/http"
)

type Handlers interface {
	GetPost() http.HandlerFunc
	DetailPost() http.HandlerFunc
	Register() http.HandlerFunc
	Login() http.HandlerFunc
	Logout() http.HandlerFunc
}
