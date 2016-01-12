package ipam

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/docker/go-plugins-helpers/sdk"
)

type TestDriver struct {
	Driver
}

func (t *TestDriver) RequestPool(r *CreateIpamPoolRequest) (*CreateIpamPoolResponse, error) {
	return &CreateIpamPoolResponse{}, nil
}

func (t *TestDriver) ReleasePool(r *ReleaseIpamPoolRequest) error {
	return nil
}

func (t *TestDriver) ReleaseAddress(r *ReleaseIpamAddressRequest) error {
	return nil
}

func (t *TestDriver) RequestAddress(r *CreateIpamAddressRequest) (*CreateIpamAddressResponse, error) {
	return &CreateIpamAddressResponse{}, nil
}

func (t *TestDriver) GetDefaultAddressSpaces() (*GetAddressSpacesResponse, error) {
	return &GetAddressSpacesResponse{}, nil
}

type ErrDriver struct {
	Driver
}

func (e *ErrDriver) RequestPool(r *CreateIpamPoolRequest) (*CreateIpamPoolResponse, error) {
	return nil, errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) ReleasePool(r *ReleaseIpamPoolRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) RequestAddress(r *CreateIpamAddressRequest) (*CreateIpamAddressResponse, error) {
	return nil, errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) ReleaseAddress(r *ReleaseIpamAddressRequest) error {
	return errors.New("I CAN HAZ ERRORZ")
}

func (e *ErrDriver) GetDefaultAddressSpaces() (*GetAddressSpacesResponse, error) {
	return nil, errors.New("I CAN HAZ ERRORZ")
}

func TestMain(m *testing.M) {
	d := &TestDriver{}
	h1 := NewHandler(d)
	go h1.ServeTCP("test", ":8234")

	e := &ErrDriver{}
	h2 := NewHandler(e)
	go h2.ServeTCP("err", ":8567")

	m.Run()
}

func TestActivate(t *testing.T) {
	response, err := http.Get("http://localhost:8234/Plugin.Activate")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if string(body) != manifest+"\n" {
		t.Fatalf("Expected %s, got %s\n", manifest+"\n", string(body))
	}
}

func TestCreateIpamPoolSuccess(t *testing.T) {
	request := `{"AddressSpace":"172.18.0.1/16","Pool":"172.18.0.1/24"}`
	response, err := http.Post("http://localhost:8234/IpamDriver.RequestPool",
		sdk.DefaultContentTypeV1_1,
		strings.NewReader(request),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d\n", response.StatusCode)
	}
	if string(body) != "{\"PoolID\":\"\",\"Pool\":\"\",\"Data\":null}\n" {
		t.Fatalf("Expected %s, got %s\n", "{}\n", string(body))
	}
}
