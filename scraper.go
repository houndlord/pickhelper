package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapePatchInfo() (PatchInfo, error) {
	url := "https://www.op.gg/champions"
	filename := "op_gg_champions.html"

	// Download the page using wget
	cmd := exec.Command("wget", "-O", filename, url)
	err := cmd.Run()
	if err != nil {
		return PatchInfo{}, fmt.Errorf("error downloading page: %v", err)
	}
	defer os.Remove(filename)

	// Open the HTML file
	file, err := os.Open(filename)
	if err != nil {
		return PatchInfo{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Parse the HTML file
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return PatchInfo{}, fmt.Errorf("error parsing HTML: %v", err)
	}

	patchVersion := doc.Find(".css-17jvkpw").Text()
	patchVersion = strings.TrimPrefix(patchVersion, "Version: ")
	return PatchInfo{Version: patchVersion}, nil
}

func ScrapeChampions() ([]Champion, error) {
	url_ := "https://www.op.gg/champions"
	filename := "op_gg_champions.html"

	// Download the page using wget
	cmd := exec.Command("wget", "-O", filename, url_)
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error downloading page: %v", err)
	}
	defer os.Remove(filename)

	// Open the HTML file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Parse the HTML file
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	var champions []Champion

	doc.Find(".css-1hw6gn9").Each(func(i int, s *goquery.Selection) {
		avatarURL, _ := s.Find("img").Attr("src")
		championName := s.Text()
		parsedURL, _ := url.Parse(avatarURL)
		avatarURL = fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, parsedURL.Path)

		champion := Champion{
			Name:      championName,
			AvatarURL: avatarURL,
		}
		champions = append(champions, champion)
	})

	return champions, nil
}

func ScrapeMatchups(champName string) (map[string][]Matchup, error) {
	roles := []string{"top", "jungle", "mid", "adc", "support"}
	matchups := make(map[string][]Matchup)

	for _, role := range roles {
		url := fmt.Sprintf("https://www.op.gg/champions/%s/counters/%s", strings.ToLower(champName), role)
		filename := fmt.Sprintf("%s_%s_matchups.html", champName, role)

		// Download the page using wget
		cmd := exec.Command("wget", "-O", filename, url)
		err := cmd.Run()
		if err != nil {
			log.Printf("Error downloading page for %s %s: %v", champName, role, err)
			continue
		}
		defer os.Remove(filename)

		// Open the HTML file
		file, err := os.Open(filename)
		if err != nil {
			log.Printf("Error opening file for %s %s: %v", champName, role, err)
			continue
		}
		defer file.Close()

		// Parse the HTML file
		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			log.Printf("Error parsing HTML for %s %s: %v", champName, role, err)
			continue
		}

		var roleMatchups []Matchup
		doc.Find(".css-12a3bv1").Each(func(i int, s *goquery.Selection) {
			opponent := s.Find(".css-72rvq0").Text()
			winRate := s.Find(".css-ekbdas").Text()
			sampleSize := s.Find(".css-1nfew2i").Text()

			// Remove the '%' symbol from the win rate
			winRate = strings.TrimSuffix(winRate, "%")

			roleMatchups = append(roleMatchups, Matchup{
				Champion:   opponent,
				WinRate:    winRate,
				SampleSize: sampleSize,
			})
		})

		matchups[role] = roleMatchups
	}

	return matchups, nil
}
