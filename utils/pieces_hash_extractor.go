package utils

import ()

func PiecesHashExtractor(info map[string]any) [][]byte {
	pieces, ok := info["pieces"].([]byte)
	if !ok {
		println("Error: Pieces info invalid")
		return nil
	}

	var piecesHash [][]byte
	for i := 0; i < len(pieces); {
		curr_piece := pieces[i : i+20]
		piecesHash = append(piecesHash, curr_piece)
		i += 20
	}

	return piecesHash
}
