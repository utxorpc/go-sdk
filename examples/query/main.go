package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/blinklabs-io/gouroboros/ledger/common"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	"github.com/utxorpc/go-codegen/utxorpc/v1alpha/query"
	utxorpc "github.com/utxorpc/go-sdk"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func main() {
	baseUrl := os.Getenv("UTXORPC_URL")
	if baseUrl == "" {
		baseUrl = "https://preview.utxorpc-v0.demeter.run"
	}
	client := utxorpc.NewClient(utxorpc.WithBaseUrl(baseUrl))
	dmtrApiKey := os.Getenv("DMTR_API_KEY")
	// set API key for demeter
	if dmtrApiKey != "" {
		client.SetHeader("dmtr-api-key", dmtrApiKey)
	}

	// Run them all
	readParams(client)
	readUtxo(
		client,
		"24efe5f12d1d93bb419cfb84338d6602dfe78c614b489edb72df0594a077431c",
		0,
	)
	// https://preprod.cexplorer.io/asset/asset1tvkt35str8aeepuflxmnjzcdj87em8xrlx4ehz
	// Use policy ID and asset name in hex format (https://cips.cardano.org/cip/CIP-68/)
	// Hunt
	searchUtxos(
		client,
		"addr_test1qptfy9zhaeuqfptcu79q6gm9l3r6cfp5gnlqc7m42qwln0lsvex239qmryg4yh3pda3rh3rnce4wd46gdyqlscrq7s4shekqrt",
		"63f9a5fc96d4f87026e97af4569975016b50eef092a46859b61898e5",
		"0014df1048554e54",
	)
	// Dedi
	searchUtxos(
		client,
		"addr_test1qptfy9zhaeuqfptcu79q6gm9l3r6cfp5gnlqc7m42qwln0lsvex239qmryg4yh3pda3rh3rnce4wd46gdyqlscrq7s4shekqrt",
		"63f9a5fc96d4f87026e97af4569975016b50eef092a46859b61898e5",
		"0014df1044454449",
	)
	// No assets
	searchUtxos(
		client,
		"addr_test1qzrkvcfvd7k5jx54xxkz87p8xn88304jd2g4jsa0hwwmg20k3c7k36lsg8rdupz6e36j5ctzs6lzjymc9vw7djrmgdnqff9z6j",
		"63f9a5fc96d4f87026e97af4569975016b50eef092a46859b61898e5",
		"0014df1044454449",
	)
	getUtxosByAddress(
		client,
		"addr_test1qptfy9zhaeuqfptcu79q6gm9l3r6cfp5gnlqc7m42qwln0lsvex239qmryg4yh3pda3rh3rnce4wd46gdyqlscrq7s4shekqrt",
	)
}

func readParams(client *utxorpc.UtxorpcClient) {
	fmt.Println("Connecting to utxorpc host:", client.URL())
	resp, err := client.ReadParams()
	if err != nil {
		utxorpc.HandleError(err)
	}
	fmt.Printf("Response: %+v\n", resp)

	if resp.Msg.GetLedgerTip() != nil {
		fmt.Printf(
			"Ledger Tip: Slot: %d, Hash: %x\n",
			resp.Msg.GetLedgerTip().GetSlot(),
			resp.Msg.GetLedgerTip().GetHash(),
		)
	}
	if resp.Msg.GetValues() != nil {
		fmt.Printf("Cardano: %+v\n", resp.Msg.GetValues())
	}
}

func readUtxo(
	client *utxorpc.UtxorpcClient,
	txHashStr string,
	txIndex uint32,
) {
	resp, err := client.ReadUtxo(txHashStr, txIndex)
	if err != nil {
		utxorpc.HandleError(err)
		return
	}

	// Process the response
	fmt.Printf("Response: %+v\n", resp)

	if resp.Msg.GetLedgerTip() != nil {
		fmt.Printf(
			"Ledger Tip:\n  Slot: %d\n  Hash: %x\n",
			resp.Msg.GetLedgerTip().GetSlot(),
			resp.Msg.GetLedgerTip().GetHash(),
		)
	}

	for _, item := range resp.Msg.GetItems() {
		fmt.Println("UTxO Data:")
		fmt.Printf("  Tx Hash: %x\n", item.GetTxoRef().GetHash())
		fmt.Printf("  Output Index: %d\n", item.GetTxoRef().GetIndex())
		fmt.Printf("  Native Bytes: %x\n", item.GetNativeBytes())
		if cardano := item.GetCardano(); cardano != nil {
			fmt.Println("  Cardano UTxO:")
			fmt.Printf("    Address: %x\n", cardano.GetAddress())
			fmt.Printf("    Coin: %d\n", cardano.GetCoin())
			if cardano.GetDatum() != nil {
				fmt.Printf("    Datum Hash: %x\n", cardano.GetDatum().GetHash())
			}
		}
	}
}

func searchUtxos(
	client *utxorpc.UtxorpcClient,
	rawAddress string,
	policyID string,
	assetName string,
) {
	// Use to support bech32/base58 addresses
	addr, err := common.NewAddress(rawAddress)
	if err != nil {
		log.Fatalf("failed to create address: %v", err)
	}
	addrCbor, err := addr.MarshalCBOR()
	if err != nil {
		log.Fatalf("failed to marshal address to CBOR: %v", err)
	}

	var txOutputPattern *cardano.TxOutputPattern
	if policyID != "" && assetName != "" {
		// Convert policyID from hex to bytes
		policyIDBytes, err := hex.DecodeString(policyID)
		if err != nil {
			log.Fatalf("failed to decode policy ID: %v", err)
		}

		// Convert assetName to bytes
		assetNameBytes, err := hex.DecodeString(assetName)
		if err != nil {
			log.Fatalf("failed to decode asset name: %v", err)
		}

		// Define the asset pattern with policy ID and asset name
		assetPattern := &cardano.AssetPattern{
			PolicyId:  policyIDBytes,
			AssetName: assetNameBytes,
		}

		// Define the TxOutput pattern including the asset filter
		txOutputPattern = &cardano.TxOutputPattern{
			Address: &cardano.AddressPattern{
				ExactAddress: addrCbor,
			},
			Asset: assetPattern,
		}
	} else {
		// Define the TxOutput pattern with only the address filter
		txOutputPattern = &cardano.TxOutputPattern{
			Address: &cardano.AddressPattern{
				ExactAddress: addrCbor,
			},
		}
	}

	// Wrap the TxOutput pattern in AnyUtxoPattern for Cardano
	anyUtxoPattern := &query.AnyUtxoPattern{
		UtxoPattern: &query.AnyUtxoPattern_Cardano{
			Cardano: txOutputPattern,
		},
	}

	// Define the UtxoPredicate with the pattern
	utxoPredicate := &query.UtxoPredicate{
		Match: anyUtxoPattern,
	}

	// Define the field mask
	fieldMask := &fieldmaskpb.FieldMask{
		Paths: []string{
			// "native_bytes",
		},
	}

	// Prepare the search request
	searchRequest := &query.SearchUtxosRequest{
		Predicate:  utxoPredicate,
		FieldMask:  fieldMask,
		MaxItems:   100, // Adjust based on your requirements
		StartToken: "",  // For pagination; empty for the first page
	}

	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.SearchUtxos(searchRequest)
	if err != nil {
		utxorpc.HandleError(err)
	}

	// Uncomment to print the full response for debugging
	// fmt.Printf("Response: %+v\n", resp)

	if resp.Msg.GetLedgerTip() != nil {
		fmt.Printf(
			"Ledger Tip:\n  Slot: %d\n  Hash: %x\n",
			resp.Msg.GetLedgerTip().GetSlot(),
			resp.Msg.GetLedgerTip().GetHash(),
		)
	}

	for _, item := range resp.Msg.GetItems() {
		fmt.Println("UTxO Data:")
		fmt.Printf("  Tx Hash: %x\n", item.GetTxoRef().GetHash())
		fmt.Printf("  Output Index: %d\n", item.GetTxoRef().GetIndex())
		fmt.Printf("  Native Bytes: %x\n", item.GetNativeBytes())
		if cardano := item.GetCardano(); cardano != nil {
			fmt.Println("  Cardano UTxO:")
			fmt.Printf("    Address: %x\n", cardano.GetAddress())
			fmt.Printf("    Coin: %d\n", cardano.GetCoin())
			fmt.Println("    Assets:")
			for _, multiasset := range cardano.GetAssets() {
				fmt.Printf("      Policy ID: %x\n", multiasset.GetPolicyId())
				for _, asset := range multiasset.GetAssets() {
					fmt.Printf("        Asset Name: %s\n", string(asset.GetName()))
					fmt.Printf("        Output Coin: %d\n", asset.GetOutputCoin())
					fmt.Printf("        Mint Coin: %d\n", asset.GetMintCoin())
				}
			}
		}
	}
}

func getUtxosByAddress(
	client *utxorpc.UtxorpcClient,
	rawAddress string,
) {
	// Use to support bech32/base58 addresses
	addr, err := common.NewAddress(rawAddress)
	if err != nil {
		log.Fatalf("failed to create address: %v", err)
	}
	addrCbor, err := addr.MarshalCBOR()
	if err != nil {
		log.Fatalf("failed to marshal address to CBOR: %v", err)
	}

	fmt.Println("connecting to utxorpc host:", client.URL())
	resp, err := client.GetUtxosByAddress(addrCbor)
	if err != nil {
		utxorpc.HandleError(err)
	}

	// Uncomment to print the full response for debugging
	// fmt.Printf("Response: %+v\n", resp)

	if resp.Msg.GetLedgerTip() != nil {
		fmt.Printf(
			"Ledger Tip:\n  Slot: %d\n  Hash: %x\n",
			resp.Msg.GetLedgerTip().GetSlot(),
			resp.Msg.GetLedgerTip().GetHash(),
		)
	}

	for _, item := range resp.Msg.GetItems() {
		fmt.Println("UTxO Data:")
		fmt.Printf("  Tx Hash: %x\n", item.GetTxoRef().GetHash())
		fmt.Printf("  Output Index: %d\n", item.GetTxoRef().GetIndex())
		fmt.Printf("  Native Bytes: %x\n", item.GetNativeBytes())
		if cardano := item.GetCardano(); cardano != nil {
			fmt.Println("  Cardano UTxO:")
			fmt.Printf("    Address: %x\n", cardano.GetAddress())
			fmt.Printf("    Coin: %d\n", cardano.GetCoin())
			fmt.Println("    Assets:")
			for _, multiasset := range cardano.GetAssets() {
				fmt.Printf("      Policy ID: %x\n", multiasset.GetPolicyId())
				for _, asset := range multiasset.GetAssets() {
					fmt.Printf("        Asset Name: %s\n", string(asset.GetName()))
					fmt.Printf("        Output Coin: %d\n", asset.GetOutputCoin())
					fmt.Printf("        Mint Coin: %d\n", asset.GetMintCoin())
				}
			}
		}
	}
}
