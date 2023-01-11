package models

import (
	"context"
	"encoding/json"
	"errors"
	"redistest/database"
	"strings"

	geojson "github.com/paulmach/go.geojson"
)

type Block struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	ParentID string           `json:"parentID"`
	CentroID geojson.Geometry `json:"centroID"`
	Value    string           `json:"value"`
}

type Tree struct {
	Block Block  `json:"block"`
	Son   []Tree `json:"son"`
}

var ErrNotFound = errors.New("404NotFound")
var PassingWrong = errors.New("passingWrongParameter")
var ParameterNotFound = errors.New("parameterNotFound")
var ParameterNotID = errors.New("parameterNotID")
var BlockTree = errors.New("blocktree")

func (i *Block) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func (i *Block) UnmarshalBinary(marsh []byte) error {
	return json.Unmarshal(marsh, i)
}

func BringAllBlocks() ([]Block, error) {
	ctx := context.Background()

	keys, err := database.DB.Keys(ctx, "*").Result()
	var block Block
	var blocks []Block
	if len(keys) == 0 {
		return blocks, ErrNotFound
	}
	for _, key := range keys {
		value, _ := database.DB.Get(ctx, key).Result()
		block.UnmarshalBinary([]byte(value))
		blocks = append(blocks, block)

	}
	return blocks, err

}

func GetBlockID(id string) (Block, error) {
	ctx := context.Background()
	var block Block

	if strings.Contains(id, "*") {
		return block, PassingWrong
	}

	key, err := database.DB.Keys(ctx, id+":*").Result()
	if err != nil {
		return block, ParameterNotID
	}
	if len(key) == 0 {
		return block, ParameterNotFound
	}
	value, err := database.DB.Get(ctx, key[0]).Result()
	if err != nil {
		return block, ParameterNotID
	}
	block.UnmarshalBinary([]byte(value))

	return block, nil

}

func GetBlockTreedID(id string) (Block, error) {
	var block Block
	ctx := context.Background()
	key, err := database.DB.Get(ctx, id).Result()
	if err != nil {
		return block, ParameterNotID
	}

	block.UnmarshalBinary([]byte(key))

	return block, err

}

func PutBlock(block Block) (Block, error) {
	ctx := context.Background()
	verifyBlock, err := GetBlockID(block.ID)
	if err == nil {
		return block, ParameterNotID
	}
	blockBinary, err := block.MarshalBinary()
	if err != nil {
		return block, ParameterNotID
	}

	if block.ParentID != verifyBlock.ParentID {
		return block, ParameterNotID
	}
	_, err = database.DB.Set(ctx, block.ID+":"+block.ParentID, blockBinary, 0).Result()
	if err != nil {
		return block, ParameterNotID
	}

	return block, nil

}

func CreateBlock(block Block) (Block, error) {
	ctx := context.Background()
	_, err := GetBlockID(block.ID)
	if err == nil {
		return block, ParameterNotID
	}
	blockBinary, err := block.MarshalBinary()
	if err != nil {
		return block, ParameterNotID
	}
	_, err = database.DB.Set(ctx, block.ID+":"+block.ParentID, blockBinary, 0).Result()
	if err != nil {
		return block, ParameterNotID
	}

	return block, nil

}

func DeleteBlockID(id string) error {
	ctx := context.Background()
	var err error
	var block Block

	block, err = GetBlockID(id)
	if err != nil || len(block.ID) == 0 {
		return ParameterNotFound
	}

	parentID, err := GetBlockTreedID(block.ID)
	if err != nil || len(parentID.ID) >= 1 {
		return BlockTree
	}

	value, err := database.DB.Del(ctx, block.ID+":"+block.ParentID).Result()

	if err != nil || value == 0 {
		return ErrNotFound
	}

	return nil

}

func GetTreeID(id string) (Tree, error) {
	ctx := context.Background()
	var tree Tree
	var err error

	tree.Block, err = GetBlockID(id)
	if err != nil {
		return tree, err
	}

	if strings.Contains(id, "*") {
		return tree, PassingWrong
	}

	keys, err := database.DB.Keys(ctx, "*:"+id).Result()
	if err != nil {
		return tree, ParameterNotID
	}
	if len(keys) == 0 {
		return tree, ParameterNotFound
	}

	var valueTree Tree
	for _, key := range keys {
		value, err := database.DB.Get(ctx, key).Result()
		if err != nil {
			return tree, ParameterNotID
		}
		valueTree.Block.UnmarshalBinary([]byte(value))
		valueTree, _ = GetTreeID(valueTree.Block.ID)
		tree.Son = append(tree.Son, valueTree)
	}

	return tree, nil

}
