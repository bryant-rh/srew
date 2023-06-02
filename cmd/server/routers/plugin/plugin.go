package plugin

import (
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/bryant-rh/srew/cmd/server/global"
	"github.com/bryant-rh/srew/pkg/index"
	"github.com/bryant-rh/srew/pkg/util"
	jsoniter "github.com/json-iterator/go"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/bryant-rh/srew/internal/model"
	"github.com/bryant-rh/srew/internal/query"

	"github.com/gin-gonic/gin"
)

func PluginRouter(r *gin.RouterGroup) {
	r.GET("/plugin/list", ListPlugin)
	r.GET("/plugin/listall", ListPluginAllVersion)
	r.GET("/plugin/search", SearchPlugin)
	r.GET("/plugin/searchall", SearchPluginAllVersion)
	r.POST("/plugin/create", CreatePlugin)
	r.PUT("/plugin/update", UpdatePlugin)
	r.DELETE("plugin/delete", DeletePlugin)
}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary ListPlugin
// @Schemes
// @Description List Plugin
// @Tags ListPlugin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param plugin_name query string false "Plugin Name"
// @Param version query string false "Plugin Version"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /plugin/list [get]
// @ID ListPlugin
func ListPlugin(ctx *gin.Context) {
	plugin_name := ctx.Query("plugin_name")
	plugin_version := ctx.Query("version")

	plugin_Detail := []model.Detail{}

	if plugin_name == "" && plugin_version == "" {
		// err := q.Cluster.WithContext(ctx).Scan(&cluster)
		// if err != nil {
		// 	error_msg := fmt.Sprintf("err: %s ", err)
		// 	util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		// 	return
		// }
		global.Config.DB.DB().
			Raw("SELECT  * FROM  detail  where detail.version IN (SELECT plugin.latest_version FROM plugin)").
			Scan(&plugin_Detail)

	} else if plugin_name != "" && plugin_version == "" {
		global.Config.DB.DB().
			Raw(" SELECT  * FROM  detail  where detail.version IN (SELECT plugin.latest_version FROM plugin) and detail.plugin_name = ? ", plugin_name).
			Scan(&plugin_Detail)

	} else if plugin_version != "" {
		q := query.Use(global.Config.DB.DB()).Detail
		PluginList, err := q.WithContext(ctx).Select(q.PluginID).Where(q.PluginName.Eq(plugin_name), q.Version.Eq(plugin_version)).First()

		if err != nil {
			if PluginList == nil {
				error_msg := fmt.Sprintf("The PluginName: [%s] is not Found! for Version: [%s]", plugin_name, plugin_version)
				util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			}
			return
		}
		global.Config.DB.DB().
			Raw(" SELECT  * FROM  detail  where detail.version = ? and detail.plugin_name = ? ", plugin_version, plugin_name).
			Scan(&plugin_Detail)

	}
	if len(plugin_Detail) == 0 {
		util.ReturnMsg(ctx, http.StatusNotFound, "", "The Plugin is not Found!")
		return
	} else {
		util.ReturnMsg(ctx, http.StatusOK, plugin_Detail, "success")
	}
}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary ListPluginAllVersion
// @Schemes
// @Description List Plugin ALL Version
// @Tags ListPluginAllVersion
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param plugin_name query string false "Plugin Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /plugin/listall [get]
// @ID ListPluginAllVersion
func ListPluginAllVersion(ctx *gin.Context) {
	plugin_name := ctx.Query("plugin_name")
	//result := make(map[string]interface{})
	result := make(map[string][]*model.Detail)
	ok_msg := ""
	q := query.Use(global.Config.DB.DB())

	if plugin_name == "" {
		plugin_list, err := q.Plugin.WithContext(ctx).Select(q.Plugin.PluginName).Find()
		if err != nil {
			error_msg := fmt.Sprintf("ListPluginAllVersion err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		if len(plugin_list) != 0 {
			for _, v := range plugin_list {
				plugin_Detail, err := q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Eq(plugin_name)).Find()
				if err != nil {
					error_msg := fmt.Sprintf("List Plugin Detail err: %s ", err)
					util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
					return
				}
				result[v.PluginName] = plugin_Detail

			}

		}
		ok_msg = "The ListPluginAllVersion successfully!"

	} else {
		plugin_Detail, err := q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Eq(plugin_name)).Find()
		if err != nil {
			error_msg := fmt.Sprintf("List Plugin Detail err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		result[plugin_name] = plugin_Detail
		ok_msg = fmt.Sprintf("The PluginName: [%s] ListALlVersion successfully!", plugin_name)

	}

	util.ReturnMsg(ctx, http.StatusOK, result, ok_msg)

}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary SearchPlugin
// @Schemes
// @Description Search Plugin
// @Tags SearchPlugin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param plugin_name query string false "Plugin Name"
// @Param version query string false "Plugin Version"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /plugin/search [get]
// @ID SearchPlugin
func SearchPlugin(ctx *gin.Context) {
	plugin_name := ctx.Query("plugin_name")
	plugin_version := ctx.Query("version")

	plugin_Detail := []model.Detail{}

	if plugin_name == "" && plugin_version == "" {
		// err := q.Cluster.WithContext(ctx).Scan(&cluster)
		// if err != nil {
		// 	error_msg := fmt.Sprintf("err: %s ", err)
		// 	util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
		// 	return
		// }
		global.Config.DB.DB().
			Raw("SELECT  * FROM  detail  where detail.version IN (SELECT plugin.latest_version FROM plugin)").
			Scan(&plugin_Detail)

	} else if plugin_name != "" && plugin_version == "" {
		global.Config.DB.DB().
			Raw(" SELECT  * FROM  detail  where detail.version IN (SELECT plugin.latest_version FROM plugin) and detail.plugin_name like ? ", "%"+plugin_name+"%").
			Scan(&plugin_Detail)

	} else if plugin_version != "" {
		q := query.Use(global.Config.DB.DB()).Detail
		PluginList, err := q.WithContext(ctx).Select(q.PluginID).Where(q.PluginName.Eq(plugin_name), q.Version.Eq(plugin_version)).First()

		if err != nil {
			if PluginList == nil {
				error_msg := fmt.Sprintf("The PluginName: [%s] is not Found! for Version: [%s]", plugin_name, plugin_version)
				util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
			}
			return
		}
		global.Config.DB.DB().
			Raw(" SELECT  * FROM  detail  where detail.version = ? and detail.plugin_name like ? ", plugin_version, "%"+plugin_name+"%").
			Scan(&plugin_Detail)

	}
	if len(plugin_Detail) == 0 {
		util.ReturnMsg(ctx, http.StatusNotFound, "", "The Plugin is not Found!")
		return
	} else {
		util.ReturnMsg(ctx, http.StatusOK, plugin_Detail, "success")
	}
}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary SearchPluginAllVersion
// @Schemes
// @Description Search Plugin ALL Version
// @Tags SearchPluginAllVersion
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param plugin_name query string false "Plugin Name"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router /plugin/searchall [get]
// @ID SearchPluginAllVersion
func SearchPluginAllVersion(ctx *gin.Context) {
	plugin_name := ctx.Query("plugin_name")
	//result := make(map[string]interface{})
	result := make(map[string][]*model.Detail)
	ok_msg := ""
	q := query.Use(global.Config.DB.DB())

	if plugin_name == "" {
		plugin_list, err := q.Plugin.WithContext(ctx).Select(q.Plugin.PluginName).Find()
		if err != nil {
			error_msg := fmt.Sprintf("ListPluginAllVersion err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		if len(plugin_list) != 0 {
			for _, v := range plugin_list {
				plugin_Detail, err := q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Eq(v.PluginName)).Find()
				if err != nil {
					error_msg := fmt.Sprintf("List Plugin Detail err: %s ", err)
					util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
					return
				}
				result[v.PluginName] = plugin_Detail

			}

		}
		ok_msg = "The ListPluginAllVersion successfully!"

	} else {
		plugin_Detail, err := q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Like("%" + plugin_name + "%")).Find()
		if err != nil {
			error_msg := fmt.Sprintf("List Plugin Detail err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		result[plugin_name] = plugin_Detail
		ok_msg = fmt.Sprintf("The PluginName: [%s] ListALlVersion successfully!", plugin_name)

	}

	util.ReturnMsg(ctx, http.StatusOK, result, ok_msg)

}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary CreatePlugin
// @Schemes
// @Description Create Plugin
// @Tags CreatePlugin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body index.Plugin true "Create Plugin"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /plugin/create [post]
// @ID CreatePlugin
func CreatePlugin(ctx *gin.Context) {
	body := index.Plugin{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}
	q := query.Use(global.Config.DB.DB())
	//查看plugin是否已存在
	_, err = q.Plugin.WithContext(ctx).Where(q.Plugin.PluginName.Eq(body.Name)).First()
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The PluginName: [%s] already exists!", body.Name)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

		return
	}

	//生成pluginId
	id, err := util.NewIdMgr(5)
	if err != nil {
		fmt.Println(err)
		return
	}

	plugin := model.Plugin{}
	plugin_id := util.Int64toString(id.ID())
	plugin.PluginName = body.Name
	plugin.PluginID = plugin_id
	plugin.LatestVersion = body.Spec.Version

	plugin_Detail := model.Detail{}
	plugin_Detail.PluginID = plugin_id
	plugin_Detail.PluginName = body.Name
	plugin_Detail.Version = body.Spec.Version
	plugin_Detail.Homepage = body.Spec.Homepage
	plugin_Detail.ShortDescription = body.Spec.ShortDescription
	plugin_Detail.Description = body.Spec.Description
	plugin_Detail.Caveats = body.Spec.Caveats
	var platform_list []string

	for _, v := range body.Spec.Platforms {
		data, err := jsoniter.Marshal(v)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
			return
		}
		platform_list = append(platform_list, string(data))

	}
	plugin_Detail.Platforms = platform_list

	// for _, v := range body.Spec.Platforms {
	// 	if _, ok := v.Selector.MatchLabels["os"]; ok {
	// 		plugin_Detail.PlatformOs = v.Selector.MatchLabels["os"]
	// 	}
	// 	if _, ok := v.Selector.MatchLabels["arch"]; ok {
	// 		plugin_Detail.PlatformArch = v.Selector.MatchLabels["arch"]
	// 	}
	// 	plugin_Detail.Sha256 = v.Sha256
	// 	plugin_Detail.URI = v.URI
	// 	plugin_Detail.Bin = v.Bin

	// }

	//不存在，直接create
	q.Transaction(func(tx *query.Query) error {
		if err := tx.Plugin.WithContext(ctx).Create(&plugin); err != nil {
			return err
		}
		if err := tx.Detail.WithContext(ctx).Create(&plugin_Detail); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		msg := fmt.Sprintf("The PluginName: [%s] Create Failed!, err: %s", body.Name, err)
		util.ReturnMsg(ctx, http.StatusNotFound, "", msg)
	} else {
		ok_msg := fmt.Sprintf("The PluginName: [%s] Create successfully!", body.Name)
		util.ReturnMsg(ctx, http.StatusOK, plugin, ok_msg)

	}
}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary UpdatePlugin
// @Schemes
// @Description Update Plugin
// @Tags UpdatePlugin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param ReqeustBody body index.Plugin true "Update Plugin"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /plugin/update [put]
// @ID UpdatePlugin
func UpdatePlugin(ctx *gin.Context) {
	body := index.Plugin{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		error_msg := fmt.Sprintf("err: %s ", err)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}
	q := query.Use(global.Config.DB.DB())
	//查看plugin是否存在
	plugin_data, err := q.Plugin.WithContext(ctx).Where(q.Plugin.PluginName.Eq(body.Name)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The PluginName: [%s] is not exists!", body.Name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)

		return
	}
	//plugin := model.Plugin{}
	//plugin.PluginName = body.Name

	plugin_Detail := model.Detail{}
	plugin_Detail.PluginName = body.Name
	plugin_Detail.Version = body.Spec.Version
	plugin_Detail.Homepage = body.Spec.Homepage
	plugin_Detail.ShortDescription = body.Spec.ShortDescription
	plugin_Detail.Description = body.Spec.Description
	plugin_Detail.Caveats = body.Spec.Caveats
	var platform_list []string

	for _, v := range body.Spec.Platforms {
		data, err := jsoniter.Marshal(v)
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
			return
		}
		platform_list = append(platform_list, string(data))

	}
	plugin_Detail.Platforms = platform_list
	//若存在，判断是否是同一版本内容
	info := gen.ResultInfo{}
	detail_data, err := q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Eq(body.Name), q.Detail.Version.Eq(body.Spec.Version)).First()

	//如果是新版本，直接create
	if errors.Is(err, gorm.ErrRecordNotFound) {
		plugin_Detail.PluginID = plugin_data.PluginID

		q.Transaction(func(tx *query.Query) error {
			if err = tx.Detail.WithContext(ctx).Create(&plugin_Detail); err != nil {
				return err
			}
			//对比当前创建版本是不是最新版本
			res := util.CompareStrVer(plugin_data.LatestVersion, body.Spec.Version)
			if res == 2 {
				//更新plugin表，更新lastversion
				//if info, err = tx.Plugin.WithContext(ctx).Where(tx.Plugin.PluginID.Eq(plugin_data.PluginID)).Updates(&plugin); err != nil {
				if info, err = tx.Plugin.WithContext(ctx).Where(tx.Plugin.PluginID.Eq(plugin_data.PluginID)).Update(tx.Plugin.LatestVersion, body.Spec.Version); err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			error_msg := fmt.Sprintf("Create Plugin_Detail err: %s ", err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

			return
		}

	} else {
		q.Transaction(func(tx *query.Query) error {
			//更新plugin表，更新lastversion
			// if info, err = tx.Plugin.WithContext(ctx).Where(tx.Plugin.PluginName.Eq(body.Name)).Updates(&plugin); err != nil {
			// 	return err
			// }
			if info, err = tx.Detail.WithContext(ctx).Where(tx.Detail.Version.Eq(detail_data.Version)).Updates(&plugin_Detail); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			error_msg := fmt.Sprintf("Update Plugin_Detail err: %s ", err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
			return
		}

	}
	ok_msg := fmt.Sprintf("The PluginName: [%s] Update successfully!", body.Name)

	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), ok_msg)
}

// @BasePath /api/v1
// PingPlugin godoc
// @Summary DeletePlugin
// @Schemes
// @Description Delete Plugin
// @Tags DeletePlugin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param plugin_name query string true "Plugin_Name"
// @Param version query string false "Plugin Version"
// @Success 200 {object} util.Res  {"code":200,"data":null,"msg":""}
// @Success 400 {object} util.Res  {"code":400,"data":null,"msg":""}
// @Success 404 {object} util.Res  {"code":404,"data":null,"msg":""}
// @Success 500 {object} util.Res  {"code":500,"data":null,"msg":""}
// @Router  /plugin/delete [delete]
// @ID DeletePlugin
func DeletePlugin(ctx *gin.Context) {
	plugin_name := ctx.Query("plugin_name")
	plugin_version := ctx.Query("version")
	detail_data := []*model.Detail{}
	q := query.Use(global.Config.DB.DB())

	plugin_data, err := q.Plugin.WithContext(ctx).Where(q.Plugin.PluginName.Eq(plugin_name)).First()
	info := gen.ResultInfo{}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		error_msg := fmt.Sprintf("The Plugin_Name: [%s] is not Found!", plugin_name)
		util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
		return

	} else {
		if err != nil {
			error_msg := fmt.Sprintf("err: %s ", err)
			util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
			return
		}
		if plugin_version != "" {
			_, err := q.Detail.WithContext(ctx).Select(q.Detail.PluginID).Where(q.Detail.PluginName.Eq(plugin_name), q.Detail.Version.Eq(plugin_version)).First()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				error_msg := fmt.Sprintf("The Plugin_Name: [%s] is not Found! for Version: [%s]", plugin_name, plugin_version)
				util.ReturnMsg(ctx, http.StatusNotFound, "", error_msg)
				return
			}

			if info, err = q.Detail.WithContext(ctx).Where(q.Detail.PluginName.Eq(plugin_name), q.Detail.Version.Eq(plugin_version)).Delete(); err != nil {
				error_msg := fmt.Sprintf("The PluginName: [%s] is deleted Failed! err: %s ", plugin_name, err)
				util.ReturnMsg(ctx, http.StatusInternalServerError, "", error_msg)
				return
			}
			q.Transaction(func(tx *query.Query) error {
				detail_data, err = q.Detail.WithContext(ctx).Select(tx.Detail.Version).Where(tx.Detail.PluginName.Eq(plugin_name)).Find()
				if len(detail_data) == 0 {
					//若该插件就一个版本，那么删除后，同时删除plugin表内容
					if info, err = tx.Plugin.WithContext(ctx).Where(tx.Plugin.PluginName.Eq(plugin_name)).Delete(); err != nil {
						return err
					}

				} else {
					//判断删除的是否是最新版本，若是，则在剩余版本中挑选最新版本，更新plugin表lastversion字段
					var version_list []string
					for _, v := range detail_data {
						version_list = append(version_list, v.Version)
					}
					sort.Sort(sort.Reverse(sort.StringSlice(version_list)))
					lastversion := version_list[0]
					fmt.Println(lastversion)

					if lastversion != plugin_data.LatestVersion {
						if info, err = tx.Plugin.WithContext(ctx).Where(tx.Plugin.PluginName.Eq(plugin_name)).Update(tx.Plugin.LatestVersion, lastversion); err != nil {
							return err
						}
					}

				}

				return nil
			})
		}
		if err != nil {
			error_msg := fmt.Sprintf("The PluginName: [%s] is deleted Failed! for Version : [%s], err: %s", plugin_name, plugin_version, err)
			util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)

			return
		}
	}
	if info.RowsAffected == 0 {
		error_msg := fmt.Sprintf("The PluginName: [%s] is deleted Failed!", plugin_name)
		util.ReturnMsg(ctx, http.StatusBadRequest, "", error_msg)
		return
	}

	util.ReturnMsg(ctx, http.StatusOK, util.Int64toString(info.RowsAffected), "Cluster deleted successfully!")
}
