package net

import (
	"errors"
	"net"
	"sync"
)

type ConnContainer struct {
	Conns []net.Conn
	mutex sync.Mutex
}

func NewConnContainer() *ConnContainer {
	c := new(ConnContainer)
	c.init()
	return c
}

func (c *ConnContainer) init() {
	c.Conns = make([]net.Conn, 0)
}

func (c *ConnContainer) Count() int {
	return len(c.Conns)
}

func (c *ConnContainer) Dial(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	c.Conns = append(c.Conns, conn)
	c.mutex.Unlock()
	return nil
}

func (c *ConnContainer) Close(index int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if index == -1 {
		for _, conn := range c.Conns {
			if err := conn.Close(); err != nil {
				return err
			}
		}
		c.init()
		return nil
	}
	if index < 0 || index >= len(c.Conns) {
		return errors.New("Index ouf of bounds")
	}
	conn := c.Conns[index]
	if err := conn.Close(); err != nil {
		return err
	}
	c.Conns = append(c.Conns[:index], c.Conns[index+1:]...)
	return nil
}
