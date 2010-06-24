include $(GOROOT)/src/Make.$(GOARCH)

TARG     = deferred
GOFILES = deferred.go

include $(GOROOT)/src/Make.pkg
