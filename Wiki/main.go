package main

import (
   "fmt" //for formatting messages to the console
  "net/http" //for web service
  "log" //logging errors
  "errors" //creating new errors
  "os" //reading and writing files
  "html/template" //for generating the page HTML
  "strconv" //converting status codes to string
)


type WikiPage struct {
  DocName string //public
  path string
  statusCode int
  Content []byte //for io libraries
  HTMLContent template.HTML //public
}

var FileNotFound *WikiPage
var ErrorOccur *WikiPage

//returns the WikiPage associated with the provided document name
//If the DocName does not coorspond to a valid file, return the 404 WikiPage
func getWikiPage( DocName string ) (*WikiPage, error) {
  var filePath string = "WebComponents/Documents/" +  DocName + ".csoc"
  body, err := os.ReadFile(filePath)
  if (err != nil) { //nil here means nothing wrong
    return FileNotFound, err
  }
  return &WikiPage{DocName: DocName, path: filePath, statusCode: 200, Content: body}, err
}

//returns the WikiPage for error files
//these exist in a different directory to not conflict with regular file names
func getErrorPage( statusCode int ) (*WikiPage, error) {
  var filePath string = "WebComponents/Errors/" + strconv.Itoa(statusCode) + ".csoc"
  body, err := os.ReadFile(filePath)
  if (err != nil) { //nil here means nothing wrong
    return ErrorOccur, err
  }
  return &WikiPage{DocName: strconv.Itoa(statusCode), path: filePath, statusCode: statusCode, Content: body, HTMLContent: template.HTML(body)}, err
}

func getFile( filePath string ) (*WikiPage, error) {
  body, err := os.ReadFile(filePath)
  if (err != nil) { //nil here means nothing wrong
    return FileNotFound, err
  }
  return &WikiPage{DocName: "temp", path: filePath, statusCode: 200, Content: body}, err
}

//uploads the Content of the page to the file structure
func writeWikiPage( page *WikiPage ) error {
  if (page == nil ) {
     return errors.New("Page is null")
  }
  return os.WriteFile(page.path, page.Content, 0666)
}

func getFileContents(writer http.ResponseWriter, request *http.Request) {
  var path string = request.URL.Path[1:]
  if ( path == "" ) { //path requested is just the root
    http.Redirect( writer, request, "/home/", 300)
    return
  }
  if ( path[ len(path) - 1 ] == '/' ) {
    //page.HTMLContent = template.HTML( page.Content ) //important later for when we create our own little language
    writer.WriteHeader( FileNotFound.statusCode )
    parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/error.html")
    err := parsedTemplate.Execute(writer, FileNotFound)
    if err != nil {
       log.Println("Error executing template :", err)
    }
    return
  }
  // page, _ := getFile(path)

  // writer.WriteHeader( page.statusCode )
  // writer.Write( page.Content )
  http.ServeFile( writer, request, path)
}

func viewHandler(writer http.ResponseWriter, request *http.Request) {
  var path string = request.URL.Path[ len("/view/"): ]
  page, _ := getWikiPage(path)
  writer.WriteHeader( page.statusCode )

  if ( page.statusCode != 200 ) {
    page.HTMLContent = template.HTML( page.Content ) //important later for when we create our own little language

    parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/error.html")
    err := parsedTemplate.Execute(writer, page)
    if err != nil {
       log.Println("Error executing template :", err)
    }
  } else {
    //page.Content = []byte( TransformStr( string(page.Content) ) )
    page.HTMLContent = template.HTML( page.Content ) //important later for when we create our own little language

    parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/view.html")
    err := parsedTemplate.Execute(writer, page)
    if err != nil {
       log.Println("Error executing template :", err)
    }

  }
}

func editHandler(writer http.ResponseWriter, request *http.Request) {
  var path string = request.URL.Path[ len("/edit/"): ]
  page, _ := getWikiPage(path)


  writer.WriteHeader( page.statusCode )

  if ( page.statusCode != 200 ) {
    page.HTMLContent = template.HTML( page.Content ) //important later for when we create our own little language

    parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/error.html")
    err := parsedTemplate.Execute(writer, page)
    if err != nil {
       log.Println("Error executing template :", err)
    }
  } else {
    //page.Content = []byte( TransformStr( string(page.Content) ) )
    page.HTMLContent = template.HTML( TransformStr( string(page.Content) ) ) //important later for when we create our own little language

    parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/edit.html")
    err := parsedTemplate.Execute(writer, page)
    if err != nil {
       log.Println("Error executing template :", err)
    }

  }
}

func newHandler(writer http.ResponseWriter, request *http.Request) {
  http.ServeFile( writer, request, "WebComponents/Templates/new.html");
}

func saveHandler( writer http.ResponseWriter, request *http.Request ) {
  //redirect if trying to use anything other than a POST
  if( request.Method != "POST" ) {
    http.Redirect( writer, request, "/home/", 300)
    return
  }

  var oldDocName string = request.PostFormValue("oldTitle")
  var newDocName string = request.PostFormValue("newTitle")
  var content string = request.PostFormValue("emitted")
  fmt.Printf("Received: %s for %s but maybe %s\n", content, newDocName, oldDocName)

  var page *WikiPage = &WikiPage{
    DocName : newDocName,
    path: "WebComponents/Documents/" + newDocName + ".csoc",
    statusCode: 200,
    Content: []byte(content) }

  var err error = writeWikiPage( page )
  if ( err != nil ) { //idk something went wrong
    writer.WriteHeader( 500 )
    log.Println("Error executing write :", err)
    writer.Write( []byte("Something went wrong, please try again later" ) )
  } else {
    if ( newDocName != oldDocName ) {
      os.Remove("WebComponents/Documents/" + oldDocName + ".csoc")
    }
    writer.Header().Set("HX-Redirect", "/view/" + newDocName)
  }


}

func homeHandler( writer http.ResponseWriter, request *http.Request ) {
  files, readErr := os.ReadDir("WebComponents/Documents")
  if readErr != nil {
    log.Fatal(readErr)
    return
  }

  var fileNames []string = make([]string, len(files), len(files) )

  for index, file := range files {
    var name string = file.Name()
    fileNames[ index ] = name[0:len(name) - 5]
  }

  parsedTemplate, _ := template.ParseFiles("WebComponents/Templates/home.html")
  err := parsedTemplate.Execute(writer, fileNames)
  if err != nil {
     log.Println("Error executing template :", err)
  }
}

func main() {
  //Initialize error files
  FileNotFound, _ = getErrorPage( 404 )
  ErrorOccur, _ = getErrorPage( 500 )


  fmt.Printf("Running...\n")

  http.HandleFunc("/", getFileContents); //catch all here
  http.HandleFunc("/view/", viewHandler);
  http.HandleFunc("/edit/", editHandler);
  http.HandleFunc("/new/", newHandler);
  http.HandleFunc("/save/", saveHandler);
  http.HandleFunc("/home/", homeHandler);

  log.Fatal(http.ListenAndServe(":8000", nil) )
}
