package client

import (
	"context"

	"encoding/json"
	"fmt"
	"time"

	"nhooyr.io/websocket"
)

// Client holds an API websocket client. From unreleased go-dvote/client
type Client struct {
	Addr string
	Conn *websocket.Conn
	Ctx  context.Context
}

// RequestMessage holds a decoded request but does not decode the body. from go-dvote
type RequestMessage struct {
	MetaRequest json.RawMessage `json:"request"`

	ID        string `json:"id"`
	Signature string `json:"signature"`
}

// MetaRequest holds a gateway api request, from go-dvote/types
type MetaRequest struct {
	CensusID   string   `json:"censusId,omitempty"`
	CensusURI  string   `json:"censusUri,omitempty"`
	ClaimData  string   `json:"claimData,omitempty"`
	ClaimsData []string `json:"claimsData,omitempty"`
	Content    string   `json:"content,omitempty"`
	Digested   bool     `json:"digested,omitempty"`
	EntityId   string   `json:"entityId,omitempty"`
	From       int64    `json:"from,omitempty"`
	FromID     string   `json:"fromId,omitempty"`
	ListSize   int64    `json:"listSize,omitempty"`
	Method     string   `json:"method"`
	Name       string   `json:"name,omitempty"`
	Nullifier  string   `json:"nullifier,omitempty"`
	// Payload    *VoteTx  `json:"payload,omitempty"`
	ProcessID string   `json:"processId,omitempty"`
	ProofData string   `json:"proofData,omitempty"`
	PubKeys   []string `json:"pubKeys,omitempty"`
	RawTx     string   `json:"rawTx,omitempty"`
	RootHash  string   `json:"rootHash,omitempty"`
	Signature string   `json:"signature,omitempty"`
	Timestamp int32    `json:"timestamp"`
	Type      string   `json:"type,omitempty"`
	URI       string   `json:"uri,omitempty"`
}

// ResponseMessage wraps an api response, from go-dvote/types
type ResponseMessage struct {
	MetaResponse json.RawMessage `json:"response"`

	ID        string `json:"id"`
	Signature string `json:"signature"`
}

// MetaResponse holds a gateway api response, from go-dvote/types
type MetaResponse struct {
	APIList        []string   `json:"apiList,omitempty"`
	BlockTime      *[5]int32  `json:"blockTime,omitempty"`
	BlockTimestamp int32      `json:"blockTimestamp,omitempty"`
	CensusID       string     `json:"censusId,omitempty"`
	CensusList     []string   `json:"censusList,omitempty"`
	ClaimsData     []string   `json:"claimsData,omitempty"`
	Content        string     `json:"content,omitempty"`
	EntityID       string     `json:"entityId,omitempty"`
	EntityIDs      []string   `json:"entityIds,omitempty"`
	Files          []byte     `json:"files,omitempty"`
	Finished       *bool      `json:"finished,omitempty"`
	Health         int32      `json:"health,omitempty"`
	Height         *int64     `json:"height,omitempty"`
	InvalidClaims  []int      `json:"invalidClaims,omitempty"`
	Message        string     `json:"message,omitempty"`
	Nullifier      string     `json:"nullifier,omitempty"`
	Nullifiers     *[]string  `json:"nullifiers,omitempty"`
	Ok             bool       `json:"ok"`
	Paused         *bool      `json:"paused,omitempty"`
	Payload        string     `json:"payload,omitempty"`
	ProcessIDs     []string   `json:"processIds,omitempty"`
	ProcessList    []string   `json:"processList,omitempty"`
	Registered     *bool      `json:"registered,omitempty"`
	Request        string     `json:"request"`
	Results        [][]uint32 `json:"results,omitempty"`
	Root           string     `json:"root,omitempty"`
	Siblings       string     `json:"siblings,omitempty"`
	Size           *int64     `json:"size,omitempty"`
	State          string     `json:"state,omitempty"`
	Timestamp      int32      `json:"timestamp"`
	Type           string     `json:"type,omitempty"`
	URI            string     `json:"uri,omitempty"`
	ValidProof     *bool      `json:"validProof,omitempty"`
}

// New starts a connection with the given endpoint address. From unreleased go-dvote/client
func New(addr string) (*Client, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return nil, cancel, err
	}
	return &Client{Addr: addr, Conn: conn, Ctx: ctx}, cancel, nil
}

// GetGatewayInfo calls gateway api and returns gateway info
func (c *Client) GetGatewayInfo() (*MetaResponse, error) {
	var req MetaRequest
	req.Method = "getGatewayInfo"
	req.Timestamp = int32(time.Now().Unix())

	resp, err := c.Request(req)
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		return nil, fmt.Errorf("cannot get keys for process")
	}
	return resp, nil
}

func (c *Client) GetEnvelopeStatus(nullifier, pid string) (bool, error) {
	var req MetaRequest
	req.Method = "getEnvelopeStatus"
	req.ProcessID = pid
	req.Nullifier = nullifier
	resp, err := c.Request(req)
	if err != nil {
		return false, err
	}
	if !resp.Ok || resp.Registered == nil {
		return false, fmt.Errorf("cannot check envelope (%s)", resp.Message)
	}
	return *resp.Registered, nil
}

// // Request makes a request to the previously connected endpoint
func (c *Client) Request(req MetaRequest) (*MetaResponse, error) {
	method := req.Method
	req.Timestamp = int32(time.Now().Unix())
	reqInner, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	reqOuter := RequestMessage{
		ID: "req-2345679",
		// ID:          fmt.Sprintf("%d", rand.Intn(1000)),
		MetaRequest: reqInner,
	}
	reqBody, err := json.Marshal(reqOuter)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}

	// Set context for request
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	fmt.Println("request: %s", reqBody)
	if err := c.Conn.Write(ctx, websocket.MessageText, reqBody); err != nil {
		return nil, fmt.Errorf("Error: %s: %v", method, err)
	}
	fmt.Println("sent request")
	_, message, err := c.Conn.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	fmt.Println("response: %s", message)
	var respOuter ResponseMessage
	if err := json.Unmarshal(message, &respOuter); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	if respOuter.ID != reqOuter.ID {
		return nil, fmt.Errorf("%s: %v", method, "request ID doesn't match")
	}
	if respOuter.Signature == "" {
		return nil, fmt.Errorf("%s: empty signature in response: %s", method, message)
	}
	var respInner MetaResponse
	if err := json.Unmarshal(respOuter.MetaResponse, &respInner); err != nil {
		return nil, fmt.Errorf("%s: %v", method, err)
	}
	return &respInner, nil
}

func (c *Client) Close() {
	err := c.Conn.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Closed websocket connection")
	}
}
