package kakaoCrawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/msyhu/GobbyIsntFree/etc"
	_struct "github.com/msyhu/GobbyIsntFree/struct"
	"net/http"
	"strconv"
	"strings"
)

var baseURL string = "https://careers.kakao.com/jobs?part=TECHNOLOGY&company=ALL"

type extractedJob = _struct.Kakao

func Crawling(kakaoC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan []extractedJob)

	totalPages := GetPages()

	for i := 1; i <= totalPages; i++ {
		go GetPage(i, c)
	}

	// TODO : waitgroup 이용해서 refactoring 해보기
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	kakaoC <- jobs
}

// 페이지 수를 가져온다
func GetPages() int {
	//lastPage := 1
	res, err := http.Get(baseURL)
	etc.CheckErr(err)
	etc.CheckCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	etc.CheckErr(err)

	pageSelection := doc.Find(".paging_list").Find("a")
	lastPageHref, _ := pageSelection.Last().Attr("href")
	lastPage := strings.Split(lastPageHref, "=")[1]
	page, err := strconv.Atoi(lastPage)
	etc.CheckErr(err)

	// 양쪽 화살표 4개 빼주고 현재 페이지 1 더해줌
	return page
}

// 하나의 페이지에서 직무를 가져와서 하나씩 채널로 넘겨준다.
func GetPage(page int, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&page=" + strconv.Itoa(page)
	fmt.Println(pageURL)
	res, err := http.Get(pageURL)
	etc.CheckErr(err)
	etc.CheckCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	etc.CheckErr(err)

	searchCards := doc.Find(".list_jobs>li")

	searchCards.Each(func(i int, card *goquery.Selection) {
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs

}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	// title
	title := card.Find(".tit_jobs").Text()

	// endDate, location
	var endDateAndLocation []string
	card.Find(".list_info>dd").Each(func(i int, s *goquery.Selection) {
		endDateAndLocation = append(endDateAndLocation, s.Text())
	})

	var endDate = ""
	var location = ""
	if len(endDateAndLocation) == 2 {
		endDate = endDateAndLocation[0]
		location = endDateAndLocation[1]
	} else {
		endDate = endDateAndLocation[0]
	}

	//jobGroups
	var jobGroups []string
	card.Find(".list_tag>a").Each(func(i int, s *goquery.Selection) {
		jobGroup, _ := s.Attr("data-code")
		jobGroups = append(jobGroups, jobGroup)
	})
	//company, jobType
	var companyAndJobType []string
	card.Find(".item_subinfo>dd").Each(func(i int, s *goquery.Selection) {
		companyAndJobType = append(companyAndJobType, s.Text())
	})
	company := companyAndJobType[0]
	jobType := companyAndJobType[1]

	//fmt.Println(company)

	c <- extractedJob{Title: title, EndDate: endDate, Location: location, JobGroups: jobGroups, Company: company, JobType: jobType}

}
