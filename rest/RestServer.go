package rest

//
// parts of listener pattern adapted from https://github.com/hydrogen18/stoppableListener
//
// Copyright (c) 2014, Eric Urban
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// End of Eric Urban Copyright
//

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/ninjasphere/app-scheduler/service"
	"github.com/ninjasphere/go-ninja/config"
	"github.com/ninjasphere/go-ninja/logger"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type listener struct {
	*net.TCPListener
	stop chan struct{}
}

// RestServer Holds stuff shared by all the rest services
type RestServer struct {
	Scheduler *service.SchedulerService
	log       *logger.Logger
	listener
	// TaskModel *TaskModel
}

var StoppedError = errors.New("Listener stopped")

func (ln listener) Accept() (c net.Conn, err error) {
	for {
		ln.SetDeadline(time.Now().Add(time.Second))
		tc, err := ln.AcceptTCP()

		select {
		case <-ln.stop:
			return nil, StoppedError
		default:
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			} else {
				return nil, err
			}
		}

		tc.SetKeepAlive(true)
		tc.SetKeepAlivePeriod(3 * time.Minute)
		return tc, nil
	}
}

func (r *RestServer) Start() error {
	r.log = logger.GetLogger("app-scheduler.rest")
	return r.listen()
}

func (r *RestServer) Stop() {
	close(r.stop)
	return
}

func (r *RestServer) listen() error {

	m := martini.Classic()

	m.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))

	// m.Map(r.TaskModel)

	task := NewTaskRouter()
	task.scheduler = r.Scheduler

	m.Group("/rest/v1/tasks", task.Register)

	listenAddress := fmt.Sprintf(":%d", config.MustInt("app-scheduler.rest.port"))

	r.log.Infof("Listening at %s", listenAddress)

	srv := &http.Server{Addr: listenAddress, Handler: m}
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return err
	}

	return srv.Serve(listener{
		TCPListener: ln.(*net.TCPListener),
		stop:        make(chan struct{}),
	})
}
