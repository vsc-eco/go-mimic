package condenser

import (
	"context"
	"log/slog"
	"mimic/mock"
	"mimic/modules/api/services"
	"mimic/modules/db/mimic/accountdb"
	"mimic/modules/db/mimic/blockdb"
	cdb "mimic/modules/db/mimic/condenserdb"
	"slices"
	"strings"
	"time"

	"github.com/chebyrash/promise"
	"github.com/sourcegraph/jsonrpc2"
)

type Condenser struct {
	Logger    *slog.Logger
	BlockDB   blockdb.BlockQuery
	AccountDB accountdb.AccountQuery
}

// Runs initialization in order of how they are passed in to `Aggregate`
func (c *Condenser) Init() error {
	return nil
}

// Runs startup and should be non blocking
func (c *Condenser) Start() *promise.Promise[any] {
	return nil
}

// Runs cleanup once the `Aggregate` is finished
func (c *Condenser) Stop() error {
	return nil
}

type GetAccountsArgs [][]string

// get_accounts
func (c *Condenser) GetAccounts(
	args *GetAccountsArgs,
	reply *[]accountdb.Account,
) {
	nameMatched := (*args)[0]
	db := accountdb.Collection()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	*reply = make([]accountdb.Account, 0)

	if err := db.QueryAccountByNames(ctx, reply, nameMatched); err != nil {
		c.Logger.Error("Failed to query for accounts.", "err", err)
		return
	}
}

// get_dynamic_global_properties
func (c *Condenser) GetDynamicGlobalProperties(
	_ *[]string,
) (*cdb.GlobalProperties, *jsonrpc2.Error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	headBlock := blockdb.HiveBlock{}
	if err := c.BlockDB.QueryHeadBlock(ctx, &headBlock); err != nil {
		c.Logger.Error("failed to query database for head block.", "err", err)
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: "failed to query database for head block",
		}
	}

	reply := &cdb.GlobalProperties{
		HeadBlockNumber: headBlock.BlockNum,
		HeadBlockID:     headBlock.BlockID,
		Time:            headBlock.Timestamp,
		CurrentWitness:  headBlock.Witness,
	}

	return reply, nil
}

// get_current_median_history_price
func (c *Condenser) GetCurrentMedianHistoryPrice(
	args *[]string,
	reply *cdb.MedianPrice,
) {
	// Fake data for now until it gets hooked up with the rest of the mock context
	reply.Base = "100.000 SBD"
	reply.Quote = "100.000 HIVE"
}

// get_reward_fund
func (c *Condenser) GetRewardFund(args *[]string, reply *cdb.RewardFund) {
	if len(*args) == 0 {
		return
	}

	var (
		rewards     []cdb.RewardFund
		mockApiData = "condenser_api.get_reward_fund"
	)

	if err := mock.GetMockData(&rewards, mockApiData); err != nil {
		c.Logger.Error("Failed to read mock data",
			"mock-json", mockApiData, "err", err)
		return
	}

	// just grab the first matched of name
	for _, reward := range rewards {
		if strings.EqualFold(reward.Name, (*args)[0]) {
			*reply = reward
			return
		}
	}
}

// get_withdraw_routes
func (c *Condenser) GetWithdrawRoutes(
	args *[]string,
	reply *[]cdb.WithdrawRoute,
) {
	var (
		routes      []cdb.WithdrawRoute
		mockApiData = "condenser_api.get_withdraw_routes"
	)

	if err := mock.GetMockData(&routes, mockApiData); err != nil {
		c.Logger.Error("Failed to read mock data",
			"mock-json", mockApiData, "err", err)
		return
	}

	*reply = make([]cdb.WithdrawRoute, 0, len(routes))

	user, transferDirection := (*args)[0], (*args)[1]

	allowedDirection := []string{"all", "incoming", "outgoing"}
	if !slices.Contains(allowedDirection, transferDirection) {
		c.Logger.Warn(
			"Invalid transfer direction query, allowed values: incoming, outgoing, all",
		)
		return
	}

	filterMap(&routes, reply, func(r *cdb.WithdrawRoute) bool {
		switch transferDirection {

		case "incoming":
			return strings.EqualFold(user, r.ToAccount)

		case "outgoing":
			return strings.EqualFold(user, r.FromAccount)

		case "all":
			return strings.EqualFold(user, r.FromAccount) ||
				strings.EqualFold(user, r.ToAccount)

		default:
			panic("invalid transfer direction")
		}
	})
}

// get_open_orders
func (c *Condenser) GetOpenOrders(args *[]string, reply *[]cdb.OpenOrder) {
	var (
		orders       []cdb.OpenOrder
		mockFilePath = "condenser_api.get_open_orders"
	)

	if err := mock.GetMockData(&orders, mockFilePath); err != nil {
		c.Logger.Error("Failed to read mock data",
			"mock-json", mockFilePath, "err", err)
		return
	}

	*reply = make([]cdb.OpenOrder, 0, len(orders))

	filterMap(&orders, reply, func(o *cdb.OpenOrder) bool {
		return slices.Contains(*args, o.Seller)
	})
}

// get_conversion_requests
// aka hbd -> hive conversion
func (c *Condenser) GetConversionRequests(
	args *[]int,
	reply *[]cdb.ConversionRequest,
) {
	var (
		conversionRequests []cdb.ConversionRequest
		mockFilePath       = "condenser_api.get_conversion_requests"
	)

	if err := mock.GetMockData(&conversionRequests, mockFilePath); err != nil {
		c.Logger.Error("Failed to read mock data",
			"mock-json", mockFilePath, "err", err)
		return
	}

	*reply = make([]cdb.ConversionRequest, 0, len(conversionRequests))

	filterMap(
		&conversionRequests,
		reply,
		func(e *cdb.ConversionRequest) bool {
			return slices.Contains(*args, int(e.ID))
		},
	)
}

// get_collateralized_conversion_requests
// aka hive -> hbd conversion
// NOTE: docs is empty right now...
// https://developers.hive.io/apidefinitions/#condenser_api.get_collateralized_conversion_requests
func (c *Condenser) GetCollateralizedConversionRequests(
	args *[]string,
	reply *[]cdb.ConversionRequest,
) {
	// For now send empty response until decided as necessary and implemented
	*reply = []cdb.ConversionRequest{}
}

// list_proposals
func (c *Condenser) ListProposals(args *[]any, reply *[]string) {
	// For now send empty response until decided as necessary and implemented
	*reply = []string{}
}

func (t *Condenser) Expose(rm services.RegisterMethod) {
	rm("get_dynamic_global_properties", "GetDynamicGlobalProperties")
	// rm("get_current_median_history_price", "GetCurrentMedianHistoryPrice")
	// rm("get_reward_fund", "GetRewardFund")
	// rm("get_withdraw_routes", "GetWithdrawRoutes")
	// rm("get_open_orders", "GetOpenOrders")
	// rm("get_conversion_requests", "GetConversionRequests")
	// rm(
	// 	"get_collateralized_conversion_requests",
	// 	"GetCollateralizedConversionRequests",
	// )
	// rm("list_proposals", "ListProposals")
	rm("broadcast_transaction", "BroadcastTransaction")
	rm("broadcast_transaction_synchronous", "BroadcastTransactionSynchronous")
	// rm("get_accounts", "GetAccounts")
	// rm("account_create", "AccountCreate")
}

// Filters elements from `data` that matches the predicate `filterFunc`, then
// writes to `buf`
func filterMap[T any](data, buf *[]T, filterFunc func(*T) bool) {
	for _, d := range *data {
		if filterFunc(&d) {
			*buf = append(*buf, d)
		}
	}
}
