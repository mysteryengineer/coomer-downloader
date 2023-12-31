package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pterm/pterm"
	"github.com/thoas/go-funk"
	"math"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func CheckUser(service string, user string) (string, error) {
	url := fmt.Sprintf("https://coomer.su/%s/user/%s", service, user)
	html := FetchUrl(url)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", fmt.Errorf("Error fetching user %s", user)
	}

	result := doc.Find("span[itemprop='name']")
	name := strings.ToLower(result.Text())

	if name == "" {
		return "", fmt.Errorf("Error fetching user %s", user)
	}

	return name, err
}

func GetMedias(service string, user string, directory string, limit int) []Media {
	pterm.Print("\n📝 Collecting files from user ", pterm.Bold.Sprintf(user), " on ", pterm.Bold.Sprintf(service), " ")

	url := fmt.Sprintf("https://coomer.su/%s/user/%s?o=0", service, user)
	numPages := countPages(url)

	posts := make([]string, 0)
	medias := make([]Media, 0)

	for i := 0; i < numPages; i++ {
		url = fmt.Sprintf("https://coomer.su/%s/user/%s?o=%d", service, user, i*50)
		postsPerPage := getPosts(url)
		posts = append(posts, postsPerPage...)

		for _, postUrl := range postsPerPage {
			mediasPerPost := getMedias(postUrl, directory)
			medias = append(medias, mediasPerPost...)
		}

		pterm.Print(".")

		if len(medias) >= limit {
			medias = medias[:limit]
			break
		}
	}

	pterm.Printf(" found %d posts/%d files\n", len(posts), len(medias))
	return medias
}

func FilterExtensions(medias []Media, extensions []string) []Media {
	ext := strings.Join(extensions, ", ")
	pterm.Print("🔍 Filtering files with extension ", pterm.Bold.Sprintf(ext), " ... ")

	filtered := funk.Filter(medias, func(media Media) bool {
		return funk.ContainsString(extensions, media.Ext())
	}).([]Media)

	pterm.Printf("%d files remained\n", len(filtered))
	return filtered
}

// region - Private functions

func countPages(url string) int {
	html := FetchUrl(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		PrintError("Error parsing HTML: %s", err)
	}

	result := doc.Find("div#paginator-top small")
	var num int

	if result.Size() > 0 {
		str := result.Text()
		re := regexp.MustCompile(`of (\d+)`)
		matches := re.FindStringSubmatch(str)
		num, _ = strconv.Atoi(matches[1])
		num = int(math.Ceil(float64(num) / 50))
	} else {
		num = 1
	}

	return num
}

func getPosts(url string) []string {
	html := FetchUrl(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		PrintError("Error parsing HTML: %s", err)
	}

	posts := doc.Find("article").
		Map(func(i int, s *goquery.Selection) string {
			id, _ := s.Attr("data-id")
			service, _ := s.Attr("data-service")
			user, _ := s.Attr("data-user")

			postUrl := fmt.Sprintf("https://coomer.su/%s/user/%s/post/%s", service, user, id)
			return postUrl
		})

	return posts
}

func getMedias(url string, directory string) []Media {
	html := FetchUrl(url)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		PrintError("Error parsing HTML: %s", err)
	}

	postId, _ := doc.Find("meta[name='id']").
		Attr("content")

	images := doc.Find("a.fileThumb").
		Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")
			return link
		})

	videos := doc.Find("a.post__attachment-link").
		Map(func(i int, s *goquery.Selection) string {
			link, _ := s.Attr("href")
			return link
		})

	links := append(images, videos...)
	medias := make([]Media, 0)

	for i, link := range links {
		fileName := fmt.Sprintf("%s-%d%s", postId, i+1, path.Ext(link))
		filePath := filepath.Join(directory, fileName)

		medias = append(medias, Media{
			Url:      link,
			FilePath: filePath,
		})
	}

	return medias
}

// endregion
