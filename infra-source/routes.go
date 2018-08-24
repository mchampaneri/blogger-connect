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

	// Home page
	router.HandleFunc("/", indexPage)

	// Google OAuth Routes
	// contains login -  callback duo
	router.Handle("/gp/login", g2.StateHandler(stateConfig, g2.LoginHandler(GpConf, nil)))
	router.Handle("/gp/callback", g2.StateHandler(stateConfig, g2.CallbackHandler(GpConf, Social.GPissueSession(), nil)))

	// Authenticated userse blog lsit
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

	// Posts inside the blog
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
		postsList = postsList.View("ADMIN")
		// getList, err := postsList.Do()
		// Taking context of current request
		ctx := r.Context()
		var posts []*blogger.Post
		dataMap := make(map[string]interface{})
		postsList.Pages(ctx, func(postlist *blogger.PostList) error {
			posts = append(posts, postlist.Items...)
			return nil
		})
		dataMap["posts"] = posts
		dataMap["blog"] = blog.Name
		dataMap["blogid"] = vars["id"]
		View(w, r, dataMap, "pages/postslist.html")
	})

	// Insert new post in the blog
	router.HandleFunc("/explore/blog/{id}/new", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["id"] == "" {
			fmt.Fprintln(w, "Invalid blog id")
			return
		}
		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["id"])
		blog, _ := getBlog.Do()

		if r.Method == "GET" {
			dataMap := make(map[string]interface{})
			dataMap["blogid"] = vars["id"]
			dataMap["blog"] = blog.Name
			View(w, r, dataMap, "pages/post.html")
			return

		}

	})
	// Fetch existing post for the blog
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
		getpost = getpost.View("ADMIN")
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

	router.HandleFunc("/explore/blog/{blogid}/post/delete/{postid}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["blogid"] == "" || vars["postid"] == "" {
			fmt.Fprintln(w, "Invalid blogId or postId")
			return
		}
		blogPostsService := blogger.NewPostsService(BloggerClient)
		deletePost := blogPostsService.Delete(vars["blogid"], vars["postid"])
		err := deletePost.Do()
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}
		redirectToPosts := fmt.Sprint("/explore/blog/", vars["blogid"])
		http.Redirect(w, r, redirectToPosts, 301)
		return
	})

	// Update route for the psot
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

			if t.Blogid == "" {
				fmt.Fprintln(w, "Blog id  is nil/empty")
				return
			}

			fmt.Println("Data of struct ", t)
			blogPostsService := blogger.NewPostsService(BloggerClient)
			post := &blogger.Post{}
			post.Content = t.Content
			post.Title = t.Title

			if t.Postid == "" {
				postRequest := blogPostsService.Insert(t.Blogid, post)
				postRequest = postRequest.IsDraft(true)
				postReturned, insertErr := postRequest.Do()
				if insertErr != nil {
					fmt.Fprintln(w, "Faild to insert new post")
					color.Red("Error during insertion : %s ", insertErr.Error())
					return
				}
				dataMap := make(map[string]interface{})
				dataMap["title"] = postReturned.Title
				dataMap["postid"] = postReturned.Id
				dataMap["content"] = postReturned.Content
				JSON(w, dataMap)
			} else {
				postRequest := blogPostsService.Patch(t.Blogid, t.Postid, post)
				postRequest = postRequest.Revert(true)
				postReturned, patcherr := postRequest.Do()
				if patcherr != nil {
					fmt.Fprintln(w, "Faild to patch existing post")
					color.Red("Error during patch : %s ", patcherr.Error())
					return
				}
				dataMap := make(map[string]interface{})
				dataMap["title"] = postReturned.Title
				dataMap["postid"] = postReturned.Id
				dataMap["content"] = postReturned.Content
				JSON(w, dataMap)
			}
		}
	})

}
