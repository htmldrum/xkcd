package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const RootURL = "https://xkcd.com/"
const DocPath = "/info.0.json"

type Comic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

func Trawl(start int) *[]*Comic {
	var comics []*Comic
	for i := start; i != nil; i++ {
		comic, _ := FetchComic(i)
		if comic != nil {
			comics = append(comics, comic)
		} else {
			i = nil
		}
	}
	return &comics
}

func FetchComic(id int) (*Comic, error) {
	var result Comic
	resp, err := http.Get(RootURL + strconv.Itoa(id) + DocPath)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("Failed to fetch document ID: %d", id)
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

// {
//		"month": "4",
//		"num": 571,
//		"link": "",
//		"year": "2009",
//		"news": "",
//		"safe_title": "Can't Sleep",
//		"transcript": "[[Someone is in bed, presumably trying to sleep. The top of each panel is a thought bubble showing sheep leaping over a fence.]]\n1 ... 2 ...\n<<baaa>>\n[[Two sheep are jumping from left to right.]]\n\n... 1,306 ... 1,307 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow.]]\n\n... 32,767 ... -32,768 ...\n<<baaa>> <<baaa>> <<baaa>> <<baaa>> <<baaa>>\n[[A whole flock of sheep is jumping over the fence from right to left. The would-be sleeper is sitting up.]]\nSleeper: ?\n\n... -32,767 ... -32,766 ...\n<<baaa>>\n[[Two sheep are jumping from left to right. The would-be sleeper is holding his pillow over his head.]]\n\n{{Title text: If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.}}",
//		"alt": "If androids someday DO dream of electric sheep, don't forget to declare sheepCount as a long int.",
//		"img": "http://imgs.xkcd.com/comics/cant_sleep.png",
//		"title": "Can't Sleep",
//		"day": "20"
// }
