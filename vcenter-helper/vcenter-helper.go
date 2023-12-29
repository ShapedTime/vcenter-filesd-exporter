package vcenterhelper

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"log"
	"net/url"
	"time"
)

type VCenterHelper struct {
	ctx       context.Context
	_client   *vim25.Client
	lastLogin *time.Time
	_username string
	_password string
	_host     string
}

func NewVCenterHelper(host, user, password string) *VCenterHelper {
	ctx := context.Background()
	c, err := newClient(ctx, host, user, password)
	if err != nil {
		panic(err)
	}

	if !c.IsVC() {
		panic("not vc")
	}

	now := time.Now()

	return &VCenterHelper{ctx: ctx, _client: c,
		lastLogin: &now, _username: user,
		_password: password, _host: host}
}

func (vc *VCenterHelper) GetVMs(f object.Reference, ps []string) ([]mo.VirtualMachine, error) {
	c, err := vc.client()
	if err != nil {
		return nil, err
	}

	m := view.NewManager(c)

	kind := []string{"VirtualMachine"}
	v, err := m.CreateContainerView(vc.ctx, f.Reference(), kind, true)
	if err != nil {
		log.Printf("error creating container view: %s", err)
		return nil, err
	}
	defer v.Destroy(vc.ctx)

	var vmsR []mo.VirtualMachine

	err = v.Retrieve(vc.ctx, kind, ps, &vmsR)
	if err != nil {
		log.Printf("error retrieving vms: %s", err)
		return nil, err
	}

	return vmsR, nil
}

func (vc *VCenterHelper) FindFolder(datacenter, vmDir string) (object.Reference, error) {
	c, err := vc.client()
	if err != nil {
		return nil, err
	}

	f, err := object.NewSearchIndex(c).FindByInventoryPath(vc.ctx, fmt.Sprintf("/%s/vm/%s", datacenter, vmDir))
	if err != nil {
		log.Printf("error finding folder: %s", err)
		return nil, err
	}

	if f == nil {
		log.Printf("folder %s not found", vmDir)
		return nil, fmt.Errorf("folder %s not found", vmDir)
	}

	return f, nil
}

func (vc *VCenterHelper) GetFolders(folder object.Reference, ps []string) ([]mo.Folder, error) {
	c, err := vc.client()
	if err != nil {
		return nil, err
	}

	m := view.NewManager(c)

	kind := []string{"Folder"}

	v, err := m.CreateContainerView(vc.ctx, folder.Reference(), kind, true)
	if err != nil {
		log.Printf("error creating container view: %s", err)
		return nil, err
	}
	defer v.Destroy(vc.ctx)

	var foldersR []mo.Folder

	err = v.Retrieve(vc.ctx, kind, ps, &foldersR)
	if err != nil {
		log.Printf("error retrieving folders: %s", err)
		return nil, err
	}

	return foldersR, nil
}

// NewClient creates a vim25.Client for use in the examples
func newClient(ctx context.Context, host, user, password string) (*vim25.Client, error) {
	// Parse URL from string
	u, err := soap.ParseURL(host)
	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword(user, password)

	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (vc *VCenterHelper) client() (*vim25.Client, error) {
	if time.Now().Sub(*vc.lastLogin) > 5*time.Minute {
		log.Println("relogin")
		log.Println("last login", *vc.lastLogin)
		c, err := newClient(vc.ctx, vc._host, vc._username, vc._password)
		if err != nil {
			return nil, err
		}

		vc._client = c
		now := time.Now()

		vc.lastLogin = &now
	}

	return vc._client, nil
}
