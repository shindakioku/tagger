package tagger

type Tagger interface {
	// Add your tag
	//    tagger.Add(tagger.New("test").InContract(nil))
	//    tagger.Add(tagger.New("test").InFunction(func() {}))
	Add(tag *Tag) Tagger
	// In for fill the struct
	// data is any value you need for fill struct (for example: json string)
	// in - struct for fill
	// tagForEmpty - your defined tag for fields which not have tags. Can be empty for to do nothing.
	// tags - list of available tags for process. You may to want to use some tags for some struct.
	//   If it's empty then will process all defined tags
	//
	//    type User struct {
	//      Username string `my_json:"user"`
	//      ID uint
	//    }
	//    tagger := NewReflectionTagger().
	//      Add(tagger.New("my_json").InFunction(func() {})).
	//      Add(tagger.New("foo").InFunction(func() {}))
	//    tagger.In("{\"username\": \"foo\"}", &User{}, "") // ID will not processed
	//    tagger.In("{\"username\": \"foo\"}", &User{}, "foo") // ID will processed by 'foo' tag
	In(data any, in any, tagForEmpty string, tags ...string) error
	// Out for fill any data from the struct
	// data - you may want to share some processed data between fields processing
	// out - full struct
	// tagForEmpty - your defined tag for fields which not have tags. Can be empty for to do nothing.
	// tags - list of available tags for process. You may to want to use some tags for some struct.
	//   If it's empty then will process all defined tags
	//
	//    type MyLoggingData struct {
	//      Fields []string
	//    }
	//    type User struct {
	//      Username string `my_json:"user"`
	//      ID uint
	//    }
	//
	//    loggingData := MyLoggingData{}
	//    tagger := NewReflectionTagger().
	//      Add(tagger.New("my_json").OutFunction(func() {})).
	//      Add(tagger.New("foo").OutFunction(func() {}))
	//    tagger.Out(&loggingData, &User{Username: "username"}, "") // ID will not processed
	//    tagger.Out(&loggingData, &User{Username: "username", ID: 1}, "foo") // ID will processed by 'foo' tag
	Out(data any, out any, tagForEmpty string, tags ...string) (any, error)
}
