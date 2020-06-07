package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// "entity_id":"774154072"
var entityRE = regexp.MustCompile(`"entity_id":"(\d+)"`)

func getUsername(url string) string {
	return fmt.Sprintf("https://www.facebook.com/%s", path.Base(url))
}

func getPresetCookies() string {
	// hackish hardcoded cookies because i dont have a lot of time to work on this
	byteCookie, err := ioutil.ReadFile("cookie.txt")
	if err != nil {
		return ""
	}
	return string(byteCookie)
}

// getUserID works for both authenticated and unauthenticated requests
func getUserID(url string) (int, error) {
	resp, err := http.Get(getUsername(url))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	body := string(byteBody)

	match := entityRE.FindString(body)
	if len(match) == 0 {
		return 0, errors.New("fbid not found")
	}

	cleanedID := strings.Replace(strings.Replace(match, `"entity_id":"`, "", 1), `"`, "", 1)

	return strconv.Atoi(cleanedID)
}

// getUserIDAuthed will only work if the request is authenticated
func getUserIDAuthed(url string) (int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", getUsername(url), nil)
	req.Header.Set("cookie", getPresetCookies())

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, err
	}

	var fbURL string
	var exists bool
	// <meta property="al:ios:url" content="fb://profile/774154072" />
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("property"); property == "al:ios:url" {
			content, ok := s.Attr("content")
			if ok {
				fbURL = content
				exists = true
			}
		}
	})

	if !exists {
		return 0, errors.New("fbid not found")
	}

	cleanedID := path.Base(fbURL)
	return strconv.Atoi(cleanedID)
}

// getApproxCreateYearAuthed will try to approximate the create year by using common sense
// can be modified to be more specific like create month but i don't have time for that
func getApproxCreateYearAuthed(url string) ([]int, error) {
	yearList := []int{}
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("cookie", getPresetCookies())

	resp, err := client.Do(req)
	if err != nil {
		return yearList, err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return yearList, err
	}

	body := string(byteBody)
	body = strings.ReplaceAll(strings.ReplaceAll(body, "<!--", ""), "-->", "")

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return yearList, err
	}

	years := doc.Find(".fbStickyHeaderBreadcrumb").Find(".lastItem")
	years.Find("select > option").Each(func(i int, s *goquery.Selection) {
		year, err := strconv.Atoi(s.Text())
		if err == nil {
			yearList = append(yearList, year)
		}
	})

	return yearList, nil
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.GET("/fb/:profile", func(c *gin.Context) {
		profile := c.Param("profile")

		fbid, err := getUserIDAuthed(profile)
		if err != nil {
			panic(err)
		}

		yearList := make(map[int][]int)
		yearCounts := make(map[int]int)

		// bruteforce find created at date
		// this attack vector is not fail safe
		for i := fbid - 10; i < fbid+10; i++ {
			years, err := getApproxCreateYearAuthed(fmt.Sprintf("https://www.facebook.com/%d", i))
			if err == nil {
				if len(years) > 0 {
					yearList[i] = years
					year := years[len(years)-1]
					_, ok := yearCounts[year]
					if !ok {
						yearCounts[year] = 0
					}
					yearCounts[year]++
				}

				// slow down to avoid getting spam blocked
				time.Sleep(2 * time.Second)
			}
		}

		if len(yearCounts) > 0 {
			type kv struct {
				Key   int
				Value int
			}

			var ss []kv
			for k, v := range yearCounts {
				ss = append(ss, kv{k, v})
			}

			sort.Slice(ss, func(i, j int) bool {
				return ss[i].Value > ss[j].Value
			})

			c.JSON(200, gin.H{
				"create_year": ss[0].Key,
			})

			return
		}

		c.JSON(400, gin.H{"error": "1000"})
	})

	r.Run()
}
