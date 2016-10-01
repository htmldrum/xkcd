// TODO: Code review. Break up functions.
// TODO: Parallelize trawl/fetchComic/Index using go routines

package xkcd

import (
	"bytes"
	"encoding/json"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

const RootURL = "https://xkcd.com/"
const DocPath = "/info.0.json"

type ComixHeader struct {
	LastFetch time.Time
	FetchCount int
	LastNum int
}

type ComicNdx struct {
	Num int
	Terms []string
}

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

// Write support files for indexing
func Write() {
	var comixB bytes.Buffer
	var header bytes.Buffer
	comix := Trawl(1)
	enc := gob.NewEncoder(&comixB)
	err := enc.Encode(comix)
	if err != nil {
		fmt.Errorf("Failed to encode comix")
	}
	err = ioutil.WriteFile("comix.gob", comixB.Bytes(), 0644)
	if err != nil {
		fmt.Errorf("Failed to write comix.gob")
	}

	ch := ComixHeader{time.Now(), len(*comix), (*comix)[len((*comix))-1].Num}
	// ch := ComixHeader{time.Now(), len(*comix)}
	enc = gob.NewEncoder(&header)
	err = enc.Encode(ch)
	if err != nil {
		fmt.Errorf("Failed to encode comixHeader")
	}
	err = ioutil.WriteFile("comixHeader.gob", header.Bytes(), 0644)
	if err != nil {
		fmt.Errorf("Failed to write comixHeader.gob")
	}
}

// Read data type that represents an index for each comic
func Read() *[]*Comic {
	var comix []*Comic
	rdisk, err := ioutil.ReadFile("comix.gob")
	if err != nil {
		fmt.Errorf("Failed to read comix.gob")
	}
	dec := gob.NewDecoder(bytes.NewBuffer(rdisk))

	err = dec.Decode(&comix)
	if err != nil {
		fmt.Errorf("Failed to decode comix")
	}
	return &comix
}

// Write data type that represents an index for each comic
func WriteIndex(){
	var b bytes.Buffer
	var comicNdcs []ComicNdx
	comix := Read()

	for _, c := range *comix {
		comicNdcs = append(comicNdcs, IndexComic(c))
	}

	enc := gob.NewEncoder(&b)
	err := enc.Encode(comicNdcs)
	if err != nil {
		fmt.Errorf("Failed to encode comicNdcs")
	}
	err = ioutil.WriteFile("comicNdcs.gob", b.Bytes(), 0644)
	if err != nil {
		fmt.Errorf("Failed to write comicNdcs.gob")
	}
}

func IndexComic(c *Comic) ComicNdx{
	var terms []string
	num := c.Num
	terms = append(terms,
		c.Month,
		c.Link,
		c.Year,
		c.News,
		c.SafeTitle,
		c.Transcript,
		c.Alt,
		c.Img,
		c.Title,
		c.Day)
	return ComicNdx{num, terms}
}

// Search for key words in the index and return matches
func Search(){
	searchTerm := "sheep"
	fmt.Printf("Searching for %q.\n", searchTerm)
	m := Matches(searchTerm)
	if len(m) == 0 {
		fmt.Printf("No matches for %q.\n", searchTerm)
	} else {
		comix := Read()
		for num, matchCount := range m {
			fmt.Printf("%d matches for %q in doc #: %d.\n", matchCount, searchTerm, num)
			fmt.Printf("%q", (*comix)[num-1])
			fmt.Printf("=====\n\n")
		}
	}
}

func Matches(term string) map[int]int {
	matches := make(map[int]int)
	var comicNdcs []ComicNdx
	f, err := ioutil.ReadFile("comicNdcs.gob")
	if err != nil {
		fmt.Errorf("Failed to read comicNdcs.gob")
	}
	dec := gob.NewDecoder(bytes.NewBuffer(f))
	err = dec.Decode(&comicNdcs)
	if err != nil {
		fmt.Errorf("Failed to decode comicNdcs")
	}

	for _, c := range comicNdcs {
		for _, t := range c.Terms {
			if term == t {
				matches[c.Num]++
			}
		}
	}

	return matches
}

func Trawl(start int) *[]*Comic {
	var comics []*Comic
	count := 0
	for i := start; i != -1; i++ {
		comic, _ := FetchComic(i)
		if comic != nil {
			comics = append(comics, comic)
		} else {
			i = -1
		}
		count++
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
