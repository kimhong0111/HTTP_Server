package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)
type Request interface{
  GetUserBalance(name string) string
  Deposit(name string , amount string)
  Withdraw(name string, amount float64)
}

type BankingServer struct {
   Request Request
}

type StoreInformation struct {
	Balance map[string] string

}



func (b *BankingServer) ServeHTTP(w http.ResponseWriter, r *http.Request){
	 
     router:=http.NewServeMux()

	 router.Handle("/user/balance/",http.HandlerFunc(b.balanceHandler))
	 router.Handle("/user/deposit/",http.HandlerFunc(b.depositHandler))
	 router.Handle("/user/withdraw/",http.HandlerFunc(b.withdrawHandler))	


	 router.ServeHTTP(w,r)
}

func (b *BankingServer) balanceHandler(w http.ResponseWriter, r *http.Request){
	 name:=""
	 amount:=""
	 name,_=getPrefix(r,name,amount)
	 b.showBalance(w,name)
}


func (b *BankingServer) depositHandler(w http.ResponseWriter, r *http.Request){
	 name:=""
	 amount:=""
	 name,amount=getPrefix(r,name,amount)
	 b.processDeposit(w,r,name,amount)
}

func (b *BankingServer) withdrawHandler(w http.ResponseWriter, r *http.Request){
	 name:=""
	 amount:=""
	 name,amount=getPrefix(r,name,amount)
	 b.processWithdraw(w,r,name,amount)
}

func (s *StoreInformation) GetUserBalance(name string) string {
	 balance:=s.Balance[name]
	 return balance
}

func (s *StoreInformation) Deposit(name string, amount string){
         s.Balance[name]=amount
}


 func (s *StoreInformation) Withdraw(name string, amount float64) {
	  convert,_:=strconv.ParseFloat(s.Balance[name],64)
	  if convert < amount{
		log.Fatal("Unsufficent Balance")
		return
	  }
      result:=convert-amount
	  convStringRes:=fmt.Sprintf("%.2f",result)
	  s.Balance[name]=convStringRes

}


func (b *BankingServer) showBalance(w http.ResponseWriter,name string){
	 balance:=b.Request.GetUserBalance(name)
	  if balance==""{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w,"Need to make a deposit request")
	  }
	 fmt.Fprint(w,balance)
}

func (b *BankingServer) processDeposit(w http.ResponseWriter,r *http.Request,name string , amount string){
	 if r.Method!="POST"{
		fmt.Fprint(w,"Request not accepted")
		return
	 }
     b.Request.Deposit(name,amount)
	 fmt.Fprint(w,"Processing Deposit\n")
	 w.WriteHeader(http.StatusAccepted)

}

func (b *BankingServer) processWithdraw(w http.ResponseWriter,r *http.Request,name string , amount string){
	 if r.Method!="POST"{
		fmt.Fprint(w,"Request not accepted")
        return
	 }
	 convert,err:=strconv.ParseFloat(amount,64)
	 if err!=nil{
		fmt.Println("error converting to float")
	 }
	 b.Request.Withdraw(name,convert)
     fmt.Fprint(w,"Processing Withdraw")
	 w.WriteHeader(http.StatusAccepted)
}


func  getPrefix(r *http.Request,name string, amount string) (string,string){
	if r.Method=="GET"{
         name=strings.TrimPrefix(r.URL.Path,"/user/balance/")
	}
	if(r.Method=="POST"){
	 if contains:=strings.Contains(r.URL.Path,"deposit"); contains{
	  trimmed:=strings.TrimPrefix(r.URL.Path,"/user/deposit/")
	  parts:=strings.SplitN(trimmed,"=",2)
      name=parts[0]
	  amount=parts[1]
	}else if contains:=strings.Contains(r.URL.Path,"withdraw"); contains{
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


