package main

import ("fmt"
		"net"
		"net/rpc"
		"net/rpc/jsonrpc"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"strings"
		"strconv")

type Server struct{}
var id int 
var mapSave map[int]string
type Request struct{
	StockSymbolAndPercentage []InnerRequest `json:"stockSymbolAndPercentage"`
	Budget float32 `json:"budget"`
}

type RequestTwo struct{
	Tradeid int `json:"tradeid"`
}
type InnerRequest struct{
	Fields ReqFields `json:"fields"`
}

type ReqFields struct{
	Name string `json:"name"`
	Percentage int `json:"perecentage"`
}

type Response struct{
	Stocks []InnerResponse `json:"stocks"`
	TradeId int `json:"tradeid"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}

type InnerResponse struct{
	ResponseFields ReqResponseFields `json:"fields"`
}

type ReqResponseFields struct{
	Name string `json:"name"`
	Number int `json:"number"`
	Price string `json:"price"`
}

type ResponseTwo struct{
	Stocks []InnerResponse `json:"stocks"`
	CurrentMarketValue float32 `json:"currentMarketValue"`
	UnvestedAmount float32 `json:"unvestedAmount"`
}

func (this *Server) GetInfo(message string,reply *string) error{
	var jsonReq RequestTwo
	var jsonMsg Response
	var jsonInt interface{}
	var company string
	var price []float64
	var structResponseTwo ResponseTwo
	json.Unmarshal([]byte(message),&jsonReq)
	tradeid:= jsonReq.Tradeid
	dataval:= mapSave[tradeid]
	json.Unmarshal([]byte(dataval),&jsonMsg)
	for _,index:= range jsonMsg.Stocks{
		company += index.ResponseFields.Name +","
	}
	company=strings.Trim(company,",")
	response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+company+"/quote?format=json")
	if(err!=nil){
		fmt.Println(err)
	}else{
		defer response.Body.Close()
		contents,_:= ioutil.ReadAll(response.Body)
		json.Unmarshal(contents,&jsonInt)
		for _,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price1,_ := strconv.ParseFloat((index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]).(string),64)
				price = append(price,price1)
			}
		var value float32=0.0
		var strprice string
		for i,index := range jsonMsg.Stocks{
				temp,_:= strconv.ParseFloat(index.ResponseFields.Price,64)
				fmt.Println(price[i],temp)
				if price[i] > temp{
					strprice = "$+"+strconv.FormatFloat(price[i],'f',-1,64)

				}
				if price[i] < temp {
					strprice = "$-"+strconv.FormatFloat(price[i],'f',-1,64)
				}else {
					strprice = "$"+strconv.FormatFloat(price[i],'f',-1,64)
				}
				structReqResponseFields:=ReqResponseFields{Name:index.ResponseFields.Name,Number:index.ResponseFields.Number,Price:strprice}
				structInnerResponse := InnerResponse{ResponseFields:structReqResponseFields}
				structResponseTwo.Stocks = append(structResponseTwo.Stocks,structInnerResponse)
				value = value + (float32(index.ResponseFields.Number) * float32(price[i]))
		}
		rone := &ResponseTwo{
    	CurrentMarketValue:value,
        Stocks: structResponseTwo.Stocks,
        UnvestedAmount:jsonMsg.UnvestedAmount}
    	rtwo, _ := json.Marshal(rone)
    	*reply = string(rtwo)
	
	}		
	return nil
}

func (this *Server) GetMsg(message string,reply *string) error{
		var jsonInt interface{}
		var structResponse Response
		var jsonMsg Request
		var company string
		var remainder float32=0.0
		json.Unmarshal([]byte(message),&jsonMsg)
		for _, i:= range jsonMsg.StockSymbolAndPercentage{
			company += i.Fields.Name +","
		}
		company=strings.Trim(company,",")
		response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+company+"/quote?format=json")
		if(err!=nil){
			fmt.Println(err)
		}else{
			defer response.Body.Close()
			contents,err:= ioutil.ReadAll(response.Body)
			json.Unmarshal(contents,&jsonInt)
			for i,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]
				price1,_ := strconv.ParseFloat(price.(string),64)
				Remainder1:=(float64(jsonMsg.StockSymbolAndPercentage[i].Fields.Percentage) * float64(jsonMsg.Budget))/100
				name := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["symbol"]
				number := int( Remainder1/price1)
				remainder = remainder + (float32(price1)*float32(number))
				structReqResponseFields:=ReqResponseFields{Name:name.(string),Number:number,Price:strconv.FormatFloat(price1,'f',-1,64)}
				structInnerResponse := InnerResponse{ResponseFields:structReqResponseFields}
				structResponse.Stocks = append(structResponse.Stocks,structInnerResponse)
			}
			remainder=jsonMsg.Budget-remainder
			rone := &Response{
    		TradeId:id,
        	Stocks: structResponse.Stocks,
        	UnvestedAmount:remainder}
    		rtwo, _ := json.Marshal(rone)
    		*reply = string(rtwo)
			mapSave[id]=string(rtwo)
			id++
			if(err!=nil){
				fmt.Println(err)
			}
				
		}
		
		return nil
}

func main(){
	id++
	mapSave=make(map[int]string)
	rpc.Register(new(Server))
	hear,err:= net.Listen("tcp",":8080")
	if(err!=nil){
		fmt.Println(err)
		return
	}
	for{
		c,error:= hear.Accept()
		if(error!=nil){
			continue
		}
		go jsonrpc.ServeConn(c)
	}

}
