
import (
	//"encoding/json"
	//"errors"
	"fmt"
	//"github.com/go-goodies/go_oops"
	//"github.com/l3x/jsoncfgo"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	//"regexp"
	"github.com/skycoin/skycoin/src/cipher"
	"os"
	"github.com/skycoin/skycoin-exchange"
)


func main() {

	server := skycoin_exchange.Server

	server.Init()
	server.Run()
}

