package rest

import (
	"context"
	"flag"
	"fmt"
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
	ListenAddr   string
	DiscountAddr string
	CpuProf      bool
	MemProf      bool
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
		listenAdr    = fs.String("listen-addr", ":8000", "Hospedare listen addr:port")
		discountAddr = fs.String("discount-addr", ":50051", "Discount service addr:port")
		cpuprof      = fs.Bool("cpuprof", false, "Enable to perform cpu profiling saving on cpu.prof")
		memprof      = fs.Bool("memprof", false, "Enable to perform memory profiling saving on mem.prof")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarPrefix("PRODUCTLIST"),
		//ff.WithConfigFile("product.conf"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		log.Println("FlagsFirst err:", err)
	}

	return AppParams{*listenAdr, *discountAddr, *cpuprof, *memprof}, nil
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Handler:      router,
		Addr:         params.ListenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe ", err.Error())
		}
	}()

	log.Printf(fmt.Sprintf("Server started %s discount-service %s", params.ListenAddr, params.DiscountAddr))

	<-done
	log.Printf("Sig received, shutdown http server, 5s timeout...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v\n", err)
	} else {
		log.Printf("gracefully stopped\n")
	}

	return err
}
