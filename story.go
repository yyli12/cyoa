package cyoa

import (
	"encoding/json"
	"net/http"
	"io"
	"html/template"
)

type StoryChapter struct {
	Title string
	Story []string
	Options []struct{
		Text, Arc string
	}
}

const InitPath = "intro"
const storyChapterHtmlTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{ .Title }}</title>
	</head>
	<body>
		<h1>{{ .Title }}</h1>
		<div>
			{{ range .Story }}
			<p>{{ . }}</p>
			{{ end }}
		</div>
		{{ range .Options }}
		<p><a href="/{{ .Arc }}">{{ .Text }}</a></p>
		{{end}}
	</body>
</html>`

func LoadStory(jsonBytes []byte) (map[string]StoryChapter, error) {
	story := make(map[string]StoryChapter)
	err := json.Unmarshal(jsonBytes, &story)
	if err != nil {
		return nil, err
	}
	return story, nil
}

func StoryHandler(story map[string]StoryChapter, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		if path == "/" {
			http.Redirect(writer, request, InitPath, http.StatusFound)
		} else if storyChapter, ok := story[path[1:]]; ok {
			storyChapter.RenderHtml(writer)
		} else {
			fallback.ServeHTTP(writer, request)
		}
	}
}

func (story * StoryChapter) RenderHtml(writer io.Writer) {
	t, err := template.New("story").Parse(storyChapterHtmlTemplate)
	if err != nil {
		panic(err)
	}
	err = t.Execute(writer, story)
	if err != nil {
		panic(err)
	}
}