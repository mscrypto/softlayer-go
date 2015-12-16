package dns_resource_record

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	softlayer "github.com/maximilien/softlayer-go/softlayer"
	testhelpers "github.com/maximilien/softlayer-go/test_helpers"
)

var _ = Describe("SoftLayer DNS Resource Records", func() {
	var (
		err								 error
		dnsDomainResourceRecordService 	 softlayer.SoftLayer_Dns_Domain_Resource_Record_Service
	)

	BeforeEach(func() {
		dnsDomainResourceRecordService, err = testhelpers.CreateDnsDomainResourceRecordService()
		Expect(err).ToNot(HaveOccurred())

		testhelpers.TIMEOUT = 30 * time.Second
		testhelpers.POLLING_INTERVAL = 10 * time.Second
	})

	Context("SoftLayer_Dns_Domain_ResourceRecord", func() {
		It("creates a DNS Domain resource record, update it, and delete it", func() {
			domainId := 123456
			createdDnsDomainResourceRecord, _ := testhelpers.CreateTestDnsDomainResourceRecord(domainId)

			testhelpers.WaitForCreatedDnsDomainResourceRecordToBePresent(createdDnsDomainResourceRecord.Id)

			result, err := dnsDomainResourceRecordService.GetObject(createdDnsDomainResourceRecord.Id)
			Expect(err).ToNot(HaveOccurred())

			Expect(result.Data).To(Equal("127.0.0.1"))
			Expect(result.Host).To(Equal("test.example.com"))
			Expect(result.ResponsiblePerson).To(Equal("testemail@sl.com"))
			Expect(result.Ttl).To(Equal(900))
			Expect(result.Type).To(Equal("a"))

			oldHost := result.Host
			oldResponsiblePerson := result.ResponsiblePerson

			result.Host = "edited.test.example.com"
			result.ResponsiblePerson = "editedtestemail@sl.com"
			dnsDomainResourceRecordService.EditObject(createdDnsDomainResourceRecord.Id, result)

			result2, err := dnsDomainResourceRecordService.GetObject(createdDnsDomainResourceRecord.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(result2.Host).To(Equal("edited.test.example.com"))
			Expect(result2.ResponsiblePerson).To(Equal("editedtestemail@sl.com"))

			deleted, err := dnsDomainResourceRecordService.DeleteObject(createdDnsDomainResourceRecord.Id)
			Expect(err).ToNot(HaveOccurred())
			Expect(deleted).To(BeTrue())

			testhelpers.WaitForDeletedDnsDomainResourceRecordToNoLongerBePresent(createdDnsDomainResourceRecord.Id)
		})
	})
})

