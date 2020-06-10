package client

import (
	"context"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/proxy"
	"github.com/p4gefau1t/trojan-go/tunnel"
	"github.com/p4gefau1t/trojan-go/tunnel/mux"
	"github.com/p4gefau1t/trojan-go/tunnel/router"
	"github.com/p4gefau1t/trojan-go/tunnel/simplesocks"
	"github.com/p4gefau1t/trojan-go/tunnel/socks"
	"github.com/p4gefau1t/trojan-go/tunnel/transport"
	"github.com/p4gefau1t/trojan-go/tunnel/trojan"
	"github.com/p4gefau1t/trojan-go/tunnel/websocket"
)

const Name = "CLIENT"

// GenerateClientTree generate general outbound protocol stack
func GenerateClientTree(isMux bool, isWebsocket bool, isRouter bool) []string {
	clientStack := []string{transport.Name}
	if isWebsocket {
		clientStack = append(clientStack, websocket.Name)
	}
	clientStack = append(clientStack, trojan.Name)
	if isMux {
		clientStack = append(clientStack, []string{mux.Name, simplesocks.Name}...)
	}
	if isRouter {
		clientStack = append(clientStack, router.Name)
	}
	return clientStack
}

func init() {
	proxy.RegisterProxyCreator(Name, func(ctx context.Context) (*proxy.Proxy, error) {
		cfg := config.FromContext(ctx, Name).(*Config)
		serverStack := []string{socks.Name}
		clientStack := GenerateClientTree(cfg.Mux.Enabled, cfg.Websocket.Enabled, cfg.Router.Enabled)
		c, err := proxy.CreateClientStack(ctx, clientStack)
		if err != nil {
			return nil, err
		}
		s, err := proxy.CreateServerStack(ctx, serverStack)
		if err != nil {
			return nil, err
		}
		return proxy.NewProxy(ctx, []tunnel.Server{s}, c), nil
	})
}
