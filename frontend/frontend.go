package frontend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/ioki-mobility/TicTxtToe/backend"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Backend interface {
	GameRequest(string) (backend.Game, error)
	MoveRequest(string, int) (backend.Game, error)
	Register(func(backend.Game))
}

type Frontend struct {
	bk       Backend
	tw       *twilio.RestClient
	svcPhone string
	log      hclog.Logger
}

func New(tw *twilio.RestClient, bk Backend, ph string, lg hclog.Logger) *Frontend {
	fr := &Frontend{
		bk:       bk,
		tw:       tw,
		log:      lg,
		svcPhone: ph,
	}
	fr.bk.Register(fr.gameReceiver)
	return fr
}

func (fr *Frontend) gameReceiver(g backend.Game) {
	fr.log.Info("backend responded with", "game", g)

	to := g.PhoneNumber

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(fr.svcPhone)
	params.SetBody(GameStateString(g))
	resp, err := fr.tw.Api.CreateMessage(params)
	if err != nil {
		fr.log.Error("failed sending sms", "error", err)
	} else {
		response, _ := json.Marshal(*resp)
		fr.log.Info("twilio response", "resp", string(response))
	}
}

func (fr *Frontend) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		fr.log.Error("cannot read response body", "error", err)
		return
	}

	s := strings.TrimSpace(string(b))
	fr.log.Info("http body", s)

	vals, err := url.ParseQuery(s)
	if err != nil {
		fr.log.Error("unexpected sms data")
	}

	msg := strings.ToLower(strings.TrimSpace(vals["Body"][0]))
	phone := strings.ToLower(strings.TrimSpace(vals["From"][0]))

	fr.log.Info("extracted values", "msg", msg, "phone", phone)

	if msg == "join" {
		g, err := fr.bk.GameRequest(phone)
		if err != nil {
			fr.log.Error("server error", "error", err)
			return
		}

		fr.log.Info("backend responded with", "game", g)
		fmt.Println(GameStateString(g))

	} else if strings.HasPrefix(msg, "move") {
		parts := strings.Split(msg, " ")
		if len(parts) != 2 {
			fr.log.Error("to user: malformed input")
			return
		}
		i, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fr.log.Error("to user: move is not a number")
			return
		}
		g, err := fr.bk.MoveRequest(phone, int(i))
		if err != nil {
			fr.log.Error("server error", "error", err)
		}

		fr.log.Info("backend responded with", "game", g)
		fmt.Println(GameStateString(g))
	} else {
		fr.log.Error("to user: malformed input")
	}
}

func GameStateString(g backend.Game) string {
	bl := strings.Builder{}
	for i, v := range g.Board {
		if i%3 == 0 {
			bl.WriteString("\n")
		}
		bl.WriteString(v + " ")
	}
	return fmt.Sprintf("State: %s\nYou: %s\n%s\n", g.State, g.Symbol, bl.String())
}
