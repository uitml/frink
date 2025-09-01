package cmd
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/spf13/cobra"
	"github.com/uitml/frink/internal/cli"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type gpuContext struct {
	cli.CommandContext
	ServiceURL string        // URL to the gpu-viewer service; if empty, use API server proxy
	Format     string        // "table" or "oneline"
	Color      bool          // enable ANSI colors
	Limit      int           // only used for oneline
	Timeout    time.Duration // HTTP timeout
}
func newGPUCmd() *cobra.Command {
	ctx := &gpuContext{
		Format:  "table",
		Timeout: 5 * time.Second,
	}
	cmd := &cobra.Command{
		Use:   "gpu",
		Short: "Show cluster GPU availability",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return ctx.Initialize(cmd)
		},
		RunE: ctx.Run,
	}
	flags := cmd.Flags()
	flags.StringVar(&ctx.ServiceURL, "service-url", "", "GPU viewer service URL (env GPU_VIEWER_URL or API proxy if empty)")
	flags.StringVar(&ctx.Format, "format", "table", "Output format: table|oneline")
	flags.BoolVar(&ctx.Color, "color", false, "Enable ANSI colors")
	flags.IntVar(&ctx.Limit, "limit", 0, "Limit number of nodes (only for --format=oneline)")
	flags.DurationVar(&ctx.Timeout, "timeout", 5*time.Second, "HTTP request timeout")
	return cmd
}
func (ctx *gpuContext) Run(cmd *cobra.Command, args []string) error {
	// 1) If user provided a direct URL (flag or env), use it (works with port-forward or Ingress)
	if ctx.ServiceURL != "" || os.Getenv("GPU_VIEWER_URL") != "" {
		serviceURL := ctx.ServiceURL
		if serviceURL == "" {
			serviceURL = os.Getenv("GPU_VIEWER_URL")
		}
		return fetchDirect(cmd, serviceURL, ctx.Format, ctx.Color, ctx.Limit, ctx.Timeout)
	}
	// 2) Otherwise, go through the API server service proxy using kubeconfig
	return fetchViaAPIServerProxy(cmd, ctx.Format, ctx.Color, ctx.Limit, ctx.Timeout)
}
// Direct HTTP call (port-forward or Ingress)
func fetchDirect(cmd *cobra.Command, baseURL, format string, color bool, limit int, timeout time.Duration) error {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return fmt.Errorf("invalid service-url: %w", err)
	}
	q := u.Query()
	switch strings.ToLower(format) {
	case "oneline":
		q.Set("format", "oneline")
		if limit > 0 {
			q.Set("limit", strconv.Itoa(limit))
		}
	case "table", "":
		// default
	default:
		return fmt.Errorf("unknown format %q (use table or oneline)", format)
	}
	if color {
		q.Set("color", "1")
	}
	u.RawQuery = q.Encode()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("service error: %s: %s", resp.Status, string(b))
	}
	_, err = io.Copy(cmd.OutOrStdout(), resp.Body)
	return err
}
// API server service proxy (no port-forward; requires RBAC to get services/proxy)
func fetchViaAPIServerProxy(cmd *cobra.Command, format string, color bool, limit int, timeout time.Duration) error {
	// Pick up the same --context flag your root command defines
	ctxFlag, _ := cmd.InheritedFlags().GetString("context")
	// Load kubeconfig with overrides (so --context is honored)
	loading := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	if ctxFlag != "" {
	  overrides.CurrentContext = ctxFlag
	}
	cfg, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loading, overrides).ClientConfig()
	if err != nil {
	  return fmt.Errorf("kubeconfig load failed: %w", err)
	}
	// Authenticated transport to the API server
	transport, err := rest.TransportFor(cfg)
	if err != nil {
	  return fmt.Errorf("transport build failed: %w", err)
	}
	host := strings.TrimRight(cfg.Host, "/")
	ns := "gpu-availability"
	svc := "gpu-viewer"
	scheme := "http"
	portName := "http"
	// proxy URL that WORKS in your cluster:
	proxyURL := fmt.Sprintf("%s/api/v1/namespaces/%s/services/%s:%s:%s/proxy/",
	  host, ns, scheme, svc, portName,
	)
	// Query params
	u, _ := url.Parse(proxyURL)
	q := u.Query()
	switch strings.ToLower(format) {
	case "oneline":
	  q.Set("format", "oneline")
	  if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	  }
	case "table", "":
	  // default
	default:
	  return fmt.Errorf("unknown format %q (use table or oneline)", format)
	}
	if color {
	  q.Set("color", "1")
	}
	u.RawQuery = q.Encode()
	// Do request via API server proxy
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, u.String(), nil)
	client := &http.Client{Transport: transport, Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
	  return fmt.Errorf("proxy request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
	  b, _ := io.ReadAll(resp.Body)
	  return fmt.Errorf("service/proxy error: %s: %s", resp.Status, string(b))
	}
	_, err = io.Copy(cmd.OutOrStdout(), resp.Body)
	return err
  }