package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {

	target := "https://www.xsbiquge.com/91_91546/"
	resp, err := http.Get(target)
	if err != nil {
		fmt.Println("get err http", err)
		return
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("html parse err", err)
		return
	}
	parseList(doc)
}

func parseList(doc *html.Node) {

	node := findSingleNode(doc)
	if node == nil {
		return
	}
	fmt.Println("Found div list")
	chapters := collectLinks(node)
	downloads(chapters)
	store(chapters)
}

func store(chapters []Chapter) {
	dir := "e:\\novel-test\\111\\"
	_ = os.MkdirAll(dir, os.ModePerm)
	for i := 0; i < len(chapters); i++ {
		_ = ioutil.WriteFile(dir+fmt.Sprint(i), []byte(chapters[i].Text), os.ModePerm)
	}
}

func downloads(chapters []Chapter) {
	download(&chapters[0])
	//for i := 0; i < len(chapters); i++ {
	//	download(chapters[i])
	//}
}

func download(chapter *Chapter) {
	resp, err := http.Get("https://www.xsbiquge.com" + chapter.Link)
	if err != nil {
		fmt.Println("get err http", err)
		return
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("html parse err", err)
		return
	}
	text := getText(doc)
	chapter.Text = text
}

func getText(doc *html.Node) string {
	node := findContentNode(doc)
	return patchText(node)
}

func patchText(node *html.Node) string {
	text := ""
	n := node.FirstChild
	for n != nil {
		if n.Type == html.ElementNode && n.Data == "br" {
			text += "\n"
		} else {
			text += n.Data
		}
		n = n.NextSibling
	}
	return text
}

func findContentNode(node *html.Node) *html.Node {

	if find(node, "div", html.Attribute{Key: "id", Val: "content"}) {
		return node
	}

	n := node.FirstChild
	for n != nil {
		r := findContentNode(n)
		if r != nil {
			return r
		}
		n = n.NextSibling
	}
	return nil
}
func collectLinks(node *html.Node) []Chapter {

	anchors := make([]*html.Node, 0)
	findMultipleNodes(node, &anchors)

	chapters := make([]Chapter, 0)
	for i := 0; i < len(anchors); i++ {
		chapters = append(chapters, Chapter{
			Title: anchors[i].FirstChild.Data,
			Link:  anchors[i].Attr[0].Val,
			Text:  "",
		})
	}
	return chapters
}

func findMultipleNodes(node *html.Node, anchors *[]*html.Node) {
	if findAnchor(node) {
		*anchors = append(*anchors, node)
	}

	n := node.FirstChild
	for n != nil {
		findMultipleNodes(n, anchors)
		n = n.NextSibling
	}
}

type Chapter struct {
	Title string
	Link  string
	Text  string
}

func findSingleNode(node *html.Node) *html.Node {

	if find(node, "div", html.Attribute{Key: "id", Val: "list"}) {
		return node
	}

	n := node.FirstChild
	for n != nil {
		r := findSingleNode(n)
		if r != nil {
			return r
		}
		n = n.NextSibling
	}
	return nil
}

func findAnchor(node *html.Node) bool {

	if node.Data != "a" {
		return false
	}
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == "href" && strings.HasSuffix(node.Attr[i].Val, ".html") {
			return true
		}
	}
	return false
}

func find(node *html.Node, targetData string, targetAttribute html.Attribute) bool {

	if node.Data != targetData {
		return false
	}
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == targetAttribute.Key && node.Attr[i].Val == targetAttribute.Val {
			return true
		}
	}
	return false
}
