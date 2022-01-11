package nftport

type ContractResponse struct {
	Response        string `json:"response"`
	Chain           string `json:"chain"`
	ContractAddress string `json:"contract_address"`
	TransactionHash string `json:"transaction_hash"`
	Error           string `json:"error"`
}

type DeployContractRequest struct {
	Chain             string `json:"chain,omitempty"`
	Name              string `json:"name,omitempty"`
	Symbol            string `json:"symbol,omitempty"`
	OwnerAddress      string `json:"owner_address,omitempty"`
	MetadataUpdatable bool   `json:"metadata_updatable,omitempty"`
	BaseUri           string `json:"base_uri,omitempty"`
}

type ContractRequest struct {
	Hash  string `json:"hash"`
	Chain string `json:"chain"`
}

type ErrorResponse struct {
	Response string `json:"response"`
	Error    string `json:"error"`
}
