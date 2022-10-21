package csdk

// #cgo darwin,arm64 LDFLAGS: -L/usr/local/lib/bcos-c-sdk/libs/darwin -lbcos-c-sdk-arm64
// #cgo darwin,amd64 LDFLAGS: -L/usr/local/lib/bcos-c-sdk/libs/darwin -lbcos-c-sdk-x86_64
// #cgo linux,amd64 LDFLAGS: -L/usr/local/lib/bcos-c-sdk/libs/linux -lbcos-c-sdk
// #cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs/win -lbcos-c-sdk
// #cgo CFLAGS: -I./
// #include "../../../bcos-c-sdk/bcos_sdk_c.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_error.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_rpc.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_uti_tx.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_common.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_amop.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_event_sub.h"
// #include "../../../bcos-c-sdk/bcos_sdk_c_uti_keypair.h"
// void on_recv_resp_callback(struct bcos_sdk_c_struct_response *);
// void on_recv_event_resp_callback(struct bcos_sdk_c_struct_response *);
// void on_recv_amop_publish_resp(struct bcos_sdk_c_struct_response *);
// void on_recv_amop_subscribe_resp(char* ,char* , struct bcos_sdk_c_struct_response *);
// void on_recv_notify_resp_callback(char* , int64_t , void* );
import "C"

import (
	"github.com/sirupsen/logrus"
	"strconv"
	"unsafe"
	"errors"
)

type CSDK struct {
	Sdk           unsafe.Pointer
	SMCrypto      C.int
	WASM          bool
	ChainID       *C.char
	GroupID       *C.char
	NodeID        *C.char
	PrivateKey    *C.char
	PrivateKeyLen C.uint
	Callback      *C.bcos_sdk_c_struct_response_cb
}

type ChanData struct {
	Sdk  unsafe.Pointer
	Data chan string
}

const (
	timeout = 1000
	proof   = 1 //是否返回交易回执证明, 0：不返回，1：返回
)


func CreateSignedTransaction(groupId string, chainId string, to string, data string, privateKey string, blockNumber int64) (error, string) {
    cPrivateKey := C.CString(privateKey)
	cPrivateKeyLen := C.uint(len(privateKey))
    cBlockNumber := C.int64_t(blockNumber)
	cTo := C.CString(to)
	cData := C.CString(data)
	cNull := C.CString("")
	cGroupId := C.CString(groupId)
	cChainId := C.CString(chainId)
	var tx_hash *C.char
	var signed_tx *C.char
	defer C.free(unsafe.Pointer(cTo))
	defer C.free(unsafe.Pointer(cData))
	defer C.free(unsafe.Pointer(cNull))
	defer C.free(unsafe.Pointer(tx_hash))
	defer C.free(unsafe.Pointer(signed_tx))
	defer C.free(unsafe.Pointer(cPrivateKey))
	defer C.free(unsafe.Pointer(cGroupId))
	defer C.free(unsafe.Pointer(cChainId))

	key_pair := C.bcos_sdk_create_keypair_by_private_key(0, unsafe.Pointer(cPrivateKey), cPrivateKeyLen)

	C.bcos_sdk_create_signed_transaction(key_pair, cGroupId, cChainId, cTo, cData, cNull, cBlockNumber, 0, &tx_hash, &signed_tx)

	if C.bcos_sdk_is_last_opr_success() == 0 {
	    errMsg := C.GoString(C.bcos_sdk_get_last_error_msg())
		logrus.Errorf("bcos_sdk_create_signed_transaction_with_signed_data failed, error: %s\n", errMsg)
		C.bcos_sdk_destroy_keypair(key_pair)
		return errors.New(errMsg), ""
	}

	C.bcos_sdk_destroy_keypair(key_pair)
	return nil, C.GoString(&signed_tx)
}
