package middleware

import(
	"log"
	"time"
	"errors"
	"regexp"
	"net/http"

	"github.com/govwa/util/config"
	"github.com/govwa/user/session"
	"github.com/julienschmidt/httprouter"
)


type Class struct{}

func New()*Class{
	return &Class{}
}

func(self *Class) LoggingMiddleware(h httprouter.Handle) httprouter.Handle{
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
		
	}
}

func (this *Class) AuthCheck(h httprouter.Handle) httprouter.Handle {
	var sess = session.New()
 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !sess.IsLoggedIn(r) {
			redirect := config.Fullurl + "login"
			http.Redirect(w, r, redirect, http.StatusSeeOther)
			return
		}

		h(w, r, ps)
	}
}

func (this *Class)CapturePanic(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var err error
		defer func() {
            r := recover()
            if r != nil {
                switch t := r.(type) {
                case string:
					err = errors.New(t)
                case error:
					err = t
                default:
                    err = errors.New("Unknown error")
				}
				log.Println(err.Error())
                http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h(w, r, ps)
	}
}

func(this *Class)DetectSQLMap(h httprouter.Handle)httprouter.Handle{
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		userAgent := r.Header.Get("User-Agent")
		sqlmapDetected, _ := regexp.MatchString("sqlmap*", userAgent)
		if sqlmapDetected{
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			log.Printf("sqlmap detect ")
			return
		}else{
			h(w, r, ps)
		}
	}
}
