package telegram

import (
	"container/list"
	"sync"
)

// MessageStore 是一个内存中的消息队列实现，用于替代 Redis 的 BLPOP 和 RPUSH 操作
type MessageStore struct {
	queue    *list.List
	mutex    sync.Mutex
	cond     *sync.Cond
	closed   bool
	closedMu sync.RWMutex
}

// NewMessageStore 创建一个新的消息存储实例
func NewMessageStore() *MessageStore {
	ms := &MessageStore{
		queue: list.New(),
	}
	ms.cond = sync.NewCond(&ms.mutex)
	return ms
}

// RPush 将消息推入队列的右侧（尾部），相当于 Redis 的 RPUSH 命令
func (ms *MessageStore) RPush(value string) error {
	ms.closedMu.RLock()
	defer ms.closedMu.RUnlock()

	if ms.closed {
		return &MessageStoreClosedError{}
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.queue.PushBack(value)
	ms.cond.Signal() // 通知等待的消费者有新消息
	return nil
}

// BLPop 阻塞式地从队列左侧（头部）弹出消息，相当于 Redis 的 BLPOP 命令
func (ms *MessageStore) BLPop() (string, error) {
	ms.closedMu.RLock()
	if ms.closed {
		ms.closedMu.RUnlock()
		return "", &MessageStoreClosedError{}
	}
	ms.closedMu.RUnlock()

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// 如果队列为空，等待直到超时或有新元素
	if ms.queue.Len() == 0 {

		// 等待直到有元素或者超时
		for ms.queue.Len() == 0 {
			ms.cond.Wait()

			// 检查是否已经关闭
			ms.closedMu.RLock()
			isClosed := ms.closed
			ms.closedMu.RUnlock()
			if isClosed {
				return "", &MessageStoreClosedError{}
			}
		}

		// 如果超时了仍然没有元素，则返回nil
		if ms.queue.Len() == 0 {
			return "", nil // Redis BLPOP 在超时时返回 nil
		}
	}

	// 弹出第一个元素
	element := ms.queue.Front()
	if element == nil {
		return "", nil
	}

	ms.queue.Remove(element)
	value := element.Value.(string)

	// 返回格式与 Redis BLPOP 一致: [key, value]
	return value, nil
}

// Close 关闭消息存储，释放资源并唤醒所有阻塞的操作
func (ms *MessageStore) Close() error {
	ms.closedMu.Lock()
	ms.closed = true
	ms.closedMu.Unlock()

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.cond.Broadcast() // 唤醒所有等待的协程
	return nil
}

// Size 返回队列中消息的数量
func (ms *MessageStore) Size() int {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	return ms.queue.Len()
}

// MessageStoreClosedError 表示消息存储已关闭的错误
type MessageStoreClosedError struct{}

func (e *MessageStoreClosedError) Error() string {
	return "message store is closed"
}
