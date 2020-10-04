package rest

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/peterbourgon/ff/v3"
)

//AppParams banco e listenAddr
type AppParams struct {
	ListenAddr string
	CpuProf    bool
	MemProf    bool
}

func cpuProf() func() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Fatal(err)
	}

	stopProfilingFunc := func() {
		f.Close()
		pprof.StopCPUProfile()
	}
	return stopProfilingFunc
}

func memProf() func() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}

	stopProfilingFunc := func() {
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	return stopProfilingFunc
}

//Setup flags de CLI e conexao banco
func Setup() (AppParams, error) {
	fs := flag.NewFlagSet("productList", flag.ExitOnError)
	var (
		listenAdr = fs.String("listen-addr", ":8000", "Hospedare listen addr")
		cpuprof   = fs.Bool("cpuprof", false, "Enable to perform cpu profiling saving on cpu.prof")
		memprof   = fs.Bool("memprof", false, "Enable to perform memory profiling saving on mem.prof")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("PRODUCTLIST"),
		//ff.WithConfigFile("product.conf"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		log.Println("FlagsFirst err:", err)
	}

	return AppParams{*listenAdr, *cpuprof, *memprof}, nil
}

//Listen inicia serviço
func Listen() error {
	params, err := Setup()
	if err != nil {
		return err
	}

	if params.CpuProf {
		stopProfFn := cpuProf()
		defer stopProfFn()
	}

	if params.MemProf {
		stopProfFn := memProf()
		defer stopProfFn()
	}

	router := NewRouter(params)

	srv := &http.Server{
		Handler:      router,
		Addr:         params.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err = srv.ListenAndServe()
	}()
	<-done

	return err
}
