package ipam

import (
	"net/http"

	"github.com/docker/go-plugins-helpers/sdk"
)

const (
	manifest                    = `{"Implements": ["IpamDriver"]}`
	capabilitiesPath            = "/IpamDriver.GetCapabilities"
	requestPoolPath             = "/IpamDriver.RequestPool"
	releasePoolPath             = "/IpamDriver.ReleasePool"
	requestAddressPath          = "/IpamDriver.RequestAddress"
	releaseAddressPath          = "/IpamDriver.ReleaseAddress"
	getDefaultAddressSpacesPath = "/IpamDriver.GetDefaultAddressSpaces"
)

// Driver represent the interface a driver must fulfill.
type Driver interface {
	Capabilities() (*CapabilitiesResponse, error)
	RequestPool(*CreateIpamPoolRequest) (*CreateIpamPoolResponse, error)
	ReleasePool(*ReleaseIpamPoolRequest) error
	RequestAddress(*CreateIpamAddressRequest) (*CreateIpamAddressResponse, error)
	ReleaseAddress(*ReleaseIpamAddressRequest) error
	GetDefaultAddressSpaces() (*GetAddressSpacesResponse, error)
}

// CreateIpamPoolRequest returns an address pool along with its unique id.
type CreateIpamPoolRequest struct {
	AddressSpace string
	Pool         string
	SubPool      string
	Options      map[string]string
	V6           bool
}

// CreateIpamPoolResponse the response from CreateIpamPoolRequest
type CreateIpamPoolResponse struct {
	PoolID string
	Pool   string // CIDR format
	Data   map[string]string
}

// ReleaseIpamPoolRequest releases the address pool identified by the passed id
type ReleaseIpamPoolRequest struct {
	PoolID string
}

// CreateIpamAddressRequest address from the specified pool ID. Input options or preferred IP can be passed.
type CreateIpamAddressRequest struct {
	PoolID  string
	Address string
	Options map[string]string
}

// CreateIpamAddressResponse the response from CreateIpamAddressRequest
type CreateIpamAddressResponse struct {
	Address string // in CIDR format
	Data    map[string]string
}

// ReleaseIpamAddressRequest release the address from the specified pool ID
type ReleaseIpamAddressRequest struct {
	PoolID  string
	Address string
}

// GetAddressSpacesResponse returns the default local and global address spaces for this ipam
type GetAddressSpacesResponse struct {
	LocalDefaultAddressSpace  string
	GlobalDefaultAddressSpace string
}

// ErrorResponse is a formatted error message that libnetwork can understand
type ErrorResponse struct {
	Err string
}

// NewErrorResponse creates an ErrorResponse with the provided message
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{Err: msg}
}

// Handler forwards requests and responses between the docker daemon and the plugin.
type Handler struct {
	driver Driver
	sdk.Handler
}

// NewHandler initializes the request handler with a driver implementation.
func NewHandler(driver Driver) *Handler {
	h := &Handler{driver, sdk.NewHandler(manifest)}
	h.initMux()
	return h
}

// CapabilitiesResponse response to Capabilities request
type CapabilitiesResponse struct {
	RequiresMacAddress bool
}

func (h *Handler) initMux() {
	h.HandleFunc(capabilitiesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.driver.Capabilities()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(requestPoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &CreateIpamPoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.RequestPool(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(releasePoolPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ReleaseIpamPoolRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.ReleasePool(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(requestAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &CreateIpamAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		res, err := h.driver.RequestAddress(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
	h.HandleFunc(releaseAddressPath, func(w http.ResponseWriter, r *http.Request) {
		req := &ReleaseIpamAddressRequest{}
		err := sdk.DecodeRequest(w, r, req)
		if err != nil {
			return
		}
		err = h.driver.ReleaseAddress(req)
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
			return
		}
		sdk.EncodeResponse(w, make(map[string]string), "")
	})
	h.HandleFunc(getDefaultAddressSpacesPath, func(w http.ResponseWriter, r *http.Request) {
		res, err := h.driver.GetDefaultAddressSpaces()
		if err != nil {
			msg := err.Error()
			sdk.EncodeResponse(w, NewErrorResponse(msg), msg)
		}
		sdk.EncodeResponse(w, res, "")
	})
}
