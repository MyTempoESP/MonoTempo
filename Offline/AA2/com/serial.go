package com

import (
	"fmt"
	"log"
	"time"

	"go.bug.st/serial"
)

// SerialSender represents a serial communication handler that manages sending
// and receiving data over a serial port.
type SerialSender struct {
	port     serial.Port // The serial port instance
	dataCh   chan string // Channel for sending data
	recvCh   chan string // Channel for receiving data
	BaudRate int         // Baud rate for the serial communication
	portName string
}

// NewSerialSender initializes a new SerialSender instance and opens the serial port.
//
// Parameters:
//   - baudRate: The baud rate for the serial communication.
//
// Returns:
//   - sender: A pointer to the initialized SerialSender instance.
//   - err: An error if the initialization or port opening fails.
func NewSerialSender(baudRate int, portName string) (sender *SerialSender, err error) {
	sender = &SerialSender{
		dataCh:   make(chan string),
		recvCh:   make(chan string, 10),
		BaudRate: baudRate,
		portName: portName,
	}

	err = sender.Open()
	if err != nil {
		close(sender.dataCh)
		close(sender.recvCh)
		return
	}

	// Start a goroutine to handle data sending and receiving
	go sender.listenAndSend()
	go sender.recvAndSend()

	return
}

// Open attempts to open the first available serial port with the configured baud rate.
// It retries multiple times with exponential backoff if the port cannot be opened.
//
// Returns:
//   - err: An error if the port cannot be opened after retries.
func (s *SerialSender) Open() (err error) {
	var portName string = s.portName
	var newPort serial.Port

	backoff := time.Millisecond * 100 // Initial backoff duration
	maxRetries := 5                   // Maximum number of retries
	retries := 0

	for retries < maxRetries {
		<-time.After(backoff) // Wait for the backoff duration

		log.Println("Attempting to open the serial port...")

		if portName == "" {
			portName, err = GetFirstAvailablePortName()
			if err != nil {
				log.Printf("Failed to get available port: %v\n", err)
				retries++
				backoff *= 2 // Exponential backoff
				continue
			}
		}

		mode := &serial.Mode{
			BaudRate: s.BaudRate,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		}

		newPort, err = serial.Open(portName, mode)
		if err != nil {
			log.Printf("Failed to open serial port: %v\n", err)
			retries++
			backoff *= 2 // Exponential backoff
			continue
		}

		s.port = newPort
		log.Println("Serial port opened successfully.")
		return
	}

	log.Println("Max retries reached. Unable to open the serial port.")
	return
}

// listenAndSend listens for data on the send channel and writes it to the serial port.
// It also reads incoming data from the serial port and sends it to the receive channel.
func (s *SerialSender) listenAndSend() {
	for data := range s.dataCh {
		_, err := s.port.Write(append([]byte(data), '\n'))
		if err != nil {
			log.Printf("Error writing to serial port: %v\n", err)
			s.port.Close()
			s.Open()
			continue
		}
	}
}

func (s *SerialSender) recvAndSend() {
	t := time.NewTicker(300 * time.Millisecond)

	for range t.C {
		buf := make([]byte, 13)
		c, err := s.port.Read(buf)
		if err != nil {
			log.Printf("Error reading from serial port: %v\n", err)
			continue
		}

		if c > 0 {
			s.recvCh <- string(buf[:c]) // Send the received data to the receive channel
		}
	}
}

// SendData sends the provided data string through the serial port.
//
// Parameters:
//   - data: The string data to send.
func (s *SerialSender) SendData(data string) {
	s.dataCh <- data // Send data to the channel
}

// Recv retrieves data from the receive channel if available.
//
// Returns:
//   - ok: A boolean indicating whether data was successfully received.
//   - data: The received string data.
func (s *SerialSender) Recv() (data string, ok bool) {
	select {
	case data, ok = <-s.recvCh:
	default:
	}
	return
}

// Close closes the serial port and associated channels.
func (s *SerialSender) Close() {
	close(s.dataCh) // Close the send channel
	close(s.recvCh) // Close the receive channel
	if s.port != nil {
		s.port.Close() // Close the serial port
	}
}

// GetFirstAvailablePortName retrieves the name of the first available serial port.
//
// Returns:
//   - port: The name of the first available serial port.
//   - err: An error if no ports are found or if the retrieval fails.
func GetFirstAvailablePortName() (port string, err error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return
	}

	if len(ports) == 0 {
		err = fmt.Errorf("no serial ports found")
		return
	}

	log.Println("Available serial ports:", ports)

	port = ports[0]
	return
}
