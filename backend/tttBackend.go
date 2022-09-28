package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type TTTOpenGames struct {
	Name string `json:"name"`
}

type TTTGame struct {
	Name          string   `json:"name"`
	State         string   `json:"state"`
	Board         []string `json:"board"`
	PlayerToken   string   `json:"player_token"`
	PlayerRole    string   `json:"player_role"`
	NextMoveToken string   `json:"next_move_token"`
}

const baseUrl = "https://tik-tak-tioki.fly.dev/api"

type ticTacTiokiClient struct {
	httpClient *http.Client
}

func (tttClient ticTacTiokiClient) OpenGames() ([]TTTOpenGames, error) {
	r, err := tttClient.httpClient.Get(baseUrl + "/join")
	if err != nil {
		return []TTTOpenGames{}, err
	}
	defer r.Body.Close()

	openGames := []TTTOpenGames{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return []TTTOpenGames{}, err
	}
	err = json.Unmarshal(bodyBytes, &openGames)
	if err != nil {
		return []TTTOpenGames{}, err
	}
	return openGames, nil
}

func (tttClient ticTacTiokiClient) JoinGame(gameName string) (TTTGame, error) {
	values := map[string]string{"name": gameName}
	jsonData, err := json.Marshal(values)
	r, err := tttClient.httpClient.Post(
		baseUrl+"/join",
		"application/json", bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return TTTGame{}, err
	}
	defer r.Body.Close()

	game := TTTGame{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return TTTGame{}, err
	}
	err = json.Unmarshal(bodyBytes, &game)
	if err != nil {
		return TTTGame{}, err
	}
	return game, nil
}

func (tttClient ticTacTiokiClient) NewGame() (TTTGame, error) {
	values := map[string]string{}
	jsonData, err := json.Marshal(values)
	r, err := tttClient.httpClient.Post(
		baseUrl+"/game",
		"application/json", bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return TTTGame{}, err
	}
	defer r.Body.Close()

	game := TTTGame{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return TTTGame{}, err
	}
	err = json.Unmarshal(bodyBytes, &game)
	if err != nil {
		return TTTGame{}, err
	}
	return game, nil
}

func (tttClient ticTacTiokiClient) Move(playerToken string, nextMoveToken string, field int) (TTTGame, error) {
	values := map[string]string{"next_move_token": nextMoveToken, "field": strconv.Itoa(field)}
	jsonData, err := json.Marshal(values)
	r, err := tttClient.httpClient.Post(
		baseUrl+"/move?player_token="+playerToken,
		"application/json", bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return TTTGame{}, err
	}
	defer r.Body.Close()

	game := TTTGame{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return TTTGame{}, err
	}
	err = json.Unmarshal(bodyBytes, &game)
	if err != nil {
		return TTTGame{}, err
	}
	return game, nil
}

func (tttClient ticTacTiokiClient) GameStatus(playerToken string) (TTTGame, error) {
	r, err := tttClient.httpClient.Get(baseUrl + "/game?player_token=" + playerToken)
	if err != nil {
		return TTTGame{}, err
	}
	defer r.Body.Close()

	game := TTTGame{}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return TTTGame{}, err
	}
	err = json.Unmarshal(bodyBytes, &game)
	if err != nil {
		return TTTGame{}, err
	}
	return game, nil
}
