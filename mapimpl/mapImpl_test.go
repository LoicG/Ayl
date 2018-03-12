package mapimpl

import (
	"campaigns"
	"encoding/json"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

const (
	campaignsfile = `
	{
		"campaigns":{
			"A":{
				"price":2.5,
				"content":{
					"title":"A title",
					"description":"A description",
					"landing":"A.fr/"
				},
				"countries":[
					"FRA",
					"BEL"
				],
				"devices":[
					"DESKTOP"
				],
				"placements":[
					"f64551fcd6f07823cb87971cfb914464",
					"3946ca64ff78d93ca61090a437cbb6b3",
					"3946ca64ff78d93ca610vva437cbb6b3"
				]
			},
			"B":{
				"price":1.3,
				"content":{
					"title":"B title",
					"description":"B description",
					"landing":"B.fr/"
				},
				"countries":[
					"FRA",
					"ALL"
				],
				"devices":[
					"DESKTOP",
					"MOBILE"
				],
				"placements":[
					"3946ca64ff78d93ca61090a437cbb6b3",
					"3946ca64ff78d93ca610vva437cbb6b3"
				]
			},
			"C":{
				"price":4.3,
				"content":{
					"title":"C title",
					"description":"C description",
					"landing":"C.fr/"
				},
				"devices":[
					"MOBILE"
				],
				"placements":[
					"3946ca64ff78d93ca61090a437cbb6b3"
				]
			}
		}
	}
`
)

func check(c *C, controller *Controller, placement, country, device string,
	expectedID string, expectedContent *data.Content) {

	id, content := controller.FindCampaign(placement, country, device)
	c.Assert(id, Equals, expectedID)
	c.Assert(content, DeepEquals, expectedContent)
}

func (s *TestSuite) TestMapImpl(c *C) {

	datas := &data.Campaigns{}
	decoder := json.NewDecoder(strings.NewReader(campaignsfile))
	err := decoder.Decode(datas)
	c.Assert(err, IsNil)

	controller := MakeController(datas)

	// invalid parameters
	check(c, controller, "", "", "", "", nil)
	check(c, controller, "invalid", "", "", "", nil)
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "invalid", "", "", nil)
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "", "invalid", "", nil)

	// simple placement
	check(c, controller, "f64551fcd6f07823cb87971cfb914464", "FRA", "MOBILE", "", nil)
	check(c, controller, "f64551fcd6f07823cb87971cfb914464", "FRA", "DESKTOP", "A", &data.Content{
		Title:       "A title",
		Description: "A description",
		Landing:     "A.fr/",
	})

	// no country
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "", "DESKTOP", "", nil)
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "", "MOBILE", "C", &data.Content{
		Title:       "C title",
		Description: "C description",
		Landing:     "C.fr/",
	})

	// best choice
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "FRA", "MOBILE", "C", &data.Content{
		Title:       "C title",
		Description: "C description",
		Landing:     "C.fr/",
	})
	check(c, controller, "3946ca64ff78d93ca61090a437cbb6b3", "FRA", "DESKTOP", "A", &data.Content{
		Title:       "A title",
		Description: "A description",
		Landing:     "A.fr/",
	})

	check(c, controller, "3946ca64ff78d93ca610vva437cbb6b3", "ALL", "DESKTOP", "B", &data.Content{
		Title:       "B title",
		Description: "B description",
		Landing:     "B.fr/",
	})
}
