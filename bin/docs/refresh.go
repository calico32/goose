package main

import (
	"net/rpc"

	"github.com/labstack/echo/v4"
)

type PageRefresh struct {
	conn   *rpc.Client
	logger echo.Logger
}

func (r *PageRefresh) Refresh() {
	if r.conn != nil {
		err := r.conn.Call("PageRefreshRPC.Refresh", struct{}{}, nil)
		if err != nil {
			r.logger.Error(err)
		} else {
			r.logger.Info("refreshed browser")
		}
	} else if mode != "release" {
		r.logger.Warn("rpc connection not available (is the page refresh extension running?)")
	}
}

func (r *PageRefresh) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}

func NewPageRefresh(logger echo.Logger) *PageRefresh {
	return &PageRefresh{conn: setupRPC(logger), logger: logger}
}

func setupRPC(logger echo.Logger) *rpc.Client {
	conn, err := rpc.DialHTTP("tcp", "localhost:41923")
	if err != nil {
		logger.Warnf("could not connect to rpc server: %v", err)
		return nil
	}
	return conn
}
