package url

import (
    //"fmt"
    "net/http"
    "net/http/httptrace"
    "time"
    "crypto/tls"
    "io"
    "io/ioutil"
    "github.com/valyala/fasthttp"
)

type UrlStats struct {
	Url string
	StatusCode int
	Error string
    DNS time.Duration
	TLSHandshake time.Duration
	ConnectTime time.Duration
	FirstByteTime time.Duration
	TotalTime time.Duration
}

func GetUrlMetrics(url string) *UrlStats{
	//fmt.Printf("URL ============================= : %v\n", url)
	var urlmetrics UrlStats
	urlmetrics.Url = url
    req, _ := http.NewRequest("GET", url, nil)
	//var dns time.Time
    var start, connect, tlsHandshake time.Time
    
    timeout := time.Duration(30 * time.Second)
    
    trTls11 := &http.Transport{
		DisableKeepAlives:true,
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS11,
			InsecureSkipVerify: true,
		},
	}
	
	client := http.Client{
		Transport: trTls11,
		Timeout:   timeout,
	}


    trace := &httptrace.ClientTrace{
//         DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
//         DNSDone: func(ddi httptrace.DNSDoneInfo) {
//             //fmt.Printf("DNS Done: %v\n", time.Since(dns))
//             urlmetrics.DNS = time.Since(dns)
//         },

        TLSHandshakeStart: func() { tlsHandshake = time.Now() },
        TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
            //fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
            urlmetrics.TLSHandshake = time.Since(tlsHandshake)
        },

        ConnectStart: func(network, addr string) { connect = time.Now() },
        ConnectDone: func(network, addr string, err error) {
            //fmt.Printf("Connect time: %v\n", time.Since(connect))
            urlmetrics.ConnectTime = time.Since(connect)
        },

        GotFirstResponseByte: func() {
            //fmt.Printf("Time from start to first byte: %v\n", time.Since(start))
            urlmetrics.FirstByteTime = time.Since(start)
        },
    }

    req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
    req.Header.Set("User-Agent", "Bot/3.0")
    start = time.Now()
    /*
    if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
        log.Fatal(err)
    }
    */
    /*
    if err == nil && (resp.StatusCode == 301 || resp.StatusCode == 302) {
		headers := resp.Header
		fmt.Printf("headers: %v\n", headers)
		
	}
	*/
	resp, err := client.Do(req)
	if err != nil{
		urlmetrics.Error = err.Error()
		
	} else{
		urlmetrics.StatusCode = resp.StatusCode
		io.Copy(ioutil.Discard, resp.Body)
	}
	if resp != nil {
	    resp.Body.Close() // MUST CLOSED THIS
	} 
	
	urlmetrics.TotalTime = time.Since(start)
    return &urlmetrics
}

func GetUrlMetricsFastHttp(websiteJob *WebsiteJob) *UrlStats{
	//fmt.Printf("URL ============================= : %v\n", url)
	var url string
	var urlmetrics UrlStats
	url = websiteJob.Website
	urlmetrics.Url = url
	//var dns time.Time
    //var start, connect, tlsHandshake time.Time
    
    timeout := time.Duration(10 * time.Second)

    
    start := time.Now()
    
    req := fasthttp.AcquireRequest()
    req.SetRequestURI(url)
    req.Header.Add("User-Agent", "Bot/3.0")

    
    req.Header.SetMethod("GET")

    resp := fasthttp.AcquireResponse()
    client := fasthttp.Client{
		TLSConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
			InsecureSkipVerify: true,
		},
		ReadTimeout:   timeout,
	}
    err := client.Do(req, resp);
    
    if err != nil {
        urlmetrics.Error = err.Error()
    } else {
    	urlmetrics.StatusCode = resp.StatusCode()
    	/*
    	if(resp.StatusCode() == 301){
    		fmt.Println("URL ============================= : ", resp.StatusCode(),string(resp.Header.Peek("Location")))
    	}else if(resp.StatusCode() == 302){
    		fmt.Println("URL ============================= : ", resp.StatusCode(),string(resp.Header.Peek("Location")))
    	}else if(resp.StatusCode() == 400){
    		fmt.Println("URL ============================= : ", resp.StatusCode())
    	}
    	*/
    	
    }
	ParseFastHttpResponse(err, resp, &urlmetrics, websiteJob)
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	urlmetrics.TotalTime = time.Since(start)
    return &urlmetrics
}