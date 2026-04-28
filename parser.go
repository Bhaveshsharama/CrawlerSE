package main

import(
	"strings"
	"golang.org/x/net/html"
	"net/url"
	
)

func ParsePage(htmlStr string,baseURL string) (string,string,[]string){

	var contentRoot *html.Node 

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "","",nil
	}
	title := ""
	text := ""
	links:=[]string{}
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "div" { 
				for _, a := range n.Attr {
					if a.Key == "id" && a.Val == "mw-content-text" {
						contentRoot = n
					}
				}
			}

			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						link := attr.Val

						base, err := url.Parse(baseURL)
						if err!=nil {
							return
						}
						parsedLink, err := url.Parse(link) 
						if err!=nil {
							continue
						}
						absoluteLink := base.ResolveReference(parsedLink)

						if absoluteLink.Scheme == "http" || absoluteLink.Scheme == "https" {
							normalized := normalizeURL(absoluteLink.String())
							if normalized != "" {
								links = append(links, normalized)
							}
						}
					}
				}
			}
		}

		if n.Type == html.TextNode && isVisibleText(n) {
			clean := strings.TrimSpace(n.Data)
			if clean != "" {
				text += " " + clean
			}
		}


		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	
	if contentRoot != nil {
    var extractParagraphs func(*html.Node)
    extractParagraphs = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "p" {
            for c := n.FirstChild; c != nil; c = c.NextSibling {
                if c.Type == html.TextNode {
                    clean := strings.TrimSpace(c.Data)
                    if clean != "" {
                        text += " " + clean
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            extractParagraphs(c)
        }
    }
    text = "" // reset junk text
    extractParagraphs(contentRoot)
}
	return title,text,links

}