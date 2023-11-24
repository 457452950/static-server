package service

import (
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"

	"static-server/config"
	"static-server/files"
)

type IndexFileItem struct {
	Path string
	Info os.FileInfo
}

type FileServiceHandler struct {
	Config config.FileServiceConfig

	fileTransformer files.FileTransformer
	indexes         []IndexFileItem
	muxRouter       *mux.Router
	bufPool         sync.Pool // use sync.Pool caching buf to reduce gc ratio
}

func (s *FileServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.muxRouter.ServeHTTP(w, r)
}

func CreateFileStaticService(conf config.FileServiceConfig) (handler *FileServiceHandler) {
	if err := conf.CheckPrefix(); err != nil {
		panic(err)
	}

	if err := conf.CheckRoot(); err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	handler = &FileServiceHandler{
		Config:          conf,
		fileTransformer: files.CreateFileTransformer(conf.Prefix, conf.Root),
		muxRouter:       router,
		bufPool: sync.Pool{
			New: func() interface{} { return make([]byte, 32*1024) },
		},
	}

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// 	for {
	// 		startTime := time.Now()
	// 		log.Println("Started making search index")
	// 		s.makeIndex()
	// 		log.Printf("Completed search index in %v", time.Since(startTime))
	// 		//time.Sleep(time.Second * 1)
	// 		time.Sleep(time.Minute * 10)
	// 	}
	// }()

	// routers for Apple *.ipa
	// m.HandleFunc("/-/ipa/plist/{path:.*}", s.hPlist)
	// m.HandleFunc("/-/ipa/link/{path:.*}", s.hIpaLink)

	// init router
	// router.HandleFunc("/{path:.*}", handler.handleIndex).Methods("GET", "HEAD")
	// router.HandleFunc("/{path:.*}", handler.hUploadOrMkdir).Methods("POST")
	// router.HandleFunc("/{path:.*}", handler.handleDelete).Methods("DELETE")

	return
}
