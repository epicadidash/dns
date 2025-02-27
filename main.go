package main

// import (
// 	"encoding/binary"
// 	"net"
//  "fmt"
//   "strings"
// )
// func main() {
//   // Bind to UDP port 53
//   conn, _ := net.ListenUDP("udp", &net.UDPAddr{Port: 53})
//   defer conn.Close()

//   // Read incoming requests
//   for {
//     buf := make([]byte, 512)
//     n, addr, _ := conn.ReadFromUDP(buf)
//     go handleRequest(conn, addr, buf[:n]) // Handle in a goroutine
//   }
// }
// func parseHeader(data []byte) (id uint16, flags uint16, qdcount uint16) {
//   id = binary.BigEndian.Uint16(data[0:2])
//   flags = binary.BigEndian.Uint16(data[2:4])
//   qdcount = binary.BigEndian.Uint16(data[4:6])
//   return
// }

// func parseQuestion(data []byte, offset int) (name string, offset int) {
//   // Decode domain name (e.g., 3www7example3com0 -> "www.example.com")
//   for {
//     length := int(data[offset])
//     if length == 0 { break }
//     name += string(data[offset+1:offset+1+length]) + "."
//     offset += length + 1
//   }
//   return name[:len(name)-1], offset + 4 // Skip QTYPE/QCLASS
// }
// func handleRequest(conn *net.UDPConn, addr *net.UDPAddr, data []byte) {
//   id, flags, qdcount := parseHeader(data)
//   name, _ := parseQuestion(data, 12) // Start after header (12 bytes)

//   // Hardcoded response for "example.com"
//   if name == "example.com" {
//     buildResponse(conn, addr, id, name, "192.0.2.1")
//   }
// }
// func buildResponse(conn *net.UDPConn, addr *net.UDPAddr, id uint16, name string, ip string) {
//   // Header (ID, Flags, QDCOUNT=1, ANCOUNT=1)
//   response := make([]byte, 12)
//   binary.BigEndian.PutUint16(response[0:2], id)
//   binary.BigEndian.PutUint16(response[2:4], 0x8180) // QR=1, AA=1
//   binary.BigEndian.PutUint16(response[4:6], 1)       // QDCOUNT=1
//   binary.BigEndian.PutUint16(response[6:8], 1)       // ANCOUNT=1

//   // Question Section (name, A record, IN class)
//   response = append(response, encodeName(name)...)
//   response = append(response, []byte{0x00, 0x01, 0x00, 0x01}...)

//   // Answer Section (name, A record, TTL=300, IP)
//   response = append(response, encodeName(name)...)
//   response = append(response, []byte{
//     0x00, 0x01,        // TYPE=A
//     0x00, 0x01,        // CLASS=IN
//     0x00, 0x00, 0x01, 0x2c, // TTL=300
//     0x00, 0x04,        // RDLENGTH=4 (IPv4)
//   }...)
//   response = append(response, net.ParseIP(ip).To4()...)

//   conn.WriteToUDP(response, addr)
// }

// // encodeName converts a domain name to DNS-encoded format.
// // Example: "example.com" → []byte{7, 'e','x','a','m','p','l','e', 3, 'c','o','m', 0}
// func encodeName(name string) ([]byte, error) {
//   // Trim trailing dot if present (e.g., "example.com." → "example.com")
//   name = strings.TrimSuffix(name, ".")
//   if name == "" {
//     // Special case: root domain "." → single 0-byte
//     return []byte{0}, nil
//   }

//   var encoded []byte
//   labels := strings.Split(name, ".")
//   for _, label := range labels {
//     if len(label) == 0 {
//       return nil, fmt.Errorf("empty label in domain name")
//     }
//     if len(label) > 63 {
//       return nil, fmt.Errorf("label '%s' exceeds 63 characters", label)
//     }
//     encoded = append(encoded, byte(len(label)))
//     encoded = append(encoded, label...)
//   }
//   encoded = append(encoded, 0) // Null terminator
//   return encoded, nil
// }
