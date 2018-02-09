package routers

import (
	"../dataProcess"
	"database/sql"
	"github.com/gin-gonic/gin"
)

var RouterMap = map[string]func(*gin.Context, *sql.DB){
	"tags":        dataProcess.QueryTagsUrl,
	"tagsinfo":    dataProcess.TgUrl,
	"count":       dataProcess.LineTagsUrl,
	"domains":     dataProcess.DomainUrl,
	"select":      dataProcess.GetSelect,
	"tagsbar":     dataProcess.GetTagsBarData,
	"barcount":    dataProcess.GetBarCountData,
	"tagtotal":    dataProcess.TotalData,
	"allflow":     dataProcess.GetAllFlow,
	"getdomains":  dataProcess.GetDomains,
	"getsiteflow": dataProcess.GetDFlow,
	"sitedetail":  dataProcess.GetSDetail,
	"flowtotal":   dataProcess.FlowTotal,
	"gettags":     dataProcess.GetTags,
	"arrival":     dataProcess.UpdateArrival,
	"cheat":       dataProcess.HandleCheat,
	"pvrate":      dataProcess.GetPvRate,
	"dimensions":  dataProcess.Dimensions,
	"custom":      dataProcess.GetCustomData,
}
