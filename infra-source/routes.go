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

	// Clearing gorilla session
	// created on time on Oauth ....
	router.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		usersession, err := UserSession.Get(r, "mvc-user-session")
		if err == nil {
			for k := range usersession.Values {
				delete(usersession.Values, k)
			}
			usersession.Options.MaxAge = -1
			usersession.Save(r, w)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	})

	// Authenticated userse blog lsit
	router.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}

		userblogservice := blogger.NewBlogsService(BloggerClient)
		list := userblogservice.ListByUser("self")
		getList, err := list.Do()
		if err != nil {
			http.Redirect(w, r, "/logout", 301)
			return
		}

		dataMap := make(map[string]interface{})
		dataMap["blogs"] = getList.Items
		View(w, r, dataMap, "pages/bloglist.html")

	})

	// Posts belongs to this blog
	router.HandleFunc("/explore/blog/{id}", func(w http.ResponseWriter, r *http.Request) {

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}

		vars := mux.Vars(r)
		if vars["id"] == "" {
			fmt.Fprintln(w, "Invalid blog id")
			return
		}
		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["id"])
		blog, err := getBlog.Do()
		if err != nil {
			http.Redirect(w, r, "/logout", 301)
			return
		}

		blogPostsService := blogger.NewPostsService(BloggerClient)
		postsList := blogPostsService.List(vars["id"])

		postsList = postsList.View("ADMIN").Status("draft", "live")

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
	// Just front end :
	// Thouhg process of saving is still
	// exectude by /save/post
	router.HandleFunc("/explore/blog/{id}/new", func(w http.ResponseWriter, r *http.Request) {

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}

		vars := mux.Vars(r)
		if vars["id"] == "" {
			fmt.Fprintln(w, "Invalid blog id")
			return
		}
		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["id"])
		blog, err := getBlog.Do()
		if err != nil {
			http.Redirect(w, r, "/logout", 301)
			return
		}

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

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}

		vars := mux.Vars(r)
		if vars["blogid"] == "" || vars["postid"] == "" {
			fmt.Fprintln(w, "Invalid blogId or postId")
			return
		}

		userblogservice := blogger.NewBlogsService(BloggerClient)
		getBlog := userblogservice.Get(vars["blogid"])
		blog, err := getBlog.Do()
		if err != nil {
			http.Redirect(w, r, "/logout", 301)
			return
		}

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

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
			return
		}

		vars := mux.Vars(r)
		if vars["blogid"] == "" || vars["postid"] == "" {
			fmt.Fprintln(w, "Invalid blogId or postId")
			return
		}
		blogPostsService := blogger.NewPostsService(BloggerClient)
		deletePost := blogPostsService.Delete(vars["blogid"], vars["postid"])
		err := deletePost.Do()
		if err != nil {
			http.Redirect(w, r, "/logout", 301)
			return
		}
		redirectToPosts := fmt.Sprint("/explore/blog/", vars["blogid"])
		http.Redirect(w, r, redirectToPosts, 301)
		return
	})

	// Update route for the psot
	router.HandleFunc("/save/post", func(w http.ResponseWriter, r *http.Request) {

		if IsNotAuthenticated() {
			http.Redirect(w, r, "/", 301)
		}

		if r.Method != "POST" {
			fmt.Fprintln(w, "Unsupported HTTP method")
			return
		} else {

			var updateErr, patchErr error
			var t *BEditorData

			decoder := json.NewDecoder(r.Body)

			err := decoder.Decode(&t)
			if err != nil {
				fmt.Println("Error During Unmarshling ", err.Error())
				return
			}

			post := &blogger.Post{}
			postReturned := &blogger.Post{}
			post.Content = t.Content
			post.Title = t.Title

			if t.Blogid == "" {
				fmt.Fprintln(w, "blog id is nil/empty")
				return
			}

			if err != nil {
				fmt.Println("failed to get post before updating it", err.Error())
				return
			}

			blogPostsService := blogger.NewPostsService(BloggerClient)

			// ... When post is new [ literaly ... ]

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
				return
			} else {
				// When post status is DRAFT we need to
				// make update call
				if t.Status == "DRAFT" {
					getPostCall := blogPostsService.Get(t.Blogid, t.Postid)
					getPostCall = getPostCall.View("ADMIN")
					post, err = getPostCall.Do()
					post.Content = t.Content
					post.Title = t.Title

					postRequest := blogPostsService.Update(t.Blogid, t.Postid, post)
					postReturned, updateErr = postRequest.Do()
					if updateErr != nil {
						fmt.Fprintln(w, "Faild to update  existing post")
						color.Red("Error during patch : %s ", updateErr.Error())
						return
					}
				} else if t.Status == "LIVE" {
					postRequest := blogPostsService.Patch(t.Blogid, t.Postid, post)
					postReturned, patchErr = postRequest.Do()
					if patchErr != nil {
						fmt.Fprintln(w, "Faild to patch existing post")
						color.Red("Error during patch : %s ", patchErr.Error())
						return
					}

				}

				dataMap := make(map[string]interface{})
				dataMap["title"] = postReturned.Title
				dataMap["postid"] = postReturned.Id
				dataMap["content"] = postReturned.Content
				JSON(w, dataMap)
			}
		}
	})

	router.HandleFunc("/chage-state/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprintln(w, "wrong request method")
		} else {
			var t BEditorData

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&t)
			if err != nil {
				fmt.Fprintln(w, " Failed to decode the struct provided by input")
				return
			}

			var revertErr, publishErr error
			post := &blogger.Post{}

			blogPostsService := blogger.NewPostsService(BloggerClient)
			if t.Status == "LIVE" {
				// Rever blog from live status to
				// draft stage
				revertBlog := blogPostsService.Revert(t.Blogid, t.Postid)
				revertBlog = revertBlog.Fields("status")
				post, revertErr = revertBlog.Do()
				if revertErr != nil {
					fmt.Fprintln(w, " Failed to revert blog post")
					return
				}

			} else if t.Status == "DRAFT" {
				publishBlog := blogPostsService.Publish(t.Blogid, t.Postid)
				publishBlog = publishBlog.Fields("status")
				post, publishErr = publishBlog.Do()
				if publishErr != nil {
					fmt.Fprintln(w, "Failed to publish blog post")
					return
				}
			}

			dataMap := make(map[string]string)
			dataMap["status"] = post.Status
			JSON(w, dataMap)
		}
	})

}

// Just checking if user is authenticated
// and connection is still alive
func IsNotAuthenticated() bool {
	return BloggerClient == nil
}
