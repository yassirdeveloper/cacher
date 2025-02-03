// command_handler.go

package main

type Result[T any] interface {
	String() string
}

type result[T any] struct {
	data T
}

type ExecutableCommand[K comparable, V any] interface {
	Command
	Run(input CommandInput, cache Cache[K, V]) (Result[V], Error)
}

type Executor[K comparable, V any] interface {
	Execute(ExecutableCommand[K, V], CommandInput) (Result[V], Error)
}

type executor[K comparable, V any] struct {
	cacheManager CacheManager[K, V]
}

func NewExecutor[K comparable, V any](cacheManager CacheManager[K, V]) Executor[K, V] {
	return &executor[K, V]{
		cacheManager: cacheManager,
	}
}

func (ch *executor[K, V]) Execute(command ExecutableCommand[K, V], input CommandInput) (Result[V], Error) {
	// Get the appropriate cache from the CacheManager
	frequentAccessOption := input.GetOption(*FrequentAccessOption)
	useSyncCache := frequentAccessOption != nil
	cache := ch.cacheManager.Get(useSyncCache)

	return command.Run(input, cache)
}
