package blockchain

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hsibAD/payment-service/internal/domain"
)

var (
	ErrInvalidWalletAddress = errors.New("invalid wallet address")
	ErrInvalidTransaction   = errors.New("invalid transaction")
	ErrTransactionFailed   = errors.New("transaction failed")
)

type MetaMaskProcessor struct {
	client         *ethclient.Client
	contractAddr   common.Address
	contractABI    string
	minConfirmations uint64
}

func NewMetaMaskProcessor(ethNodeURL, contractAddress, contractABI string, minConfirmations uint64) (*MetaMaskProcessor, error) {
	client, err := ethclient.Dial(ethNodeURL)
	if err != nil {
		return nil, err
	}

	return &MetaMaskProcessor{
		client:         client,
		contractAddr:   common.HexToAddress(contractAddress),
		contractABI:    contractABI,
		minConfirmations: minConfirmations,
	}, nil
}

func (p *MetaMaskProcessor) InitiateTransaction(ctx context.Context, payment *domain.Payment, walletAddress string) (*domain.MetaMaskInfo, error) {
	if !common.IsHexAddress(walletAddress) {
		return nil, ErrInvalidWalletAddress
	}

	// Convert payment amount to Wei
	amountWei := p.convertToWei(payment.Amount)

	// Create MetaMask info with transaction details
	info := &domain.MetaMaskInfo{
		WalletAddress:   walletAddress,
		ContractAddress: p.contractAddr.Hex(),
		AmountWei:      amountWei.String(),
	}

	return info, nil
}

func (p *MetaMaskProcessor) VerifyTransaction(ctx context.Context, payment *domain.Payment, transactionHash string) error {
	if len(transactionHash) != 66 || transactionHash[:2] != "0x" {
		return ErrInvalidTransaction
	}

	txHash := common.HexToHash(transactionHash)

	// Get transaction receipt
	receipt, err := p.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if err == ethereum.NotFound {
			return ErrInvalidTransaction
		}
		return err
	}

	// Check if transaction was successful
	if receipt.Status != types.ReceiptStatusSuccessful {
		return ErrTransactionFailed
	}

	// Get current block number
	currentBlock, err := p.client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	// Check confirmations
	confirmations := currentBlock - receipt.BlockNumber.Uint64()
	if confirmations < p.minConfirmations {
		return errors.New("insufficient confirmations")
	}

	return nil
}

func (p *MetaMaskProcessor) GetTransactionStatus(ctx context.Context, transactionHash string) (string, error) {
	if len(transactionHash) != 66 || transactionHash[:2] != "0x" {
		return "", ErrInvalidTransaction
	}

	txHash := common.HexToHash(transactionHash)

	// Get transaction receipt
	receipt, err := p.client.TransactionReceipt(ctx, txHash)
	if err != nil {
		if err == ethereum.NotFound {
			return "PENDING", nil
		}
		return "", err
	}

	// Get current block number
	currentBlock, err := p.client.BlockNumber(ctx)
	if err != nil {
		return "", err
	}

	// Check confirmations
	confirmations := currentBlock - receipt.BlockNumber.Uint64()

	if receipt.Status != types.ReceiptStatusSuccessful {
		return "FAILED", nil
	}

	if confirmations < p.minConfirmations {
		return "CONFIRMING", nil
	}

	return "CONFIRMED", nil
}

func (p *MetaMaskProcessor) convertToWei(amount float64) *big.Int {
	// Convert amount to Wei (1 ETH = 10^18 Wei)
	amountStr := big.NewFloat(amount)
	multiplier := big.NewFloat(1e18)
	amountStr.Mul(amountStr, multiplier)

	amountWei := new(big.Int)
	amountStr.Int(amountWei)

	return amountWei
}

// Smart Contract Interface
const PaymentContractABI = `[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"internalType": "string",
				"name": "orderID",
				"type": "string"
			},
			{
				"indexed": true,
				"internalType": "address",
				"name": "payer",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "PaymentReceived",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "string",
				"name": "orderID",
				"type": "string"
			}
		],
		"name": "makePayment",
		"outputs": [],
		"stateMutability": "payable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "newOwner",
				"type": "address"
			}
		],
		"name": "transferOwnership",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "withdraw",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]` 