package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AbenezerWork/AASTU-CPC/models"
)

//TODO: check user handle against the handle of the submission

func GetAndCheckAdmission(problem models.Problem, submissionNo string, cfusername string) (error, bool) {
	contestID := problem.ContestID
	submissionID := submissionNo

	req := fmt.Sprintf("https://codeforces.com/api/contest.status?contestId=%s&handle=%s", contestID, cfusername)
	fmt.Println("request: ", req)
	res, err := http.Get(req)
	if err != nil {
		return err, false
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err, false
	}
	var apiResponse APIResponse

	err = json.Unmarshal(body, &apiResponse)

	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err, false
	}

	fmt.Println(apiResponse)
	for _, submission := range apiResponse.Result {
		subid, _ := strconv.Atoi(submissionID)
		if subid == submission.ID {

			contid, _ := strconv.Atoi(problem.ContestID)
			fmt.Println(contid, submission.ContestID, problem.Index, submission.Problem.Index, submission.Verdict)
			if contid != submission.ContestID || problem.Index != submission.Problem.Index || submission.Verdict != "OK" {
				continue
			}
			return nil, true
		}

	}
	return errors.New("your submission is not correct"), false

}
