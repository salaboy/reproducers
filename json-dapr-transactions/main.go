package main

import (
	"context"
	"encoding/json"
	"log"

	dapr "github.com/dapr/go-sdk/client"
)

var (
	STATESTORE_NAME = "statestore"
)

type Vote struct {
	Type    string `json:"type"`
	VoterId string `json:"voterId"`
	Option  string `json:"option"`
}

func main() {
	daprClient, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	vote1 := &Vote{
		Type:    "vote",
		VoterId: "123",
		Option:  "winner",
	}

	vote2 := &Vote{
		Type:    "vote",
		VoterId: "456",
		Option:  "looser",
	}

	vote3 := &Vote{
		Type:    "vote",
		VoterId: "789",
		Option:  "draw",
	}

	jsonData, err := json.Marshal(vote1)
	if err != nil {
		log.Printf("An error occured while marshalling vote to json: %v", err)
	}

	jsonData2, err := json.Marshal(vote2)
	if err != nil {
		log.Printf("An error occured while marshalling vote to json: %v", err)
	}

	jsonData3, err := json.Marshal(vote3)
	if err != nil {
		log.Printf("An error occured while marshalling vote to json: %v", err)
	}

	err = daprClient.SaveState(ctx, STATESTORE_NAME, "voter-"+vote1.VoterId, jsonData, map[string]string{
		"contentType": "application/json",
	})
	if err != nil {
		log.Printf("An error occured while storing the vote: %v", err)
	}

	ops := make([]*dapr.StateOperation, 0)

	op1 := &dapr.StateOperation{
		Type: dapr.StateOperationTypeUpsert,
		Item: &dapr.SetStateItem{
			Key:   "voter-" + vote2.VoterId,
			Value: jsonData2,
		},
	}

	ops = append(ops, op1)

	err = daprClient.ExecuteStateTransaction(ctx, STATESTORE_NAME, map[string]string{
		"contentType": "application/json",
	}, ops)
	if err != nil {
		log.Printf("An error occured while storing transactionally the vote: %v", err)
	}

	item1 := &dapr.SetStateItem{
		Key:   "voter-" + vote3.VoterId,
		Value: jsonData3,
		Metadata: map[string]string{
			"contentType": "application/json",
		},
	}

	if err := daprClient.SaveBulkState(ctx, STATESTORE_NAME, item1); err != nil {
		panic(err)
	}

}
