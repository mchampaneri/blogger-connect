package main

import (
	"fmt"

	"github.com/fatih/color"
	sciter "github.com/sciter-sdk/go-sciter"
	blogger "google.golang.org/api/blogger/v3"
)

// NavTo function
// responseds to the url change
// on sciter frontend
func NavTo(vals ...*sciter.Value) *sciter.Value {
	color.Yellow("NavTo called ")
	for _, val := range vals {
		switch val.String() {
		case "Blogs":
			fmt.Println(fetchBlogs())
			return fetchBlogs()
		}
	}
	return nil
}

// fethcing blog for authenticate user
func fetchBlogs() *sciter.Value {
	if BloggerClient == nil {
		color.Red("Blogger Client  is not beign initialized ")
		return nil
	}
	userblogservice := blogger.NewBlogsService(BloggerClient)
	list := userblogservice.ListByUser("self")
	getList, err := list.Do()
	if err != nil {
		store := sciter.NewValue()
		for index, item := range getList.Items {
			store.SetIndex(index, item.Name)
		}
		return store
	} else {
		return nil
	}
}
