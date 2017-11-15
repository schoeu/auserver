package routers

import (
	"../dataProcess"
	"github.com/gin-gonic/gin"
	"database/sql"
)

var RouterMap = map[string]func(*gin.Context, *sql.DB){
	"tags": dataProcess.QueryTagsUrl,
	"tagsinfo": dataProcess.TgUrl,
	"count": dataProcess.LineTagsUrl,
	"domains": dataProcess.DomainUrl,
	"select": dataProcess.GetSelect,
	"tagsbar": dataProcess.GetTagsBarData,
	"barcount": dataProcess.GetBarCountData,
	"tagtotal": dataProcess.TotalData,
	"allflow": dataProcess.GetAllFlow,
	"getdomains": dataProcess.GetDomains,
	"getsiteflow": dataProcess.GetDFlow,
	"sitedetail": dataProcess.GetSDetail,
	"getnewer": dataProcess.GetNewer,
}
