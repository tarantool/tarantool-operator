package reconciliation

type CartridgeConfigController interface {
	Controller
}

type CommonCartridgeConfigController struct {
	*CommonController
}
