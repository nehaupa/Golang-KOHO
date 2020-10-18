package velocitylimit

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)
func Process(inputFile string, outputFile string, db *Store) error {
	in, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer in.Close()
	defer out.Close()

	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)
	for scanner.Scan() {
		req, err := NewRequest(scanner.Text())
		if err != nil {
			return err
		}
		if db.IsDupTxn(req.ID, req.CustID) {
			log.Println("Ignoring duplicate txn: ", req.ID)
			continue
		}
		account := db.GetAccount(req.CustID)
		accepted := account.LoadFunds(req)
		response := NewResponse(req.ID, req.CustID, accepted)
		resBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		if _, err = writer.WriteString(string(resBytes) + "\n"); err != nil {
			return err
		}
		db.AddTxn(req.ID, req.CustID)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	writer.Flush()
	return nil
}
