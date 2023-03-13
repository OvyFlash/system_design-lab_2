package service

import (
	"bytes"
	"fmt"
	"io"
	"lab_2/config"
	"lab_2/pkg/models/settings"
	"lab_2/pkg/services/isw_reader/models"
	"lab_2/pkg/services/isw_reader/repository"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog"
)

type Service struct {
	*settings.Essential
	repo *repository.Repository
	l    zerolog.Logger
}

func NewService(e *settings.Essential, l zerolog.Logger) *Service {
	return &Service{
		Essential: e,
		l:         l,
	}
}

func (s *Service) Start() (err error) {
	if err = s.init(); err != nil {
		return
	}
	if err = s.tryToGetOld(); err != nil {
		return
	}
	if err = s.tryToParseOld(); err != nil {
		return
	}
	go s.scheduleReadAndParseNew()
	return nil
}

func (s *Service) init() (err error) {
	if err = s.initDatabases(); err != nil {
		return err
	}
	return
}

func (s *Service) initDatabases() (err error) {
	s.repo, err = repository.NewRepository(s.Storage)
	return
}

//
// Scheduled getting page
//

func (s *Service) scheduleReadAndParseNew() {
	ticker := time.NewTicker(time.Hour)
	ctx := s.GetContext()
	page, err := s.repo.GetLastPage()
	if err != nil {
		s.l.Error().Msg(err.Error())
	}
	for now := time.Now(); ; now = time.Now() {
		select {
		case <-ticker.C:
			if page.Date.Format(dateLayout) == now.Format(dateLayout) {
				s.l.Info().Msg("No new pages")
				continue
			}
			s.requestAndAddPage(getLink(now), now)
			page, err = s.repo.GetLastPage()
			if err != nil {
				s.l.Error().Msg(err.Error())
				continue
			}
			err = s.parseRawPage(page)
			if err != nil {
				s.l.Error().Msg(err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

//
// Fill missing raw pages in database
//

/*
Missing:
	2022-05-05
	2022-07-11
	2022-11-24
	2022-12-25
	2023-01-01
	2023-02-05
	2023-03-05
*/

const (
	dateLayout = "2006-01-02"
	baseUrl    = "https://www.understandingwar.org/backgrounder/"
)

func getLink(date time.Time) string {
	pattern := baseUrl + "russian-offensive-campaign-assessment-%s-%d"
	if date.Year() > 2022 {
		pattern = baseUrl + "russian-offensive-campaign-assessment-%s-%d-2023"
	}
	return fmt.Sprintf(pattern, strings.ToLower(date.Month().String()), date.Day())
}

func (s *Service) generateDateToLinkPattern() (map[time.Time]string, error) {
	var dateToLink = map[time.Time]string{
		time.Date(2022, time.February, 24, 1, 0, 0, 0, config.GetLocation()): baseUrl + "russia-ukraine-warning-update-initial-russian-offensive-campaign-assessment",
		time.Date(2022, time.February, 25, 1, 0, 0, 0, config.GetLocation()): baseUrl + "russia-ukraine-warning-update-russian-offensive-campaign-assessment-february-25-2022",
		time.Date(2022, time.February, 26, 1, 0, 0, 0, config.GetLocation()): baseUrl + "russia-ukraine-warning-update-russian-offensive-campaign-assessment-february-26",
		time.Date(2022, time.February, 27, 1, 0, 0, 0, config.GetLocation()): baseUrl + "russia-ukraine-warning-update-russian-offensive-campaign-assessment-february-27",
		time.Date(2022, time.February, 28, 1, 0, 0, 0, config.GetLocation()): baseUrl + "russian-offensive-campaign-assessment-february-28-2022",
	}

	startDate := time.Date(2022, time.March, 1, 1, 0, 0, 0, config.GetLocation())
	now := time.Now()
	var ok bool
	for ; now.After(startDate); startDate = startDate.Add(time.Hour * 24) {
		_, ok = dateToLink[startDate]
		if !ok {
			dateToLink[startDate] = getLink(startDate)
		}
	}
	err := s.clearPresentPages(dateToLink)
	return dateToLink, err
}

func (s *Service) clearPresentPages(m map[time.Time]string) (err error) {
	presentPages, err := s.getPresentPages()
	if err != nil {
		return
	}
	var ok bool
	for date := range m {
		_, ok = presentPages[date.Format(dateLayout)]
		if ok {
			delete(m, date)
		}
	}
	return
}

func (s *Service) getPresentPages() (map[string]struct{}, error) {
	pages, err := s.repo.GetPresentPagesDates()
	if err != nil {
		s.l.Error().Msg(err.Error())
		return nil, err
	}
	presentPages := map[string]struct{}{}
	for _, v := range pages {
		presentPages[v.Date.Format(dateLayout)] = struct{}{}
	}
	return presentPages, nil
}

func (s *Service) tryToGetOld() error {
	links, err := s.generateDateToLinkPattern()
	if err != nil {
		return err
	}
	for k, v := range links {
		<-time.After(time.Second / 30)
		s.requestAndAddPage(v, k)
	}
	return nil
}

func (s *Service) requestAndAddPage(link string, date time.Time) {
	resp, err := http.Get(link)
	if err != nil {
		s.printGetErr(date, err)
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		s.printGetErr(date, fmt.Errorf("page not found"))
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.printGetErr(date, err)
		return
	}
	err = s.repo.WriteRawData(models.RawISWPage{
		Date:    date,
		RawPage: string(body),
	})
	if err != nil {
		s.printGetErr(date, err)
		return
	}
	s.l.Info().Msgf("Added page for date %s", date.Format(dateLayout))
}

func (s *Service) printGetErr(date time.Time, err error) {
	s.l.Error().Msgf("Could not get page for date: %s, err: %s",
		date.Format(dateLayout), err.Error())
}

//
// Parse old pages
//

func (s *Service) tryToParseOld() (err error) {
	res, err := s.repo.GetUnprocessedPages()
	if err != nil {
		s.l.Error().Msg(err.Error())
		return
	}
	for _, v := range res {
		err = s.parseRawPage(v)
		if err != nil {
			s.l.Error().Msg(err.Error())
		}
	}
	return nil
}

func (s *Service) parseRawPage(raw models.RawISWPage) (err error) {
	var parsed models.ParsedISWPage
	parsed.Date = raw.Date
	parsed.ParsedPage, err = s.parse(raw.RawPage)
	if err != nil {
		return
	}
	parsed.ParsedPage = s.processParsedText(parsed.ParsedPage)
	if len(parsed.ParsedPage) == 0 {
		return fmt.Errorf("%s: zero parsed page length", parsed.Date.Format(dateLayout))
	}
	return s.repo.WriteParsedPage(parsed)
}

var (
	timeZone   *regexp.Regexp
	references *regexp.Regexp
	links      *regexp.Regexp
)

func init() {
	var err error
	timeZone, err = regexp.Compile(`\b(January|February|March|April|May|June|July|August|September|October|November|December)\s+\d{1,2}(,|)\s+\d{1,2}(:\d{2})?\s*(am|pm)?\s*[A-Z]{2,4}\b`)
	if err != nil {
		panic(err)
	}
	references, err = regexp.Compile(`\[\d+\]`)
	if err != nil {
		panic(err)
	}
	links, err = regexp.Compile(`((\[\d+\](| )+)|)(http|https):\/\/[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?`)
	if err != nil {
		panic(err)
	}
}

func (s *Service) parse(text string) (data []string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(text)))
	if err != nil {
		return
	}
	allowSaveText := false
	var counter int
	for i := 0; i < 2; i++ {
		doc.Find("div.field-item.even > p").Each(func(i int, s *goquery.Selection) {
			counter++
			if timeZone.MatchString(s.Text()) {
				allowSaveText = true
				return
			}
			if !allowSaveText {
				return
			}
			if links.MatchString(s.Text()) {
				return
			}
			finalText := references.ReplaceAllString(s.Text(), "")
			if finalText == "" {
				return
			}
			data = append(data, finalText)
		})
		if len(data) == 1 {
			data = []string{}
		}
		//для випадку, коли дати немає
		if len(data) == 0 && counter > 1 {
			allowSaveText = true
		} else {
			break
		}
	}
	return
}

func (s *Service) processParsedText(text []string) (res []string) {
	//todo? call python
	return text
}
