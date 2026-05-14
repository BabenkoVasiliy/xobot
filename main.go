package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GameState struct {
	Board   []string `json:"board"`
	Xmoves  []int    `json:"xMoves"`
	Omoves  []int    `json:"oMoves"`
	LastMove int     `json:"lastMove"`
}

var gameState = make(map[int64]*GameState)
var botToken = "8669066608:AAHQeaPBVCT_khKKTWSEZsDXWe7pAWjuoMo"
var webAppURL = "https://BabenkoVasiliy.github.io/"

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	http.HandleFunc("/move", handleMove)

	go func() {
		log.Println("HTTP server started on 0.0.0.0:8080")
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	kb := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Найти соперника"),
			tgbotapi.NewKeyboardButton("Играть с ботом"),
		),
	)

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			text := update.Message.Text

			if text == "Играть с ботом" {
				gameState[chatID] = &GameState{
					Board:  make([]string, 9),
					Xmoves: []int{},
					Omoves: []int{},
				}

				btn := tgbotapi.NewInlineKeyboardButtonURL("🎮 Играть", webAppURL)
				kbInline := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(btn))

				msg := tgbotapi.NewMessage(chatID, "Нажми кнопку ниже чтобы начать игру:")
				msg.ReplyMarkup = kbInline
				bot.Send(msg)
			} else if text == "Найти соперника" {
				msg := tgbotapi.NewMessage(chatID, "Поиск соперника временно недоступен")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(chatID, "Выбери режим игры:")
				msg.ReplyMarkup = kb
				bot.Send(msg)
			}
		}

		if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			bot.Request(callback)
		}
	}
}

type MoveRequest struct {
	ChatID   int64   `json:"chatId"`
	Board    []string `json:"board"`
	Xmoves   []int   `json:"xMoves"`
	Omoves   []int   `json:"oMoves"`
	LastMove int     `json:"lastMove"`
}

type MoveResponse struct {
	Index    int    `json:"index"`
	Remove   int    `json:"remove"`
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	move := calculateBotMove(req.Board, req.Xmoves, req.Omoves)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(move)
}

func calculateBotMove(board []string, xMoves, oMoves []int) MoveResponse {
	emptyCells := []int{}
	for i, v := range board {
		if v == "" {
			emptyCells = append(emptyCells, i)
		}
	}

	if len(emptyCells) == 0 {
		return MoveResponse{Index: -1}
	}

	bestMove := -1
	bestScore := -1000

	for _, idx := range emptyCells {
		board[idx] = "O"
		testMoves := append(oMoves, idx)
		score := minimax(board, testMoves, xMoves, 0, false)
		board[idx] = ""

		if score > bestScore {
			bestScore = score
			bestMove = idx
		}
	}

	removeIdx := -1
	if len(oMoves) >= 3 {
		removeIdx = oMoves[0]
	}

	return MoveResponse{Index: bestMove, Remove: removeIdx}
}

func minimax(board []string, oMoves, xMoves []int, depth int, isMaximizing bool) int {
	winner := checkWinnerBot(board, oMoves, xMoves)
	if winner == "O" {
		return 100 - depth
	}
	if winner == "X" {
		return depth - 100
	}
	if isDrawBot(board) {
		return 0
	}

	if isMaximizing {
		maxScore := -1000
		for i := 0; i < 9; i++ {
			if board[i] == "" {
				board[i] = "O"
				score := minimax(board, append(oMoves, i), xMoves, depth+1, false)
				board[i] = ""
				if score > maxScore {
					maxScore = score
				}
			}
		}
		return maxScore
	} else {
		minScore := 1000
		for i := 0; i < 9; i++ {
			if board[i] == "" {
				board[i] = "X"
				score := minimax(board, oMoves, append(xMoves, i), depth+1, true)
				board[i] = ""
				if score < minScore {
					minScore = score
				}
			}
		}
		return minScore
	}
}

func checkWinnerBot(board []string, oMoves, xMoves []int) string {
	lines := [][]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		{0, 4, 8}, {2, 4, 6},
	}

	for _, line := range lines {
		if board[line[0]] != "" && board[line[0]] == board[line[1]] && board[line[0]] == board[line[2]] {
			return board[line[0]]
		}
	}
	return ""
}

func isDrawBot(board []string) bool {
	for _, v := range board {
		if v == "" {
			return false
		}
	}
	return true
}

var _ = strconv.Atoi