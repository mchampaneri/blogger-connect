# blogger-connect
Alternative editor for google blogger

#### Have look at demo deployment on YouTube https://youtu.be/CoXymNBPtV4

## Setup

### Download Source code

### Add your client credentails in `gpload()` functions in `infra-source\authHandlres.go`

```
func gpLoad() {

	GpConf = &oauth2.Config{
		ClientID:     "Your-Client-Id",
		ClientSecret: "Your-Client-Secret",
		Scopes:       []string{"profile", "email", blogger.BloggerScope},
		Endpoint:     googleOAuth2.Endpoint,
	}

	GpConf.RedirectURL = fmt.Sprintf("%s/gp/callback", Config.AppURL)
	fmt.Println("App url is ", GpConf.RedirectURL)
}
```

### Compile Golang Code & Frontend Code
- Run ` go build -o ../app/server ` from  `/infra-source folder `
- Run ` npm run prod ` from `/app folder`

### Copy app folder to VPS and start server
- Copy your app folder with ` server  ` executable
- execute server

## Done!

##### [www.mchampaneri.in](https://www.mchampaneri.in)
