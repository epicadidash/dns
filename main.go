package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// DNS header structure
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNS question structure
type DNSQuestion struct {
	Name  string
	Type  uint16
	Class uint16
}

// DNS answer structure
type DNSAnswer struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

// Simple DNS record store
var dnsRecords = map[string]string{
	"example.com":    "93.184.216.34",
	"test.local":     "127.0.0.1",
	"myserver.local": "192.168.1.100",
	"google.com":     "8.8.8.8",
}

func main() {
	// Listen on UDP port 53
	addr, err := net.ResolveUDPAddr("udp", ":53")
	if err != nil {
		log.Fatal("Error resolving address:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer conn.Close()

	fmt.Println("DNS Server started on port 53")
	fmt.Println("Configured records:")
	for domain, ip := range dnsRecords {
		fmt.Printf("  %s -> %s\n", domain, ip)
	}

	// Handle incoming requests
	for {
		buffer := make([]byte, 512)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("Error reading UDP message: %v", err)
			continue
		}

		go handleDNSRequest(conn, clientAddr, buffer[:n])
	}
}

func handleDNSRequest(conn *net.UDPConn, clientAddr *net.UDPAddr, data []byte) {
	if len(data) < 12 {
		log.Printf("DNS request too short")
		return
	}

	// Parse DNS header
	header := parseDNSHeader(data)

	// Parse question section
	question, err := parseDNSQuestion(data[12:])
	if err != nil {
		log.Printf("Error parsing DNS question: %v", err)
		return
	}

	fmt.Printf("DNS Query: %s (Type: %d, Class: %d) from %s\n",
		question.Name, question.Type, question.Class, clientAddr)

	// Create response
	response := createDNSResponse(header, question)

	// Send response
	_, err = conn.WriteToUDP(response, clientAddr)
	if err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

func parseDNSHeader(data []byte) DNSHeader {
	return DNSHeader{
		ID:      uint16(data[0])<<8 | uint16(data[1]),
		Flags:   uint16(data[2])<<8 | uint16(data[3]),
		QDCount: uint16(data[4])<<8 | uint16(data[5]),
		ANCount: uint16(data[6])<<8 | uint16(data[7]),
		NSCount: uint16(data[8])<<8 | uint16(data[9]),
		ARCount: uint16(data[10])<<8 | uint16(data[11]),
	}
}

func parseDNSQuestion(data []byte) (DNSQuestion, error) {
	var question DNSQuestion
	var name strings.Builder
	i := 0

	// Parse domain name
	for i < len(data) {
		length := int(data[i])
		if length == 0 {
			i++
			break
		}

		if name.Len() > 0 {
			name.WriteByte('.')
		}

		i++
		if i+length > len(data) {
			return question, fmt.Errorf("invalid name length")
		}

		name.Write(data[i : i+length])
		i += length
	}

	if i+4 > len(data) {
		return question, fmt.Errorf("insufficient data for type and class")
	}

	question.Name = name.String()
	question.Type = uint16(data[i])<<8 | uint16(data[i+1])
	question.Class = uint16(data[i+2])<<8 | uint16(data[i+3])

	return question, nil
}

func createDNSResponse(header DNSHeader, question DNSQuestion) []byte {
	response := make([]byte, 0, 512)

	// DNS Header
	response = append(response, byte(header.ID>>8), byte(header.ID))
	response = append(response, 0x81, 0x80) // Standard query response, no error
	response = append(response, byte(header.QDCount>>8), byte(header.QDCount))

	// Check if we have a record for this domain
	ip, exists := dnsRecords[strings.ToLower(question.Name)]
	if exists && question.Type == 1 { // A record
		response = append(response, 0x00, 0x01) // Answer count = 1
	} else {
		response = append(response, 0x00, 0x00) // Answer count = 0
	}

	response = append(response, 0x00, 0x00) // Authority RRs
	response = append(response, 0x00, 0x00) // Additional RRs

	// Question section (echo back)
	response = append(response, encodeDomainName(question.Name)...)
	response = append(response, byte(question.Type>>8), byte(question.Type))
	response = append(response, byte(question.Class>>8), byte(question.Class))

	// Answer section
	if exists && question.Type == 1 {
		// Name (pointer to question)
		response = append(response, 0xC0, 0x0C)

		// Type (A record)
		response = append(response, 0x00, 0x01)

		// Class (IN)
		response = append(response, 0x00, 0x01)

		// TTL (300 seconds)
		response = append(response, 0x00, 0x00, 0x01, 0x2C)

		// Data length (4 bytes for IPv4)
		response = append(response, 0x00, 0x04)

		// IP address
		ipBytes := net.ParseIP(ip).To4()
		if ipBytes != nil {
			response = append(response, ipBytes...)
			fmt.Printf("Responding with: %s -> %s\n", question.Name, ip)
		}
	} else {
		fmt.Printf("No record found for: %s\n", question.Name)
	}

	return response
}

func encodeDomainName(name string) []byte {
	if name == "" {
		return []byte{0}
	}

	parts := strings.Split(name, ".")
	encoded := make([]byte, 0, len(name)+2)

	for _, part := range parts {
		if len(part) > 0 {
			encoded = append(encoded, byte(len(part)))
			encoded = append(encoded, []byte(part)...)
		}
	}
	encoded = append(encoded, 0) // Null terminator

	return encoded
}
