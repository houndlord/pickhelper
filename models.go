package main

type Champion struct {
	Name      string
	AvatarURL string
}

type Matchup struct {
	Champion   string
	WinRate    string
	SampleSize string
}

type PatchInfo struct {
	Version string
}

type ScrapingStatus struct {
	CurrentPatch     string
	LastScrapedPatch string
	IsUpdating       bool
}
