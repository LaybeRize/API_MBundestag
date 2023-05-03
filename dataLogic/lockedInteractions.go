package dataLogic

import (
	"sync"
)

var documentLock = &sync.Mutex{}
var userLock = &sync.Mutex{}
var orgnisationLock = &sync.Mutex{}
var titleLock = &sync.Mutex{}
