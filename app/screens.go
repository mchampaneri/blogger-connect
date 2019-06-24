package main

import (
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
			fetchBlogs()
		}
	}
	return nil
}

// fethcing blog for authenticate user
func fetchBlogs() {
	if BloggerClient == nil {
		color.Red("Blogger Client  is not beign initialized")
	}
	color.Green("Fetching user's blog list ")
	userblogservice := blogger.NewBlogsService(BloggerClient)
	list := userblogservice.ListByUser("self")
	getList, err := list.Do()
	if err == nil {
		// Instead of returing single elemtns
		// we are going ot insert dom elements directly
		// on the screen for ease ...
		ListItemContainer, ListContainerSelectErr := RootElement.SelectById("#blogslist")
		if ListContainerSelectErr != nil {
			color.Red("failed to selected list contaienr %s", ListContainerSelectErr.Error())
			return
		}
		for index, item := range getList.Items {
			currentItem := &sciter.Element{}
			currentItem.SetHtml(ListItem("list-success", string(index), item.Name), sciter.SOH_INSERT_AFTER)
			ListItemContainer.Insert(currentItem, index)
			color.Yellow("element should ne attached")
		}
	}
}
