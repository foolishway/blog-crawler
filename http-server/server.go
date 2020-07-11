package http_server

import (
	"blog-crawler/models"
	"blog-crawler/robot"
	"blog-crawler/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type HttpServer struct {
	Port string
}

const sessionKey string = "blog-crawler"

var tpl *template.Template
var store = sessions.NewCookieStore([]byte(os.Getenv("BLOG_CRAWLER_FOR_TAL")))

func (hs *HttpServer) StartServer() {
	//static server
	serveFile()
	tpl = template.Must(
		template.ParseGlob(
			"./static/views/*.html",
		),
	)

	//index
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(w, r)
	})

	//duty
	http.HandleFunc("/duty/", authWapper(dutyHandler))

	//login
	http.HandleFunc("/login/", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r)
	})

	//edit duty
	http.HandleFunc("/dutyedit/", authWapper(editDutyHandler))
	//share
	http.HandleFunc("/share", shareHandler)
	err := http.ListenAndServe(":"+hs.Port, nil)
	if err != nil {
		log.Fatalf("Http server listen error %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" && r.RequestURI != "/index.html" && r.RequestURI != "/index.htm" {
		return
	}
	articles := models.GetAllArticles()
	randomArticles := make([]models.Article, len(articles))
	if itfs := utils.RandomSlice(articles); len(itfs) > 0 {
		for index, itf := range itfs {
			if newArticle, ok := itf.(models.Article); ok {
				randomArticles[index] = newArticle
			}
		}
	}
	// render template with tplName index
	_ = tpl.ExecuteTemplate(
		w,
		"index.html",
		randomArticles,
	)
}
func shareHandler(w http.ResponseWriter, r *http.Request) {
	//ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	io.Copy(ioutil.Discard, r.Body)
	//}))
	if r.Method != http.MethodPost {
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rs := models.Article{}
	err = json.Unmarshal(b, &rs)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	basePath, basePathSet := os.LookupEnv("ROBOT_BASE_PATH")
	var errMsg string
	if !basePathSet {
		errMsg = "ROBOT_BASE_PATH is required."
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	accessToken, accessTokenSet := os.LookupEnv("ROBOT_ACCESS_TOKEN")
	if !accessTokenSet {
		errMsg = "ROBOT_ACCESS_TOKEN is required."
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	accessKey, accessKeySet := os.LookupEnv("ROBOT_ACCESS_KEY")
	if !accessKeySet {
		errMsg = "ROBOT_ACCESS_KEY is required."
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	rb := &robot.Robot{BasePath: basePath, AccessToken: accessToken, AccessKey: accessKey}
	//n, err := rb.Write([]byte("hellorobot"))

	msg := fmt.Sprintf("大家好我是机器人小库，推荐给大家一篇文章；\n题目：%s；\n地址：%s；\n作者：%s；\n发布时间：%s；", rs.Title, rs.Address, rs.Author, rs.PublishTime)
	n, err := rb.Write([]byte(msg))
	if err != nil || n == 0 {
		log.Printf("Robot write error: %v", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	//w.WriteHeader(http.StatusOK)
	//set article's is_shared field "1"
	err = models.UpdateShareFeild(rs.ArticleId)
	if err != nil {
		log.Printf("Update share feild error: %v", err)
	}
}
func dutyHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/duty/" {
		return
	}

	duty := models.GetAllDuty()
	_ = tpl.ExecuteTemplate(
		w,
		"duty.html",
		duty,
	)
}
func editDutyHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/dutyedit/" {
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parse form error.", http.StatusInternalServerError)
		return
	}
	//id := r.PostFormValue("id")
	//name := r.PostFormValue("name")
	//employeeNum := r.PostFormValue("employeeNum")
	//phone := r.PostFormValue("phone")
	////u := models.Duty{Name: name, EmployeesNum: employeeNum, PhoneNum: phone}
	////update
	//if id != "0" {
	//
	//}
	////insert
	//if id == "0" {
	//
	//}
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/login/" && r.RequestURI != "/login" {
		return
	}
	if r.Method == http.MethodGet {
		_ = tpl.ExecuteTemplate(
			w,
			"login.html",
			nil,
		)
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Parse form error.", http.StatusInternalServerError)
			return
		}
		name, pwd := r.PostFormValue("admin-user"), r.PostFormValue("admin-pwd")
		u := models.User{name, pwd}
		if models.CheckUser(u) {
			session, err := store.Get(r, sessionKey)
			session.Options = &sessions.Options{
				MaxAge: 86400,
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set some session values.
			session.Values["isLogin"] = true
			// Save it before we write to the response/return from the handler.
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/duty", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

//static server
func serveFile() {
	//absolue path
	fs := http.Dir("./static")
	handler := http.StripPrefix("/static", http.FileServer(fs))
	http.Handle("/static/", handler)
}

func authWapper(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, sessionKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check whether is login
		isLogin, ok := session.Values["isLogin"].(bool)
		if !isLogin || !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		handler(w, r)
	}
}
