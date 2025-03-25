package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Kline struct {
	OpenTime                 int64  `json:"open_time"`
	OpenPrice                string `json:"open_price"`
	HighPrice                string `json:"high_price"`
	LowPrice                 string `json:"low_price"`
	ClosePrice               string `json:"close_price"`
	Volume                   string `json:"volume"`
	CloseTime                int64  `json:"close_time"`
	QuoteAssetVolume         string `json:"quote_asset_volume"`
	NumberOfTrades           int64  `json:"number_of_trades"`
	TakerBuyBaseAssetVolume  string `json:"taker_buy_base_asset_volume"`
	TakerBuyQuoteAssetVolume string `json:"taker_buy_quote_asset_volume"`
}

func main() {
	// req, err := http.NewRequest("GET", "https://api.binance.com/api/v3/klines?symbol=ETHUSDT&interval=1h", nil)
	// if err != nil {
	// 	panic(err)
	// }

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }

	// var data [][]interface{}
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	panic(err)
	// }

	// var klines []Kline
	// for _, d := range data {
	// 	fmt.Println(d[4])

	// 	kline := Kline{
	// 		OpenTime:                 int64(d[0].(float64)),
	// 		OpenPrice:                d[1].(string),
	// 		HighPrice:                d[2].(string),
	// 		LowPrice:                 d[3].(string),
	// 		ClosePrice:               d[4].(string),
	// 		Volume:                   d[5].(string),
	// 		CloseTime:                int64(d[6].(float64)),
	// 		QuoteAssetVolume:         d[7].(string),
	// 		NumberOfTrades:           int64(d[8].(float64)),
	// 		TakerBuyBaseAssetVolume:  d[9].(string),
	// 		TakerBuyQuoteAssetVolume: d[10].(string),
	// 	}
	// 	klines = append(klines, kline)
	// }

	// fmt.Println(klines)

	s := server.NewMCPServer(
		"Demo",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
	)

	calculatorTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Hello World Tool"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Says hello to the name"),
		),
	)

	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)

		return mcp.NewToolResultText("Hello " + name), nil
	})

	spotTickerPrice := mcp.NewTool("spot_ticker_price",
		mcp.WithDescription("Get the current price of a coin"),
		mcp.WithString("symbol",
			mcp.Required(),
			mcp.Description("The symbol of the coin to get the price"),
		),
	)

	s.AddTool(spotTickerPrice, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		symbol := request.Params.Arguments["symbol"].(string)

		req, err := http.NewRequest("GET", "https://api.binance.com/api/v3/ticker/price?symbol="+symbol, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(body)), nil
	})

	spotKline := mcp.NewTool("spot_kline",
		mcp.WithDescription("Get the current kline of a coin"),
		mcp.WithString("symbol",
			mcp.Required(),
			mcp.Description("The symbol of the coin to get the kline data"),
		),
		mcp.WithString("interval",
			mcp.Required(),
			mcp.Description("The interval of the kline data"),
		),
		mcp.WithString("startTime",
			mcp.Description("The start time (milliseconds) of the kline data"),
		),
		mcp.WithString("endTime",
			mcp.Description("The end time (milliseconds) of the kline data"),
		),
	)

	s.AddTool(spotKline, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		symbol := request.Params.Arguments["symbol"].(string)
		interval := request.Params.Arguments["interval"].(string)
		startTime := request.Params.Arguments["startTime"].(string)
		endTime := request.Params.Arguments["endTime"].(string)

		var url string
		if len(startTime) != 0 && len(endTime) != 0 {
			url = "https://api.binance.com/api/v3/klines?symbol=" + symbol + "&interval=" + interval + "&startTime=" + startTime + "&endTime=" + endTime
		} else {
			url = "https://api.binance.com/api/v3/klines?symbol=" + symbol + "&interval=" + interval
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var data [][]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		var klines []Kline
		for _, d := range data {
			kline := Kline{
				OpenTime:                 int64(d[0].(float64)),
				OpenPrice:                d[1].(string),
				HighPrice:                d[2].(string),
				LowPrice:                 d[3].(string),
				ClosePrice:               d[4].(string),
				Volume:                   d[5].(string),
				CloseTime:                int64(d[6].(float64)),
				QuoteAssetVolume:         d[7].(string),
				NumberOfTrades:           int64(d[8].(float64)),
				TakerBuyBaseAssetVolume:  d[9].(string),
				TakerBuyQuoteAssetVolume: d[10].(string),
			}
			klines = append(klines, kline)
		}

		jsonData, err := json.Marshal(klines)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	})

	srv := server.NewStdioServer(s)
	srv.Listen(context.Background(), os.Stdin, os.Stdout)

	// srv := server.NewSSEServer(s)
	// log.Printf("SSE server listening on localhost:8081\n")
	// srv.Start("localhost:8081")

	// if err := server.ServeStdio(s); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }
}
