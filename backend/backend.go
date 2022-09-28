package backend

import (
	"net/http"
	"time"
)

type Backend struct {
	players   map[string]string
	listener  func(Game)
	tttClient ticTacTiokiClient
}

func NewBackend() *Backend {
	return &Backend{
		players:   map[string]string{},
		listener:  func(g Game) {},
		tttClient: ticTacTiokiClient{httpClient: &http.Client{}},
	}
}

type State string

const (
	AwaitingJoin State = "awaiting_join"
	YourTurn     State = "your_turn"
	TheirTurn    State = "their_turn"
	YouWon       State = "you_won"
	TheyWon      State = "they_won"
	Draw         State = "draw"
)

type Game struct {
	PhoneNumber string
	Name        string
	Board       []string
	State       State
	Symbol      string
}

func (b *Backend) GameRequest(phoneNumber string) (Game, error) {
	openGames, err := b.tttClient.OpenGames()
	if err != nil {
		return Game{}, err
	}
	if len(openGames) > 0 {
		firstAvailableGame := openGames[0]
		tttGame, err := b.tttClient.JoinGame(firstAvailableGame.Name)
		if err != nil {
			return Game{}, err
		}
		b.players[phoneNumber] = tttGame.PlayerToken
		go b.maybePoll(phoneNumber, State(tttGame.State))
		return Game{
			PhoneNumber: phoneNumber,
			Name:        tttGame.Name,
			Board:       tttGame.Board,
			State:       State(tttGame.State),
			Symbol:      tttGame.PlayerRole,
		}, nil
	}

	tttGame, err := b.tttClient.NewGame()
	if err != nil {
		return Game{}, err
	}
	b.players[phoneNumber] = tttGame.PlayerToken
	go b.maybePoll(phoneNumber, State(tttGame.State))
	return Game{
		PhoneNumber: phoneNumber,
		Name:        tttGame.Name,
		Board:       tttGame.Board,
		State:       State(tttGame.State),
		Symbol:      tttGame.PlayerRole,
	}, nil
}

func (b *Backend) MoveRequest(phoneNumber string, fieldOnBoard int) (Game, error) {
	gameStatus, err := b.tttClient.GameStatus(b.players[phoneNumber])
	if err != nil {
		return Game{}, err
	}
	tttGame, err := b.tttClient.Move(b.players[phoneNumber], gameStatus.NextMoveToken, fieldOnBoard)
	if err != nil {
		return Game{}, err
	}
	go b.maybePoll(phoneNumber, State(tttGame.State))
	return Game{
		PhoneNumber: phoneNumber,
		Name:        tttGame.Name,
		Board:       tttGame.Board,
		State:       State(tttGame.State),
		Symbol:      tttGame.PlayerRole,
	}, nil
}

func (b *Backend) Register(listener func(Game)) {
	b.listener = listener
}

func (b *Backend) maybePoll(phoneNumber string, state State) {
	if state == AwaitingJoin || state == TheirTurn {
		for {
			tttGame, err := b.tttClient.GameStatus(b.players[phoneNumber])
			if err == nil && State(tttGame.State) != AwaitingJoin && State(tttGame.State) != TheirTurn {
				updatedgame := Game{
					PhoneNumber: phoneNumber,
					Name:        tttGame.Name,
					Board:       tttGame.Board,
					State:       State(tttGame.State),
					Symbol:      tttGame.PlayerRole,
				}
				b.listener(updatedgame)
				break
			}

			time.Sleep(5 * time.Second)
		}
	}
}
