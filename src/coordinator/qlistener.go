package coordinator
import(
  "github.com/streadway/amqp"
  "qutils/qutils"
  "bytes"
  "encoding/gob"
  "dto/dto"
  "fmt"
  )                        

const url = "amqp://guest:guest@localhost:5672"
type QueueListener                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       struct(
  conn *amqp.Connection
  ch *amqp.Channel
  sources map[string]<-chan amqp.Delivery
  )               
                                                                                                                                                                                                                                                                                                                                                                                                       
func NewQueueListener *QueueListener {
  ql:= QueueListener{
    sources: make(map[string]<-chan amqp.Delivery),
  }

  ql.conn, ql.ch = qutils.GetChannel(url)
  return &ql
}

func (ql *QueueListener) ListenForNewSource(){
  q := qutils.GetQueue("",qql.ch)

  ql.ch.QueueBind(
    q.Name,
    "",
    "amq.fanout",
    false,
    nil)

  msgs,_ := ql.ch.Consume(
    q.Name,
    "",
    true,
    false,
    false,
    false,
    nil)

  for msg := range msgs {
    SourceChan,_ := ql.ch.Consume(
      string(msg.Body),
      "",
      true,
      false,
      false,
      false,
      nil)

    if ql.sources[string(msg.Body)] == nil{
      ql.sources[sring(msg.Body)] = SourceChan

      go ql.AddListener(sourceChan)
    }

  }
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery){
  for msg:= range msgs{
    r:= bytes.NewReader(msg.Body)
    d:= gob.NewDecoder(r)
    sd:=  new(dto.SensorMessage)
    d.Decode(sd)

    fmt.Printf("Recieved messages: %v\n",sd)
  }
}
