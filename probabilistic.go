package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9/internal/proto"
)

type ProbabilisticCmdble interface {
	BFAdd(ctx context.Context, key, item interface{}) *IntCmd
	BFCard(ctx context.Context, key string) *IntCmd
	BFExists(ctx context.Context, key, item interface{}) *IntCmd
	BFInfo(ctx context.Context, key string) *BFInfoCmd
	BFInfoArg(ctx context.Context, key string, option BFInfo) *BFInfoCmd
	BFInsert(ctx context.Context, key string, options *BFReserveOptions, items ...interface{}) *IntSliceCmd
	BFMAdd(ctx context.Context, key string, items ...interface{}) *IntSliceCmd
	BFMExists(ctx context.Context, key string, items ...interface{}) *IntSliceCmd
	BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) *StatusCmd
	BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity, expansion int64) *StatusCmd
	BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) *StatusCmd
	BFReserveArgs(ctx context.Context, key string, options *BFReserveOptions) *StatusCmd
	//TODO Loadchunk and scandump missing

	CFAdd(ctx context.Context, key, item interface{}) *IntCmd
	CFAddNX(ctx context.Context, key, item interface{}) *IntCmd
	CFCount(ctx context.Context, key, item interface{}) *IntCmd
	CFDel(ctx context.Context, key string) *IntCmd
	CFExists(ctx context.Context, key, item interface{}) *IntCmd
	CFInfo(ctx context.Context, key string) *CFInfoCmd
	CFReserve(ctx context.Context, key string, capacity int64) *StatusCmd
	CFInsert(ctx context.Context, key string, options *CFInsertOptions, items ...interface{}) *IntSliceCmd
	CFInsertNx(ctx context.Context, key string, options *CFInsertOptions, items ...interface{}) *IntSliceCmd
	CFMExists(ctx context.Context, key string, items ...interface{}) *IntSliceCmd

	CMSIncrBy(ctx context.Context, key string, items ...interface{}) *IntSliceCmd
	CMSInfo(ctx context.Context, key string) *CMSInfoCmd
	CMSInitByDim(ctx context.Context, key string, width, height int64) *StatusCmd
	CMSInitByProb(ctx context.Context, key string, errorRate, probability float64) *StatusCmd
	CMSMerge(ctx context.Context, destKey string, sourceKeys ...string) *StatusCmd
	CMSMergeWithWeight(ctx context.Context, destKey string, sourceKeys map[string]int) *StatusCmd
	CMSQuery(ctx context.Context, key string, items ...interface{}) *IntSliceCmd

	TOPKAdd(ctx context.Context, key string, items ...interface{}) *StringSliceCmd
	TOPKReserve(ctx context.Context, key string, k int) *StatusCmd
	TOPKReserveWithOptions(ctx context.Context, key string, k int, width, depth int64, decay float64) *StatusCmd
	TOPKInfo(ctx context.Context, key string) *TOPKInfoCmd
	TOPKQuery(ctx context.Context, key string, items ...interface{}) *BoolSliceCmd
	TOPKCount(ctx context.Context, key string, items ...interface{}) *IntSliceCmd
	TOPKIncrBy(ctx context.Context, key string, items ...interface{}) *StringSliceCmd
}

type BFReserveOptions struct {
	Capacity   int64
	Error      float64
	Expansion  int64
	NonScaling bool
}

type CFReserveOptions struct {
	Capacity      int64
	BucketSize    int64
	MaxIterations int64
	Expansion     int64
}

type CFInsertOptions struct {
	Capacity int64
	NoCreate bool
}

type BFInfoArgs int

const (
	BFCAPACITY BFInfoArgs = iota
	BFSIZE
	BFFILTERS
	BFITEMS
	BFEXPANSION
)

func (b BFInfoArgs) String() string {
	switch b {
	case BFCAPACITY:
		return "capacity"
	case BFSIZE:
		return "size"
	case BFFILTERS:
		return "filters"
	case BFITEMS:
		return "items"
	case BFEXPANSION:
		return "expansion"
	}
	return ""
}

// -------------------------------------------
// Bloom filter commands
//-------------------------------------------

func (c cmdable) BFReserve(ctx context.Context, key string, errorRate float64, capacity int64) *StatusCmd {
	args := []interface{}{"bf.reserve", key, errorRate, capacity}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFReserveExpansion(ctx context.Context, key string, errorRate float64, capacity, expansion int64) *StatusCmd {
	args := []interface{}{"bf.reserve", key, errorRate, capacity, "expansion", expansion}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFReserveNonScaling(ctx context.Context, key string, errorRate float64, capacity int64) *StatusCmd {
	args := []interface{}{"bf.reserve", key, errorRate, capacity, "nonscaling"}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFReserveArgs(ctx context.Context, key string, options *BFReserveOptions) *StatusCmd {
	args := []interface{}{"bf.reserve", key}
	if options != nil {
		if options.Error != 0 {
			args = append(args, options.Error)
		}
		if options.Capacity != 0 {
			args = append(args, options.Capacity)
		}
		if options.Expansion != 0 {
			args = append(args, "expansion", options.Expansion)
		}
		if options.NonScaling {
			args = append(args, "nonscaling")
		}
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFAdd(ctx context.Context, key, item interface{}) *BoolCmd {
	args := []interface{}{"bf.add", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFCard(ctx context.Context, key string) *IntCmd {
	args := []interface{}{"bf.card", key}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFExists(ctx context.Context, key, item string) *BoolCmd {
	args := []interface{}{"bf.exists", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFInfo(ctx context.Context, key string) *BFInfoCmd {
	args := []interface{}{"bf.info", key}
	cmd := NewBFInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type BFInfo struct {
	Capacity         int64
	Size             int64
	NumFilters       int64
	NumItemsInserted int64
	ExpansionRate    int64
}

type BFInfoCmd struct {
	baseCmd

	val BFInfo
}

func NewBFInfoCmd(ctx context.Context, args ...interface{}) *BFInfoCmd {
	return &BFInfoCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *BFInfoCmd) SetVal(val BFInfo) {
	cmd.val = val
}

func (cmd *BFInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *BFInfoCmd) Val() BFInfo {
	return cmd.val
}

func (cmd *BFInfoCmd) Result() (BFInfo, error) {
	return cmd.val, cmd.err
}

func (cmd *BFInfoCmd) readReply(rd *proto.Reader) (err error) {
	n, err := rd.ReadMapLen()
	if err != nil {
		return err
	}

	var key string
	var result BFInfo
	for f := 0; f < n; f++ {
		key, err = rd.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "Capacity":
			result.Capacity, err = rd.ReadInt()
		case "Size":
			result.Size, err = rd.ReadInt()
		case "Number of filters":
			result.NumFilters, err = rd.ReadInt()
		case "Number of items inserted":
			result.NumItemsInserted, err = rd.ReadInt()
		case "Expansion rate":
			result.ExpansionRate, err = rd.ReadInt()
		default:
			return fmt.Errorf("redis: bloom.info unexpected key %s", key)
		}

		if err != nil {
			return err
		}
	}

	cmd.val = result
	return nil
}

func (c cmdable) BFInfoArg(ctx context.Context, key string, option BFInfoArgs) *BFInfoCmd {
	args := []interface{}{"bf.info", key, option.String()}
	cmd := NewBFInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFInsert(ctx context.Context, key string, options *BFReserveOptions, items ...string) *BoolSliceCmd {
	args := []interface{}{"bf.insert", key}
	if options != nil {
		if options.Error != 0 {
			args = append(args, "error", options.Error)
		}
		if options.Capacity != 0 {
			args = append(args, "capacity", options.Capacity)
		}
		if options.Expansion != 0 {
			args = append(args, "expansion", options.Expansion)
		}
		if options.NonScaling {
			args = append(args, "nonscaling")
		}
	}
	args = append(args, "items")
	for _, s := range items {
		args = append(args, s)
	}

	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFMAdd(ctx context.Context, key string, items ...string) *BoolSliceCmd {
	args := []interface{}{"bf.madd", key}
	for _, s := range items {
		args = append(args, s)
	}
	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) BFMExists(ctx context.Context, key string, items ...string) *BoolSliceCmd {
	args := []interface{}{"bf.mexists", key}
	for _, s := range items {
		args = append(args, s)
	}
	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// -------------------------------------------
// Cuckoo filter commands
//-------------------------------------------

func (c cmdable) CFReserve(ctx context.Context, key string, capacity int64) *StatusCmd {
	args := []interface{}{"cf.reserve", key, capacity}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFReserveArgs(ctx context.Context, key string, options *CFReserveOptions) *StatusCmd {
	args := []interface{}{"cf.reserve", key, options.Capacity}
	if options.BucketSize != 0 {
		args = append(args, "bucketsize", options.BucketSize)
	}
	if options.MaxIterations != 0 {
		args = append(args, "maxiterations", options.MaxIterations)
	}
	if options.Expansion != 0 {
		args = append(args, "expansion", options.Expansion)
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFAdd(ctx context.Context, key, item string) *BoolCmd {
	args := []interface{}{"cf.add", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFAddNX(ctx context.Context, key, item string) *BoolCmd {
	args := []interface{}{"cf.addnx", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFCount(ctx context.Context, key, item string) *IntCmd {
	args := []interface{}{"cf.count", key, item}
	cmd := NewIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFDel(ctx context.Context, key string, item string) *BoolCmd {
	args := []interface{}{"cf.del", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFExists(ctx context.Context, key, item string) *BoolCmd {
	args := []interface{}{"cf.exists", key, item}
	cmd := NewBoolCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type CFInfo struct {
	Size             int64
	NumBuckets       int64
	NumFilters       int64
	NumItemsInserted int64
	NumItemsDeleted  int64
	BucketSize       int64
	ExpansionRate    int64
	MaxIteration     int64
}

type CFInfoCmd struct {
	baseCmd

	val CFInfo
}

func NewCFInfoCmd(ctx context.Context, args ...interface{}) *CFInfoCmd {
	return &CFInfoCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *CFInfoCmd) SetVal(val CFInfo) {
	cmd.val = val
}

func (cmd *CFInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *CFInfoCmd) Val() CFInfo {
	return cmd.val
}

func (cmd *CFInfoCmd) Result() (CFInfo, error) {
	return cmd.val, cmd.err
}

func (cmd *CFInfoCmd) readReply(rd *proto.Reader) (err error) {
	n, err := rd.ReadMapLen()
	if err != nil {
		return err
	}

	var key string
	var result CFInfo
	for f := 0; f < n; f++ {
		key, err = rd.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "Size":
			result.Size, err = rd.ReadInt()
		case "Number of buckets":
			result.NumBuckets, err = rd.ReadInt()
		case "Number of filters":
			result.NumFilters, err = rd.ReadInt()
		case "Number of items inserted":
			result.NumItemsInserted, err = rd.ReadInt()
		case "Number of items deleted":
			result.NumItemsDeleted, err = rd.ReadInt()
		case "Bucket size":
			result.BucketSize, err = rd.ReadInt()
		case "Expansion rate":
			result.ExpansionRate, err = rd.ReadInt()
		case "Max iterations":
			result.MaxIteration, err = rd.ReadInt()

		default:
			return fmt.Errorf("redis: cf.info unexpected key %s", key)
		}

		if err != nil {
			return err
		}
	}

	cmd.val = result
	return nil
}

func (c cmdable) CFInfo(ctx context.Context, key string) *CFInfoCmd {
	args := []interface{}{"cf.info", key}
	cmd := NewCFInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFInsert(ctx context.Context, key string, options *CFInsertOptions, items ...string) *BoolSliceCmd {
	args := []interface{}{"cf.insert", key}
	args = c.getCfInsertArgs(args, options, items...)

	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CFInsertNx(ctx context.Context, key string, options *CFInsertOptions, items ...string) *IntSliceCmd {
	args := []interface{}{"cf.insertnx", key}
	args = c.getCfInsertArgs(args, options, items...)

	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) getCfInsertArgs(args []interface{}, options *CFInsertOptions, items ...string) []interface{} {
	if options != nil {
		if options.Capacity != 0 {
			args = append(args, "capacity", options.Capacity)
		}
		if options.NoCreate {
			args = append(args, "nocreate")
		}
	}
	args = append(args, "items")
	for _, s := range items {
		args = append(args, s)
	}
	return args
}

func (c cmdable) CFMExists(ctx context.Context, key string, items ...string) *BoolSliceCmd {
	args := []interface{}{"cf.mexists", key}
	for _, s := range items {
		args = append(args, s)
	}
	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// -------------------------------------------
// CMS commands
//-------------------------------------------

func (c cmdable) CMSIncrBy(ctx context.Context, key string, items ...interface{}) *IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "cms.incrby"
	args[1] = key
	args = appendArgs(args, items)

	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type CMSInfo struct {
	Width int64
	Depth int64
	Count int64
}

type CMSInfoCmd struct {
	baseCmd

	val CMSInfo
}

func NewCMSInfoCmd(ctx context.Context, args ...interface{}) *CMSInfoCmd {
	return &CMSInfoCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *CMSInfoCmd) SetVal(val CMSInfo) {
	cmd.val = val
}

func (cmd *CMSInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *CMSInfoCmd) Val() CMSInfo {
	return cmd.val
}

func (cmd *CMSInfoCmd) Result() (CMSInfo, error) {
	return cmd.val, cmd.err
}

func (cmd *CMSInfoCmd) readReply(rd *proto.Reader) (err error) {
	n, err := rd.ReadMapLen()
	if err != nil {
		return err
	}

	var key string
	var result CMSInfo
	for f := 0; f < n; f++ {
		key, err = rd.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "width":
			result.Width, err = rd.ReadInt()
		case "depth":
			result.Depth, err = rd.ReadInt()
		case "count":
			result.Count, err = rd.ReadInt()
		default:
			return fmt.Errorf("redis: cms.info unexpected key %s", key)
		}

		if err != nil {
			return err
		}
	}

	cmd.val = result
	return nil
}

func (c cmdable) CMSInfo(ctx context.Context, key string) *CMSInfoCmd {
	args := []interface{}{"cms.info", key}
	cmd := NewCMSInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CMSInitByDim(ctx context.Context, key string, width, depth int64) *StatusCmd {
	args := []interface{}{"cms.initbydim", key, width, depth}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CMSInitByProb(ctx context.Context, key string, errorRate, probability float64) *StatusCmd {
	args := []interface{}{"cms.initbyprob", key, errorRate, probability}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CMSMerge(ctx context.Context, destKey string, sourceKeys ...string) *StatusCmd {
	args := []interface{}{"cms.merge", destKey, len(sourceKeys)}
	for _, s := range sourceKeys {
		args = append(args, s)
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CMSMergeWithWeight(ctx context.Context, destKey string, sourceKeys map[string]int) *StatusCmd {
	args := make([]interface{}, 0, 4+(len(sourceKeys)*2+1))
	args = append(args, "cms.merge", destKey, len(sourceKeys))

	if len(sourceKeys) > 0 {
		sk := make([]interface{}, len(sourceKeys))
		sw := make([]interface{}, len(sourceKeys))

		i := 0
		for k, w := range sourceKeys {
			sk[i] = k
			sw[i] = w
			i++
		}

		args = append(args, sk...)
		args = append(args, "weights")
		args = append(args, sw...)
	}

	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) CMSQuery(ctx context.Context, key string, items ...interface{}) *IntSliceCmd {
	args := []interface{}{"cms.query", key}
	for _, s := range items {
		args = append(args, s)
	}
	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

// -------------------------------------------
// TOPK commands
//-------------------------------------------

func (c cmdable) TOPKAdd(ctx context.Context, key string, items ...interface{}) *StringSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "topk.add"
	args[1] = key
	args = appendArgs(args, items)

	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKReserve(ctx context.Context, key string, k int) *StatusCmd {
	args := []interface{}{"topk.reserve", key, k}

	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKReserveWithOptions(ctx context.Context, key string, k int, width, depth int64, decay float64) *StatusCmd {
	args := []interface{}{"topk.reserve", key, k, width, depth, decay}

	cmd := NewStatusCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

type TOPKInfo struct {
	K     int64
	Width int64
	Depth int64
	Decay float64
}

type TOPKInfoCmd struct {
	baseCmd

	val TOPKInfo
}

func NewTOPKInfoCmd(ctx context.Context, args ...interface{}) *TOPKInfoCmd {
	return &TOPKInfoCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			args: args,
		},
	}
}

func (cmd *TOPKInfoCmd) SetVal(val TOPKInfo) {
	cmd.val = val
}

func (cmd *TOPKInfoCmd) String() string {
	return cmdString(cmd, cmd.val)
}

func (cmd *TOPKInfoCmd) Val() TOPKInfo {
	return cmd.val
}

func (cmd *TOPKInfoCmd) Result() (TOPKInfo, error) {
	return cmd.val, cmd.err
}

func (cmd *TOPKInfoCmd) readReply(rd *proto.Reader) (err error) {
	n, err := rd.ReadMapLen()
	if err != nil {
		return err
	}

	var key string
	var result TOPKInfo
	for f := 0; f < n; f++ {
		key, err = rd.ReadString()
		if err != nil {
			return err
		}

		switch key {
		case "k":
			result.K, err = rd.ReadInt()
		case "width":
			result.Width, err = rd.ReadInt()
		case "depth":
			result.Depth, err = rd.ReadInt()
		case "decay":
			result.Decay, err = rd.ReadFloat()
		default:
			return fmt.Errorf("redis: topk.info unexpected key %s", key)
		}

		if err != nil {
			return err
		}
	}

	cmd.val = result
	return nil
}

func (c cmdable) TOPKInfo(ctx context.Context, key string) *TOPKInfoCmd {
	args := make([]interface{}, 2, 2)
	args[0] = "topk.info"
	args[1] = key

	cmd := NewTOPKInfoCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKQuery(ctx context.Context, key string, items ...interface{}) *BoolSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "topk.query"
	args[1] = key
	args = appendArgs(args, items)

	cmd := NewBoolSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKCount(ctx context.Context, key string, items ...interface{}) *IntSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "topk.count"
	args[1] = key
	args = appendArgs(args, items)

	cmd := NewIntSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKIncrBy(ctx context.Context, key string, items ...interface{}) *StringSliceCmd {
	args := make([]interface{}, 2, 2+len(items))
	args[0] = "topk.incrby"
	args[1] = key
	args = appendArgs(args, items)

	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKList(ctx context.Context, key string) *StringSliceCmd {
	args := []interface{}{"topk.list", key}

	cmd := NewStringSliceCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}

func (c cmdable) TOPKListWithCount(ctx context.Context, key string) *MapStringIntCmd {
	args := []interface{}{"topk.list", key, "withcount"}

	cmd := NewMapStringIntCmd(ctx, args...)
	_ = c(ctx, cmd)
	return cmd
}