package main

/*
import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/taeduard/PI-SMS/rabbitmq_rpc"
	"github.com/taeduard/PI-SMS/utils"
)

type Job struct {
	id        int
	command   string
	StopWorkr bool
}
type Result struct {
	job    Job
	output string
}

var jobs = make(chan Job, 100)
var results = make(chan Result, 100)
var test = make(chan string, 100)

func worker(wg *sync.WaitGroup) {
	x := false
	for {
		if x {
			break
		}
		for job := range jobs {
			if job.stopWorkr {
				wg.Done()
				x = true
			} else {
				test <- "test"
				res, err := rabbitmq_rpc.RPCcommand(job.command)
				utils.FailOnError(err, "Failed to handle RPC request")
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
func allocate(lastjobid int, noOfJobs int, command string, stopchan bool) {
	if stopchan {
		job := Job{0, "", true}
		jobs <- job
		time.Sleep(3 * time.Second)
		close(jobs)
		close(results)
	} else {
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

	go allocate(0, noOfJobs, "df -h /", false)
	go allocate(10, noOfJobs, "df -h /", false)

	done := make(chan bool)
	go result(done)
	noOfWorkers := 1

	createWorkerPool(noOfWorkers)

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")
	go allocate(0, 0, "", true)

	a := <-test
	fmt.Println(a)

	<-done

	// }
}
*/
