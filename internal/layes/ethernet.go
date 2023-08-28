package layes

import "fmt"

type EthernetLayer struct {
	Origem  [6]byte
	Destino [6]byte
	Tipo    [2]byte
	Data    []byte
}

// type ClearEthernet struct {
// 	Origem  string
// 	Destino string
// 	Tipo    string
// 	Data    []byte
// }

// func WrapEthernet() {

// }

func UnWrapEthernet(bytes *[]byte) EthernetLayer {
	var raw EthernetLayer
	copy(raw.Destino[:], (*bytes)[0:6])
	copy(raw.Origem[:], (*bytes)[6:12])
	copy(raw.Tipo[:], (*bytes)[12:14])
	raw.Data = (*bytes)[14:]
	return raw
}

func (e EthernetLayer) String() string {
	str := fmt.Sprintf("MAC Destino: %x:%x:%x:%x:%x:%x \n", e.Destino[0], e.Destino[1], e.Destino[2], e.Destino[3], e.Destino[4], e.Destino[5])
	str += fmt.Sprintf("MAC Origen: %x:%x:%x:%x:%x:%x \n", e.Origem[0], e.Origem[1], e.Origem[2], e.Origem[3], e.Origem[4], e.Origem[5])
	str += fmt.Sprintf("Tipo: %x \n", e.Tipo)
	str += fmt.Sprintf("Data: %d Bytes \n\n", len(e.Data))

	return str
}
