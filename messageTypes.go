package main

type Message struct {
	// From    string
	Payload any
}

type MessageStoreFile struct {
	Key string
	ID  string
	// Ext  string
	Size      int64
	Signature []byte
}

type MessageGetFile struct {
	Key string
	ID  string
}

type MessageDeleteFile struct {
	Key string
	ID  string
}

type MessageGetFileNotFound struct {
	Key string
	ID  string
}

type MessageStoreAck struct {
	Key  string
	From string
	Err  string
}

type MessageDeleteAck struct {
	Key  string
	From string
	Err  string
}

type MessageDuplicateCheck struct {
	Key string
	ID  string
}

type MessageDuplicateResponse struct {
    Key   string
    ID    string
    HasIt bool
    From  string
}
