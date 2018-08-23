package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dghubble/gologin"
	g2 "github.com/dghubble/gologin/google"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	blogger "google.golang.org/api/blogger/v3"
)

func dynamicRoutes(router *mux.Router) {

	// Social Controller //
	stateConfig := gologin.DebugOnlyCookieConfig

	router.HandleFunc("/", indexPage)
	router.Handle("/gp/login", g2.StateHandler(stateConfig, g2.LoginHandler(GpConf, nil)))
	router.Handle("/gp/callback", g2.StateHandler(stateConfig, g2.CallbackHandler(GpConf, Social.GPissueSession(), nil)))

	router.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {

		userblogservice := blogger.NewBlogsService(BloggerClient)
		// list := userblogservice.List("self") // for v2
		list := userblogservice.ListByUser("self")
		getList, err := list.Do()
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}

		dataMap := make(map[string]interface{})
		dataMap["blogs"] = getList.Items
		View(w, r, dataMap, "pages/bloglist.html")
	})

	router.HandleFunc("/explore/blog/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] == "" {
			fmt.Fprintln(w, "Invalid blog id")
			return
		}
		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["id"])
		blog, _ := getBlog.Do()

		blogPostsService := blogger.NewPostsService(BloggerClient)
		postsList := blogPostsService.List(vars["id"])
		getList, err := postsList.Do()
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		dataMap := make(map[string]interface{})
		dataMap["posts"] = getList.Items
		dataMap["blog"] = blog.Name
		View(w, r, dataMap, "pages/postslist.html")
	})

	router.HandleFunc("/save/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprintln(w, "Unsupported HTTP method")
			return
		} else {

			decoder := json.NewDecoder(r.Body)
			var t BEditorData
			err := decoder.Decode(&t)
			if err != nil {
				panic(err)
			}

			if t.Blogid == "" || t.Postid == "" {
				fmt.Fprintln(w, "Blog id or postid is nil/empty")
				return
			}

			fmt.Println("Data of struct ", t)
			blogPostsService := blogger.NewPostsService(BloggerClient)
			post := &blogger.Post{}
			post.Content = t.Content
			post.Title = t.Title
			postRequest := blogPostsService.Patch(t.Blogid, t.Postid, post)
			post, patcherr := postRequest.Do()
			if patcherr != nil {
				fmt.Fprintln(w, "Faild to patch existing post")
				color.Red("Error during patch : %s ", patcherr.Error())
				return
			}

			dataMap := make(map[string]interface{})
			dataMap["title"] = post.Title
			dataMap["content"] = post.Content
			JSON(w, dataMap)

		}
	})

	router.HandleFunc("/explore/blog/{blogid}/post/{postid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["blogid"] == "" || vars["postid"] == "" {
			fmt.Fprintln(w, "Invalid blogId or postId")
			return
		}

		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["blogid"])
		blog, _ := getBlog.Do()

		blogPostsService := blogger.NewPostsService(BloggerClient)
		getpost := blogPostsService.Get(vars["blogid"], vars["postid"])
		post, err := getpost.Do()
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		dataMap := make(map[string]interface{})
		dataMap["post"] = post
		dataMap["blog"] = blog.Name
		dataMap["blogid"] = vars["blogid"]
		dataMap["postid"] = vars["postid"]
		View(w, r, dataMap, "pages/post.html")

	})
}
