package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

type GenesisAlloc map[string]GenesisAccount

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       string            `json:"code,omitempty"`
	Storage    map[string]string `json:"storage,omitempty"`
	Balance    string            `json:"balance"             // gencodec:"required"`
	Nonce      string            `json:"nonce,omitempty"`
	PrivateKey string            `json:"secretKey,omitempty"` // for tests
}

type Genesis struct {
	Config     ChainConfig               `json:"config"`
	Nonce      string                    `json:"nonce"`
	Timestamp  string                    `json:"timestamp"`
	ExtraData  string                    `json:"extraData"`
	GasLimit   string                    `json:"gasLimit"   gencodec:"required"`
	Difficulty string                    `json:"difficulty" gencodec:"required"`
	Alloc      map[string]GenesisAccount `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
	BaseFee    *big.Int    `json:"baseFeePerGas"`
}

type ChainConfig struct {
	Type    string   `json:"type,omitempty"     yaml:"type"`     // type indicates which processor to use for transaction applying
	ChainID *big.Int `json:"chain_id,omitempty" yaml:"chain_id"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homestead_block,omitempty" yaml:"homestead_block,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"dao_fork_block,omitempty"   yaml:"dao_fork_block,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"dao_fork_support,omitempty" yaml:"dao_fork_support,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int `json:"eip_150_block,omitempty" yaml:"eip_150_block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  string   `json:"eip_150_hash,omitempty"  yaml:"eip_150_hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip_155_block,omitempty" yaml:"eip_155_block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip_158_block,omitempty" yaml:"eip_158_block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantium_block,omitempty"       yaml:"byzantium_block,omitempty"`       // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinople_block,omitempty"  yaml:"constantinople_block,omitempty"`  // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburg_block,omitempty"      yaml:"petersburg_block,omitempty"`      // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbul_block,omitempty"        yaml:"istanbul_block,omitempty"`        // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muir_glacier_block,omitempty"    yaml:"muir_glacier_block,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
	BerlinBlock         *big.Int `json:"berlin_block,omitempty"          yaml:"berlin_block,omitempty"`          // Berlin switch block (nil = no fork, 0 = already on berlin)
	LondonBlock         *big.Int `json:"london_block,omitempty"          yaml:"london_block,omitempty"`          // London switch block (nil = no fork, 0 = already on london)
	ArrowGlacierBlock   *big.Int `json:"arrow_glacier_block,omitempty"   yaml:"arrow_glacier_block,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeNetSplitBlock  *big.Int `json:"merge_net_split_block,omitempty" yaml:"merge_net_split_block,omitempty"` // Virtual fork after The Merge to use as a network splitter

	// Fork scheduling was switched from blocks to timestamps here
	ShanghaiTime *uint64 `json:"shanghai_time,omitempty" yaml:"shanghai_time,omitempty"` // Shanghai switch time (nil = no fork, 0 = already on shanghai)

	YoloV3Block *big.Int `json:"yolo_v_3_block,omitempty" yaml:"yolo_v_3_block,omitempty"` // YOLO v3: Gas repricings
	EWASMBlock  *big.Int `json:"ewasm_block,omitempty"    yaml:"ewasm_block,omitempty"`    // EWASM switch block (nil = no fork, 0 = already activated)

	// Polygon network
	JaipurBlock *big.Int `json:"jaipur_block,omitempty" yaml:"jaipur_block,omitempty"`

	// Avalanche Network Upgrades
	ApricotPhase1BlockTimestamp     *big.Int `json:"apricot_phase_1_block_timestamp,omitempty"      yaml:"apricot_phase_1_block_timestamp,omitempty"`      // Apricot Phase 1 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhase2BlockTimestamp     *big.Int `json:"apricot_phase_2_block_timestamp,omitempty"      yaml:"apricot_phase_2_block_timestamp,omitempty"`      // Apricot Phase 2 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhase3BlockTimestamp     *big.Int `json:"apricot_phase_3_block_timestamp,omitempty"      yaml:"apricot_phase_3_block_timestamp,omitempty"`      // Apricot Phase 3 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhase4BlockTimestamp     *big.Int `json:"apricot_phase_4_block_timestamp,omitempty"      yaml:"apricot_phase_4_block_timestamp,omitempty"`      // Apricot Phase 4 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhase5BlockTimestamp     *big.Int `json:"apricot_phase_5_block_timestamp,omitempty"      yaml:"apricot_phase_5_block_timestamp,omitempty"`      // Apricot Phase 5 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhasePre6BlockTimestamp  *big.Int `json:"apricot_phase_pre_6_block_timestamp,omitempty"  yaml:"apricot_phase_pre_6_block_timestamp,omitempty"`  // Apricot Phase Pre 6 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhase6BlockTimestamp     *big.Int `json:"apricot_phase_6_block_timestamp,omitempty"      yaml:"apricot_phase_6_block_timestamp,omitempty"`      // Apricot Phase 6 Block Timestamp (nil = no fork, 0 = already activated)
	ApricotPhasePost6BlockTimestamp *big.Int `json:"apricot_phase_post_6_block_timestamp,omitempty" yaml:"apricot_phase_post_6_block_timestamp,omitempty"` // Apricot Phase Post 6 Block Timestamp (nil = no fork, 0 = already activated)
	BanffBlockTimestamp             *big.Int `json:"banff_block_timestamp,omitempty"                yaml:"banff_block_timestamp,omitempty"`
	CortinaBlockTimestamp           *big.Int `json:"cortina_block_timestamp,omitempty"              yaml:"cortina_block_timestamp,omitempty"`

	// Optimism Bedrock network
	BedrockBlock          *big.Int        `json:"bedrock_block,omitempty"    yaml:"bedrock_block,omitempty"`
	RegolithTime          *uint64         `json:"regolith_time,omitempty"    yaml:"regolith_time,omitempty"`
	OptimismBedrockConfig *OptimismConfig `json:"optimism_bedrock,omitempty" yaml:"optimism_bedrock,omitempty"`

	// Various consensus engines
	Ethash *params.EthashConfig `json:"ethash,omitempty" yaml:"ethash,omitempty"`
	Clique *params.CliqueConfig `json:"clique,omitempty" yaml:"clique,omitempty"`

	ExtraEIPs []int `json:"extra_eips,omitempty" yaml:"extra_eips,omitempty"`
}

type OptimismConfig struct {
	EIP1559Elasticity  uint64 `json:"eip_1559_elasticity"  yaml:"eip_1559_elasticity"`
	EIP1559Denominator uint64 `json:"eip_1559_denominator" yaml:"eip_1559_denominator"`
}
