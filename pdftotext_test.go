package pdftotext

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pdfOkRhel          = "testdata/rhel_whatsnewrhel7beta_techoverview_.pdf"
	pdfProtectedLetter = "testdata/protectedLetter.pdf"
	pdfEncrypted       = "testdata/Hello-World-password_is_password.pdf"
	notAPdf            = "model.go"
	pdfMalformed       = "testdata/malformed.pdf"
)

func TestBasicUsage(t *testing.T) {
	pdf, err := os.ReadFile(pdfOkRhel)
	assert.NoError(t, err)

	pages, err := Extract(pdf)
	assert.NoError(t, err)

	assert.Equal(t, 9, len(pages))
	assert.Equal(t, 1, pages[0].Number)
	assert.Equal(t, 2, pages[1].Number)
	assert.Equal(t, 9, pages[8].Number)

	assert.Equal(t,
		"WHAT’S NEW IN RED HAT ENTERPRISE LINUX 7 RED HAT ENTERPRISE LINUX 7 BETA TECHNOLOGY OVERVIEW Learn more about Red Hat Enterprise Linux 7 Download Red Hat Enterprise Linux 7 beta 1 and access documentation 2 in the Red Hat Customer Portal. Introduction Red Hat’s latest release of its flagship platform delivers dramatic improvements in reliability, performance, and scalability. A wealth of new features provides the architect, system administrator, and developer with the resources necessary to innovate and manage more efficiently. Architects: Red Hat ® Enterprise Linux ® 7 beta is ready for whatever infrastructure choices you make, efficiently integrating with other operating environment, authentication, and management systems. Whether your primary goal is to build network-intensive applications, massively scalable data repositories, or a build-once-deploy-often solution that performs well in physical, virtual, and cloud environments, Red Hat Enterprise Linux 7 beta has functionality to support your project. System administrators: Red Hat Enterprise Linux 7 beta has new features that help you do your job better. You’ll have better insights into what the system is doing and more controls to optimize it, with unified management tools and system-wide resource management that reduce the administrative burden. Container-based isolation and enhanced performance tools allow you to see and adjust resource allocation to each application. And, of course, there are continued improvements to scalability, reliability, and security. Developers and dev-ops: Red Hat Enterprise Linux 7 beta has more than just operating system functionality; it provides a rich application infrastructure with built-in mechanisms for security, identity management, resource allocation, and performance optimization. In addition to well-tuned default behaviors, you can take advantage of controls for application resources so you don’t leave performance up to chance. Red Hat Enterprise Linux 7 beta includes the latest stable versions of the most in-demand programming languages, databases, and runtime environments. The Linux container architecture in Red Hat Enterprise Linux 7 beta covers four technology areas: • Process isolation—Namespaces • Resource management—Cgroups Whichever role (or roles) apply to you, the Red Hat Enterprise Linux team hopes that you will find features and enhancements in Red Hat Enterprise Linux 7 beta that you want to test and try out, and then share your feedback. Linux containers Linux containers provide a method of isolating a process and simulating its environment inside a single host. It provides application sandboxing technology to run applications in a secure container environment, isolated from other applications running in the same host operating system environment. Linux containers are useful when multiple copies of an application or workload need to be run in isolation, but share environments and resources. • Security—SELinux • Management—Libvirt facebook.com/redhatinc @redhatnews linkedin.com/company/red-hat redhat.com 1 https://access.redhat.com/site/products/Red_Hat_Enterprise_Linux/Get-Beta 2 https://access.redhat.com/site/documentation/Red_Hat_Enterprise_Linux/ ",
		pages[0].Content)

	assert.True(t, strings.Contains(
		pages[4].Content,
		"OpenLMI is a common infrastructure for automating system management operations across physical and virtual deployments."))
}

func TestBasicUsageAndFailWithProblematicPdf(t *testing.T) {
	pdf, err := os.ReadFile(pdfMalformed)
	assert.NoError(t, err)

	_, err = ExtractOrError(pdf)
	assert.Error(t, err)
}

func TestAdvancedTsvFormat(t *testing.T) {
	pdf, err := os.ReadFile(pdfOkRhel)
	assert.NoError(t, err)

	tsv, err := ExtractInPopplerTsv(pdf)
	assert.NoError(t, err)
	assert.True(t, len(tsv) > 0)

	assert.Equal(t, tsv[0].Text, "###PAGE###")
	assert.Equal(t, tsv[1].Text, "###FLOW###")
	assert.Equal(t, tsv[2].Text, "###LINE###")

	firstText := tsv[3]
	assert.Equal(t, firstText.Level, 5)
	assert.Equal(t, firstText.PageNum, 1)
	assert.Equal(t, firstText.ParNum, 0)
	assert.Equal(t, firstText.BlockNum, 0)
	assert.Equal(t, firstText.LineNum, 0)
	assert.Equal(t, firstText.WordNum, 0)
	assert.Equal(t, firstText.Left, 41.64)
	assert.Equal(t, firstText.Top, 93.12)
	assert.Equal(t, firstText.Width, 70.67)
	assert.Equal(t, firstText.Height, 19.85)
	assert.Equal(t, firstText.Conf, 100)
	assert.Equal(t, firstText.Text, "WHAT’S")
}

func TestBasicUsageProtectedPdf(t *testing.T) {
	pdf, err := os.ReadFile(pdfProtectedLetter)
	assert.NoError(t, err)

	pages, err := Extract(pdf)
	assert.NoError(t, err)

	assert.True(t, len(pages) > 0)
	assert.Equal(t, 1, pages[0].Number)

	assert.Equal(t,
		"C O M PA N Y NA M E 4. June 2023 Recipient Name Company Name 123 New Street Anytown, County, Postcode Dear Recipient Name, To get started, just tap or click this placeholder text and begin typing. You can view and edit this document on your Mac, iPad, iPhone, or on iCloud.com. It’s easy to edit text, change fonts and add beautiful graphics. Use paragraph styles to get a consistent look throughout your document. For example, this paragraph uses Body style. You can change it in the Text tab of the Format controls. To add photos, movies, audio and other objects, tap or click one of the buttons in the toolbar or drag and drop the objects onto the page. Yours sincerely, Sender Name 123 High Street, Anytown, County, Postcode 01234 567 890 www.example.com ",
		pages[0].Content)
}

func TestNotAPdf(t *testing.T) {
	pdf, err := os.ReadFile(notAPdf)
	assert.NoError(t, err)

	_, err = ExtractInPopplerTsv(pdf)
	assert.Error(t, err)
}

func TestEncryptedPdf(t *testing.T) {
	// Encrypted PDFs are not supported, this is indicated by return code 1 from poppler
	pdf, err := os.ReadFile(pdfEncrypted)
	assert.NoError(t, err)

	_, err = ExtractInPopplerTsv(pdf)
	assert.Error(t, err)
}
