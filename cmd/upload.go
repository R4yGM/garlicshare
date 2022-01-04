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
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Creates an .onion link where you can download a file or directory ",
	Run: func(cmd *cobra.Command, args []string) {
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


	http.Handle("/upload", http.FileServer(http.Dir(".")))

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

		tmpl := template.Must(template.ParseFiles("static/index.html"))
		tmpl.Execute(w, vars)
	})


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))


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


