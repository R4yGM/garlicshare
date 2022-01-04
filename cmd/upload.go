package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"path/filepath"

	"context"
	"log"
	"net/http"	
	"time"
	"github.com/cretz/bine/tor"
	"github.com/innix/shrek"
	"html/template"
	"os"
	"math/rand"
    "strings"
	"strconv"
	gs "github.com/R4yGM/garlicshare/size"
)



var (
	key string
	path string
	templateString string
)


var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Creates an .onion link where you can download a file or directory ",
	Run: func(cmd *cobra.Command, args []string) {

		templateString = `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<meta http-equiv="X-UA-Compatible" content="ie=edge" />
				<title>GarlicShare : Secret file sharing</title>
				<meta name="description" content="" />
				<meta name="keywords" content="" />
				<meta name="author" content="" />

			<link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.10.2/css/all.css">

				<link rel="stylesheet" href="https://unpkg.com/tailwindcss@2.2.19/dist/tailwind.min.css"/>
				<link
				href="https://unpkg.com/@tailwindcss/custom-forms/dist/custom-forms.min.css"
				rel="stylesheet"
				/>

				<style>
					@import url("https://rsms.me/inter/inter.css");
					html {
					font-family: "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI",
						Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif,
						"Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol",
						"Noto Color Emoji";
					}

					.gradient {
					background: #16222a; /* fallback for old browsers */
					background: -webkit-linear-gradient(to bottom, #16222a, #3a6073); /* Chrome 10-25, Safari 5.1-6 */
					background: linear-gradient(to bottom, #16222a, #3a6073); /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */
					}

					button{
					background-color: white;
					}
					.gradient2 {
					background-color: #f39f86;
					background-image: linear-gradient(315deg, #f39f86 0%, #f9d976 74%);
					}
				</style>

				<script type="text/javascript">
						function download(){
							if(document.getElementById("password")){
							document.location.assign(("/download?Key="+document.getElementById("password").value))
							}else{
							document.location.assign("/download")
							}
						}
												
						var navMenuDiv = document.getElementById("nav-content");
						var navMenu = document.getElementById("nav-toggle");
						
						document.onclick = check;
						function check(e) {
							var target = (e && e.target) || (event && event.srcElement);
						
							//Nav Menu
							if (!checkParent(target, navMenuDiv)) {
							// click NOT on the menu
							if (checkParent(target, navMenu)) {
								// click on the link
								if (navMenuDiv.classList.contains("hidden")) {
								navMenuDiv.classList.remove("hidden");
								} else {
								navMenuDiv.classList.add("hidden");
								}
							} else {
								// click both outside link and outside menu, hide menu
								navMenuDiv.classList.add("hidden");
							}
							}
						}
						function checkParent(t, elm) {
							while (t.parentNode) {
							if (t == elm) {
								return true;
							}
							t = t.parentNode;
							}
							return false;
						}
				</script>
			</head>

			<body class="gradient leading-relaxed tracking-wide flex flex-col">
				<!--Nav-->
				<nav id="header" class="w-full z-30 top-0 text-white py-1 lg:py-6">
				<div
					class="w-full container mx-auto flex flex-wrap items-center justify-between mt-0 px-2 py-2 lg:py-6"
				>
				<img src="https://i.imgur.com/bfUnl8P.png" width=50 height=50>
					<div class="pl-4 flex items-center">
					<a
						class="text-white no-underline hover:no-underline font-bold text-2xl lg:text-4xl"
						href="#"
					>
						
						GarlicShare
					</a>
					</div>

					<div class="block lg:hidden pr-4">
					<button
						id="nav-toggle"
						class="flex items-center px-3 py-2 border rounded text-gray-500 border-gray-600 hover:text-gray-800 hover:border-green-500 appearance-none focus:outline-none"
					>
						<svg
						class="fill-current h-3 w-3"
						viewBox="0 0 20 20"
						xmlns="http://www.w3.org/2000/svg"
						>
						<title>Menu</title>
						<path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z" />
						</svg>
					</button>
					</div>

					<div
					class="w-full flex-grow lg:flex lg:items-center lg:w-auto hidden lg:block mt-2 lg:mt-0 text-black p-4 lg:p-0 z-20"
					id="nav-content"
					>
					<ul class="list-reset lg:flex justify-end flex-1 items-center">
						<li class="mr-3">
						<a
							class="inline-block text-black no-underline hover:text-gray-600 text-gray-500 hover:text-underline py-2 px-4"
							href="https://github.com/R4yGM/GarlicShare"
							>Github</a
						>
						</li>
						<li class="mr-3">
						<a
							class="inline-block text-black no-underline hover:text-gray-600 text-gray-500 hover:text-underline py-2 px-4"
							href="https://r4ygm.github.io/garlicshare/"
							>Website</a
						>
						</li>
					</ul>
					</div>
				</div>
				</nav>

				<div class="container mx-auto h-screen">
				<div class="text-center px-3 lg:px-0">
					<p
					class="leading-normal text-gray-400 text-base md:text-xl lg:text-2xl mb-8"
				>
				Simple, Secure and private file sharing over the Tor network
				</p>
					<h1 id="filename"
					class="my-4 text-2xl md:text-3xl lg:text-5xl font-black leading-tight" style="color:white"
					>
					{{.FileName}}
					</h1>

					

					{{.PasswordForm}}



					<button id="download" onClick="download()"
					class="mx-auto lg:mx-0 hover:underline text-gray-800 font-extrabold rounded my-2 md:my-6 py-4 px-8 shadow-lg w-68"
					>
					<i class="fas fa-download"></i> Download ({{.FileSize}})
					</button>
					<img id="qrcode" style="margin:auto;margin-top:10%;border: 3px solid black;  border-radius: 6px;">

					<p class="text-center text-gray-500 text-xs" style="margin-top:4px">
					Use this QR code to share this page
					</p>

					
				</div>
				<script>document.getElementById('qrcode').src="https://api.qrserver.com/v1/create-qr-code/?size=150x150&data="+window.location.href;</script>

				<div class="flex items-center w-full mx-auto content-end">
				</div>
				</div>
				
			</body>
			</html>`

		Upload();
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.PersistentFlags().StringVarP(&key, "key", "k", "", "Password to download the files")
	uploadCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Path")
	uploadCmd.MarkPersistentFlagRequired("path")
}



func Upload(){

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Panicf("Path '%s' does not exist", path)
	}
	dt := time.Now()
	fmt.Printf(strings.Repeat("=", 60))
    fmt.Println("\nGarlicShare starting",dt.Format("01-02-2006 15:04:05"))
	fmt.Println("Starting and registering onion service, please wait a couple of minutes...")
	t, err := tor.Start(nil, nil)
	if err != nil {
		log.Panicf("Unable to start Tor: %v", err)
	}
	defer t.Close()

	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()


	addr, err := shrek.MineOnionHostName(context.Background(), nil, shrek.StartEndMatcher{
		Start: []byte("gar"),
	})

	if err != nil {
		panic(err)
	}
	onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, Key: addr.SecretKey, RemotePorts: []int{80}})


	if err != nil {
		log.Panicf("Unable to create onion service: %v", err)
	}
	defer onion.Close()
	fmt.Printf("Open Tor browser and navigate to http://%v.onion\n", onion.ID)

	errCh := make(chan error, 1)

	/*var staticFS = fs.FS(Static)
	htmlContent, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(htmlContent))
	
	http.Handle("/static/", fs)
*/
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vars := make(map[string]interface{})
		

		info, err := os.Stat(path)
		if err != nil {
			panic(err)
		}

		if info.IsDir() {
			if(path[len(path)-1:] == "/"){
				vars["FileName"] = template.HTML(`<i class="fas fa-folder" style="margin-right:20px"></i>`+filepath.Base(path))
			}else{
				vars["FileName"] = template.HTML(`<i class="fas fa-folder" style="margin-right:20px"></i>`+filepath.Base(path)+"/")
			}
			
		}else{
			vars["FileName"] = template.HTML(`<i class="fas fa-file" style="margin-right:20px"></i>`+filepath.Base(path))
		}

		vars["FileSize"] = gs.HumanFileSize(float64(info.Size()))

		if(key != ""){
			vars["PasswordForm"] = template.HTML(`
						<div class="" style="margin:auto;width:30%">
							<div class="mb-6">
							<label class="block text-gray-700 text-sm font-bold mb-2" style="color:white" for="password">
								Key
							</label>
							<input onKeyDown="if(event.keyCode==13) download()"class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************">
							</div>
						<p class="text-center text-gray-500 text-s">
							This file is protected by a password, insert the password and download the file
						</p>
						</div>`)
		}else{
			vars["PasswordForm"] = "";
		}
		tmpl := template.Must(template.New("").Parse(templateString))
		tmpl.Execute(w, vars)
	})




	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < 30; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() 

	if fileInfo.IsDir() {
		http.Handle("/download/"+str+"/", http.StripPrefix("/download/"+str+"/", http.FileServer(http.Dir(path))))
	}

	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		if(key != ""){	
			param, ok := r.URL.Query()["Key"]
			if !ok || len(param[0]) < 1 {
				fmt.Fprintf(w, "Key parameter is missing!")
				return
			}
			defer r.Body.Close()
			
			if(key == param[0]){
				fileInfo, err := os.Stat(path)
				if err != nil {
					// error handling
				}
				if fileInfo.IsDir() {
					http.Redirect(w, r, "/download/"+str, http.StatusSeeOther)
				} else {
						w.Header().Set("Content-Disposition", "attachment; filename="+path)
						w.Header().Set("Content-Type", "application/octet-stream")
						fi, err := os.Stat(path)
						if err != nil {
						}
						w.Header().Set("Content-Lenght", strconv.FormatInt(int64(fi.Size()), 10))
						http.ServeFile(w, r, path)
				}
			}else{
				fmt.Fprintf(w, "Wrong Key!")
			}
		}else{
			fileInfo, err := os.Stat(path)
			if err != nil {
				// error handling
			}
			if fileInfo.IsDir() {
				http.Redirect(w, r, "/download/"+str, http.StatusSeeOther)
			} else {
				w.Header().Set("Content-Disposition", "attachment; filename="+path)
				w.Header().Set("Content-Type", "application/octet-stream")
				fi, err := os.Stat(path)
				if err != nil {
				}
				w.Header().Set("Content-Lenght", strconv.FormatInt(int64(fi.Size()), 10))
				http.ServeFile(w, r, path)
			}
		}
	})

	if err != nil {
		panic(err)
	}


	http.Serve(onion, nil)
	// End when enter is pressed
	go func() {
		fmt.Scanln()
		errCh <- nil
	}()
	if err = <-errCh; err != nil {
		log.Panicf("Failed serving: %v", err)
	}
}
