// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"fmt"
	"testing"

	"github.com/tsuru/tsuru/event"
	"github.com/tsuru/tsuru/provision"

	check "gopkg.in/check.v1"
)

type S struct{}

type customPlatformBuilder struct {
	customBehavior func(PlatformOptions, string) error
}

var _ = check.Suite(S{})
var _ PlatformBuilder = &customPlatformBuilder{}
var _ Builder = &customPlatformBuilder{}

func (b *customPlatformBuilder) PlatformAdd(opts PlatformOptions) error {
	if b.customBehavior == nil {
		return nil
	}
	return b.customBehavior(opts, "")
}

func (b *customPlatformBuilder) PlatformUpdate(opts PlatformOptions) error {
	if b.customBehavior == nil {
		return nil
	}
	return b.customBehavior(opts, "")
}

func (b *customPlatformBuilder) PlatformRemove(name string) error {
	if b.customBehavior == nil {
		return nil
	}
	return b.customBehavior(PlatformOptions{}, name)
}

func (b *customPlatformBuilder) Build(provision.BuilderDeploy, provision.App, *event.Event, BuildOpts) (string, error) {
	return "", nil
}

func Test(t *testing.T) {
	check.TestingT(t)
}

func (s S) SetUpTest(c *check.C) {
	builders = make(map[string]Builder)
}

func (s S) TestRegisterAndGetBuilder(c *check.C) {
	var b Builder
	Register("my-builder", b)
	got, err := Get("my-builder")
	c.Assert(err, check.IsNil)
	c.Check(got, check.DeepEquals, b)
	_, err = Get("unknown-builder")
	c.Check(err, check.NotNil)
	expectedMessage := `unknown builder: "unknown-builder"`
	c.Assert(err.Error(), check.Equals, expectedMessage)
}

func (s S) TestGetDefaultBuilder(c *check.C) {
	var b1, b2 Builder
	DefaultBuilder = "default-builder"
	Register("default-builder", b1)
	Register("other-builder", b2)
	got, err := GetDefault()
	c.Check(err, check.IsNil)
	c.Check(got, check.DeepEquals, b1)
}

func (s S) TestRegistry(c *check.C) {
	var b1, b2, b3 Builder
	Register("my-builder", b1)
	Register("your-builder", b2)
	Register("default-builder", b3)
	builders, err := Registry()
	c.Assert(err, check.IsNil)
	c.Assert(builders, check.HasLen, 3)
}

func (s S) TestPlatformAdd(c *check.C) {
	platformAddCalls := 0
	callCounter := func(PlatformOptions, string) error {
		platformAddCalls++
		return nil
	}
	b1 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b2 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b3 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	Register("my-builder", &b1)
	Register("your-builder", &b2)
	Register("default-builder", &b3)
	err := PlatformAdd(PlatformOptions{})
	c.Assert(err, check.IsNil)
	c.Assert(platformAddCalls, check.Equals, 3)
}

func (s S) TestPlatformAddError(c *check.C) {
	errMsg := "error adding platform"
	var b1 = customPlatformBuilder{
		customBehavior: func(PlatformOptions, string) error {
			return fmt.Errorf(errMsg)
		},
	}
	Register("my-builder", &b1)
	err := PlatformAdd(PlatformOptions{})
	c.Assert(err, check.ErrorMatches, errMsg)
}

func (s S) TestPlatformUpdate(c *check.C) {
	platformUpdateCalls := 0
	callCounter := func(PlatformOptions, string) error {
		platformUpdateCalls++
		return nil
	}
	b1 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b2 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b3 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	Register("my-builder", &b1)
	Register("your-builder", &b2)
	Register("default-builder", &b3)
	err := PlatformUpdate(PlatformOptions{})
	c.Assert(err, check.IsNil)
	c.Assert(platformUpdateCalls, check.Equals, 3)
}

func (s S) TestPlatformUpdateError(c *check.C) {
	errMsg := "error updating platform"
	var b1 = customPlatformBuilder{
		customBehavior: func(PlatformOptions, string) error {
			return fmt.Errorf(errMsg)
		},
	}
	Register("my-builder", &b1)
	err := PlatformUpdate(PlatformOptions{})
	c.Assert(err, check.ErrorMatches, errMsg)
}

func (s S) TestPlatformRemove(c *check.C) {
	platformRemoveCalls := 0
	callCounter := func(PlatformOptions, string) error {
		platformRemoveCalls++
		return nil
	}
	b1 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b2 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	b3 := customPlatformBuilder{
		customBehavior: callCounter,
	}
	Register("my-builder", &b1)
	Register("your-builder", &b2)
	Register("default-builder", &b3)
	err := PlatformRemove("platform-name")
	c.Assert(err, check.IsNil)
	c.Assert(platformRemoveCalls, check.Equals, 3)
}

func (s S) TestPlatformRemoveError(c *check.C) {
	errMsg := "error removing platform"
	var b1 = customPlatformBuilder{
		customBehavior: func(PlatformOptions, string) error {
			return fmt.Errorf(errMsg)
		},
	}
	Register("my-builder", &b1)
	err := PlatformRemove("platform-name")
	c.Assert(err, check.ErrorMatches, errMsg)
}
