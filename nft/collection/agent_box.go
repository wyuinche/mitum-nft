package collection

import (
	"bytes"
	"sort"

	extensioncurrency "github.com/ProtoconNet/mitum-currency-extension/currency"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
	"github.com/spikeekips/mitum/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	AgentBoxType   = hint.Type("mitum-nft-agent-box")
	AgentBoxHint   = hint.NewHint(AgentBoxType, "v0.0.1")
	AgentBoxHinter = AgentBox{BaseHinter: hint.NewBaseHinter(AgentBoxHint)}
)

type AgentBox struct {
	hint.BaseHinter
	collection extensioncurrency.ContractID
	agents     []base.Address
}

func NewAgentBox(symbol extensioncurrency.ContractID, agents []base.Address) AgentBox {
	if agents == nil {
		return AgentBox{BaseHinter: hint.NewBaseHinter(AgentBoxHint), collection: symbol, agents: []base.Address{}}
	}
	return AgentBox{BaseHinter: hint.NewBaseHinter(AgentBoxHint), collection: symbol, agents: agents}
}

func (abx AgentBox) Bytes() []byte {
	bas := make([][]byte, len(abx.agents))
	for i := range abx.agents {
		bas[i] = abx.agents[i].Bytes()
	}

	return util.ConcatBytesSlice(bas...)
}

func (abx AgentBox) Hint() hint.Hint {
	return AgentBoxHint
}

func (abx AgentBox) Hash() valuehash.Hash {
	return abx.GenerateHash()
}

func (abx AgentBox) GenerateHash() valuehash.Hash {
	return valuehash.NewSHA256(abx.Bytes())
}

func (abx AgentBox) IsEmpty() bool {
	return len(abx.agents) < 1
}

func (abx AgentBox) IsValid([]byte) error {
	for i := range abx.agents {
		if err := abx.agents[i].IsValid(nil); err != nil {
			return err
		}
	}
	return nil
}

func (abx AgentBox) Collection() extensioncurrency.ContractID {
	return abx.collection
}

func (abx AgentBox) Equal(b AgentBox) bool {
	abx.Sort(true)
	b.Sort(true)
	for i := range abx.agents {
		if !abx.agents[i].Equal(b.agents[i]) {
			return false
		}
	}
	return true
}

func (abx *AgentBox) Sort(ascending bool) {
	sort.Slice(abx.agents, func(i, j int) bool {
		if ascending {
			return bytes.Compare(abx.agents[j].Bytes(), abx.agents[i].Bytes()) > 0
		}
		return bytes.Compare(abx.agents[j].Bytes(), abx.agents[i].Bytes()) < 0
	})
}

func (abx AgentBox) Exists(ag base.Address) bool {
	if abx.IsEmpty() {
		return false
	}
	for i := range abx.agents {
		if ag.Equal(abx.agents[i]) {
			return true
		}
	}
	return false
}

func (abx AgentBox) Get(ag base.Address) (base.Address, error) {
	for i := range abx.agents {
		if ag.Equal(abx.agents[i]) {
			return abx.agents[i], nil
		}
	}
	return currency.Address{}, errors.Errorf("agent not found in owner's agent box; %v", ag)
}

func (abx *AgentBox) Append(ag base.Address) error {
	if err := ag.IsValid(nil); err != nil {
		return err
	}
	if abx.Exists(ag) {
		return errors.Errorf("agent %v already exists in agent box", ag)
	}
	if len(abx.agents) >= MaxAgents {
		return errors.Errorf("max agents; %v", ag)
	}

	abx.agents = append(abx.agents, ag)
	return nil
}

func (abx *AgentBox) Remove(ag base.Address) error {
	if !abx.Exists(ag) {
		return errors.Errorf("agent %v not found in agent box", ag)
	}
	for i := range abx.agents {
		if ag.String() == abx.agents[i].String() {
			abx.agents[i] = abx.agents[len(abx.agents)-1]
			abx.agents[len(abx.agents)-1] = currency.Address{}
			abx.agents = abx.agents[:len(abx.agents)-1]
			return nil
		}
	}
	return nil
}

func (abx AgentBox) Agents() []base.Address {
	return abx.agents
}

type AgentBoxJSONPacker struct {
	jsonenc.HintedHead
	CL extensioncurrency.ContractID `json:"collection"`
	AG []base.Address               `json:"agents"`
}

func (abx AgentBox) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(AgentBoxJSONPacker{
		HintedHead: jsonenc.NewHintedHead(abx.Hint()),
		CL:         abx.collection,
		AG:         abx.agents,
	})
}

type AgentBoxJSONUnpacker struct {
	CL string                `json:"collection"`
	AG []base.AddressDecoder `json:"agents"`
}

func (abx *AgentBox) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var ubox AgentBoxJSONUnpacker
	if err := enc.Unmarshal(b, &ubox); err != nil {
		return err
	}

	return abx.unpack(enc, ubox.CL, ubox.AG)
}

type AgentBoxBSONPacker struct {
	CL extensioncurrency.ContractID `bson:"collection"`
	AG []base.Address               `bson:"agents"`
}

func (abx AgentBox) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(abx.Hint()),
		bson.M{
			"collection": abx.collection,
			"agents":     abx.agents,
		}),
	)
}

type AgentBoxBSONUnpacker struct {
	CL string                `bson:"collection"`
	AG []base.AddressDecoder `bson:"agents"`
}

func (abx *AgentBox) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubox AgentBoxBSONUnpacker
	if err := bsonenc.Unmarshal(b, &ubox); err != nil {
		return err
	}

	return abx.unpack(enc, ubox.CL, ubox.AG)
}

func (abx *AgentBox) unpack(
	enc encoder.Encoder,
	cl string,
	bags []base.AddressDecoder, // base.Addresss
) error {

	abx.collection = extensioncurrency.ContractID(cl)

	agents := make([]base.Address, len(bags))
	for i := range agents {
		agent, err := bags[i].Encode(enc)
		if err != nil {
			return err
		}

		agents[i] = agent
	}

	abx.agents = agents

	return nil
}
