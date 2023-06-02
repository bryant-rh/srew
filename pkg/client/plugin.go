package client

import (
	"github.com/bryant-rh/srew/internal/model"
	"github.com/bryant-rh/srew/pkg/index"
)

type pluginDetailData struct {
	GlobalRes
	Data []model.Detail `json:"data"`
}

//ListPlugin /plugin/list [get]
func (c *SrewClient) ListPlugin(plugin_name, plugin_version string) (*pluginDetailData, error) {
	res := &pluginDetailData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"plugin_name": plugin_name,
			"version":     plugin_version,
		}).
		SetSuccessResult(res).Get("/plugin/list")
	return res, err
}

//SearchPlugin /plugin/search [get]
func (c *SrewClient) SearchPlugin(plugin_name, plugin_version string) (*pluginDetailData, error) {
	res := &pluginDetailData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"plugin_name": plugin_name,
			"version":     plugin_version,
		}).
		SetSuccessResult(res).Get("/plugin/search")
	return res, err
}

type pluginAllData struct {
	GlobalRes
	Data map[string][]*model.Detail `json:"data"`
}

//ListPluginAllVersion /plugin/listall [get]
func (c *SrewClient) ListPluginAllVersion(plugin_name string) (*pluginAllData, error) {
	res := &pluginAllData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"plugin_name": plugin_name,
		}).
		SetSuccessResult(res).Get("/plugin/listall")
	return res, err
}

//SearchPluginAllVersion /plugin/searchall [get]
func (c *SrewClient) SearchPluginAllVersion(plugin_name string) (*pluginAllData, error) {
	res := &pluginAllData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"plugin_name": plugin_name,
		}).
		SetSuccessResult(res).Get("/plugin/searchall")
	return res, err
}

type pluginCreateData struct {
	GlobalRes
	Data model.Detail `json:"data"`
}

//CreatePlugin /plugin/create
func (c *SrewClient) CreatePlugin(plugin index.Plugin) (*pluginCreateData, error) {
	res := &pluginCreateData{}

	_, err := c.R().
		SetBody(plugin).
		SetSuccessResult(res).Post("/plugin/create")
	return res, err
}

type pluginData struct {
	GlobalRes
	Data string `json:"data"`
}

//UpdatePlugin /plugin/update
func (c *SrewClient) UpdatePlugin(plugin index.Plugin) (*pluginData, error) {
	res := &pluginData{}

	_, err := c.R().
		SetBody(plugin).
		SetSuccessResult(res).Put("/plugin/update")
	return res, err
}

//DeletePlugin /plugin/delete
func (c *SrewClient) DeletePlugin(project_name, plugin_version string) (*pluginData, error) {
	res := &pluginData{}

	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"project_name": project_name,
			"version":      plugin_version,
		}).
		SetSuccessResult(res).Delete("/plugin/delete")
	return res, err
}
