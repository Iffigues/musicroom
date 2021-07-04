package config

import (
	"errors"
)

type ConfState map[string]interface{}

type ConfType map[string]ConfState

type Conf struct {
	a ConfType;
}

func (c *Conf)NewConfType(a string, override bool) (err error) {
	if _, ok := c.a[a];  !ok{
		c.a[a] = make(map[string]interface{})
		return nil;
	}
	if !override {
		return errors.New("key already exists");
	}
	c.a[a] = make(map[string]interface{})
	return
}


func NewConf() (c *Conf) {
	c = new(Conf)
	c.a = make(map[string]ConfState)
	return
}

func (c *Conf)AddState(types string, name string, value interface{}, o bool) (err error) {
	if _, ok := c.a[types]; ok {
		if _, yes := c.a[types][name]; !yes {
			c.a[types][name] = value;
			return nil
		}
		if o {
			c.a[types][name] = value
			return nil
		}
		return errors.New("key already exists")
	}
	return errors.New("conf type don't exists")
}

func (c *Conf) GetValue(a string, key string) (cc interface{})  {
	if z, ok := c.a[a]; ok {
		if zz, yes := z[key]; yes {
			return zz
		}
		return nil
	}
	return nil
}
