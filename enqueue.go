package workers

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"time"
)

const (
	NanoSecondPrecision = 1000000000.0
)

type EnqueueData struct {
	Queue      string      `json:"queue,omitempty"`
	Class      string      `json:"class"`
	Args       interface{} `json:"args"`
	Jid        string      `json:"jid"`
	EnqueuedAt float64     `json:"enqueued_at"`
	EnqueueOptions
}

const (
	fifo = iota + 1
	lifo
)

type EnqueueOptions struct {
	RetryCount int     `json:"retry_count,omitempty"`
	Retry      bool    `json:"retry,omitempty"`
	At         float64 `json:"at,omitempty"`
	Policy     int     `json:"policy,omitempty"`
}

func Enqueue(queue, class string, args interface{}) (string, error) {
	return EnqueueWithOptions(queue, class, args, EnqueueOptions{Policy: fifo})
}

func EnqueueIn(queue, class string, in float64, args interface{}) (string, error) {
	return EnqueueWithOptions(queue, class, args, EnqueueOptions{At: nowToSecondsWithNanoPrecision() + in, Policy: fifo})
}

func EnqueueAt(queue, class string, at time.Time, args interface{}) (string, error) {
	return EnqueueWithOptions(queue, class, args, EnqueueOptions{At: timeToSecondsWithNanoPrecision(at), Policy: fifo})
}

func EnqueueWithOptions(queue, class string, args interface{}, opts EnqueueOptions) (string, error) {
	now := nowToSecondsWithNanoPrecision()
	data := EnqueueData{
		Queue:          queue,
		Class:          class,
		Args:           args,
		EnqueueOptions: opts,
	}
	argsBytes, _ := json.Marshal(args)

	checksum := sha1.New()
	_, _ = checksum.Write(argsBytes)
	data.Jid = BytesToHex(checksum.Sum(nil))

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	if now < opts.At {
		err := enqueueAt(data.At, bytes)
		return data.Jid, err
	}

	conn := Config.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("sadd", Config.Namespace+"queues", queue)
	if err != nil {
		return "", err
	}
	queue = Config.Namespace + "queue:" + queue
	_, err = conn.Do("zadd", queue, now, bytes)
	if err != nil {
		return "", err
	}

	return data.Jid, nil
}

func enqueueAt(at float64, bytes []byte) error {
	conn := Config.Pool.Get()
	defer conn.Close()

	_, err := conn.Do(
		"zadd",
		Config.Namespace+SCHEDULED_JOBS_KEY, at, bytes,
	)
	if err != nil {
		return err
	}

	return nil
}

func timeToSecondsWithNanoPrecision(t time.Time) float64 {
	return float64(t.UnixNano()) / NanoSecondPrecision
}

func durationToSecondsWithNanoPrecision(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / NanoSecondPrecision
}

func nowToSecondsWithNanoPrecision() float64 {
	return timeToSecondsWithNanoPrecision(time.Now())
}

func BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return string(c)
}
