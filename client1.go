package main

import ("fmt"
        "net/rpc/jsonrpc"
        "strings"
        "bufio"
        "os"
        "strconv"
        "encoding/json")

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

func PurchaseStocks(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:8080")
	if err!=nil{
		fmt.Println(err)
		return
	}
	var reply string
	var structRequest Request
	var message,dataval,newData []string
	fmt.Println("Enter the request")
	in := bufio.NewReader(os.Stdin)
	inp, err := in.ReadString('\n')
	message = strings.SplitN(inp," ",-1)
	dataval = strings.SplitN(message[0],":",2)
	newData = strings.SplitN(message[1],":",2)
	bValue,err:=strconv.ParseFloat(strings.TrimSpace(newData[1]),64)
	dataval[1]= strings.Replace(dataval[1],"\"","",-1)
	dataval[1]= strings.Replace(dataval[1],"%","",-1)
	fields := strings.SplitN(dataval[1],",",-1)
	for _,index:=range fields{
			c:= strings.SplitN(index,":",-1)
			a,_:=strconv.Atoi(c[1])
			structFields := ReqFields{Name:c[0],Percentage:a} 
			structInnerRequest := InnerRequest {Fields:structFields}
			structRequest.StockSymbolAndPercentage =append(structRequest.StockSymbolAndPercentage,structInnerRequest)
	}
	rone := &Request{
    	Budget:float32(bValue),
        StockSymbolAndPercentage: structRequest.StockSymbolAndPercentage}
    rtwo, _ := json.Marshal(rone)
	err = c.Call("Server.GetMsg",string(rtwo),&reply)
	var jsonMsg Response
	var output string
	output = "\"tradeid\":"
	json.Unmarshal([]byte(reply),&jsonMsg)
	output+=strconv.Itoa(jsonMsg.TradeId)+"\n"+"\"stocks\":\""
	for _, i:= range jsonMsg.Stocks{
		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+"$"+i.ResponseFields.Price+","
	}
	output=strings.Trim(output,",")
	output+="\"\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg.UnvestedAmount),'f',-1,32)		
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nResponse:\n")
		fmt.Println(output)
	}
}

func Portfolio(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:8080")
	if err!=nil{
		fmt.Println(err)
		return
	}
	structSRequest:=new(RequestTwo)
	fmt.Println("Please Enter:")
	var sRequest string
	fmt.Scanf("%s",&sRequest)
	sRequest= strings.Replace(sRequest,"\"","",-1)
	newsRequest:=strings.SplitN(sRequest,":",-1)
	structSRequest.Tradeid,_= strconv.Atoi(newsRequest[1])
	result3 := &RequestTwo{
		Tradeid: structSRequest.Tradeid}
	result4,_:= json.Marshal(result3)
	var jsonMsg2 ResponseTwo
	var reply string
	err = c.Call("Server.GetInfo",string(result4),&reply)
	var output string
	output = "\"stocks\":"
	json.Unmarshal([]byte(reply),&jsonMsg2)
	for _, i:= range jsonMsg2.Stocks{
		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+i.ResponseFields.Price+","
	}
	output=strings.Trim(output,",")
	output+="\"\n\"currentMarketValue\":$"+strconv.FormatFloat(float64(jsonMsg2.CurrentMarketValue),'f',-1,32)
	output+="\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg2.UnvestedAmount),'f',-1,32)
	if err!=nil {
		fmt.Println(err)
	}else{
		fmt.Println("\nResponse:\n")
		fmt.Println(output)
	}
}

func main(){
	fmt.Println("Select Your Option\n1.Purchase Stocks\n2.My portfolio")
	var choice int64 
	fmt.Scanf("%d\n",&choice)
	switch choice{
		case 1:
			PurchaseStocks()
			break
		case 2:
			Portfolio()
			break
		default:
			fmt.Println("Select One of the above options")
			break
		}
}
