package mapimpl

import (
	"campaigns"
	"sort"
)

type Price struct {
	Price    float32
	Campaign string
}

type Campaigns []*Price

func (items Campaigns) Len() int {
	return len(items)
}

func (items Campaigns) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items Campaigns) Less(i, j int) bool {
	return items[i].Price > items[j].Price
}

type Devices map[string]Campaigns
type Countries map[string]Devices

func makeDevicesMap(data map[string]Countries, v string) {
	if data[v] == nil {
		data[v] = make(map[string]Devices)
	}
}

func makeCampaignsMap(data map[string]Devices, v string) {
	if data[v] == nil {
		data[v] = make(map[string]Campaigns)
	}
}

func convertData(c *data.Campaigns) (map[string]*data.Content, map[string]Countries, map[string]Devices) {
	contents := make(map[string]*data.Content)
	locals := make(map[string]Countries)
	globals := make(map[string]Devices)
	for k, v := range c.Elements {
		contents[k] = v.Content
		for p := range v.Placements.Elements {
			if len(v.Countries.Elements) == 0 {
				makeCampaignsMap(globals, p)
				for d := range v.Devices.Elements {
					globals[p][d] = append(globals[p][d], &Price{
						Price:    v.Price,
						Campaign: k,
					})
				}
				continue
			}
			makeDevicesMap(locals, p)
			for c := range v.Countries.Elements {
				makeCampaignsMap(locals[p], c)
				for d := range v.Devices.Elements {
					locals[p][c][d] = append(locals[p][c][d], &Price{
						Price:    v.Price,
						Campaign: k,
					})
				}
			}
		}
	}
	for _, v := range locals {
		for _, v1 := range v {
			for _, v2 := range v1 {
				sort.Sort(v2)
			}
		}
	}
	for _, v := range globals {
		for _, v1 := range v {
			sort.Sort(v1)
		}
	}
	return contents, locals, globals
}

type Controller struct {
	contents map[string]*data.Content
	locals   map[string]Countries
	globals  map[string]Devices
}

func MakeController(data *data.Campaigns) *Controller {
	contents, locals, globals := convertData(data)
	controller := &Controller{
		contents: contents,
		locals:   locals,
		globals:  globals,
	}
	return controller
}

func (c *Controller) makePriceResult(price *Price) (string, *data.Content) {
	return price.Campaign, c.contents[price.Campaign]
}

func (c *Controller) makeResult(prices []*Price) (string, *data.Content) {
	if len(prices) != 0 {
		return c.makePriceResult(prices[0])
	}
	return "", nil
}

func (c *Controller) FindCampaign(placement, country, device string) (string, *data.Content) {
	global, ok := c.globals[placement][device]
	if country == "" {
		return c.makeResult(global)
	}
	local, ok2 := c.locals[placement][country][device]
	if !ok2 {
		return c.makeResult(global)
	}
	if !ok {
		return c.makeResult(local)
	}
	if local[0].Price >= global[0].Price {
		return c.makePriceResult(local[0])
	}
	return c.makePriceResult(global[0])
}
