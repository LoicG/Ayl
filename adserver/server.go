package main

import (
	"campaigns"
	"fmt"
	"mapimpl"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

type Finder interface {
	FindCampaign(string, string, string) (string, *data.Content)
}

type Server struct {
	ws     *restful.WebService
	port   int
	finder Finder
}

func newServer(port int, campaigns string) (*Server, error) {
	data, err := data.LoadJSONFile(campaigns)
	if err != nil {
		return nil, err
	}
	controller := mapimpl.MakeController(data)
	server := &Server{
		ws:     &restful.WebService{},
		port:   port,
		finder: controller,
	}
	server.AttachRoute()
	return server, nil
}

type RestHandler func(*restful.Request) (interface{}, error)

func wrap(handler RestHandler) restful.RouteFunction {
	return func(request *restful.Request, response *restful.Response) {
		result, err := handler(request)
		if err != nil {
			response.WriteError(http.StatusBadRequest, err)
			return
		}
		response.WriteEntity(result)
	}
}

type Input struct {
	Country string
	Device  string
}

type Output struct {
	Campaign string
	Content  *data.Content
}

func (s *Server) getPlacement(request *restful.Request) (interface{}, error) {
	placement := request.QueryParameter("placement")
	if placement == "" {
		return nil, fmt.Errorf("invalid placement parameter")
	}
	input := &Input{}
	err := request.ReadEntity(input)
	if err != nil {
		return nil, err
	}
	campaign, content := s.finder.FindCampaign(placement, input.Country, input.Device)
	if campaign != "" {
		return &Output{
			Campaign: campaign,
			Content:  content,
		}, nil
	}
	return nil, nil
}

func (s *Server) AttachRoute() {
	s.ws.
		Path("/ad").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	s.ws.
		Route(s.ws.POST("/").
			To(wrap(s.getPlacement)))
	restful.Add(s.ws)
}

func (s *Server) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(s.port), nil)
}
