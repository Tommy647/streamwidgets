package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/Tommy647/streamwidgets/internal/plugins/twitch"
	"github.com/Tommy647/streamwidgets/internal/widgets"
)

type config struct {
	Watchers []widgets.Widget `yaml:"watchers"`
}

// port just a const for now until we add viper
const port = `:8080`

func main() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name) // @todo: implement this it would be nice
	})
	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	log.Printf("config is : %v\n", viper.Get("watchers"))

	var cfg config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err.Error())
	}
	_ = cfg

	srv := http.NewServeMux()

	for k, v := range cfg.Watchers {
		log.Printf("setting up handlers for : %s %v", k, v)
		_, err := os.Open(v.Path)
		if err != nil {
			panic(err.Error())
		}
		fs := http.FileServer(http.Dir(v.Path))
		srv.Handle("/widgets/"+v.Name, http.StripPrefix("/widgets/"+v.Name, fs))
	}

	// @todo: explore, uses the wrong default config :/
	// if err := viper.SafeWriteConfig(); err != nil {
	// 	panic(err.Error())
	// }

	// switch x := viper.Get("bots").(type) {
	// case map[string]interface{}:
	// 	log.Printf("found config as map[string]interface{}: %s", x)
	// default:
	// 	panic(fmt.Sprintf("unknown type in config: %s", x))
	// }
	//
	// for k, v := range viper.Get("bots").(map[string]interface{}) {
	// 	log.Printf("setting up bots for : %s %v", k, v)
	// 	var c map[string]interface{}
	// 	switch x := v.(type) {
	// 	case map[string]interface{}:
	// 		c = x
	// 		log.Printf("found %q config as map[string]interface{}: %s", k, c)
	// 	default:
	// 		panic(fmt.Sprintf("unknown type in %q config: %s", k, x))
	// 	}
	// }

	//go twitch.New()
	// log.Printf("TWITCH DONE AND LISTENING")
	srv.HandleFunc(`/twitch`, twitch.Handle)
	// http.HandleFunc("/", FolderHandler(fmt.Sprintf("%sexamples", rootFolder)))

	fs := http.FileServer(http.Dir("admin"))
	srv.Handle("/admin", http.StripPrefix("/admin", fs))

	log.Printf("Starting up on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, srv))
}

// FolderHandler handles folder paths
func FolderHandler(folder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir, err := filepath.Abs(filepath.Dir(folder))
		if err != nil {
			log.Fatal(err)
		}
		path := r.URL.Path[1:]
		file := fmt.Sprintf("%s/%s/index.html", dir, path) // @todo: just testing code
		log.Printf("using: %s", file)
		f, err := os.Open(file)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		_, err = fmt.Fprint(w, string(data))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
	}
}
