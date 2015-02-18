go run GSM.go sequence

go run GSM2csv.go/100GSM2csv.go testinput/sequence/ testinput/read.fq sequenceGSM ReadGSM

go run readGSM2csv.go testinput/read.fq testinput/sequence.csv ReadGSM

go run topK/topK.go testinput/sequence.csv topK
