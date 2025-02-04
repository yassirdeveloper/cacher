package main

// Command arguments
var (
	KeyCommandArgument = &commandArgument{label: "key", position: 0, valueType: TypeString, description: "a unique identifier for quickly storing and retrieving specific data"}
)

// Command options
var (
	FrequentAccessOption = &commandOption{label: "frequent access cache", letter: 'f', name: "frequent-access", valueType: NoType, description: "pass this option for frequently accessed values"}
	ExpirationOption     = &commandOption{label: "expiration time", letter: 'e', name: "expires-in", valueType: TypeInt, description: "period in seconds until the key value pair are deleted"}
)
