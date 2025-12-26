package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// SupplierData — это полезная нагрузка (Payload).
// Здесь хранятся данные о конкретном событии с поставщиком.
type SupplierData struct {
	INN          string  `json:"inn"`           // ИНН компании (уникальный идентификатор)
	CompanyName  string  `json:"company_name"`  // Название юрлица
	Action       string  `json:"action"`        // Событие: "REGISTER", "CONTRACT_FAIL", "RATING_UPDATE"
	RatingChange float64 `json:"rating_change"` // Изменение рейтинга (например, -0.5 или +1.0)
	Details      string  `json:"details"`       // Комментарий (например, номер госконтракта)
}

// Block — основная единица блокчейна.
type Block struct {
	Index     int          `json:"index"`     // Порядковый номер блока (высота)
	Timestamp int64        `json:"timestamp"` // Время создания (Unix timestamp)
	Data      SupplierData `json:"data"`      // Информация о поставщике
	PrevHash  string       `json:"prev_hash"` // Хеш предыдущего блока (связь цепочки)
	Hash      string       `json:"hash"`      // Хеш текущего блока (цифровой отпечаток)
	Nonce     int          `json:"nonce"`     // Случайное число (для имитации Proof-of-Work, если понадобится)
}

// NewBlock создает новый блок.
// Принимает данные, хеш предыдущего блока и текущий индекс.
func NewBlock(data SupplierData, prevHash string, index int) *Block {
	block := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Data:      data,
		PrevHash:  prevHash,
		Nonce:     0,
	}

	// Сразу вычисляем хеш для этого блока
	block.Hash = block.CalculateHash()
	return block
}

// NewGenesisBlock создает самый первый блок в цепочке.
// У него нет предыдущего хеша (он равен "0").
func NewGenesisBlock() *Block {
	genesisData := SupplierData{
		INN:         "0000000000",
		CompanyName: "System Genesis",
		Action:      "INIT",
		Details:     "Genesis Block - начало цепочки",
	}
	return NewBlock(genesisData, "0", 0)
}

// CalculateHash создает SHA-256 хеш блока на основе его содержимого.
// Если изменить хоть одну букву в Data, хеш полностью изменится.
func (b *Block) CalculateHash() string {
	// Преобразуем данные поставщика в JSON строку для хеширования
	dataBytes, _ := json.Marshal(b.Data)

	// Собираем все поля блока в одну строку
	record := fmt.Sprintf("%d%d%s%s%d", b.Index, b.Timestamp, string(dataBytes), b.PrevHash, b.Nonce)

	// Вычисляем SHA-256
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)

	return hex.EncodeToString(hashed)
}
