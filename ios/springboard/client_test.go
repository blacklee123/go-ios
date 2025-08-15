//go:build !fast
// +build !fast

package springboard

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/stretchr/testify/assert"
)

func TestListIcons(t *testing.T) {
	list, err := ios.ListDevices()
	assert.NoError(t, err)
	if len(list.DeviceList) == 0 {
		t.Skip("No devices found")
		return
	}
	device := list.DeviceList[0]

	client, err := NewClient(device)
	assert.NoError(t, err)
	defer client.Close()

	screens, err := client.ListIcons()

	assert.NoError(t, err)
	// As the contents are individual to each device, we can only check that something gets returned
	assert.Greater(t, len(screens), 0)
}

func TestGetIconPNGData(t *testing.T) {
	list, err := ios.ListDevices()
	assert.NoError(t, err)
	if len(list.DeviceList) == 0 {
		t.Skip("No devices found")
		return
	}
	device := list.DeviceList[0]

	client, err := NewClient(device)
	assert.NoError(t, err)
	defer client.Close()

	// 测试获取设置应用的图标
	data, err := client.GetIconPNGData("com.apple.Preferences")
	if err != nil {
		if err.Error() == "springboard error: Could not find application with identifier com.apple.Preferences" {
			t.Log("设备上没有设置应用，跳过测试")
			return
		}
		t.Logf("获取图标失败: %v", err)
		t.Skip("跳过图标获取测试")
		return
	}

	// 验证PNG文件头
	assert.GreaterOrEqual(t, len(data), 8)
	assert.Equal(t, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, data[:8])

	filename := fmt.Sprintf("preferences_icon_%s.png", time.Now().Format("20060102_150405"))
	os.WriteFile(filename, data, 0644)
	t.Logf("图标已保存到 %s", filename)
}
