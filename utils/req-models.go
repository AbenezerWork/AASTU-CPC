package utils

// Represents the problem details
type Problem struct {
	ContestID int      `json:"contestId"`
	Index     string   `json:"index"`
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Points    float32  `json:"points"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

// Represents the author details
type Author struct {
	ContestID int `json:"contestId"`
	Members   []struct {
		Handle          string `json:"handle"`
		ParticipantType string `json:"participantType"`
		Ghost           bool   `json:"ghost"`
		Room            int    `json:"room"`
		StartTime       int    `json:"startTimeSeconds"`
	} `json:"members"`
}

// Represents a single submission result
type Submission struct {
	ID                  int     `json:"id"`
	ContestID           int     `json:"contestId"`
	CreationTimeSeconds int     `json:"creationTimeSeconds"`
	RelativeTimeSeconds int     `json:"relativeTimeSeconds"`
	Problem             Problem `json:"problem"`
	Author              Author  `json:"author"`
	Verdict             string  `json:"verdict"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
	PassedTestCount     int     `json:"passedTestCount"`
	TimeConsumedMillis  int     `json:"timeConsumedMillis"`
}

// Represents the full API response
type APIResponse struct {
	Status string       `json:"status"`
	Result []Submission `json:"result"`
}
