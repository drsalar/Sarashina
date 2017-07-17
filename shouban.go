package main

/*手办维基爬厂商数据：http://www.hpoi.net/
1.厂商名称
2.logo
3. 中文名
4. 别名
5. 所在地
6. 厂商介绍

手办酱爬ip数据：https://www.shoubanjiang.com/anime
1. ip名称
2. iplogo
3. 中文名
4. 别名
5. ip介绍
*/

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx"
	"strings"
)

var file *xlsx.File
var sheet *xlsx.Sheet
var id = 205
var hpoiUrl = "http://www.hpoi.net/"
var shoubanjiangUrl = "https://www.shoubanjiang.com/anime/company"
var shoubanjiangPageCount = 36
var companyName = `名称`
var chineseName = `中文名`
var alias = `别名`
var location = `地区`
var location2 = `所在地`
var location3 = `总部地点`
var location4 = `本社地址`
var description = ""

type HpoiInfo struct {
	companyName string
	logo        string
	chineseName string
	alias       string
	location    string
	description string
}

type ShoubanjiangInfo struct {
	companyName string
	logo        string
	chineseName string
	description string
}

func main() {
	var err error
	file, err = xlsx.OpenFile("/home/salar/Documents/ips_test.xlsx")
	if err != nil {
		log.WithFields(log.Fields{
			"func": "main",
		}).Error(err.Error())
		return
	}
	sheet = file.Sheets[0]
	//getHpoi()
	getShoubanjiang()
	file.Save("/home/salar/Documents/ips_test2.xlsx")
}

func getShoubanjiang() {
	for i := 0; i < shoubanjiangPageCount; i++ {
		doc, err := goquery.NewDocument(fmt.Sprintf("%s?page=%d", shoubanjiangUrl, i))
		if err != nil {
			log.WithFields(log.Fields{
				"func": "getShoubanjiang",
				"url":  shoubanjiangUrl,
			}).Error(err.Error())
			return
		}
		doc.Find(".companyboard-brand").Each(func(i int, s *goquery.Selection) {
			url, isExist := s.Attr("href")
			if isExist {
				img, _ := s.Children().Attr("src")
				getShoubanjiangCompanyDetail(url, img)
			}
		})
	}
}

func getShoubanjiangCompanyDetail(url string, img string) {
	var si ShoubanjiangInfo
	si.logo = img
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "getShoubanjiangCompanyDetail",
			"url":  url,
		}).Error(err.Error())
	}
	si.chineseName = doc.Find(".list-header-name").Text()
	si.companyName = doc.Find(".list-header-jpname").Text()
	si.description = doc.Find(".list-header-tab-pane").Text()

	log.WithFields(log.Fields{
		"func":        "getShoubanjiangCompanyDetail",
		"url":         url,
		"companyName": si.companyName,
		"logo":        si.logo,
		"csineseName": si.chineseName,
		"desc":        si.description,
	}).Info()
	insert2Excel2(si)
}

func getHpoi() {
	doc, err := goquery.NewDocument(hpoiUrl + "company")
	if err != nil {
		log.WithFields(log.Fields{
			"func": "getHpoi",
			"url":  hpoiUrl,
		}).Error(err.Error())
		return
	}
	//fmt.Println(doc.Text())
	doc.Find(".bs-glyphicons li .caption").Each(func(i int, s *goquery.Selection) {
		company := s.Text()
		href, isExist := s.Attr("href")
		if isExist {
			log.WithFields(log.Fields{
				"func":    "getHpoi",
				"url":     hpoiUrl,
				"index":   i,
				"company": company,
				"href":    href,
			}).Info()
			getHpoiCompanyDetail(hpoiUrl + href)
		} else {
			log.WithFields(log.Fields{
				"func":    "getHpoi",
				"url":     hpoiUrl,
				"index":   i,
				"company": company,
				"href":    href,
			}).Info()
		}
	})
}

func getHpoiCompanyDetail(url string) {
	var hi HpoiInfo
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.WithFields(log.Fields{
			"func": "getHpoiCompanyDetail",
			"url":  url,
		}).Error(err.Error())
		return
	}
	doc.Find(".col-xs-24 .thumbnail").Each(func(i int, s *goquery.Selection) {
		img, isExist := s.Children().Attr("src")
		if isExist {
			// log.WithFields(log.Fields{
			// 	"func":  "getHpoiCompanyDetail",
			// 	"url":   url,
			// 	"index": i,
			// 	"src":   img,
			// }).Info()
			hi.logo = img
		} else {
			// log.WithFields(log.Fields{
			// 	"func":  "getHpoiCompanyDetail",
			// 	"url":   url,
			// 	"index": i,
			// 	"src":   img,
			// }).Info()
		}
	})
	desc := ""
	doc.Find(".detail-content").Children().Each(func(i int, s *goquery.Selection) {
		if desc == "" {
			desc = strings.Replace(s.Text(), "&nbsp;", "", -1)
		} else {
			desc = desc + "\n" + strings.Replace(s.Text(), "&nbsp;", "", -1)
		}
	})
	hi.description = desc
	var t []string
	n := 0
	doc.Find(".col-xs-24 .table").Children().Children().Each(func(i int, s *goquery.Selection) {
		s.Children().Each(func(i int, s *goquery.Selection) {
			x := strings.Split(s.Text(), `:`)[0]
			x = strings.Replace(x, "\n", "", -1)
			x = strings.Replace(x, ` `, "", -1)
			t = append(t, x)
			n++
		})
	})
	for i := 0; i < n; i += 2 {
		switch t[i] {
		case companyName:
			hi.companyName = t[i+1]
			break
		case chineseName:
			hi.chineseName = t[i+1]
			break
		case alias:
			hi.alias = t[i+1]
			break
		case location, location2, location3, location4:
			hi.location = t[i+1]
			break
		default:
			fmt.Println(t[i])
			break
		}
	}
	log.WithFields(log.Fields{
		"func":        "getHpoiCompanyDetail",
		"url":         url,
		"companyName": hi.companyName,
		"logo":        hi.logo,
		"chineseName": hi.chineseName,
		"alias":       hi.alias,
		"location":    hi.location,
		"desc":        hi.description,
	}).Info()
	insert2Excel(hi)
}

func insert2Excel(hi HpoiInfo) {
	id++
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = fmt.Sprintf("%d", id)
	cell = row.AddCell()
	cell.Value = hi.companyName
	cell = row.AddCell()
	cell.Value = hi.chineseName
	cell = row.AddCell()
	cell.Value = hi.alias
	cell = row.AddCell()
	cell.Value = hi.location
	cell = row.AddCell()
	cell.Value = hi.logo
	cell = row.AddCell()
	cell.Value = hi.description
	cell = row.AddCell()
	cell.Value = "手办维基"
}

func insert2Excel2(si ShoubanjiangInfo) {
	id++
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = fmt.Sprintf("%d", id)
	cell = row.AddCell()
	cell.Value = si.companyName
	cell = row.AddCell()
	cell.Value = si.chineseName
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = si.logo
	cell = row.AddCell()
	cell.Value = si.description
	cell = row.AddCell()
	cell.Value = "手办酱"
}
