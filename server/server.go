package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)
type Request interface{
  GetUserBalance(name string) string
  Deposit(name string , amount string)
}

type BankingServer struct {
   Request Request
}

type StoreInformation struct {
	Balance map[string] string

}



func (b *BankingServer) ServeHTTP(w http.ResponseWriter, r *http.Request){
	 name:=""
	 amount:=""
	 name,amount=getPrefix(r,name,amount)

	switch r.URL.Path{
   	case fmt.Sprintf("/user/balance/%s",name):{
        b.showBalance(w,name)
	}
    case fmt.Sprintf("/user/deposit/%s=%s",name,amount):{
        b.processDeposit(w,name,amount)
    }
	/*case fmt.Sprintf("/user/withdraw/%s=%s",name,amount):{
        amount,err:=strconv.ParseFloat(amount,64)
		if err!=nil{
			fmt.Printf("Fail converting")
		}
        b.request.Withdraw(name,amount)

    }*/
 }
}

func (s *StoreInformation) GetUserBalance(name string) string {
	 balance:=s.Balance[name]
	 return balance
}

func (s *StoreInformation) Deposit(name string, amount string){
         s.Balance[name]=amount
}


func (b *BankingServer) showBalance(w http.ResponseWriter,name string){
	 balance:=b.Request.GetUserBalance(name)
	  if balance==""{
		w.WriteHeader(http.StatusNotFound)
	  }
	 fmt.Fprint(w,balance)
}

func (b *BankingServer) processDeposit(w http.ResponseWriter,name string , amount string){
     b.Request.Deposit(name,amount)
	 w.WriteHeader(http.StatusAccepted)

}



 func (s *StoreInformation) Withdraw(name string, amount float64) {
	  convert,_:=strconv.ParseFloat(s.Balance[name],64)
      result:=convert-amount
	  convStringRes:=fmt.Sprint(result)
	  fmt.Printf("%s\n",convStringRes)
	  s.Balance[name]=convStringRes

}


func  getPrefix(r *http.Request,name string, amount string) (string,string){
	if r.Method=="GET"{
         name=strings.TrimPrefix(r.URL.Path,"/user/balance/")
	}
	if(r.Method=="POST"){
	 if contains:=strings.ContainsAny(r.URL.Path,"deposit"); contains{
	  trimmed:=strings.TrimPrefix(r.URL.Path,"/user/deposit/")
	  parts:=strings.SplitN(trimmed,"=",2)
      name=parts[0]
	  amount=parts[1]
	}else if contains:=strings.ContainsAny(r.URL.Path,"withdraw"); contains{
      trimmed:=strings.TrimPrefix(r.URL.Path,"/user/withdraw/")
	  parts:=strings.SplitN(trimmed,"=",2)
      name=parts[0]
	  amount=parts[1]
	}
}
	return name,amount

}



func NewMemoryStorage() *StoreInformation{
	return &StoreInformation{map[string]string{}}
}


