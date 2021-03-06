package main

import "fmt"

// msr is used to parse and print information that comes from debug
// information.
type msr struct {
	V1      int `json:"auth.v1"`
	Authed  int `json:"auth.authenticated"`
	Plain   int `json:"auth.plain"`
	Expired int `json:"auth.expired"`
	Error   int `json:"auth.error"`
}

func (m msr) String() string {
	return fmt.Sprintf("v1: %8d, auth: %8d, plain: %8d, expired: %8d, error: %8d",
		m.V1,
		m.Authed,
		m.Plain,
		m.Expired,
		m.Error,
	)
}

// Add adds members from o into m.
func (m *msr) Add(o msr) {
	m.Authed += o.Authed
	m.Error += o.Error
	m.Expired += o.Expired
	m.Plain += o.Plain
	m.V1 += o.V1
}

// Delta returns an msr with the values from o subtracted from m.
func (m msr) Delta(o msr) msr {
	m.Authed -= o.Authed
	m.Error -= o.Error
	m.Expired -= o.Expired
	m.Plain -= o.Plain
	m.V1 -= o.V1
	return m
}

func (m msr) Positive() msr {
	if m.V1 < 0 {
		m.V1 = 0
	}
	if m.Authed < 0 {
		m.Authed = 0
	}
	if m.Plain < 0 {
		m.Plain = 0
	}
	if m.Expired < 0 {
		m.Expired = 0
	}
	if m.Error < 0 {
		m.Error = 0
	}
	return m
}

type msrErr struct {
	m msr
	e error
}
