package service

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func SolCallBack(address string) {}

const (
	// Mint token
	USDT_Mint = "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB"
	USDC_Mint = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
)

const (
	TokenProgramID     = "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	Token2022ProgramID = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"
)

func FindATAAddress(owner, mint string) (string, error) {
	ownerPubKey, err := solana.PublicKeyFromBase58(owner)
	if err != nil {
		return "", fmt.Errorf("invalid owner public key: %w", err)
	}

	mintPubKey, err := solana.PublicKeyFromBase58(mint)
	if err != nil {
		return "", fmt.Errorf("invalid mint public key: %w", err)
	}

	ata, _, err := solana.FindAssociatedTokenAddress(ownerPubKey, mintPubKey)
	if err != nil {
		return "", fmt.Errorf("find associated token address failed: %w", err)
	}

	return ata.String(), nil
}

const (
	InstructionTransfer        = 3
	InstructionTransferChecked = 12
)

type TransferInfo struct {
	ProgramID string
	Type      string

	Source      string
	Destination string
	Mint        string
	Authority   string

	Amount   uint64
	Decimals *uint8
}

func isTokenProgram(programID solana.PublicKey) bool {
	p := programID.String()
	return p == TokenProgramID || p == Token2022ProgramID
}

func DecodeTransfer(data []byte, accountKeys []solana.PublicKey, accountIndexes []uint16, programID solana.PublicKey) (*TransferInfo, error) {
	if len(data) < 9 {
		return nil, fmt.Errorf("invalid transfer data length: %d", len(data))
	}

	if data[0] != InstructionTransfer {
		return nil, fmt.Errorf("not transfer instruction, got: %d", data[0])
	}

	if len(accountIndexes) < 3 {
		return nil, fmt.Errorf("transfer accounts not enough: %d", len(accountIndexes))
	}

	amount := binary.LittleEndian.Uint64(data[1:9])

	sourceIdx := accountIndexes[0]
	destIdx := accountIndexes[1]
	authIdx := accountIndexes[2]

	if int(sourceIdx) >= len(accountKeys) || int(destIdx) >= len(accountKeys) || int(authIdx) >= len(accountKeys) {
		return nil, errors.New("account index out of range")
	}

	info := &TransferInfo{
		ProgramID:   programID.String(),
		Type:        "transfer",
		Source:      accountKeys[sourceIdx].String(),
		Destination: accountKeys[destIdx].String(),
		Authority:   accountKeys[authIdx].String(),
		Amount:      amount,
	}

	return info, nil
}

type TransferCheckedInfo struct {
	ProgramID string
	Type      string

	Source      string
	Destination string
	Mint        string
	Authority   string

	Amount   uint64
	Decimals uint8
}

func DecodeTransferCheck(data []byte, accountKeys []solana.PublicKey, accountIndexes []uint16, programID solana.PublicKey) (*TransferInfo, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("invalid transferChecked data length: %d", len(data))
	}

	if data[0] != InstructionTransferChecked {
		return nil, fmt.Errorf("not transferChecked instruction, got: %d", data[0])
	}

	if len(accountIndexes) < 4 {
		return nil, fmt.Errorf("transferChecked accounts not enough: %d", len(accountIndexes))
	}

	amount := binary.LittleEndian.Uint64(data[1:9])
	decimals := uint8(data[9])

	sourceIdx := accountIndexes[0]
	mintIdx := accountIndexes[1]
	destIdx := accountIndexes[2]
	authIdx := accountIndexes[3]

	if int(sourceIdx) >= len(accountKeys) ||
		int(mintIdx) >= len(accountKeys) ||
		int(destIdx) >= len(accountKeys) ||
		int(authIdx) >= len(accountKeys) {
		return nil, errors.New("account index out of range")
	}

	info := &TransferInfo{
		ProgramID:   programID.String(),
		Type:        "transferChecked",
		Source:      accountKeys[sourceIdx].String(),
		Mint:        accountKeys[mintIdx].String(),
		Destination: accountKeys[destIdx].String(),
		Authority:   accountKeys[authIdx].String(),
		Amount:      amount,
		Decimals:    &decimals,
	}

	return info, nil
}

// ParseTransactionTransfers
func ParseTransactionTransfers(ctx context.Context, client *rpc.Client, sig string) ([]*TransferInfo, error) {
	signature, err := solana.SignatureFromBase58(sig)
	if err != nil {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	tx, err := client.GetTransaction(
		ctx,
		signature,
		&rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: func(v uint64) *uint64 { return &v }(0),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get transaction failed: %w", err)
	}
	if tx == nil {
		return nil, errors.New("transaction not found")
	}

	decodedTx, err := tx.Transaction.GetTransaction()
	if err != nil {
		return nil, fmt.Errorf("decode transaction failed: %w", err)
	}

	msg := decodedTx.Message
	accountKeys := msg.AccountKeys
	results := make([]*TransferInfo, 0)

	for _, ix := range msg.Instructions {
		if int(ix.ProgramIDIndex) >= len(accountKeys) {
			continue
		}

		programID := accountKeys[ix.ProgramIDIndex]
		if !isTokenProgram(programID) {
			continue
		}

		data := ix.Data
		if len(data) == 0 {
			continue
		}

		switch data[0] {
		case InstructionTransfer:
			info, err := DecodeTransfer(data, accountKeys, toUint16Slice(ix.Accounts), programID)
			if err == nil {
				results = append(results, info)
			}
		case InstructionTransferChecked:
			info, err := DecodeTransferCheck(data, accountKeys, toUint16Slice(ix.Accounts), programID)
			if err == nil {
				results = append(results, info)
			}
		}
	}

	return results, nil
}

func toUint16Slice(in []uint16) []uint16 {
	out := make([]uint16, len(in))
	for i, v := range in {
		out[i] = uint16(v)
	}
	return out
}

func DecodeBase64InstructionData(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
