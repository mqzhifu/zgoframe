package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
	httpresponse "zgoframe/http/response"
)

// @Tags Base
// @Summary 牛博网- 某人题库资料目录结构
// @Description 某人题库资料目录结构
// @Param X-Source-Type header string true "来源" default(11)
// @Param X-Project-Id header string true "项目ID"  default(6)
// @Param X-Access header string true "访问KEY" default(imzgoframe)
// @Success 200 {bool} bool "true:成功 false:否"
// @Router /base/niuke/question/dir/list [get]
func NiukeQuestionDirList(c *gin.Context) {

	catalogList := InitCata()

	catalogNum := 0
	for k2, v2 := range catalogList {
		//levelTwo := make(map[string][]NiukeQuestion)
		for k, v := range v2.Sub {
			var logList []NiukeQuestion
			for i := 0; i < 2; i++ {
				title := "chapter_title_" + v2.Name + "_" + v.Name + "_" + strconv.Itoa(i)
				oneData := NiukeQuestion{ChapterTitle: title, Content: "我是内容啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊啊 " + title}
				logList = append(logList, oneData)
				catalogNum++
			}
			v.CataLog = logList
			catalogList[k2].Sub[k].CataLog = logList
		}
	}

	niukeQuestionDirList := NiukeQuestionDirListStruct{
		Author:         "小z",
		ZhuanlanId:     1111,
		Title:          "BAT面试宝典之实践与总结",
		TypeA:          2222,
		AuthorId:       3333,
		IsSubscription: false,
		IsFree:         false,
		Price:          9999,
		Tutorial:       4444,
		CatalogList:    catalogList,
		CatalogNum:     catalogNum,
	}

	httpresponse.NiuKeOkWithDetailed(niukeQuestionDirList, "成功", c)
}

func InitCata() []CataInfo {
	cataInfoA := CataInfo{Id: 1, Key: "A", Name: "1级A", ParentKey: "root"}

	var subCateA []CataInfo
	subCateA = append(subCateA, CataInfo{Id: 11, Key: "level_2_A", Name: "2级A", ParentKey: "A"})
	subCateA = append(subCateA, CataInfo{Id: 11, Key: "level_2_B", Name: "2级B", ParentKey: "A"})
	subCateA = append(subCateA, CataInfo{Id: 11, Key: "level_2_C", Name: "2级C", ParentKey: "A"})
	cataInfoA.Sub = subCateA
	//subCate = append(subCate , )

	cataInfoB := CataInfo{Id: 2, Key: "B", Name: "1级B", ParentKey: "root"}

	var subCateB []CataInfo
	subCateB = append(subCateB, CataInfo{Id: 21, Key: "level_2_A", Name: "2级A", ParentKey: "B"})
	subCateB = append(subCateB, CataInfo{Id: 21, Key: "level_2_B", Name: "2级B", ParentKey: "B"})
	subCateB = append(subCateB, CataInfo{Id: 23, Key: "level_2_C", Name: "2级C", ParentKey: "B"})
	cataInfoB.Sub = subCateB

	cataInfoC := CataInfo{Id: 3, Key: "C", Name: "1级C", ParentKey: "root"}

	var subCateC []CataInfo
	subCateC = append(subCateA, CataInfo{Id: 11, Key: "level_2_A", Name: "2级A", ParentKey: "C"})
	subCateC = append(subCateA, CataInfo{Id: 11, Key: "level_2_B", Name: "2级B", ParentKey: "C"})
	subCateC = append(subCateA, CataInfo{Id: 11, Key: "level_2_C", Name: "2级C", ParentKey: "C"})
	cataInfoC.Sub = subCateC

	cataList := []CataInfo{cataInfoA, cataInfoB, cataInfoC}
	return cataList
}

type NiukeQuestionDirListStruct struct {
	Author         string     `json:"author"`
	CatalogList    []CataInfo `json:"catalog_list"`
	ZhuanlanId     int        `json:"zhuanlan_id"`
	Title          string     `json:"title"`
	TypeA          int        `json:"type"`
	AuthorId       int        `json:"author_id"`
	IsSubscription bool       `json:"is_subscription"`
	IsFree         bool       `json:"is_free"`
	Price          int        `json:"price"`
	Tutorial       int        `json:"tutorial"`
	CatalogNum     int        `json:"catalog_num"`
	//CataDirIndex   map[string]map[string]string          `json:"cata_dir_index"`
}

type NiukeQuestion struct {
	ChapterId    int    `json:"chapter_id"`
	ChapterTitle string `json:"chapter_title"`
	Content      string `json:"content"`
	HasPurchased bool   `json:"has_purchased"`
	Id           int    `json:"id"`
	SectionId    int    `json:"section_id"`
	SourceType   int    `json:"source_type"`
	Status       int    `json:"status"`
	Title        string `json:"title"`
	Uuid         int    `json:"uuid"`
	WordCount    int    `json:"word_count"`
}

type CataInfo struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	Key       string          `json:"key"`
	Sub       []CataInfo      `json:"sub"`
	CataLog   []NiukeQuestion `json:"cataLog"`
	ParentKey string          `json:"parent_key"`
}
