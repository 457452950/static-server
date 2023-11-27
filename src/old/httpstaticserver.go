package old

// const YAMLCONF = ".ghs.yml"

// type IndexFileItem struct {
// 	Path string
// 	Info os.FileInfo
// }

// type Directory struct {
// 	size  map[string]int64
// 	mutex *sync.RWMutex
// }

// type HTTPStaticServer struct {
// 	Config config.FileServiceConfig

// 	fileTransformer files.FileTransformer
// 	indexes         []IndexFileItem
// 	muxRouter       *mux.Router
// 	bufPool         sync.Pool // use sync.Pool caching buf to reduce gc ratio
// }

// Return real path with Seperator(/)
// func (s *HTTPStaticServer) getRealPath(r *http.Request) string {
// 	path := mux.Vars(r)["path"]
// 	mp, err := s.fileTransformer.TransformPath(path)
// 	log.Printf("transform path {%s}  {%s} \n", mp, err)
// 	rp, err := s.fileTransformer.IsolationPath(string(mp))
// 	log.Printf("Isolation path {%s} {%s} {%s} \n", rp, path, err)

// 	if !strings.HasPrefix(path, "/") {
// 		path = "/" + path
// 	}
// 	path = filepath.Clean(path) // prevent .. for safe issues
// 	relativePath, err := filepath.Rel(s.Config.Prefix, path)
// 	if err != nil {
// 		relativePath = path
// 	}
// 	realPath := filepath.Join(s.Config.Root, relativePath)
// 	return filepath.ToSlash(realPath)
// }

// func combineURL(r *http.Request, path string) *url.URL {
// 	return &url.URL{
// 		Scheme: r.URL.Scheme,
// 		Host:   r.Host,
// 		Path:   path,
// 	}
// }

// func (s *HTTPStaticServer) hPlist(w http.ResponseWriter, r *http.Request) {
// 	path := mux.Vars(r)["path"]
// 	// rename *.plist to *.ipa
// 	if filepath.Ext(path) == ".plist" {
// 		path = path[0:len(path)-6] + ".ipa"
// 	}

// 	relPath := s.getRealPath(r)
// 	plinfo, err := parseIPA(relPath)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}

// 	scheme := "http"
// 	if r.TLS != nil {
// 		scheme = "https"
// 	}
// 	baseURL := &url.URL{
// 		Scheme: scheme,
// 		Host:   r.Host,
// 	}
// 	data, err := generateDownloadPlist(baseURL, path, plinfo)
// 	if err != nil {
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "text/xml")
// 	w.Write(data)
// }

// func (s *HTTPStaticServer) hIpaLink(w http.ResponseWriter, r *http.Request) {
// 	path := mux.Vars(r)["path"]
// 	var plistUrl string

// 	if r.URL.Scheme == "https" {
// 		plistUrl = combineURL(r, "/-/ipa/plist/"+path).String()
// 	} else if s.PlistProxy != "" {
// 		httpPlistLink := "http://" + r.Host + "/-/ipa/plist/" + path
// 		url, err := s.genPlistLink(httpPlistLink)
// 		if err != nil {
// 			http.Error(w, err.Error(), 500)
// 			return
// 		}
// 		plistUrl = url
// 	} else {
// 		http.Error(w, "500: Server should be https:// or provide valid plistproxy", 500)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/html")
// 	log.Println("PlistURL:", plistUrl)
// 	renderHTML(w, "assets/ipa-install.html", map[string]string{
// 		"Name":      filepath.Base(path),
// 		"PlistLink": plistUrl,
// 	})
// }

// func (s *HTTPStaticServer) genPlistLink(httpPlistLink string) (plistUrl string, err error) {
// 	// Maybe need a proxy, a little slowly now.
// 	pp := s.PlistProxy
// 	if pp == "" {
// 		pp = defaultPlistProxy
// 	}
// 	resp, err := http.Get(httpPlistLink)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()

// 	data, _ := ioutil.ReadAll(resp.Body)
// 	retData, err := http.Post(pp, "text/xml", bytes.NewBuffer(data))
// 	if err != nil {
// 		return
// 	}
// 	defer retData.Body.Close()

// 	jsonData, _ := ioutil.ReadAll(retData.Body)
// 	var ret map[string]string
// 	if err = json.Unmarshal(jsonData, &ret); err != nil {
// 		return
// 	}
// 	plistUrl = pp + "/" + ret["key"]
// 	return
// }

// func (s *HTTPStaticServer) hFileOrDirectory(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, s.getRealPath(r))
// }

// TODO: I need to read more abouthtml/template
// var (
// 	funcMap template.FuncMap
// )

// func init() {
// 	funcMap = template.FuncMap{
// 		"title": strings.Title,
// 		// "urlhash": func(path string) string {
// 		// 	httpFile, err := Assets.Open(path)
// 		// 	log.Printf("assets open %s err %s\n", path, err)
// 		// 	if err != nil {
// 		// 		return path + "#no-such-file"
// 		// 	}
// 		// 	info, err := httpFile.Stat()
// 		// 	if err != nil {
// 		// 		return path + "#stat-error"
// 		// 	}
// 		// 	return fmt.Sprintf("%s?t=%d", path, info.ModTime().Unix())
// 		// },
// 	}
// }
