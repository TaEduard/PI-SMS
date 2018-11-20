package main

import (
	"fmt"
	"log"
	"math/rand"
        "os"
        "bufio"
        "time"
        "sync"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func RPCcommand(n string) (res string, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(n),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			break
		}
	}

	return
}


type Job struct {  
        id       int
        command  string
        stopWorkr bool
    }
    type Result struct {  
        job         Job
        output string
    }
    
    var jobs = make(chan Job, 100)  
    var results = make(chan Result, 100)
    var test = make(chan string, 100)


    func worker(wg *sync.WaitGroup) {  
        x:=false;
        for{
        if x {
                break
        }
        for job := range jobs {
                if job.stopWorkr{
                        wg.Done()
                        x=true;
                }else{
                        test <- "test"
            res, err := RPCcommand(job.command)
	    failOnError(err, "Failed to handle RPC request")
            output := Result{job, res}
            results <- output
        }
        }
        }
        // wg.Done()
    }
    func createWorkerPool(noOfWorkers int) {  
        var wg sync.WaitGroup
        for i := 0; i < noOfWorkers; i++ {
            wg.Add(1)
            go worker(&wg)
        }
        wg.Wait()
    }
    func allocate(lastjobid int ,noOfJobs int, command string,stopchan bool) {  
        if stopchan {
                job := Job{0, "", true}
                jobs <- job
                time.Sleep(3 * time.Second)
                close(jobs)
                close(results)
        }else{
        for i := lastjobid; i < lastjobid+noOfJobs; i++ {
            job := Job{i, command, false}
            jobs <- job
        }
}
    }

    func result(done chan bool) {  
        f, _ := os.Create("/tmp/dat1")
        defer f.Close()
        w := bufio.NewWriter(f)
        for result := range results {
                _, _ = w.WriteString(fmt.Sprintf("Job id %d, input random no %s , sum of digits %s\n", result.job.id, result.job.command, result.output))
                w.Flush()
        }
        done <- true
    }

    func main() {  
        // for {

        startTime := time.Now()
        noOfJobs := 10

        go allocate(0,noOfJobs,"df -h /",false)
        go allocate(10,noOfJobs,"df -h /",false)
        

        done := make(chan bool)
       
        go result(done)
        
        noOfWorkers := 1
        
        createWorkerPool(noOfWorkers)

        // <-done
        // close(done)
        endTime := time.Now()
        diff := endTime.Sub(startTime)
        fmt.Println("total time taken ", diff.Seconds(), "seconds")	
        // go allocate(20,noOfJobs,"df -h /",false)
        // go result(done)
        go allocate(0,0,"",true)
        
        a:= <-test
        fmt.Println( a )


        <-done

// }
}

    
