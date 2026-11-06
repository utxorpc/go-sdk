package cardano

import (
	"bytes"
	"encoding/csv"
	"encoding/hex"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	v1alpha "github.com/utxorpc/go-codegen/utxorpc/v1alpha/cardano"
	v1beta "github.com/utxorpc/go-codegen/utxorpc/v1beta/cardano"
	"google.golang.org/protobuf/proto"
)

type vectorHeader interface {
	GetSlot() uint64
	GetHash() []byte
}

type vectorTx interface {
	GetHash() []byte
}

type decoded struct {
	message proto.Message
	header  vectorHeader
	txs     []vectorTx
}

func TestVectors(t *testing.T) {
	dir := filepath.Join("testdata", "vectors")

	index, err := os.Open(filepath.Join(dir, "vectors.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	defer index.Close()

	reader := csv.NewReader(index)
	reader.Comma = '\t'
	rows, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) < 2 {
		t.Fatal("no vectors")
	}

	decoders := map[string]func([]byte) (decoded, error){
		"v1alpha": func(raw []byte) (decoded, error) {
			block := &v1alpha.Block{}
			err := proto.Unmarshal(raw, block)
			out := decoded{message: block, header: block.GetHeader()}
			for _, tx := range block.GetBody().GetTx() {
				out.txs = append(out.txs, tx)
			}
			return out, err
		},
		"v1beta": func(raw []byte) (decoded, error) {
			block := &v1beta.Block{}
			err := proto.Unmarshal(raw, block)
			out := decoded{message: block, header: block.GetHeader()}
			for _, tx := range block.GetBody().GetTx() {
				out.txs = append(out.txs, tx)
			}
			return out, err
		},
	}

	for version, decode := range decoders {
		for _, row := range rows[1:] {
			name, slot, hash, txCount, firstTx := row[0], row[3], row[4], row[5], row[6]

			t.Run(version+"/"+name, func(t *testing.T) {
				raw, err := os.ReadFile(filepath.Join(dir, version, name+".pb"))
				if err != nil {
					t.Fatal(err)
				}

				block, err := decode(raw)
				if err != nil {
					t.Fatal(err)
				}

				if got := strconv.FormatUint(block.header.GetSlot(), 10); got != slot {
					t.Errorf("slot = %s, want %s", got, slot)
				}
				if got := hex.EncodeToString(block.header.GetHash()); got != hash {
					t.Errorf("hash = %s, want %s", got, hash)
				}
				if got := strconv.Itoa(len(block.txs)); got != txCount {
					t.Errorf("tx count = %s, want %s", got, txCount)
				}
				if firstTx != "-" {
					if got := hex.EncodeToString(block.txs[0].GetHash()); got != firstTx {
						t.Errorf("first tx = %s, want %s", got, firstTx)
					}
				}

				encoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(block.message)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(encoded, raw) {
					t.Error("re-encoded bytes differ")
				}
			})
		}
	}
}
