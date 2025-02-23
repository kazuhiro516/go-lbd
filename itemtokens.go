package lbd

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Meta struct {
	data map[string]string
}

func NewMeta() *Meta {
	return &Meta{
		data: map[string]string{},
	}
}

func (m *Meta) Set(key, value string) (err error) {
	if len(key) < 1 && 15 < len(key) {
		return fmt.Errorf("Invalid key length")
	}
	if len(value) < 1 && 15 < len(value) {
		return fmt.Errorf("Invalid value length")
	}
	m.data[key] = value
	return nil
}

func (m *Meta) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.data)
}

func UnmarshalMeta(data []byte) (*Meta, error) {
	m := new(Meta)
	return m, json.Unmarshal(data, &m.data)
}

func (m *Meta) String() string {
	b, _ := m.MarshalJSON()
	return string(b)
}

type ItemTokenContractInformation struct {
	ContractID   string `json:"contractId"`
	BaseImgURI   string `json:"baseImgUri"`
	OwnerAddress string `json:"ownerAddress"`
	CreatedAt    int64  `json:"createdAt"`
	ServiceID    string `json:"serviceId"`
}

func (l LBD) RetrieveItemTokenContractInformation(contractId string) (*ItemTokenContractInformation, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s", contractId)
	r := NewGetRequest(path)
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}

	ret := &ItemTokenContractInformation{}
	return ret, json.Unmarshal(resp.ResponseData, &ret)
}

type TokenType struct {
	TokenType   string   `json:"tokenType"`
	Name        string   `json:"name"`
	Meta        string   `json:"meta"`
	CreatedAt   int64    `json:"createdAt"`
	TotalSupply string   `json:"totalSupply"`
	TotalMint   string   `json:"totalMint"`
	TotalBurn   string   `json:"totalBurn"`
	Token       []*Token `json:"token"`
}

type Token struct {
	TokenIndex string `json:"tokenIndex"`
	Name       string `json:"name"`
	Meta       string `json:"meta"`
	CreatedAt  int64  `json:"createdAt"`
	BurnedAt   int64  `json:"burnedAt"`
}

func (l LBD) ListAllNonFungibles(contractId string) ([]*TokenType, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles", contractId)

	all := []*TokenType{}
	page := 1
	for {
		r := NewGetRequest(path)
		r.pager.Page = page
		r.pager.OrderBy = "asc"
		resp, err := l.Do(r, true)
		if err != nil {
			return nil, err
		}
		ret := []*TokenType{}
		err = json.Unmarshal(resp.ResponseData, &ret)
		if err != nil {
			return nil, err
		}
		if len(ret) == 0 {
			break
		}
		all = append(all, ret...)
		page++
	}
	return all, nil
}

type CreateNonFungibleRequest struct {
	*Request
	OwnerAddress string `json:"ownerAddress"`
	OwnerSecret  string `json:"ownerSecret"`
	Name         string `json:"name"`
	Meta         string `json:"meta"`
}

func (r CreateNonFungibleRequest) Encode() string {
	base := r.Request.Encode()
	return fmt.Sprintf("%s?meta=%s&name=%s&ownerAddress=%s&ownerSecret=%s", base, r.Meta, r.Name, r.OwnerAddress, r.OwnerSecret)
}

func (l *LBD) CreateNonFungible(contractId, name, meta string) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles", contractId)
	r := CreateNonFungibleRequest{NewPostRequest(path), l.Owner.Address, l.Owner.Secret, name, meta}
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}

type UpdateNonFungibleTokenTypeRequest struct {
	*Request
	OwnerAddress string `json:"ownerAddress"`
	OwnerSecret  string `json:"ownerSecret"`
	Name         string `json:"name"`
	Meta         string `json:"meta"`
}

func (r UpdateNonFungibleTokenTypeRequest) Encode() string {
	base := r.Request.Encode()
	return fmt.Sprintf("%s?meta=%s&name=%s&ownerAddress=%s&ownerSecret=%s", base, r.Meta, r.Name, r.OwnerAddress, r.OwnerSecret)
}

func (l *LBD) UpdateNonFungibleTokenType(contractId, tokenType, name, meta string) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s", contractId, tokenType)
	r := CreateNonFungibleRequest{NewPutRequest(path), l.Owner.Address, l.Owner.Secret, name, meta}
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}

type NonFungibleTokenType struct {
	*Request
	OwnerAddress string `json:"ownerAddress"`
	OwnerSecret  string `json:"ownerSecret"`
	Name         string `json:"name"`
	Meta         string `json:"meta"`
}

func (l *LBD) RetrieveNonFungibleTokenType(contractId, tokenType string, pager *Pager) (*TokenType, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s", contractId, tokenType)
	if pager == nil {
		pager = &Pager{
			Limit:   10,
			Page:    1,
			OrderBy: "desc",
		}
	}

	r := NewGetRequest(path)
	r.pager = pager
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	ret := new(TokenType)
	return ret, json.Unmarshal(resp.ResponseData, ret)
}

type NonFungibleInformation struct {
	Name      string      `json:"name"`
	TokenID   string      `json:"tokenId"`
	Meta      string      `json:"meta"`
	CreatedAt int64       `json:"createdAt"`
	BurnedAt  interface{} `json:"burnedAt"`
}

func (l *LBD) RetrieveNonFungibleInformation(contractId, tokenType, tokenIndex string) (*NonFungibleInformation, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/%s", contractId, tokenType, tokenIndex)
	r := NewGetRequest(path)
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	ret := new(NonFungibleInformation)
	return ret, json.Unmarshal(resp.ResponseData, ret)
}

// Holders Response Struct
type Holder struct {
	WalletAddress *string `json:"walletAddress"`
	UserID        *string `json:"userId"`
	NumberOfIndex string  `json:"numberOfIndex"`
}

func (l LBD) RetrieveHolderOfSpecificNonFungible(contractId, tokenType string) ([]*Holder, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/holders", contractId, tokenType)

	all := []*Holder{}
	page := 1
	for {
		r := NewGetRequest(path)
		r.pager.Page = page
		r.pager.OrderBy = "asc"
		resp, err := l.Do(r, true)
		if err != nil {
			return nil, err
		}
		ret := []*Holder{}
		err = json.Unmarshal(resp.ResponseData, &ret)
		if err != nil {
			return nil, err
		}
		all = append(all, ret...)
		page++
		if len(ret) < r.pager.Limit {
			break
		}
	}
	return all, nil
}

type ItemTokenHolder struct {
	WalletAddress *string `json:"walletAddress"`
	UserID        *string `json:"userId"`
	TokenID       *string `json:"tokenId"`
	Amount        string  `json:"amount"`
}

func (l LBD) RetrieveTheHolderOfSpecificNonFungible(contractId, tokenType, tokenIndex string) (*ItemTokenHolder, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/%s/holder", contractId, tokenType, tokenIndex)

	r := NewGetRequest(path)
	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	ret := new(ItemTokenHolder)
	return ret, json.Unmarshal(resp.ResponseData, ret)
}

type MintNonFungibleRequest struct {
	*Request
	OwnerAddress string `json:"ownerAddress"`
	OwnerSecret  string `json:"ownerSecret"`
	Name         string `json:"name"`
	Meta         string `json:"meta"`
	ToUserId     string `json:"toUserId,omitempty"`
	ToAddress    string `json:"toAddress,omitempty"`
}

func (r MintNonFungibleRequest) Encode() string {
	base := r.Request.Encode()
	if r.ToUserId != "" {
		return fmt.Sprintf("%s?meta=%s&name=%s&ownerAddress=%s&ownerSecret=%s&toUserId=%s", base, r.Meta, r.Name, r.OwnerAddress, r.OwnerSecret, r.ToUserId)
	}
	return fmt.Sprintf("%s?meta=%s&name=%s&ownerAddress=%s&ownerSecret=%s&toAddress=%s", base, r.Meta, r.Name, r.OwnerAddress, r.OwnerSecret, r.ToAddress)
}

func (l *LBD) MintNonFungible(contractId, tokenType, name, meta, to string) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/mint", contractId, tokenType)

	r := MintNonFungibleRequest{
		Request:      NewPostRequest(path),
		OwnerAddress: l.Owner.Address,
		OwnerSecret:  l.Owner.Secret,
		Name:         name,
		Meta:         meta,
	}

	if l.IsAddress(to) {
		r.ToAddress = to
	} else {
		r.ToUserId = to
	}

	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}

type MintMultipleNonFungibleRequest struct {
	*Request
	OwnerAddress string      `json:"ownerAddress"`
	OwnerSecret  string      `json:"ownerSecret"`
	MintList     []*MintList `json:"mintList"`
	ToUserId     string      `json:"toUserId,omitempty"`
	ToAddress    string      `json:"toAddress,omitempty"`
}

type MintList struct {
	TokenType string `json:"tokenType"`
	Name      string `json:"name"`
	Meta      string `json:"meta"`
}

func (r MintMultipleNonFungibleRequest) Encode() string {
	base := r.Request.Encode()
	names := make([]string, len(r.MintList))
	metas := make([]string, len(r.MintList))
	TokenTypes := make([]string, len(r.MintList))
	for i, m := range r.MintList {
		names[i] = m.Name
		metas[i] = m.Meta
		TokenTypes[i] = m.TokenType
	}
	mintList := fmt.Sprintf("mintList.meta=%s&mintList.name=%s&mintList.tokenType=%s",
		strings.Join(metas, ","),
		strings.Join(names, ","),
		strings.Join(TokenTypes, ","),
	)

	if r.ToUserId != "" {
		return fmt.Sprintf("%s?%s&ownerAddress=%s&ownerSecret=%s&toUserId=%s", base, mintList, r.OwnerAddress, r.OwnerSecret, r.ToUserId)
	}
	ret := fmt.Sprintf("%s?%s&ownerAddress=%s&ownerSecret=%s&toAddress=%s", base, mintList, r.OwnerAddress, r.OwnerSecret, r.ToAddress)
	return ret
}

func (l *LBD) MintMultipleNonFungible(contractId, to string, mintList []*MintList) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/multi-mint", contractId)

	r := MintMultipleNonFungibleRequest{
		Request:      NewPostRequest(path),
		OwnerAddress: l.Owner.Address,
		OwnerSecret:  l.Owner.Secret,
		MintList:     mintList,
	}

	if l.IsAddress(to) {
		r.ToAddress = to
	} else {
		r.ToUserId = to
	}

	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}

type UpdateNonFungibleInformationRequest struct {
	*Request
	OwnerAddress string `json:"ownerAddress"`
	OwnerSecret  string `json:"ownerSecret"`
	Name         string `json:"name"`
	Meta         string `json:"meta,omitempty"`
}

func (r UpdateNonFungibleInformationRequest) Encode() string {
	base := r.Request.Encode()
	if r.Meta != "" {
		return fmt.Sprintf("%s?meta=%s&name=%s&ownerAddress=%s&ownerSecret=%s", base, r.Meta, r.Name, r.OwnerAddress, r.OwnerSecret)
	}
	return fmt.Sprintf("%s?name=%s&ownerAddress=%s&ownerSecret=%s", base, r.Name, r.OwnerAddress, r.OwnerSecret)
}

func (l *LBD) UpdateNonFungibleInformation(contractId, tokenType, tokenIndex, name, meta string) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/%s/%s", contractId, tokenType, tokenIndex)

	r := UpdateNonFungibleInformationRequest{
		Request:      NewPutRequest(path),
		OwnerAddress: l.Owner.Address,
		OwnerSecret:  l.Owner.Secret,
		Name:         name,
		Meta:         meta,
	}

	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}

type UpdateMultipleFungibleTokenIconsRequest struct {
	*Request
	UpdateList []*UpdateList `json:"updateList"`
}

type UpdateList struct {
	TokenType  string `json:"tokenType"`
	TokenIndex string `json:"tokenIndex"`
}

func (r UpdateMultipleFungibleTokenIconsRequest) Encode() string {
	base := r.Request.Encode()
	types := make([]string, len(r.UpdateList))
	indexes := make([]string, len(r.UpdateList))

	for i, m := range r.UpdateList {
		types[i] = m.TokenType
		indexes[i] = m.TokenIndex
	}
	updateList := fmt.Sprintf("updateList.tokenIndex=%s&updateList.tokenType=%s",
		strings.Join(indexes, ","),
		strings.Join(types, ","),
	)

	ret := fmt.Sprintf("%s?%s", base, updateList)
	return ret
}

func (l *LBD) UpdateMultipleFungibleTokenIcons(contactId string, updateList []*UpdateList) (*Transaction, error) {
	path := fmt.Sprintf("/v1/item-tokens/%s/non-fungibles/icon", contactId)

	r := UpdateMultipleFungibleTokenIconsRequest{
		Request:    NewPutRequest(path),
		UpdateList: updateList,
	}

	resp, err := l.Do(r, true)
	if err != nil {
		return nil, err
	}
	return UnmarshalTransaction(resp.ResponseData)
}
