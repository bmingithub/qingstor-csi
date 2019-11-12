package qbd

import (
	"bou.ke/monkey"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yunify/qingstor-csi/pkg/common"
	"testing"
)

const configFile = "/etc/neonsan/qbd.conf"

var (
	errMock = errors.New("error mock")
)

func TestListAttachVolume(t *testing.T) {
	RegisterFailHandler(Fail)

	var cmdOut string
	var cmdError error
	var patchExecCmd *monkey.PatchGuard
	poolName, volName := "csi", "foo1"
	BeforeEach(func() {
		patchExecCmd = monkey.Patch(common.ExecCommand, func(string, []string) ([]byte, error) { return []byte(cmdOut), cmdError })
	})

	AfterEach(func() {
		patchExecCmd.Unpatch()
	})

	Describe("exe cmd error", func() {
		poolName, volName := "csi", "foo1"
		It("", func() {
			cmdOut = ""
			cmdError = errMock
			attachInfo, err := ListVolume(configFile, poolName, volName)
			Expect(err).To(HaveOccurred())
			Expect(attachInfo).To(BeNil())
		})
	})

	Describe("List attached volumes", func() {
		It("volume has attached", func() {
			cmdOut = `dev_id  vol_id  device  volume  config  read_bps    write_bps   read_iops   write_iops
0   0x3ff7000000    qbd0    csi/foo1    /etc/neonsan/qbd.conf   0   0   0   0
1   0x3a7c000000    qbd1    csi/foo /etc/neonsan/qbd.conf   0   0   0   0`
			cmdError = nil
			attachInfo, err := ListVolume(configFile, poolName, volName)
			Expect(err).To(BeNil())
			Expect(attachInfo).NotTo(BeNil())
			Expect(attachInfo.Pool).To(Equal(poolName))
			Expect(attachInfo.Name).To(Equal(volName))
		})
	})

	Describe("volume has not attached ", func() {
		It("", func() {
			cmdOut = `dev_id  vol_id  device  volume  config  read_bps    write_bps   read_iops   write_iops
0   0x3ff7000000    qbd0    csi/foo0    /etc/neonsan/qbd.conf   0   0   0   0
1   0x3a7c000000    qbd1    csi/foo /etc/neonsan/qbd.conf   0   0   0   0`
			cmdError = nil
			attachInfo, err := ListVolume(configFile, poolName, volName)
			Expect(err).NotTo(HaveOccurred())
			Expect(attachInfo).To(BeNil())
		})
	})

	Describe("volume has not attached when same name in another pool", func() {
		It("", func() {
			cmdOut = `dev_id  vol_id  device  volume  config  read_bps    write_bps   read_iops   write_iops
0   0x3ff7000000    qbd0    csi1/foo1    /etc/neonsan/qbd.conf   0   0   0   0
1   0x3a7c000000    qbd1    csi/foo /etc/neonsan/qbd.conf   0   0   0   0`
			cmdError = nil
			attachInfo, err := ListVolume(configFile, poolName, volName)
			Expect(err).NotTo(HaveOccurred())
			Expect(attachInfo).To(BeNil())
		})
	})

	Describe("volume has two attached infos", func() {
		It("", func() {
			cmdOut = `dev_id  vol_id  device  volume  config  read_bps    write_bps   read_iops   write_iops
0   0x3ff7000000    qbd0    csi/foo1    /etc/neonsan/qbd.conf   0   0   0   0
1   0x3a7c000000    qbd1    csi/foo1 /etc/neonsan/qbd.conf   0   0   0   0`
			cmdError = nil
			attachInfo, err := ListVolume(configFile, poolName, volName)
			Expect(err).To(HaveOccurred())
			Expect(attachInfo).To(BeNil())
		})
	})

	RunSpecs(t, "CSI Sanity Test Suite")
}
