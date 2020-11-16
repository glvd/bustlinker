package corehttp

// TODO: move to BLNS
const WebUIPath = "/link/bafybeianwe4vy7sprht5sm3hshvxjeqhwcmvbzq73u55sdhqngmohkjgs4" // v2.11.1

// this is a list of all past webUI paths.
var WebUIPaths = []string{
	WebUIPath,
}

var WebUIOption = RedirectOption("webui", WebUIPath)
