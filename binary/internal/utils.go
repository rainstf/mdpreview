package internal

import (
	"encoding/json"
	"github.com/gomarkdown/markdown"
	"fmt"
)

func ToMarkdown(input string) string {
	md := []byte(input)

	html := markdown.ToHTML(md, nil, nil)
	return string(html)
}

type Obj struct {
	Event string `json:"event"`
	Content string `json:"content"`
}

func EventsJSON(jsonObj []byte) Obj {
	var obj Obj

	err := json.Unmarshal(jsonObj, &obj)
	if err != nil {
		fmt.Printf("error unmarshalling json: %s", err)
	}

	return obj
}

type Page struct {
	content string
}

func newPage() *Page {
	base := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Live Preview</title>
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  `
	return &Page{content: base}
}

func scripts() string {
  return `
  <script>
  const evtSource = new EventSource("/events");
  evtSource.addEventListener("reload", (event) => {
	  const data = event.data.replace(/<br>/g, "\n");
	  let view = document.getElementById("container");
	  view.innerHTML = data;
	  console.log(event.data);
	});
  });
  </script>
  <script>
	  evtSource.addEventListener("scroll", (event) => {
		  let percentAmount = parseFloat(event.data);
		  var body = document.body,
    	  html = document.documentElement;

		  var height = Math.max( body.scrollHeight, body.offsetHeight, 
        	  html.clientHeight, html.scrollHeight, html.offsetHeight );
		  
		  let yLoc = percentAmount * height;
		  window.scrollTo({
			  top: yLoc,
			  left: 0,
		  });

		  console.log(yLoc);
	  });
  </script>
  </head>
<body>
<div id="container">`
}

func inject(p *Page, payload string) {
	p.content += "\n" + payload
}

func LayoutPage(html string) string {
	page := newPage()
	inject(page, Css())
	inject(page, scripts())
	inject(page, html)
	ending := `</div></body>
</html>`	
	inject(page, ending)
	return page.content
}
