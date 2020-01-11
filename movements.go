package main

import (
	"context"
	"fmt"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"net/http"
	"os"
	"time"
)

type Movements struct {
	file   *os.File
	blocks *client.Blocks
}

func NewMovements() (*Movements, error) {
	file, err := os.Create(fmt.Sprintf("./scan-%d-asset-%s-height-%d", time.Now().Unix(), *flagTrackMovementsAsset, *flagTrackMovementsHeightFrom))
	if err != nil {
		return nil, err
	}

	var defaultOptions = client.Options{
		BaseUrl: *flagNodeEndpoint,
		Client:  &http.Client{Timeout: 3 * time.Second},
	}

	blocks := client.NewBlocks(defaultOptions)

	return &Movements{
		file:   file,
		blocks: blocks,
	}, nil
}

func (m *Movements) Parse() error {
	defer m.file.Close()

	ctx := context.Background()
	asset, _ := proto.NewOptionalAssetFromString(*flagTrackMovementsAsset)
	height := uint64(0)
	count := *flagTrackMovementsHeightFrom

	current, _, err := m.blocks.Height(ctx)
	if err != nil {
		return err
	}
	height = current.Height

	for {

		if count < height {
			block, _, err := m.blocks.At(ctx, count+1)
			if err != nil {
				continue
			}

			for _, tx := range block.Transactions {
				_, err := tx.MarshalBinary()
				if err != nil {
					m.WriteError(count+1, err)
				}

				switch tx.(type) {
				case *proto.TransferV1:
					transfer := tx.(*proto.TransferV1)
					if transfer.AmountAsset.ID == asset.ID {
						var sender, recipient string

						s, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, tx.GetSenderPK())
						sender = s.String()
						if transfer.Recipient.Address != nil {
							recipient = transfer.Recipient.Address.String()
						} else if transfer.Recipient.Alias != nil {
							recipient = transfer.Recipient.Alias.String()
						}

						m.WriteAssetEntry(
							count+1,
							transfer.ID.String(),
							int(tx.GetTypeVersion().Type),
							transfer.Amount,
							sender,
							recipient,
						)
					}
					break
				case *proto.TransferV2:
					transfer := tx.(*proto.TransferV2)
					if transfer.AmountAsset.ID == asset.ID {
						var sender, recipient string

						s, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, tx.GetSenderPK())
						sender = s.String()
						if transfer.Recipient.Address != nil {
							recipient = transfer.Recipient.Address.String()
						} else if transfer.Recipient.Alias != nil {
							recipient = transfer.Recipient.Alias.String()
						}

						m.WriteAssetEntry(
							count+1,
							transfer.ID.String(),
							int(tx.GetTypeVersion().Type),
							transfer.Amount,
							sender,
							recipient,
						)
					}
					break
				case *proto.MassTransferV1:
					transfer := tx.(*proto.MassTransferV1)
					if transfer.Asset.ID == asset.ID {
						var sender, recipient string

						s, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, tx.GetSenderPK())
						sender = s.String()

						for _, entry := range transfer.Transfers {
							if entry.Recipient.Address != nil {
								recipient = entry.Recipient.Address.String()
							} else if entry.Recipient.Alias != nil {
								recipient = entry.Recipient.Alias.String()
							}

							m.WriteAssetEntry(
								count+1,
								transfer.ID.String(),
								int(tx.GetTypeVersion().Type),
								entry.Amount,
								sender,
								recipient,
							)
						}
					}
					break
				case *proto.ExchangeV1:
					exchange := tx.(*proto.ExchangeV1)
					if exchange.BuyOrder.AssetPair.AmountAsset.ID == asset.ID || exchange.BuyOrder.AssetPair.PriceAsset.ID == asset.ID ||
						exchange.SellOrder.AssetPair.AmountAsset.ID == asset.ID || exchange.SellOrder.AssetPair.PriceAsset.ID == asset.ID {
						var sender, recipient string

						s, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, exchange.BuyOrder.SenderPK)
						r, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, exchange.SellOrder.SenderPK)
						sender = s.String()
						recipient = r.String()

						m.WriteAssetEntry(
							count+1,
							exchange.ID.String(),
							int(tx.GetTypeVersion().Type),
							exchange.Amount,
							sender,
							recipient,
						)
					}
					break
				case *proto.ExchangeV2:
					exchange := tx.(*proto.ExchangeV2)
					if exchange.BuyOrder.GetAssetPair().AmountAsset.ID == asset.ID || exchange.BuyOrder.GetAssetPair().PriceAsset.ID == asset.ID ||
						exchange.SellOrder.GetAssetPair().AmountAsset.ID == asset.ID || exchange.SellOrder.GetAssetPair().PriceAsset.ID == asset.ID {
						var sender, recipient string

						s, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, exchange.BuyOrder.GetSenderPK())
						r, _ := proto.NewAddressFromPublicKey(proto.MainNetScheme, exchange.SellOrder.GetSenderPK())
						sender = s.String()
						recipient = r.String()

						m.WriteAssetEntry(
							count+1,
							exchange.ID.String(),
							int(tx.GetTypeVersion().Type),
							exchange.Amount,
							sender,
							recipient,
						)
					}
					break
				}
			}

			count = count + 1
		} else {
			return nil
		}
	}
}

func (m *Movements) WriteAssetEntry(height uint64, txHash string, txType int, amount uint64, sender, recipient string) {
	message := fmt.Sprintf("%-10d %-42s %-4d %-16d %-42s %-42s\n", height, txHash, txType, amount, sender, recipient)
	fmt.Print(message)
	_, _ = m.file.WriteString(message)
}

func (m *Movements) WriteError(height uint64, err error) {
	fmt.Println(height, err)
}
